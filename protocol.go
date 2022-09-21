// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "encoding/binary"
import "fmt"
import "strconv"

const MilterVersion = 2

// Define constant for each milter message, note the constant is a byte, this
// byte is exactly the byte used by the milter protocol.
type MsgType byte
const (
	SMFIC_ABORT MsgType = 'A'
	SMFIC_BODY MsgType = 'B'
	SMFIC_CONNECT MsgType = MsgType(MS_CONNECT)
	SMFIC_MACRO MsgType = 'D'
	SMFIC_BODYEOB MsgType = 'E'
	SMFIC_HELO MsgType = MsgType(MS_HELO)
	SMFIC_HEADER MsgType = 'L'
	SMFIC_MAIL MsgType = MsgType(MS_MAIL)
	SMFIC_EOH MsgType = 'N'
	SMFIC_OPTNEG MsgType = 'O'
	SMFIC_RCPT MsgType = MsgType(MS_RCPT)
	SMFIC_QUIT MsgType = 'Q'

	SMFIR_ADDRCPT MsgType = MsgType(MC_ADDRCPT)
	SMFIR_DELRCPT MsgType = MsgType(MC_DELRCPT)
	SMFIR_ACCEPT MsgType = MsgType(AC_ACCEPT)
	SMFIR_REPLBODY MsgType = MsgType(MC_REPLBODY)
	SMFIR_CONTINUE MsgType = MsgType(AC_CONTINUE)
	SMFIR_DISCARD MsgType = MsgType(AC_DISCARD)
	SMFIR_ADDHEADER MsgType = MsgType(MC_ADDHEADER)
	SMFIR_CHGHEADER MsgType = MsgType(MC_CHGHEADER)
	SMFIR_PROGRESS MsgType = 'p'
	SMFIR_QUARANTINE MsgType = MsgType(MC_QUARANTINE)
	SMFIR_REJECT MsgType = MsgType(AC_REJECT)
	SMFIR_TEMPFAIL MsgType = MsgType(AC_TEMPFAIL)
	SMFIR_REPLYCODE MsgType = MsgType(AC_REPLYCODE)

	SMFIR_ERROR MsgType = 0xff
)

func (b *MsgType)String()(string) {
	switch *b {
	case SMFIC_ABORT:      return "ABORT"
	case SMFIC_BODY:       return "BODY"
	case SMFIC_CONNECT:    return "CONNECT"
	case SMFIC_MACRO:      return "MACRO"
	case SMFIC_BODYEOB:    return "BODYEOB"
	case SMFIC_HELO:       return "HELO"
	case SMFIC_HEADER:     return "HEADER"
	case SMFIC_MAIL:       return "MAIL"
	case SMFIC_EOH:        return "EOH"
	case SMFIC_OPTNEG:     return "OPTNEG"
	case SMFIC_RCPT:       return "RCPT"
	case SMFIC_QUIT:       return "QUIT"
	case SMFIR_ADDRCPT:    return "ADDRCPT"
	case SMFIR_DELRCPT:    return "DELRCPT"
	case SMFIR_ACCEPT:     return "ACCEPT"
	case SMFIR_REPLBODY:   return "REPLBODY"
	case SMFIR_CONTINUE:   return "CONTINUE"
	case SMFIR_DISCARD:    return "DISCARD"
	case SMFIR_ADDHEADER:  return "ADDHEADER"
	case SMFIR_CHGHEADER:  return "CHGHEADER"
	case SMFIR_PROGRESS:   return "PROGRESS"
	case SMFIR_QUARANTINE: return "QUARANTINE"
	case SMFIR_REJECT:     return "REJECT"
	case SMFIR_TEMPFAIL:   return "TEMPFAIL"
	case SMFIR_REPLYCODE:  return "REPLYCODE"
	}
	return fmt.Sprintf("UNKNOWN[%02x]", byte(*b))
}

func toMsgType(b byte)(MsgType) {
	switch b {
	case 'A': return SMFIC_ABORT
	case 'B': return SMFIC_BODY
	case 'C': return SMFIC_CONNECT
	case 'D': return SMFIC_MACRO
	case 'E': return SMFIC_BODYEOB
	case 'H': return SMFIC_HELO
	case 'L': return SMFIC_HEADER
	case 'M': return SMFIC_MAIL
	case 'N': return SMFIC_EOH
	case 'O': return SMFIC_OPTNEG
	case 'R': return SMFIC_RCPT
	case 'Q': return SMFIC_QUIT
	case '+': return SMFIR_ADDRCPT
	case '-': return SMFIR_DELRCPT
	case 'a': return SMFIR_ACCEPT
	case 'b': return SMFIR_REPLBODY
	case 'c': return SMFIR_CONTINUE
	case 'd': return SMFIR_DISCARD
	case 'h': return SMFIR_ADDHEADER
	case 'm': return SMFIR_CHGHEADER
	case 'p': return SMFIR_PROGRESS
	case 'q': return SMFIR_QUARANTINE
	case 'r': return SMFIR_REJECT
	case 't': return SMFIR_TEMPFAIL
	case 'y': return SMFIR_REPLYCODE
	}
	return SMFIR_ERROR
}

type MacroStep byte

// Macros Sending steps. these steps are defined whith each macro in order
// to send the macro as appropriate step.
const (
	MS_CONNECT MacroStep = MacroStep('C')
	MS_HELO MacroStep = MacroStep('H')
	MS_MAIL MacroStep = MacroStep('M')
	MS_RCPT MacroStep = MacroStep('R')
)

// Accept any byte as macro step
func toMacroStep(b byte)(MacroStep) {
	return MacroStep(b)
}

type FamilyCode byte

// SMFIA means Sendmail Milter Functions Internal Actions. These constants
// defines which type of connexion is used between email sender and the MTA.
//
// ▶︎ SMFIA_UNKNOWN : Connexion type is unknown
//
// ▶︎ SMFIA_UNIX : connexion type is unix socket
//
// ▶︎ SMFIA_INET : connexion type is IPv4
//
// ▶︎ SMFIA_INET6 : connexion type is IPv6
const (
	SMFIA_UNKNOWN FamilyCode = FamilyCode('U')
	SMFIA_UNIX FamilyCode = FamilyCode('L')
	SMFIA_INET FamilyCode = FamilyCode('4')
	SMFIA_INET6 FamilyCode = FamilyCode('6')

	SMFIA_ERROR FamilyCode = FamilyCode(0xff)
)

func (n *FamilyCode)String()(string) {
	switch *n {
	case SMFIA_UNKNOWN: return "unknown"
	case SMFIA_UNIX: return "unix"
	case SMFIA_INET: return "ipv4"
	case SMFIA_INET6: return "ipv6"
	}
	return "unknown"
}

