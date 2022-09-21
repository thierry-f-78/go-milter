// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "fmt"
import "reflect"
import "testing"

func expect(data []byte, kind MsgType, msg interface{})(string) {
	var l uint
	var decode_kind MsgType
	var decode_msg interface{}
	var err error

	l, err = DecodeLength(data)
	if err != nil {
		return err.Error()
	}
	if len(data[4:]) != int(l) {
		return fmt.Sprintf("Expect message of length %d, got %d", l, len(data[4:]))
	}

	decode_kind, decode_msg, err = Decode(data[4:])
	if err != nil {
		return err.Error()
	}
	if decode_kind != kind {
		return fmt.Sprintf("Expect message of type %s, got %s (%02x)",
		       kind.String(), decode_kind.String(), decode_kind)
	}
	if !reflect.DeepEqual(decode_msg, msg) {
		return fmt.Sprintf("Decoded struct %#v not match expected %#v", decode_msg, msg)
	}
	return ""
}

func Test_proto(t *testing.T) {
	var message []byte
	var verdict string
	var emailAddress string = "myemail.address@anylocation.fr"
	var helo string = "my.host.name"
	var reason string = "because"
	var msgHeader MsgHeader = MsgHeader{
		Name: "Header-Name",
		Value: "header value",
	}
	var msgAddHeader MsgAddHeader = MsgAddHeader{
		Name: "Header-Name",
		Value: "header value",
	}
	var msgChgHeader MsgChgHeader = MsgChgHeader{
		Index: 33,
		Name: "Header-Name",
		Value: "header value",
	}
	var body []byte = []byte("This is the body")
	var msgConnect MsgConnect = MsgConnect{
		Hostname: "my.host.name",
		Family: SMFIA_INET,
		Port: 25,
		Address: "127.0.0.1",
	}
	var macros []*Macro = []*Macro{
		&Macro{
			Step: MS_CONNECT,
			Name: "{macro1}",
			Value: "value 01",
		},
		&Macro{
			Step: MS_CONNECT,
			Name: "{macro2}",
			Value: "value 02",
		},
	}
	var msgMail MsgMail = MsgMail{
		Address: "myemail.address@anylocation.fr",
		Args: []string{
			"arg0",
			"arg1",
		},
	}
	var msgOptNeg MsgOptNeg = MsgOptNeg{
		Version: 2,
		Actions: SMFIF_ADDRCPT | SMFIF_CHGHDRS,
		Protocol: SMFIP_NOHELO | SMFIP_NOHDRS,
	}
	var msgReply MsgReply = MsgReply{
		Code: 405,
		Reason: "because",
	}

	message = EncodeAbort()
	verdict = expect(message, SMFIC_ABORT, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeAccept()
	verdict = expect(message, SMFIR_ACCEPT, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeAddHeader(&msgAddHeader)
	verdict = expect(message, SMFIR_ADDHEADER, &msgAddHeader)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeAddRcpt(emailAddress)
	verdict = expect(message, SMFIR_ADDRCPT, emailAddress)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeBody(body)
	verdict = expect(message, SMFIC_BODY, body)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeBodyEOB()
	verdict = expect(message, SMFIC_BODYEOB, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeChgHeader(&msgChgHeader)
	verdict = expect(message, SMFIR_CHGHEADER, &msgChgHeader)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeConnect(&msgConnect, nil)
	verdict = expect(message, SMFIC_CONNECT, &msgConnect)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeContinue()
	verdict = expect(message, SMFIR_CONTINUE, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeDelRcpt(emailAddress)
	verdict = expect(message, SMFIR_DELRCPT, emailAddress)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeDiscard()
	verdict = expect(message, SMFIR_DISCARD, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeEOH()
	verdict = expect(message, SMFIC_EOH, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeHeader(&msgHeader)
	verdict = expect(message, SMFIC_HEADER, &msgHeader)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeHelo(helo, nil)
	verdict = expect(message, SMFIC_HELO, helo)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeMacro(MS_CONNECT, macros)
	verdict = expect(message, SMFIC_MACRO, macros)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeMail(&msgMail, nil)
	verdict = expect(message, SMFIC_MAIL, &msgMail)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeOptNeg(&msgOptNeg)
	verdict = expect(message, SMFIC_OPTNEG, &msgOptNeg)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeProgress()
	verdict = expect(message, SMFIR_PROGRESS, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeQuarantine(reason)
	verdict = expect(message, SMFIR_QUARANTINE, reason)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeQuit()
	verdict = expect(message, SMFIC_QUIT, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeRcpt(&msgMail, nil)
	verdict = expect(message, SMFIC_RCPT, &msgMail)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeReject()
	verdict = expect(message, SMFIR_REJECT, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeReplBody(body)
	verdict = expect(message, SMFIR_REPLBODY, body)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeReplyCode(&msgReply)
	verdict = expect(message, SMFIR_REPLYCODE, &msgReply)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}

	message = EncodeTempfail()
	verdict = expect(message, SMFIR_TEMPFAIL, nil)
	if verdict != "" {
		t.Errorf("%s", verdict)
	}
}
