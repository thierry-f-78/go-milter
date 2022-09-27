// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "fmt"
import "net"

// Server Callbacks interface are used with Exchange() function.
// Each callback is called when the client send corresponding message.
// Each function which implements this interface must return the
// expected answer. In most cases is simply an Action. The OPTNEG
// step expect OPTNEG message and the BODYEOB step could take
// a list of Modifications. The callback OnERROR is not handled
// by the protocol, but it is called if an error occurs.
type ServerCallbacks interface {
	OnOPTNEG(*Server, *MsgOptNeg)(*MsgOptNeg, error)
	OnCONNECT(*Server, *MsgConnect)(*Action, error)
	OnHELO(*Server, string)(*Action, error)
	OnMAIL(*Server, *MsgMail)(*Action, error)
	OnRCPT(*Server, *MsgMail)(*Action, error)
	OnHEADER(*Server, *MsgHeader)(*Action, error)
	OnEOH(*Server)(*Action, error)
	OnBODY(*Server, []byte)(*Action, error)
	OnBODYEOB(*Server)([]*Modification, *Action, error)
	OnABORT(*Server)(error)
	OnQUIT(*Server)(error)
	OnERROR(*Server, error)
}

// This struct contains server things like Macros. It allow
// communication with client in Send*/Receive* mode.
type Server struct {
	buffer bufferIO
	Macros []*Macro
}

// Create new server based on network connection.
func ServerNew(conn net.Conn)(*Server) {
	var srv *Server

	/* Init new server */
	srv = &Server{}
	srv.Macros = nil
	srv.buffer.InitBufferIO(conn)

	return srv
}

// This function returns milter byte ready to be decoded
// If error is filled, the connexion should be close and processing aborted
func (srv *Server)ReceivePacket()([]byte, error) {
	return srv.buffer.ReceivePacket()
}

// This function return next decoded Milter message. See documentation of
// Decode function to understand cast between MsgType and interface{}.
// If error is filled, the connexion should be close and processing aborted
func (srv *Server)ReceiveMessage()(MsgType, interface{}, error) {
	var msg []byte
	var err error

	msg, err = srv.buffer.ReceivePacket()
	if err != nil {
		return SMFIR_ERROR, nil, err
	}

	return Decode(msg)
}

// This function is called to handle new server request. "inst" is a variable
// which implements InstanceCallbacks. The defined callbacks are called each
// time is necessary.
//
// If the fucntion returns 0, end of connection is required, the caller should
// close connection. If the function returns 1, the connection should be keep
// opened and a new request could arrive.
func Exchange(conn net.Conn, inst ServerCallbacks) {
	var srv *Server
	var msgType MsgType
	var msg interface{}
	var err error
	var macros []*Macro
	var m *Macro
	var optNeg *MsgOptNeg
	var modification *Modification
	var modifications []*Modification
	var action *Action

	/* Init new server */
	srv = ServerNew(conn)
	srv.Macros = nil

	/* Read first message, expect Optneg */
	msgType, msg, err = srv.ReceiveMessage()
	if err != nil {
		inst.OnERROR(srv, err)
		return
	}

	/* Expect negociation */
	if msgType != SMFIC_OPTNEG {
		inst.OnERROR(srv, fmt.Errorf("protocol error: expect message SMFIC_OPTNEG, got %s", msgType.String()))
		return
	}

	// Call OptNeg callback
	optNeg, err = inst.OnOPTNEG(srv, msg.(*MsgOptNeg))
	if err != nil {
		inst.OnERROR(srv, err)
		return
	}

	err = srv.buffer.Write(EncodeOptNeg(optNeg))
	if err != nil {
		inst.OnERROR(srv, err)
		return
	}

	for {

		// Read next message
		msgType, msg, err = srv.ReceiveMessage()
		if err != nil {
			inst.OnERROR(srv, err)
			return
		}

		// Call the right callback according with received message
		switch msgType {
		case SMFIC_CONNECT:

			action, err = inst.OnCONNECT(srv, msg.(*MsgConnect))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_MACRO:

			macros = msg.([]*Macro)
			for _, m = range macros {
				macroAdd(&srv.Macros, m.Step, m.Name, m.Value)
			}

		case SMFIC_HELO:

			action, err = inst.OnHELO(srv, msg.(string))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_MAIL:

			action, err = inst.OnMAIL(srv, msg.(*MsgMail))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_RCPT:

			action, err = inst.OnRCPT(srv, msg.(*MsgMail))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_HEADER:

			action, err = inst.OnHEADER(srv, msg.(*MsgHeader))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_EOH:

			action, err = inst.OnEOH(srv)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_BODY:

			action, err = inst.OnBODY(srv, msg.([]byte))
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

		case SMFIC_BODYEOB:

			modifications, action, err = inst.OnBODYEOB(srv)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			for _, modification = range modifications {
				err = srv.SendModification(modification)
				if err != nil {
					inst.OnERROR(srv, err)
					return
				}
			}

			err = srv.SendAction(action)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			/* could proces other message */
			srv.Macros = nil

		case SMFIC_ABORT:

			err = inst.OnABORT(srv)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}

			/* Abort command must reset transaction to the step HELO */
			srv.Macros = nil

		case SMFIC_QUIT:

			err = inst.OnQUIT(srv)
			if err != nil {
				inst.OnERROR(srv, err)
				return
			}
			return

		default:
			inst.OnERROR(srv, fmt.Errorf("receive unknown response code %q: %s", string(byte(msgType)), msgType.String()))
			return
		}
	}
}