func toFamily(b byte)(FamilyCode) {
	switch b {
	case byte(SMFIA_UNKNOWN): return SMFIA_UNKNOWN
	case byte(SMFIA_UNIX): return SMFIA_UNIX
	case byte(SMFIA_INET): return SMFIA_INET
	case byte(SMFIA_INET6): return SMFIA_INET6
	}
	return SMFIA_ERROR
}

type ActionFlag uint32

// SMFIF means Sendmail Milter Functions Interceptor Flags. These constants are
// used to negaciate which modification actions are allowed by the MTA and
// required by the milter.
//
// ▶︎ SMFIF_ADDHDRS : MTA offer or milter wants to add email headers
//
// ▶︎ SMFIF_CHGBODY : MTA offer or milter wants to modify email body
//
// ▶︎ SMFIF_ADDRCPT : MTA offer or milter wants to add recipient
//
// ▶︎ SMFIF_DELRCPT : MTA offer or milter wants to delete recipient
//
// ▶︎ SMFIF_CHGHDRS : MTA offer or milter wants to modify email header
//
// ▶︎ SMFIF_QUARANTINE : MTA offer or milter wants to drop email in quarantine
//
// ▶︎ SMFIF_ALL : MTA offer or milter wants all the modification actions
const (
	SMFIF_ADDHDRS ActionFlag = ActionFlag(0x01)
	SMFIF_CHGBODY ActionFlag = ActionFlag(0x02)
	SMFIF_ADDRCPT ActionFlag = ActionFlag(0x04)
	SMFIF_DELRCPT ActionFlag = ActionFlag(0x08)
	SMFIF_CHGHDRS ActionFlag = ActionFlag(0x10)
	SMFIF_QUARANTINE ActionFlag = ActionFlag(0x20)
	SMFIF_ALL ActionFlag = SMFIF_ADDHDRS | SMFIF_CHGBODY | SMFIF_ADDRCPT | SMFIF_DELRCPT | SMFIF_CHGHDRS | SMFIF_QUARANTINE
)

type ProtocolFlag uint32

// SMFIP means Sendmail Milter Functions Protocol Flags. These constants are
// used to negociate which events are not offered by the MTA and which are not
// required by the milter.
//
// ▶︎ SMFIP_NOCONNECT : MTA not offer or milter don't want CONNECT message
//
// ▶︎ SMFIP_NOHELO : MTA not offer or milter don't want HELO message
//
// ▶︎ SMFIP_NOMAIL : MTA not offer or milter don't want MAIL (MAIL FROM) message
//
// ▶︎ SMFIP_NORCPT : MTA not offer or milter don't want RCPT (RCPT TO) message
//
// ▶︎ SMFIP_NOBODY : MTA not offer or milter don't want BODY nor BODYEOB messages
//
// ▶︎ SMFIP_NOHDRS : MTA not offer or milter don't want HEADER message
//
// ▶︎ SMFIP_NOEOH : MTA not offer or milter don't want EOH messages
const (
	SMFIP_NOCONNECT ProtocolFlag = ProtocolFlag(0x01)
	SMFIP_NOHELO ProtocolFlag = ProtocolFlag(0x02)
	SMFIP_NOMAIL ProtocolFlag = ProtocolFlag(0x04)
	SMFIP_NORCPT ProtocolFlag = ProtocolFlag(0x08)
	SMFIP_NOBODY ProtocolFlag = ProtocolFlag(0x10)
	SMFIP_NOHDRS ProtocolFlag = ProtocolFlag(0x20)
	SMFIP_NOEOH ProtocolFlag = ProtocolFlag(0x40)
)

// Define milter actions
type ActionCode byte
const (
	AC_ACCEPT ActionCode    = ActionCode('a')
	AC_CONTINUE ActionCode  = ActionCode('c')
	AC_DISCARD ActionCode   = ActionCode('d')
	AC_REJECT ActionCode    = ActionCode('r')
	AC_TEMPFAIL ActionCode  = ActionCode('t')
	AC_REPLYCODE ActionCode = ActionCode('y')
)

// Define milter modifications
type ModificationCode byte
const (
	MC_ADDRCPT ModificationCode    = ModificationCode('+')
	MC_DELRCPT ModificationCode    = ModificationCode('-')
	MC_REPLBODY ModificationCode   = ModificationCode('b')
	MC_ADDHEADER ModificationCode  = ModificationCode('h')
	MC_CHGHEADER ModificationCode  = ModificationCode('m')
	MC_QUARANTINE ModificationCode = ModificationCode('q')
)

// define milter protocol negociation message
type MsgOptNeg struct {
	Version uint32 // use MilterVersion
	Actions ActionFlag // use SMFIF_* constants
	Protocol ProtocolFlag // use SMFIP_* constants
}

// define reply code message content
type MsgReply struct {
	Code int
	Reason string
}

// This struct define the Add Header modification. Name is the Name of the
// header to add, and Value is its value.
type MsgAddHeader struct {
	Name string
	Value string
}

// This struct defines an header which change the content. Index is the index
// of the occurrence of this header. Name is the name of header. Value is the
// new value of header
//
// Note that the "index" above is per-name. for example, a 3 in this field
// indicates that the modification is to be applied to the third such
// header matching the supplied "name" field. A zero length string for
// "value", indicates that the header should be deleted entirely.
type MsgChgHeader struct {
	Index uint32
	Name string
	Value string
}

// contains data for CONNECT message
type MsgConnect struct {
	Hostname string
	Family FamilyCode
	Port int
	Address string
}

// Contains data for HEADER Message
type MsgHeader struct {
	Name string
	Value string
}

// contains data for MAIL and RCPT messages
type MsgMail struct {
	Address string
	Args []string
}

// This struct contains and action. An action is sent as response by milter to
// MTA. actions are defined by constants AC_* of type ActionCode. The actions
// are:
//
// ▶︎ AC_ACCEPT : milter ask to MTA to accept message completely. This will skip
// to the end of the milter sequence, and recycle back to the state before
// SMFIC_MAIL. The MTA may, instead, close the connection at that point.
//
// ▶︎ AC_CONTINUE : milter ask to MTA to continue processing. If issued at the
// end of the milter conversation, functions the same as AC_ACCEPT.
//
// ▶︎ AC_DISCARD : milter ask to MTA to set discard flag for entire message
// processing. Note that message processing MAY continue afterwards, but the
// mail will not be delivered even if accepted with SMFIR_ACCEPT.
//
// ▶︎ AC_REJECT : milter ask to MTA to reject email with a 5xx code.
//
// ▶︎ AC_TEMPFAIL : milter ask to MTA to reject email with a temporary 4xx code.
//
// ▶︎ AC_REPLYCODE : milter ask to MTA to answer with the code and message
// specified in the fields Code and Text.
type Action struct {
	Action ActionCode
	Value *MsgReply
}

