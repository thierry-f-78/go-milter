// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

// Implement simple milter server which reject SMTP client randomly
package milter_test

import "fmt"
import "io"
import "log"
import "net"
import "time"

import "github.com/thierry-f-78/go-milter"

type IpDecision struct{}

func (id *IpDecision)OnOPTNEG(srv *milter.Server, optNeg *milter.MsgOptNeg)(*milter.MsgOptNeg, error) {

	// Inform client you want only CONNECT step and no specific actions
	return &milter.MsgOptNeg{
		Version: milter.MilterVersion,
		Protocol: milter.SMFIP_NOHELO | milter.SMFIP_NOMAIL | milter.SMFIP_NORCPT |
		          milter.SMFIP_NOBODY | milter.SMFIP_NOHDRS | milter.SMFIP_NOEOH,
		Actions: 0,
	}, nil
}

func (id *IpDecision)OnCONNECT(srv *milter.Server, connect *milter.MsgConnect)(*milter.Action, error) {

	// Randomly reject email (reject if the current microsecond id odd)
	if time.Now().UnixNano() & 0x1 == 1 {
		return milter.ActionReject(), nil
	} else {
		return milter.ActionContinue(), nil
	}
}

func (id *IpDecision)OnBODYEOB(srv *milter.Server)([]*milter.Modification, *milter.Action, error) {
	return nil, milter.ActionContinue(), nil
}

func (id *IpDecision)OnERROR(srv *milter.Server, err error)() {
	if err != nil {
		if err != io.EOF {
			log.Printf("%s\n", err.Error())
		}
	} else {
		log.Printf("Unknown error\n")
	}
}

func (id *IpDecision)OnABORT(srv *milter.Server)(error) {
	return nil
}

func (id *IpDecision)OnQUIT(srv *milter.Server)(error) {
	return nil
}

// useless functions declared to match interface

func (id *IpDecision)OnHELO(srv *milter.Server, helo string)(*milter.Action, error) {
	return nil, fmt.Errorf("Step HELO not supported")
}
func (id *IpDecision)OnMAIL(srv *milter.Server, mail *milter.MsgMail)(*milter.Action, error) {
	return nil, fmt.Errorf("Step MAIL not supported")
}
func (id *IpDecision)OnRCPT(srv *milter.Server, mail *milter.MsgMail)(*milter.Action, error) {
	return nil, fmt.Errorf("Step RCPT not supported")
}
func (id *IpDecision)OnHEADER(srv *milter.Server, hdr *milter.MsgHeader)(*milter.Action, error) {
	return nil, fmt.Errorf("Step HEADER not supported")
}
func (id *IpDecision)OnEOH(srv *milter.Server)(*milter.Action, error) {
	return nil, fmt.Errorf("Step EOH not supported")
}
func (id *IpDecision)OnBODY(srv *milter.Server, body []byte)(*milter.Action, error) {
	return nil, fmt.Errorf("Step BODY not supported")
}

// This example propose simple milter server which block email according with its IP address
func Example_exchangeIpDecision() {
	var err error
	var conn net.Conn
	var l net.Listener

	// listen port
	l, err = net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Listen for an incoming connection.
	conn, err = l.Accept()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// start session handler
	milter.Exchange(conn, &IpDecision{})
	conn.Close()

	// Output: 
}
