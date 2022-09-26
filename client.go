// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "fmt"
import "net"
import "time"

// this struct handle client connexion.
type Client struct {
	buffer bufferIO
	Macros []*Macro
	do_close bool
}

// This function process message as expected "Accept/reject action"
// it is used when the client receive server message which should be
// an action.
func AnswerToAction(msgType MsgType, value interface{})(*Action, error) {
	switch msgType {
	case SMFIR_ACCEPT,
	     SMFIR_CONTINUE,
	     SMFIR_DISCARD,
	     SMFIR_REJECT,
	     SMFIR_TEMPFAIL:
		return &Action{Action: ActionCode(msgType)}, nil

	case SMFIR_REPLYCODE:
		return &Action{
			Action: ActionCode(msgType),
			Value: value.(*MsgReply),
		}, nil
	default:
		return nil, fmt.Errorf("protocol error: unexpected command %q", msgType.String())
	}
}

// This function process message as modification message.
// It is used when client receive server message which should be
// a modification.
func AnswerToModification(msgType MsgType, value interface{})(*Modification, error) {

	switch msgType {

	// Process modification.
	case SMFIR_ADDRCPT,
	     SMFIR_DELRCPT,
	     SMFIR_QUARANTINE:
		return &Modification{
			Modification: ModificationCode(msgType),
			Value: value.(string),
		}, nil

	case SMFIR_REPLBODY:
		return &Modification{
			Modification: ModificationCode(msgType),
			Value: value.([]byte),
		}, nil

	case SMFIR_ADDHEADER:
		return &Modification{
			Modification: ModificationCode(msgType),
			Value: value.(*MsgAddHeader),
		}, nil

	case SMFIR_CHGHEADER:
		return &Modification{
			Modification: ModificationCode(msgType),
			Value: value.(*MsgChgHeader).Index,
		}, nil

	default:
		return nil, fmt.Errorf("protocol error: unexpected command %q", msgType.String())
	}
}

// This function set macro. The macro must be added before called function
// associated with corresponding step, otherwise the macro will not sent.
// step is a constant MS_*. name is the macro name avec value is the macro
// value. Note that Sendmail and Postfix defines standard macros associated
// with steps. These macro are:
//
// ▶︎ SMFIC_CONNECT : _ j {daemon_name} {if_name} {if_addr}
//
// ▶︎ SMFIC_HELO : {tls_version} {cipher} {cipher_bits cert_subject} {cert_issuer}
//
// ▶︎ SMFIC_MAIL : i {auth_type} {auth_authen} {auth_ssf} {auth_author} {mail_mailer} {mail_host} {mail_addr}
//
// ▶︎ SMFIC_RCPT : {rcpt_mailer} {rcpt_host} {rcpt_addr}
func (cli *Client)MacroAdd(step MacroStep, name string, value string)() {
	macroAdd(&cli.Macros, step, name, value)
}

// This function perform a lookup in the macro container. It returns macro value
// or empty string if none is found
func (cli *Client)MacroGet(name string)(MacroStep, string) {
	return macroGet(cli.Macros, name)
}

// This function helps to debug macro container. It dumps macro content on stdout.
func (cli *Client)MacroDebug()() {
	macroDebug(cli.Macros)
}

// Add macro "_". See MacroAdd for more information.
func (cli *Client)MacroAdd__(value string)            { macroAdd__(&cli.Macros, value) }
// Add macro "j". See MacroAdd for more information.
func (cli *Client)MacroAdd_j(value string)            { macroAdd_j(&cli.Macros, value) }
// Add macro "{daemon_name}". See MacroAdd for more information.
func (cli *Client)MacroAdd_daemon_name(value string)  { macroAdd_daemon_name(&cli.Macros, value) }
// Add macro "{if_name}". See MacroAdd for more information.
func (cli *Client)MacroAdd_if_name(value string)      { macroAdd_if_name(&cli.Macros, value) }
// Add macro "{if_addr}". See MacroAdd for more information.
func (cli *Client)MacroAdd_if_addr(value string)      { macroAdd_if_addr(&cli.Macros, value) }