// This struct contains a modification requirement. A modification is sent as
// response by the milter to the MTA. Kind of modification is defined by
// constants MC_* of type ModificationCode. The fiels Value depends of the kind
// of modification. Avalaible modifications are:
//
// ▶︎ MC_ADDRCPT : Add recipient. The Value is a simple string which contains
// recipient to add.
//
// ▶︎ MC_DELRCPT : Delete recipient. The Value is a simple string which contains
// recipient to delete.
//
// ▶︎ MC_REPLBODY : Replace email body. The Value is a simple string which
// contains the new body.
//
// ▶︎ MC_ADDHEADER : Add header in the email. The Value is a ModAddHeader
// struct. Check documentation struct to understand how use it.
//
// ▶︎ MC_CHGHEADER : change header content in the email. ModAddHeader
// ModChgHeader struct. Check documentation struct to understand how use it.
//
// ▶︎ MC_QUARANTINE : Quarantine message. This quarantines the message into
// a holding pool defined by the MTA.
type Modification struct {
	Modification ModificationCode
	Value interface{}
}

func null_terminated_string(msg []byte, pos int)(int, string) {
	var index int
	for index = pos; index < len(msg); index++ {
		if msg[index] == 0 {
			return index + 1, string(msg[pos:index])
		}
	}
	return -1, ""
}

// expect this format:
// char  smtpcode[3] Nxx code (ASCII), not NUL terminated
// char  space    ' '
// char  text[]      Text of reply message, NUL terminated
func split_code_message(text string)(int, string) {
	var i int
	var err error

	if len(text) < 4 {
		return -1, ""
	}

	if text[3] != ' ' {
		return -1, ""
	}

	i, err = strconv.Atoi(text[:3])
	if err != nil {
		return -1, ""
	}

	return i, text[4:]
}

// this function expect 4 byte from message and return the expected length
// for the next packet to decode. These 4 byte becomes useless. If the
// msg is less than 4 byte, the function fail and error is filled
func DecodeLength(msg []byte)(uint, error) {
	if len(msg) < 4 {
		return 0, fmt.Errorf("Expect 4 bytes, got %d", len(msg))
	}
	return uint(binary.BigEndian.Uint32(msg[:4])), nil
}