// This function perform a lookup in the macro container. It returns macro value
// or empty string if none is found
func (srv *Server)MacroGet(name string)(MacroStep, string) {
	return macroGet(srv.Macros, name)
}

// Display Macro summary on stdout.
func (srv *Server)MacroDebug()() {
	macroDebug(srv.Macros)
}

// Send OPTNEG message
func (srv *Server)SendOptNeg(optNeg *MsgOptNeg)(error) {
	return srv.buffer.Write(EncodeOptNeg(optNeg))
}

// Send PROGRESS message. This message is used to maintain network
// connexion alive.
func (srv *Server)SendProgress()(error) {
	return srv.buffer.Write(EncodeProgress())
}

// Send ADDRCPT modification message
func (srv *Server)ModificationAddRcpt(rcpt string)(error) {
	return srv.buffer.Write(EncodeAddRcpt(rcpt))
}

// Send DELRCP modification message
func (srv *Server)ModificationDelRcpt(rcpt string)(error) {
	return srv.buffer.Write(EncodeDelRcpt(rcpt))
}

// Send REPLBODY modification message
func (srv *Server)ModificationReplBody(body []byte)(error) {
	return srv.buffer.Write(EncodeReplBody(body))
}

// Send ADDHEADER modification message
func (srv *Server)ModificationAddHeader(addhdr *MsgAddHeader)(error) {
	return srv.buffer.Write(EncodeAddHeader(addhdr))
}

// Send CHGHEADER modification message. Note if the header content
// is empty string, the header is removed.
func (srv *Server)ModificationChgHeader(chghdr *MsgChgHeader)(error) {
	return srv.buffer.Write(EncodeChgHeader(chghdr))
}

// Send QUARANTINE modification message.
func (srv *Server)ModificationQuarantine(reason string)(error) {
	return srv.buffer.Write(EncodeQuarantine(reason))
}

// Send ACCEPT action
func (srv *Server)ActionAccept()(error) {
	return srv.buffer.Write(EncodeAccept())
}

// Send CONTINUE action
func (srv *Server)ActionContinue()(error) {
	return srv.buffer.Write(EncodeContinue())
}

// Send DISCARD action
func (srv *Server)ActionDiscard()(error) {
	return srv.buffer.Write(EncodeDiscard())
}

// Send REJECT action
func (srv *Server)ActionReject()(error) {
	return srv.buffer.Write(EncodeReject())
}

// Send TEMPFAIL action
func (srv *Server)ActionTempfail()(error) {
	return srv.buffer.Write(EncodeTempfail())
}

// Send REPLYCODE action
func (srv *Server)ActionReplyCode(reply *MsgReply)(error) {
	return srv.buffer.Write(EncodeReplyCode(reply))
}

// Build ACCEPT struct for Exchange API
func ActionAccept()(*Action) {
	return &Action{Action: AC_ACCEPT}
}

// Build CONTINUE struct for Exchange API
func ActionContinue()(*Action) {
	return &Action{Action: AC_CONTINUE}
}

// Build DISCARD struct for Exchange API
func ActionDiscard()(*Action) {
	return &Action{Action: AC_DISCARD}
}

// Build REJECT struct for Exchange API
func ActionReject()(*Action) {
	return &Action{Action: AC_REJECT}
}

// Build TEMPSFAIL struct for Exchange API
func ActionTempfail()(*Action) {
	return &Action{Action: AC_TEMPFAIL}
}

// Build REPLYCODE struct for Exchange API
func ActionReplyCode(code int, reason string)(*Action) {
	return &Action{Action: AC_REPLYCODE, Value: &MsgReply{Code: code, Reason: reason}}
}

// This function send generic *Modification
func (srv *Server)SendModification(modification *Modification)(error) {
	switch modification.Modification {
	case MC_ADDRCPT:    return srv.ModificationAddRcpt(modification.Value.(string))
	case MC_DELRCPT:    return srv.ModificationDelRcpt(modification.Value.(string))
	case MC_REPLBODY:   return srv.ModificationReplBody(modification.Value.([]byte))
	case MC_ADDHEADER:  return srv.ModificationAddHeader(modification.Value.(*MsgAddHeader))
	case MC_CHGHEADER:  return srv.ModificationChgHeader(modification.Value.(*MsgChgHeader))
	case MC_QUARANTINE: return srv.ModificationQuarantine(modification.Value.(string))
	default:            return fmt.Errorf("Unknwon modification %q", modification.Modification)
	}
}

// This function send generic *Action
func (srv *Server)SendAction(action *Action)(error) {
	switch action.Action {
	case AC_ACCEPT:    return srv.ActionAccept()
	case AC_CONTINUE:  return srv.ActionContinue()
	case AC_DISCARD:   return srv.ActionDiscard()
	case AC_REJECT:    return srv.ActionReject()
	case AC_TEMPFAIL:  return srv.ActionTempfail()
	case AC_REPLYCODE: return srv.ActionReplyCode(action.Value)
	default:           return fmt.Errorf("Unknwon action %q", action.Action)
	}
}