// Add macro "{tls_version}". See MacroAdd for more information.
func (cli *Client)MacroAdd_tls_version(value string)  { macroAdd_tls_version(&cli.Macros, value) }
// Add macro "{cipher}". See MacroAdd for more information.
func (cli *Client)MacroAdd_cipher(value string)       { macroAdd_cipher(&cli.Macros, value) }
// Add macro "{cipher_bits}". See MacroAdd for more information.
func (cli *Client)MacroAdd_cipher_bits(value string)  { macroAdd_cipher_bits(&cli.Macros, value) }
// Add macro "{cert_subject}". See MacroAdd for more information.
func (cli *Client)MacroAdd_cert_subject(value string) { macroAdd_cert_subject(&cli.Macros, value) }
// Add macro "{cert_issuer}". See MacroAdd for more information.
func (cli *Client)MacroAdd_cert_issuer(value string)  { macroAdd_cert_issuer(&cli.Macros, value) }

// Add macro "i". See MacroAdd for more information.
func (cli *Client)MacroAdd_i(value string)            { macroAdd_i(&cli.Macros, value) }
// Add macro "{auth_type}". See MacroAdd for more information.
func (cli *Client)MacroAdd_auth_type(value string)    { macroAdd_auth_type(&cli.Macros, value) }
// Add macro "{auth_authen}". See MacroAdd for more information.
func (cli *Client)MacroAdd_auth_authen(value string)  { macroAdd_auth_authen(&cli.Macros, value) }
// Add macro "{auth_ssf}". See MacroAdd for more information.
func (cli *Client)MacroAdd_auth_ssf(value string)     { macroAdd_auth_ssf(&cli.Macros, value) }
// Add macro "{auth_author}". See MacroAdd for more information.
func (cli *Client)MacroAdd_auth_author(value string)  { macroAdd_auth_author(&cli.Macros, value) }
// Add macro "{mail_mailer}". See MacroAdd for more information.
func (cli *Client)MacroAdd_mail_mailer(value string)  { macroAdd_mail_mailer(&cli.Macros, value) }
// Add macro "{mail_host}". See MacroAdd for more information.
func (cli *Client)MacroAdd_mail_host(value string)    { macroAdd_mail_host(&cli.Macros, value) }
// Add macro "{mail_addr}". See MacroAdd for more information.
func (cli *Client)MacroAdd_mail_addr(value string)    { macroAdd_mail_addr(&cli.Macros, value) }

// Add macro "{rcpt_mailer}". See MacroAdd for more information.
func (cli *Client)MacroAdd_rcpt_mailer(value string)  { macroAdd_rcpt_mailer(&cli.Macros, value) }
// Add macro "{rcpt_host}". See MacroAdd for more information.
func (cli *Client)MacroAdd_rcpt_host(value string)    { macroAdd_rcpt_host(&cli.Macros, value) }
// Add macro "{rcpt_addr}". See MacroAdd for more information.
func (cli *Client)MacroAdd_rcpt_addr(value string)    { macroAdd_rcpt_addr(&cli.Macros, value) }

// This function use connection to milter server defined in conn. It returns
// a *Client on success and bever fails. Note, the caller must close the
// connexion once its no longer used.
func ClientNewFromConn(conn net.Conn)(*Client) {
	var cli *Client

	// Create client struct
	cli = &Client{}

	// Declare connection
	cli.buffer.InitBufferIO(conn)

	return cli
}

// This function connects to milter server using proto (like "tcp"), adress
// (like "localhost:4567") and timeout in seconds. It returns a *Client
// on success or fill error on error cases.
func ClientNew(proto string, addr string, timeout int)(*Client, error) {
	var err error
	var conn net.Conn
	var cli *Client

	// Open connection
	conn, err = net.DialTimeout(proto, addr, time.Duration(timeout) * time.Second)
	if err != nil {
		return nil, err
	}

	// Create new connection
	cli = ClientNewFromConn(conn)

	// Set flag do close to indicate to the the close function its behavior
	cli.do_close = true

	return cli, nil
}

// Terminate client connexion. The connexion ois closed if the connection
// was established using ClientNew function, otherwise, the caller should
// close the connexion
func (cli *Client)Close()(error) {
	if !cli.do_close {
		return nil
	}
	return cli.buffer.Close()
}

// This function returns milter byte ready to be decoded
// If error is filled, the connexion should be close and processing aborted
func (cli *Client)ReceivePacket()([]byte, error) {
	return cli.buffer.ReceivePacket()
}