// This function expect a buffer which contains full milter message. The
// message size is obtained using function DecodeLength(). The function
// return MsgType which help you to cast interface to the right value.
// if the decoding fail, error will be filled.
//
// returned types acording with MsgType are:
//
//  SMFIC_ABORT      : nil
//  SMFIC_BODY       : []byte
//  SMFIC_CONNECT    : *MsgConnect
//  SMFIC_MACRO      : []*Macro
//  SMFIC_BODYEOB    : nil
//  SMFIC_HELO       : string
//  SMFIC_HEADER     : *MsgHeader
//  SMFIC_MAIL       : *MsgMail
//  SMFIC_EOH        : nil
//  SMFIC_OPTNEG     : *MsgOptNeg
//  SMFIC_RCPT       : *MsgMail
//  SMFIC_QUIT       : nil
//
//  SMFIR_ADDRCPT    : string
//  SMFIR_DELRCPT    : string
//  SMFIR_ACCEPT     : nil
//  SMFIR_REPLBODY   : []byte
//  SMFIR_CONTINUE   : nil
//  SMFIR_DISCARD    : nil
//  SMFIR_ADDHEADER  : *MsgAddHeader
//  SMFIR_CHGHEADER  : *MsgChgHeader
//  SMFIR_PROGRESS   : nil
//  SMFIR_QUARANTINE : nil
//  SMFIR_REJECT     : nil
//  SMFIR_TEMPFAIL   : nil
//  SMFIR_REPLYCODE  : *MsgReply
func Decode(msg []byte)(MsgType, interface{}, error) {
	var msgType MsgType
	var step MacroStep
	var text string
	var pos int
	var code int
	var addheader *MsgAddHeader
	var chgheader *MsgChgHeader
	var optneg *MsgOptNeg
	var mail *MsgMail
	var header *MsgHeader
	var connect *MsgConnect
	var macros []*Macro
	var str string
	var name string
	var value string

	// Read the message code, at least 1 byte
	if len(msg) < 1 {
		return SMFIR_ERROR, nil, fmt.Errorf("protocol error: empty message")
	}

	// Check message type
	msgType = toMsgType(msg[0])
	if msgType == SMFIR_ERROR {
		return SMFIR_ERROR, nil, fmt.Errorf("protocol error: unknown message code %q", string(msg[0]))
	}
	pos = 1

	// Choose the right decoder
	switch msgType {

	// MESSAGES WITHOUT ARGUMENTS
	// --------------------------
	//
	// Client side
	// -----------
	//
	// 'a'	SMFIR_ACCEPT	Accept message completely (accept/reject action)
	//
	// (This will skip to the end of the milter sequence, and recycle back to
	// the state before SMFIC_MAIL.  The MTA may, instead, close the connection
	// at that point.)
	//
	// 'c'	SMFIR_CONTINUE	Accept and keep processing (accept/reject action)
	//
	// (If issued at the end of the milter conversation, functions the same as
	// SMFIR_ACCEPT.)
	//
	// 'd'	SMFIR_DISCARD	Set discard flag for entire message (accept/reject action)
	//
	// (Note that message processing MAY continue afterwards, but the mail will
	// not be delivered even if accepted with SMFIR_ACCEPT.)
	//
	// 'p'	SMFIR_PROGRESS	Progress (asynchronous action)
	//
	// This is an asynchronous response which is sent to the MTA to reset the
	// communications timer during long operations.  The MTA should consume
	// as many of these responses as are sent, waiting for the real response
	// for the issued command.
	//
	// 'r'	SMFIR_REJECT	Reject command/recipient with a 5xx (accept/reject action)
	//
	// 't'	SMFIR_TEMPFAIL	Reject command/recipient with a 4xx (accept/reject action)
	//
	// Server side
	// -----------
	//
	// 'A'	SMFIC_ABORT	Abort current filter checks
	// 			Expected response:  NONE
	//
	// (Resets internal state of milter program to before SMFIC_HELO, but keeps
	// the connection open.)
	//
	// 'E'	SMFIC_BODYEOB	Final body chunk
	//			Expected response:  Zero or more modification
	//			actions, then accept/reject action
	//
	// 'N'	SMFIC_EOH	End of headers marker
	//			Expected response:  Accept/reject action
	//
	// 'Q'	SMFIC_QUIT	Quit milter communication
	//			Expected response:  Close milter connection
	case SMFIC_ABORT,
	     SMFIC_BODYEOB,
	     SMFIC_EOH,
	     SMFIC_QUIT,
	     SMFIR_PROGRESS:
		return msgType, nil, nil

	case SMFIR_ACCEPT,
	     SMFIR_CONTINUE,
	     SMFIR_DISCARD,
	     SMFIR_REJECT,
	     SMFIR_TEMPFAIL:
		return msgType, nil, nil

	// 'B'	SMFIC_BODY	Body chunk
	// 			Expected response:  Accept/reject action
	//
	// char	buf[]		Up to MILTER_CHUNK_SIZE (65535) bytes
	//
	// (These body chunks can be buffered by the milter for later replacement
	// via SMFIR_REPLBODY during the SMFIC_BODYEOB phase.)
	//
	// 'b'	SMFIR_REPLBODY	Replace body (modification action)
	//
	// char	buf[]		Full body, as a single packet
	case SMFIC_BODY,
	     SMFIR_REPLBODY:
		return msgType, msg[pos:], nil

	// 'C'	SMFIC_CONNECT	SMTP connection information
	// 			Expected response:  Accept/reject action
	//
	// char	hostname[]	Hostname, NUL terminated
	// char	family		Protocol family (see below)
	// uint16	port		Port number (SMFIA_INET or SMFIA_INET6 only)
	// char	address[]	IP address (ASCII) or unix socket path, NUL terminated
	//
	// (Sendmail invoked via the command line or via "-bs" will report the
	// connection as the "Unknown" protocol family.)
	//
	// Protocol families used with SMFIC_CONNECT in the "family" field:
	//
	// 'U'	SMFIA_UNKNOWN	Unknown (NOTE: Omits "port" and "host" fields entirely)
	// 'L'	SMFIA_UNIX	Unix (AF_UNIX/AF_LOCAL) socket ("port" is 0)
	// '4'	SMFIA_INET	TCPv4 connection
	// '6'	SMFIA_INET6	TCPv6 connection
	case SMFIC_CONNECT:

		connect = &MsgConnect{}

		// read hostname
		pos, connect.Hostname = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}

		// read 1 byte for family
		if pos + 1 > len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1 byte for protocol type", msgType.String())
		}
		connect.Family = toFamily(msg[pos])
		if connect.Family == SMFIA_ERROR {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> unknown protocol type %q", msgType.String(), string(msg[pos]))
		}
		pos++

		// Read two bytes for port
		if pos + 2 > len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2 bytes for port number", msgType.String())
		}
		connect.Port = int(binary.BigEndian.Uint16(msg[pos:]))
		pos += 2

		// read string for address
		pos, connect.Address = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
		}

		// Check full message eaten
		if pos != len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: remains some byte", msgType.String())
		}

		return msgType, connect, nil

	// 'D'	SMFIC_MACRO	Define macros
	// 			Expected response:  NONE
	//
	// char	cmdcode		Command for which these macros apply
	// char	nameval[][]	Array of NUL-terminated strings, alternating
	// 			between name of macro and value of macro.
	//
	// SMFIC_MACRO appears as a packet just before the corresponding "cmdcode"
	// (here), which is the same identifier as the following command.  The
	// names correspond to Sendmail macros, omitting the "$" identifier
	// character.
	//
	// Types of macros, and some commonly supplied macro names, used with
	// SMFIC_MACRO are as follows, organized by "cmdcode" value.
	// Implementations SHOULD NOT assume that any of these macros will be
	// present on a given connection.  In particular, communications protocol
	// information may not be present on the "Unknown" protocol type.
	//
	// 'C'	SMFIC_CONNECT	$_ $j ${daemon_name} ${if_name} ${if_addr}
	//
	// 'H'	SMFIC_HELO	${tls_version} ${cipher} ${cipher_bits}
	// 			${cert_subject} ${cert_issuer}
	//
	// 'M'	SMFIC_MAIL	$i ${auth_type} ${auth_authen} ${auth_ssf}
	// 			${auth_author} ${mail_mailer} ${mail_host}
	// 			${mail_addr}
	//
	// 'R'	SMFIC_RCPT	${rcpt_mailer} ${rcpt_host} ${rcpt_addr}
	//
	// For future compatibility, implementations MUST allow SMFIC_MACRO at any
	// time, but the handling of unspecified command codes, or SMFIC_MACRO not
	// appearing before its specified command, is currently undefined.
	case SMFIC_MACRO:

		// read step
		if len(msg) < 1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too short: require at least 1 byte", msgType.String())
		}
		step = toMacroStep(msg[pos])
		pos++

		// read alternance of name/values
		for pos < len(msg) {

			// read name
			pos, name = null_terminated_string(msg, pos)
			if pos == -1 {
				return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
			}

			// read value
			pos, value = null_terminated_string(msg, pos)
			if pos == -1 {
				return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
			}

			// Append macro
			macros = append(macros, &Macro{step, name, value})
		}

		return msgType, macros, nil

	// 'H'	SMFIC_HELO	HELO/EHLO name
	//			Expected response:  Accept/reject action
	//
	// char	helo[]		HELO string, NUL terminated
	case SMFIC_HELO:

		// read helo name
		pos, str = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}
		return msgType, str, nil

	// 'L'	SMFIC_HEADER	Mail header
	//			Expected response:  Accept/reject action
	//
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated
	case SMFIC_HEADER:

		header = &MsgHeader{}

		// read name
		pos, header.Name = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}

		// read value
		pos, header.Value = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
		}

		// Check full message eaten
		if pos != len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: remains some byte", msgType.String())
		}

		return msgType, header, nil

	// 'M'	SMFIC_MAIL	MAIL FROM: information
	//			Expected response:  Accept/reject action
	//
	// char	args[][]	Array of strings, NUL terminated (address at index 0).
	//			args[0] is sender, with <> qualification.
	//			args[1] and beyond are ESMTP arguments, if any.
	//
	// 'R'	SMFIC_RCPT	RCPT TO: information
	//			Expected response:  Accept/reject action
	//
	// char	args[][]	Array of strings, NUL terminated (address at index 0).
	//			args[0] is recipient, with <> qualification.
	//			args[1] and beyond are ESMTP arguments, if any.
	case SMFIC_MAIL,
	     SMFIC_RCPT:

		mail = &MsgMail{}

		// read name
		pos, mail.Address = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}

		// Remove brackets
		if len(mail.Address) > 1 {
			if mail.Address[0] == '<' {
				mail.Address = mail.Address[1:]
			}
			if mail.Address[len(mail.Address) - 1] == '>' {
				mail.Address = mail.Address[:len(mail.Address) - 1]
			}
		}

		// read optional args
		for pos < len(msg) {
			pos, value = null_terminated_string(msg, pos)
			if pos == -1 {
				return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
			}
			mail.Args = append(mail.Args, value)
		}

		// Check full message eaten
		if pos != len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: remains some byte", msgType.String())
		}

		return msgType, mail, nil

	// JUST RETURN ONE CHAR[] NULL TERMNATED
	// -------------------------------------
	//
	// '+'	SMFIR_ADDRCPT	Add recipient (modification action)
	//
	// char	rcpt[]		New recipient, NUL terminated
	//
	// '-'	SMFIR_DELRCPT	Remove recipient (modification action)
	//
	// char	rcpt[]		Recipient to remove, NUL terminated
	// 			(string must match the one in SMFIC_RCPT exactly)
	//
	// 'q'	SMFIR_QUARANTINE Quarantine message (modification action)
	// char	reason[]	Reason for quarantine, NUL terminated
	//
	// This quarantines the message into a holding pool defined by the MTA.
	// (First implemented in Sendmail in version 8.13; offered to the milter by
	// the SMFIF_QUARANTINE flag in "actions" of SMFIC_OPTNEG.)
	//
	// 'y'	SMFIR_REPLYCODE	Send specific Nxx reply message (accept/reject action)
	//
	//    Note: Technically reply code is only one string because it is the
	//          concatenation of two non null terminated string and one last
	//          null terminated.
	//
	// char	smtpcode[3]	Nxx code (ASCII), not NUL terminated
	// char	space		' '
	// char	text[]		Text of reply message, NUL terminated
	//
	// ('%' characters present in "text" must be doubled to prevent problems
	// with printf-style formatting that may be used by the MTA.)
	case SMFIR_ADDRCPT,
	     SMFIR_DELRCPT,
	     SMFIR_QUARANTINE,
	     SMFIR_REPLYCODE:
		if msg[len(msg) - 1] != 0 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: expect 1 byte, receive %d bytes", msgType.String(), len(msg))
		}
		text = string(msg[pos:len(msg)-1])
		switch msgType {
		case SMFIR_ADDRCPT,
		     SMFIR_DELRCPT,
		     SMFIR_QUARANTINE:
			return msgType, text, nil

		case SMFIR_REPLYCODE:
			code, text = split_code_message(text)
			if code == -1 {
				return SMFIR_ERROR, nil, fmt.Errorf("protocol error: can't decode reply code %q", text)
			}
			return msgType, &MsgReply{Code: code, Reason: text}, nil
		}
		panic("unreacheable code")

	// 'h'	SMFIR_ADDHEADER	Add header (modification action)
	//
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated
	case SMFIR_ADDHEADER:
		addheader = &MsgAddHeader{}
		pos, addheader.Name = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}
		pos, addheader.Value = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
		}
		if pos != len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: remains some byte", msgType.String())
		}
		return msgType, addheader, nil

	// 'm'	SMFIR_CHGHEADER	Change header (modification action)
	//
	// uint32	index		Index of the occurrence of this header
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated
	//
	// (Note that the "index" above is per-name--i.e. a 3 in this field
	// indicates that the modification is to be applied to the third such
	// header matching the supplied "name" field.  A zero length string for
	// "value", leaving only a single NUL byte, indicates that the header
	// should be deleted entirely.)
	case SMFIR_CHGHEADER:
		chgheader = &MsgChgHeader{}
		if len(msg) < 7 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too short: require at least 7 bytes", msgType.String())
		}
		chgheader.Index = binary.BigEndian.Uint32(msg[pos:])
		pos += 4
		pos, chgheader.Name = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 1st NULL terminated message", msgType.String())
		}
		pos, chgheader.Value = null_terminated_string(msg, pos)
		if pos == -1 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> expect 2nd NULL terminated message", msgType.String())
		}
		if pos != len(msg) {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too long: remains some byte", msgType.String())
		}
		return msgType, chgheader, nil

	// 'O'	SMFIC_OPTNEG	Option negotiation
	// 			Expected response:  SMFIC_OPTNEG packet
	//
	// uint32	version		SMFI_VERSION (2)
	// uint32	actions		Bitmask of allowed actions from SMFIF_*
	// uint32	protocol	Bitmask of possible protocol content from SMFIP_*
	case SMFIC_OPTNEG: // The C is not an error, this code is both client and server message
		if len(msg) != 13 {
			return SMFIR_ERROR, nil, fmt.Errorf("receive message <%s> too short: require at least 13 bytes", msgType.String())
		}
		optneg = &MsgOptNeg{}
		optneg.Version = binary.BigEndian.Uint32(msg[pos:])
		pos += 4
		optneg.Actions = ActionFlag(binary.BigEndian.Uint32(msg[pos:]))
		pos += 4
		optneg.Protocol = ProtocolFlag(binary.BigEndian.Uint32(msg[pos:]))
		return msgType, optneg, nil

	default:
		return SMFIR_ERROR, nil, fmt.Errorf("receive unknown message code %q: %s", string(msgType), msgType.String())
	}
}

// build milter packet like this:
// - uint32  len          Size of data to follow
// - char    cmd          Command/response code
// - char    data[len-1]  Code-specific data (may be empty)
const headerLength = 5
func fillHeader(msg []byte, cmd MsgType, length uint) {
	binary.BigEndian.PutUint32(msg[0:], uint32(length) + 1)
	msg[4] = byte(cmd)
}

// Return []byte which contains QUIT message
func EncodeQuit()([]byte) {
	var msg []byte

	// Make buffer with payload length
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIC_QUIT, 0)

	return msg
}

// Return []byte which contains ABORT message
func EncodeAbort()([]byte) {
	var msg []byte

	// Make buffer with payload length
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIC_ABORT, 0)

	return msg
}

const encodeOptNegLength = 12

// Return []byte which contains OPTNEG message
func EncodeOptNeg(optNeg *MsgOptNeg)([]byte) {
	var msg []byte

	// 'O'	SMFIC_OPTNEG	Option negotiation
	// 			Expected response:  SMFIC_OPTNEG packet
	//
	// uint32	version		SMFI_VERSION (2)
	// uint32	actions		Bitmask of allowed actions from SMFIF_*
	// uint32	protocol	Bitmask of possible protocol content from SMFIP_*

	// Make buffer with payload length
	msg = make([]byte, headerLength + encodeOptNegLength)
	fillHeader(msg, SMFIC_OPTNEG, encodeOptNegLength)

	// Forge data
	binary.BigEndian.PutUint32(msg[headerLength:], optNeg.Version)
	binary.BigEndian.PutUint32(msg[headerLength+4:], uint32(optNeg.Actions))
	binary.BigEndian.PutUint32(msg[headerLength+8:], uint32(optNeg.Protocol))

	return msg
}