// This function return next decoded Milter message. See documentation of
// Decode function to understand cast between MsgType and interface{}.
// If error is filled, the connexion should be close and processing aborted
func (cli *Client)ReceiveMessage()(MsgType, interface{}, error) {
	var msg []byte
	var err error

	msg, err = cli.buffer.ReceivePacket()
	if err != nil {
		return SMFIR_ERROR, nil, err
	}

	return Decode(msg)
}

// Client send message to quit milter communication. The server do not
// answer anything. The client should free milter protocol handler using
// cli.Close(). If the client established connection, the connection is
// closed. If the connection was establish by the caller, the caller
// shoul close the connexion.
func (cli *Client)SendQuit()(error) {
	return cli.buffer.Write(EncodeQuit())
}

// Client send message to milter to abort current filter checks. The connection
// is reset to the HELO state.
func (cli *Client)SendAbort()(error) {
	return cli.buffer.Write(EncodeAbort())
}

// Client send its protocol and modifications options and get the milter server
// requirement as return. actions is "or" between SMFIF_* constants and protocol
// is "or" between SMFIP_* constants.
func (cli *Client)SendOptNeg(optNeg *MsgOptNeg)(error) {
	return cli.buffer.Write(EncodeOptNeg(optNeg))
}

// Client send MACRO message which inform milter server about MACRO
// names and value.
func (cli *Client)SendMacro(step MacroStep, macros []*Macro)(error) {
	return cli.buffer.Write(EncodeMacro(step, macros))
}

// Client send CONNECT message which inform milter server about CONNNECT
// informations. hostname is the hostname as string if known. port is the
// client port and adress is IPv4 or IPv6 client address. Note the unix
// socket are not yet supported. If an error occurs, error is filled,
// otherwise it is nil.
func (cli *Client)SendConnect(connect *MsgConnect)(error) {
	return cli.buffer.Write(EncodeConnect(connect, cli.Macros))
}

// This function send SMTP HELO information to the milter server. HELO is just
// one string. If an error occurs, error is filled, otherwise it is nil.
func (cli *Client)SendHelo(helo string)(error) {
	return cli.buffer.Write(EncodeHelo(helo, cli.Macros))
}

// This function send the SMTP MAIL FROM command content. Its juste on string.
// If an error occurs, error is filled, otherwise it is nil.
func (cli *Client)SendMail(email *MsgMail)(error) {
	return cli.buffer.Write(EncodeMail(email, cli.Macros))
}

// This function send the SMTP RCPT TO command content. Its juste on string.
// If an error occurs, error is filled, otherwise it is nil.
func (cli *Client)SendRcpt(email *MsgMail)(error) {
	return cli.buffer.Write(EncodeRcpt(email, cli.Macros))
}

// The client send header contained in the email. This function should call one
// time per header. Its important to send email using encountered order because
// the modification function "change header" gives and index of the header to
// be modified. If an error occurs, error is filled, otherwise it is nil.
func (cli *Client)SendHeader(hdr *MsgHeader)(error) {
	return cli.buffer.Write(EncodeHeader(hdr))
}

// this message indicated to the milter server the end of headers.
// If an error occurs, error is filled, otherwise it is nil.
func (cli *Client)SendEOH()(error) {
	return cli.buffer.Write(EncodeEOH())
}

// the client send to milter server the body using chunks of 65535 bytes. This
// function must be called more than one time. If an error occurs, error is
// filled, otherwise it is nil.
func (cli *Client)SendBody(body []byte)(error) {
	return cli.buffer.Write(EncodeBody(body))
}

// This function indicated the end of body to the milter server. If an error
// occurs, error is filled, otherwise it is nil.
func (cli *Client)SendBodyEOB()(error) {
	return cli.buffer.Write(EncodeBodyEOB())
}

// Client send message to quit milter communication. The server do not
// answer anything. The client should free milter protocol handler using
// cli.Close(). If the client established connection, the connection is
// closed. If the connection was establish by the caller, the caller
// shoul close the connexion.
func (cli *Client)ExchangeQuit()(error) {
	return cli.buffer.Write(EncodeQuit())
}

// Client send message to milter to abort current filter checks. The connection
// is reset to the HELO state. The server do not answer anything.
func (cli *Client)ExchangeAbort()(error) {
	return cli.buffer.Write(EncodeAbort())
}