// 'D'	SMFIC_MACRO	Define macros
// 			Expected response:  NONE
//
// char	cmdcode		Command for which these macros apply
// char	nameval[][]	Array of NUL-terminated strings, alternating
// 			between name of macro and value of macro.
//
// SMFIC_MACRO appears as a packet just before the corresponding "cmdcode"
// (here), which is the same identifier as the following command.  The
// names correspond to Sendmail macros, omitting the "$" identifier
// character.
//
// Types of macros, and some commonly supplied macro names, used with
// SMFIC_MACRO are as follows, organized by "cmdcode" value.
// Implementations SHOULD NOT assume that any of these macros will be
// present on a given connection.  In particular, communications protocol
// information may not be present on the "Unknown" protocol type.
//
// 'C'	SMFIC_CONNECT	$_ $j ${daemon_name} ${if_name} ${if_addr}
//
// 'H'	SMFIC_HELO	${tls_version} ${cipher} ${cipher_bits}
// 			${cert_subject} ${cert_issuer}
//
// 'M'	SMFIC_MAIL	$i ${auth_type} ${auth_authen} ${auth_ssf}
// 			${auth_author} ${mail_mailer} ${mail_host}
// 			${mail_addr}
//
// 'R'	SMFIC_RCPT	${rcpt_mailer} ${rcpt_host} ${rcpt_addr}
//
// For future compatibility, implementations MUST allow SMFIC_MACRO at any
// time, but the handling of unspecified command codes, or SMFIC_MACRO not
// appearing before its specified command, is currently undefined.
func fillMacroLength(step MacroStep, macros []*Macro)(uint) {
	var length uint
	var m *Macro
	var do_send uint // store value 1 for macro command if send macros has data to send

	// Compute payload langth and make buffer with payload length
	for _, m = range macros {
		if m.Step != step {
			continue
		}
		length += uint(len(m.Name)) + 1
		length += uint(len(m.Value)) + 1
		do_send = 1
	}

	return length + do_send
}
func fillMacro(msg []byte, step MacroStep, macros []*Macro) {
	var pos uint
	var m *Macro

	if len(msg) == 0 {
		return
	}

	// Forge data
	msg[pos] = byte(step)
	pos++
	for _, m = range macros {
		if m.Step != step {
			continue
		}
		copy(msg[pos:], []byte(m.Name))
		pos += uint(len(m.Name))

		msg[pos] = 0
		pos++

		copy(msg[pos:], []byte(m.Value))
		pos += uint(len(m.Value))

		msg[pos] = 0
		pos++
	}
}

// This function encode message which contains macro definition.
// cmd is the step associated with macros. The function encode
// only macro in macros with the corresponding step. If no macros
// match, the function returns nil.
func EncodeMacro(step MacroStep, macros []*Macro)([]byte) {
	var length uint
	var msg []byte
	var pos int

	length = fillMacroLength(step, macros)
	if length == 0 {
		return nil
	}
	msg = make([]byte, headerLength + length)

	fillHeader(msg[pos:], SMFIC_MACRO, length)
	pos += headerLength
	fillMacro(msg[pos:], step, macros)

	return msg
}

// return []byte which contains CONNECT message. The CONNECT message
// is commonly associated with macro, you can give the list of macros
// to transfer. macro is nil if nothing to transfert
func EncodeConnect(connect *MsgConnect, macros []*Macro)([]byte) {
	var msg []byte
	var pos uint
	var macro_header_length uint
	var macro_length uint
	var data_length uint
	var b_hostname []byte
	var b_address []byte

	// 'C'	SMFIC_CONNECT	SMTP connection information
	// 			Expected response:  Accept/reject action
	//
	// char	hostname[]	Hostname, NUL terminated
	// char	family		Protocol family (see below)
	// uint16	port		Port number (SMFIA_INET or SMFIA_INET6 only)
	// char	address[]	IP address (ASCII) or unix socket path, NUL terminated
	//
	// (Sendmail invoked via the command line or via "-bs" will report the
	// connection as the "Unknown" protocol family.)
	//
	// Protocol families used with SMFIC_CONNECT in the "family" field:
	//
	// 'U'	SMFIA_UNKNOWN	Unknown (NOTE: Omits "port" and "host" fields entirely)
	// 'L'	SMFIA_UNIX	Unix (AF_UNIX/AF_LOCAL) socket ("port" is 0)
	// '4'	SMFIA_INET	TCPv4 connection
	// '6'	SMFIA_INET6	TCPv6 connection

	b_hostname = []byte(connect.Hostname)
	b_address = []byte(connect.Address)

	// Compute payload langth and make buffer with payload length
	macro_length = fillMacroLength(MS_CONNECT, macros)
	if macro_length > 0 {
		macro_header_length = headerLength
	}
	data_length = 3 // 1B family + 2B port
	data_length += uint(len(b_hostname)) + 1
	data_length += uint(len(b_address)) + 1

	// Create message
	msg = make([]byte, macro_header_length + macro_length + headerLength + data_length)

	// Encode macro relative to the current step
	if macro_length > 0 {
		fillHeader(msg[pos:], SMFIC_MACRO, macro_length)
		pos += headerLength
		fillMacro(msg[pos:], MS_CONNECT, macros)
		pos += macro_length
	}

	// Append header
	fillHeader(msg[pos:], SMFIC_CONNECT, data_length)
	pos += headerLength

	// Fill payload

	// copy hostname
	copy(msg[pos:], b_hostname)
	pos += uint(len(b_hostname))
	msg[pos] = 0
	pos++

	// Set network address type. type is based on the IP buffer length
	// documentation says "if ip is not an IPv4 address, To4 returns nil"
	msg[pos] = byte(connect.Family)
	pos++;

	// Set port as bigendian 16bit
	binary.BigEndian.PutUint16(msg[pos:], uint16(connect.Port))
	pos += 2

	// display IP as string and copy it
	copy(msg[pos:], b_address)
	pos += uint(len(b_address))
	msg[pos] = 0

	return msg
}

// return []byte which contains HELO message. The HELO message
// is commonly associated with macro, you can give the list of macros
// to transfer. macro is nil if nothing to transfert
func EncodeHelo(helo string, macros []*Macro)([]byte) {
	var msg []byte
	var pos uint
	var bytes []byte
	var macro_length uint
	var macro_header_length uint
	var data_length uint

	// 'H'	SMFIC_HELO	HELO/EHLO name
	// 			Expected response:  Accept/reject action
	//
	// char	helo[]		HELO string, NUL terminated

	bytes = []byte(helo)

	// Compute payload langth and make buffer with payload length
	macro_length = fillMacroLength(MS_HELO, macros)
	if macro_length > 0 {
		macro_header_length = headerLength
	}
	data_length = uint(len(bytes)) + 1

	// Create message
	msg = make([]byte, macro_header_length + macro_length + headerLength + data_length)

	// Encode macro relative to the current step
	if macro_length > 0 {
		fillHeader(msg[pos:], SMFIC_MACRO, macro_length)
		pos += headerLength
		fillMacro(msg[pos:], MS_HELO, macros)
		pos += macro_length
	}

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIC_HELO, data_length)
	pos += headerLength

	// Fill payload
	copy(msg[pos:], bytes)
	pos += uint(len(bytes))
	msg[pos] = 0

	return msg
}

func encodeMailRcpt(step MsgType, email *MsgMail, macros []*Macro)([]byte) {
	var msg []byte
	var pos uint
	var macro_length uint
	var macro_header_length uint
	var data_length uint
	var str string

	// 'M'	SMFIC_MAIL	MAIL FROM: information
	// 			Expected response:  Accept/reject action
	//
	// char	args[][]	Array of strings, NUL terminated (address at index 0).
	// 			args[0] is sender, with <> qualification.
	// 			args[1] and beyond are ESMTP arguments, if any.
	//
	// 'R'	SMFIC_RCPT	RCPT TO: information
	// 			Expected response:  Accept/reject action
	//
	// char	args[][]	Array of strings, NUL terminated (address at index 0).
	// 			args[0] is recipient, with <> qualification.
	// 			args[1] and beyond are ESMTP arguments, if any.

	// Compute payload langth and make buffer with payload length
	macro_length = fillMacroLength(MacroStep(step), macros)
	if macro_length > 0 {
		macro_header_length = headerLength
	}
	data_length = uint(len([]byte(email.Address))) + 2 /* brackets */ + 1 /* final \0 */;
	for _, str = range email.Args {
		data_length += uint(len([]byte(str))) + 1
	}

	// Create message
	msg = make([]byte, macro_header_length + macro_length + headerLength + data_length)

	// Encode macro relative to the current step
	if macro_length > 0 {
		fillHeader(msg[pos:], SMFIC_MACRO, macro_length)
		pos += headerLength
		fillMacro(msg[pos:], MacroStep(step), macros)
		pos += macro_length
	}

	// Compute payload langth and make buffer with payload length
	fillHeader(msg[pos:], step, data_length)
	pos += headerLength

	// Fill payload
	msg[pos] = '<'
	pos++
	copy(msg[pos:], []byte(email.Address))
	pos += uint(len([]byte(email.Address)))
	msg[pos] = '>'
	pos++
	msg[pos] = 0
	pos++

	// Fill args
	for _, str = range email.Args {
		copy(msg[pos:], []byte(str))
		pos += uint(len([]byte(str)))
		msg[pos] = 0
		pos++
	}

	return msg
}

// return []byte which contains MAIL message. The MAIL message
// is commonly associated with macro, you can give the list of macros
// to transfer. macro is nil if nothing to transfert
func EncodeMail(email *MsgMail, macros []*Macro)([]byte) {
	return encodeMailRcpt(SMFIC_MAIL, email, macros)
}

// return []byte which contains RCPT message. The RCPT message
// is commonly associated with macro, you can give the list of macros
// to transfer. macro is nil if nothing to transfert
func EncodeRcpt(email *MsgMail, macros []*Macro)([]byte) {
	return encodeMailRcpt(SMFIC_RCPT, email, macros)
}

// return []byte which contains HEADER message.
func EncodeHeader(hdr *MsgHeader)([]byte) {
	var msg []byte
	var pos uint
	var length uint
	var bytes []byte

	// 'L'	SMFIC_HEADER	Mail header
	// 			Expected response:  Accept/reject action
	//
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated

	// Compute payload langth and make buffer with payload length
	length = 0
	length += uint(len(hdr.Name)) + 1
	length += uint(len(hdr.Value)) + 1

	// Build buffer
	msg = make([]byte, length + headerLength)

	// Append header
	fillHeader(msg, SMFIC_HEADER, length)
	pos += headerLength

	// copy name
	bytes = []byte(hdr.Name)
	copy(msg[pos:], bytes)
	pos += uint(len(bytes))
	msg[pos] = 0
	pos++

	// copy value
	bytes = []byte(hdr.Value)
	copy(msg[pos:], bytes)
	pos += uint(len(bytes))
	msg[pos] = 0

	return msg
}

// return []byte which contains EOH message.
func EncodeEOH()([]byte) {
	var msg []byte

	// 'N'	SMFIC_EOH	End of headers marker
	//			Expected response:  Accept/reject action

	// Compute payload langth and make buffer with payload length
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIC_EOH, 0)

	return msg
}

// return []byte which contains BODY message. The body should not
// exceed 65535 bytes. This function should called more than one
// time to transfert all the body.
func EncodeBody(body []byte)([]byte) {
	var msg []byte
	var pos uint
	var bytes []byte

	// 'B'	SMFIC_BODY	Body chunk
	// 			Expected response:  Accept/reject action
	//
	// char	buf[]		Up to MILTER_CHUNK_SIZE (65535) bytes

	bytes = []byte(body)

	// build buffer
	msg = make([]byte, headerLength + len(bytes))

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIC_BODY, uint(len(bytes)))
	pos += headerLength

	// Fill payload
	copy(msg[pos:], bytes)

	return msg
}

// return []byte which contains BODYEOB message.
func EncodeBodyEOB()([]byte) {
	var msg []byte

	// 'E'	SMFIC_BODYEOB	Final body chunk
	// 			Expected response:  Zero or more modification
	// 			actions, then accept/reject action

	// Compute payload langth and make buffer with payload length
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIC_BODYEOB, 0)

	return msg
}

// return []byte which contains ADDRCPT message.
func EncodeAddRcpt(rcpt string)([]byte) {
	var msg []byte
	var pos uint
	var byte_rcpt []byte
	var len_rcpt uint

	// '+'	SMFIR_ADDRCPT	Add recipient (modification action)
	//
	// char	rcpt[]		New recipient, NUL terminated

	byte_rcpt = []byte(rcpt)
	len_rcpt = uint(len(byte_rcpt))

	msg = make([]byte, headerLength + len_rcpt + 1)

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIR_ADDRCPT, len_rcpt + 1)
	pos = headerLength

	// Copy name
	copy(msg[pos:], byte_rcpt)
	pos += len_rcpt
	msg[pos] = 0

	return msg
}

// return []byte which contains DELRCPT message.
func EncodeDelRcpt(rcpt string)([]byte) {
	var msg []byte
	var pos uint
	var len_rcpt uint
	var byte_rcpt []byte

	// '-' SMFIR_DELRCPT   Remove recipient (modification action)
	//
	// char        rcpt[]          Recipient to remove, NUL terminated
	//                     (string must match the one in SMFIC_RCPT exactly)

	byte_rcpt = []byte(rcpt)
	len_rcpt = uint(len(byte_rcpt))

	msg = make([]byte, headerLength + len_rcpt + 1)

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIR_DELRCPT, len_rcpt + 1)
	pos += headerLength

	// Copy name
	copy(msg[pos:], []byte(rcpt))
	pos += len_rcpt
	msg[pos] = 0

	return msg
}

// return []byte which contains ACCEPT message.
func EncodeAccept()([]byte) {
	var msg []byte

	// 'a' SMFIR_ACCEPT    Accept message completely (accept/reject action)
	//
	// (This will skip to the end of the milter sequence, and recycle back to
	// the state before SMFIC_MAIL.  The MTA may, instead, close the connection
	// at that point.)
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_ACCEPT, 0)
	return msg
}

// return []byte which contains CONTINUE message.
func EncodeContinue()([]byte) {
	var msg []byte

	// 'c'	SMFIR_CONTINUE	Accept and keep processing (accept/reject action)
	//
	// (If issued at the end of the milter conversation, functions the same as
	// SMFIR_ACCEPT.)
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_CONTINUE, 0)
	return msg
}