// Client send its protocol and modifications options and get the milter server
// requirement as return. actions is "or" between SMFIF_* constants and protocol
// is "or" between SMFIP_* constants. The function waits for server answer.
func (cli *Client)ExchangeOptNeg(optNeg *MsgOptNeg)(*MsgOptNeg, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	msg = EncodeOptNeg(optNeg)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	if msgType != SMFIC_OPTNEG {
		return nil, fmt.Errorf("protocol error: expect SMFIC_OPTNEG message, got %q", msgType.String())
	}

	return value.(*MsgOptNeg), nil
}

// Client send CONNECT message which inform milter server about CONNNECT
// informations. hostname is the hostname as string if known. port is the
// client port and adress is IPv4 or IPv6 client address. Note the unix
// socket are not yet supported. The milter answer an *Action. If an error
// occurs, error is filled, otherwise it is nil. The function waits for
// server answer.
func (cli *Client)ExchangeConnect(connect *MsgConnect)(*Action, error) {
	var msg []byte
	var msgType MsgType
	var value interface{}
	var err error

	// Make buffer with connect payload
	msg = EncodeConnect(connect, cli.Macros)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// This function send SMTP HELO information to the milter server. HELO is just
// one string. The milter answer an *Action. If an error occurs, error is
// filled, otherwise it is nil. The function waits for server answer.
func (cli *Client)ExchangeHelo(helo string)(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// Encode message
	msg = EncodeHelo(helo, cli.Macros)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// This function send the SMTP MAIL FROM command content. Its juste on string.
// The milter answer an *Action. If an error occurs, error is filled, otherwise
// it is nil. The function waits for server answer.
func (cli *Client)ExchangeMail(email *MsgMail)(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// encode data
	msg = EncodeMail(email, cli.Macros)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// This function send the SMTP RCPT TO command content. Its juste on string.
// The milter answer an *Action. If an error occurs, error is filled, otherwise
// it is nil. The function waits for server answer.
func (cli *Client)ExchangeRcpt(email *MsgMail)(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// encode data
	msg = EncodeRcpt(email, cli.Macros)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// The client send header contained in the email. This function should call one
// time per header. Its important to send email using encountered order because
// the modification function "change header" gives and index of the header to
// be modified. The milter answer an *Action. If an error occurs, error is
// filled, otherwise it is nil. The function waits for server answer.
func (cli *Client)ExchangeHeader(hdr *MsgHeader)(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// Encode data
	msg = EncodeHeader(hdr)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// this message indicated to the milter server the end of headers. The milter
// answer an *Action. If an error occurs, error is filled, otherwise it is nil.
// The function waits for server answer.
func (cli *Client)ExchangeEOH()(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// encode data
	msg = EncodeEOH()

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// the client send to milter server the body using chunks of 65535 bytes. This
// function must be called more than one time. The milter answer an *Action. If
// an error occurs, error is filled, otherwise it is nil. The function waits for
// server answer.
func (cli *Client)ExchangeBody(body []byte)(*Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}

	// encode data
	msg = EncodeBody(body)

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, err
	}

	// Read response and decode it
	msgType, value, err = cli.ReceiveMessage()
	if err != nil {
		return nil, err
	}

	return AnswerToAction(msgType, value)
}

// This function indicated the end of body to the milter server. The server could
// answer with a list of modification and an action. The list of modification
// could be empty. The milter answer an *Action. If an error occurs, error is
// filled, otherwise it is nil. The function waits for server answer.
func (cli *Client)ExchangeBodyEOB()([]*Modification, *Action, error) {
	var msg []byte
	var err error
	var msgType MsgType
	var value interface{}
	var mods []*Modification
	var action *Action
	var modification *Modification

	// encode data
	msg = EncodeBodyEOB()

	// Send packet
	err = cli.buffer.Write(msg)
	if err != nil {
		return nil, nil, err
	}

	// Read all responses until accept/reject action
	for {

		// Read next response and decode it
		msgType, value, err = cli.ReceiveMessage()
		if err != nil {
			return nil, nil, err
		}

		// Process modification or action
		modification, err = AnswerToModification(msgType, value)
		if err == nil {
			mods = append(mods, modification)
		} else {
			action, err = AnswerToAction(msgType, value)
			return mods, action, err
		}
	}
}