// return []byte which contains DISCARD message.
func EncodeDiscard()([]byte) {
	var msg []byte

	// 'd'	SMFIR_DISCARD	Set discard flag for entire message (accept/reject action)
	//
	// (Note that message processing MAY continue afterwards, but the mail will
	// not be delivered even if accepted with SMFIR_ACCEPT.)
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_DISCARD, 0)
	return msg
}

// return []byte which contains REJECT message.
func EncodeReject()([]byte) {
	var msg []byte

	// 'r'	SMFIR_REJECT	Reject command/recipient with a 5xx (accept/reject action)
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_REJECT, 0)
	return msg
}

// return []byte which contains TEMPFAIL message.
func EncodeTempfail()([]byte) {
	var msg []byte

	// 't'	SMFIR_TEMPFAIL	Reject command/recipient with a 4xx (accept/reject action)
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_TEMPFAIL, 0)
	return msg
}

// return []byte which contains REPLBODY message.
func EncodeReplBody(body []byte)([]byte) {
	var msg []byte
	var pos uint

	// 'b'	SMFIR_REPLBODY	Replace body (modification action)
	//
	// char	buf[]		Full body, as a single packet
	msg = make([]byte, headerLength + len(body))

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIR_REPLBODY, uint(len(body)))
	pos += headerLength

	// Copy body
	copy(msg[pos:], body)

	return msg
}

// return []byte which contains ADDHEADER message.
func EncodeAddHeader(addhdr *MsgAddHeader)([]byte) {
	var msg []byte
	var pos uint
	var byte_name []byte
	var byte_value []byte
	var len_name uint
	var len_value uint

	// 'h'	SMFIR_ADDHEADER	Add header (modification action)
	//
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated
	byte_name = []byte(addhdr.Name)
	byte_value = []byte(addhdr.Value)
	len_name = uint(len(byte_name))
	len_value = uint(len(byte_value))

	msg = make([]byte, headerLength + len_name + 1 + len_value + 1)

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIR_ADDHEADER, len_name + 1 + len_value + 1)
	pos += headerLength

	// Copy name
	copy(msg[pos:], byte_name)
	pos += len_name
	msg[pos] = 0
	pos++

	// Copy value
	copy(msg[pos:], byte_value)
	pos += len_value
	msg[pos] = 0

	return msg
}

// return []byte which contains CHGHEADER message.
func EncodeChgHeader(chghdr *MsgChgHeader)([]byte) {
	var msg []byte
	var pos uint
	var byte_name []byte
	var byte_value []byte
	var len_name uint
	var len_value uint

	// 'm'	SMFIR_CHGHEADER	Change header (modification action)
	//
	// uint32	index		Index of the occurrence of this header
	// char	name[]		Name of header, NUL terminated
	// char	value[]		Value of header, NUL terminated
	//
	// (Note that the "index" above is per-name--i.e. a 3 in this field
	// indicates that the modification is to be applied to the third such
	// header matching the supplied "name" field.  A zero length string for
	// "value", leaving only a single NUL byte, indicates that the header
	// should be deleted entirely.)
	byte_name = []byte(chghdr.Name)
	byte_value = []byte(chghdr.Value)
	len_name = uint(len(byte_name))
	len_value = uint(len(byte_value))

	msg = make([]byte, headerLength + 4 + len_name + 1 + len_value + 1)

	// Compute payload langth and make buffer with payload length
	fillHeader(msg, SMFIR_CHGHEADER, 4 + len_name + 1 + len_value + 1)
	pos += headerLength

	// Copy index as 32 bit integer
	binary.BigEndian.PutUint32(msg[pos:], uint32(chghdr.Index))
	pos += 4

	// Copy name
	copy(msg[pos:], byte_name)
	pos += len_name
	msg[pos] = 0
	pos++

	// Copy value
	copy(msg[pos:], byte_value)
	pos += len_value
	msg[pos] = 0

	return msg
}

// return []byte which contains PROGRESS message.
func EncodeProgress()([]byte) {
	var msg []byte

	// 'p'	SMFIR_PROGRESS	Progress (asynchronous action)
	//
	// This is an asynchronous response which is sent to the MTA to reset the
	// communications timer during long operations.  The MTA should consume
	// as many of these responses as are sent, waiting for the real response
	// for the issued command.
	msg = make([]byte, headerLength)
	fillHeader(msg, SMFIR_PROGRESS, 0)
	return msg
}

// return []byte which contains QUARANTINE message.
func EncodeQuarantine(reason string)([]byte) {
	var msg []byte
	var pos uint
	var byte_reason []byte
	var len_reason uint

	// 'q'	SMFIR_QUARANTINE Quarantine message (modification action)
	// char	reason[]	Reason for quarantine, NUL terminated
	//
	// This quarantines the message into a holding pool defined by the MTA.
	// (First implemented in Sendmail in version 8.13; offered to the milter by
	// the SMFIF_QUARANTINE flag in "actions" of SMFIC_OPTNEG.)
	byte_reason = []byte(reason)
	len_reason = uint(len(byte_reason))

	// Make buffer
	msg = make([]byte, headerLength + len_reason + 1)

	// Compute header
	fillHeader(msg, SMFIR_QUARANTINE, len_reason + 1)
	pos += headerLength

	// Copy message
	copy(msg[pos:], byte_reason)
	pos += len_reason


	// nul char
	msg[pos] = 0

	return msg
}

// return []byte which contains REPLYCODE message.
func EncodeReplyCode(reply *MsgReply)([]byte) {
	var msg []byte
	var pos uint
	var byte_reason []byte
	var len_reason uint
	var byte_code []byte
	var len_code uint
	var code int

	// 'y'	SMFIR_REPLYCODE	Send specific Nxx reply message (accept/reject action)
	//
	// char	smtpcode[3]	Nxx err (ASCII), not NUL terminated
	// char	space		' '
	// char	text[]		Text of reply message, NUL terminated
	//
	// ('%' characters present in "text" must be doubled to prevent problems
	// with printf-style formatting that may be used by the MTA.)
	//
	// NOTE:
	//     Nxx err (ASCII), not NUL terminated
	//   + ' '
	//   + Text of reply message, NUL terminated
	//  ------------------------------------------
	//   = NUL terminated string
	code = reply.Code
	if code < 0 || code > 999 {
		code = 0
	}
	byte_code = []byte(strconv.Itoa(reply.Code))
	len_code = uint(len(byte_code))

	byte_reason = []byte(reply.Reason)
	len_reason = uint(len(byte_reason))

	msg = make([]byte, headerLength + len_code + 1 /* space */ + len_reason + 1 /* nul */)

	// add header
	fillHeader(msg, SMFIR_REPLYCODE, len_code + 1 + len_reason + 1)
	pos += headerLength

	// copy code
	copy(msg[pos:], byte_code)
	pos += len_code

	// add space
	msg[pos] = ' '
	pos++

	// add reason
	copy(msg[pos:], byte_reason)
	pos += len_reason

	// null
	msg[pos] = 0

	return msg
}

