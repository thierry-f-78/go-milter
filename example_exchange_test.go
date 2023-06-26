// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

// This exemple shows usage of Exchange* API in the same time for server and client.
package milter_test

import "fmt"
import "log"
import "net"

import "github.com/thierry-f-78/go-milter"

type Proxy struct{
	cli *milter.Client
}

func (px *Proxy)OnOPTNEG(srv *milter.Server, optNeg *milter.MsgOptNeg)(*milter.MsgOptNeg, error) {
	var err error

	// Dump received information
	fmt.Printf("> OPTNEG %s\n", optNeg.String())

	// First connect to the target server
	px.cli, err = milter.ClientNew("tcp", "127.0.0.1:7357", 10)
	if err != nil {
		return nil, fmt.Errorf("Milter server error: %s", err.Error())
	}

	// Send negociation options
	optNeg, err = px.cli.ExchangeOptNeg(optNeg)
	if err != nil {
		px.cli.Close()
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< OPTNEG %s\n", optNeg.String())

	// Send answer
	return optNeg, nil
}

func (px *Proxy)OnCONNECT(srv *milter.Server, connect *milter.MsgConnect)(*milter.Action, error) {
	var err error
	var action *milter.Action
	var m *milter.Macro

	// Dump received information
	fmt.Printf("> CONNECT %s\n", connect.String())

	// Dump macros
	for _, m = range srv.Macros {
		if m.Step != milter.MS_CONNECT {
			continue
		}
		fmt.Printf("  MACRO %s\n", m.String())
	}

	// Copy macro
	px.cli.Macros = srv.Macros

	// Send connect info
	action, err = px.cli.ExchangeConnect(connect)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< CONNECT %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnHELO(srv *milter.Server, helo string)(*milter.Action, error) {
	var err error
	var action *milter.Action
	var m *milter.Macro

	// Dump received information
	fmt.Printf("> HELO %q\n", helo)

	// Dump macros
	for _, m = range srv.Macros {
		if m.Step != milter.MS_HELO {
			continue
		}
		fmt.Printf("  MACRO %s\n", m.String())
	}

	// Copy macro
	px.cli.Macros = srv.Macros

	// Send helo information
	action, err = px.cli.ExchangeHelo(helo)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< HELO %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnMAIL(srv *milter.Server, mail *milter.MsgMail)(*milter.Action, error) {
	var err error
	var action *milter.Action
	var m *milter.Macro

	// Dump received information
	fmt.Printf("> MAIL %s\n", mail.String())

	// Dump macros
	for _, m = range srv.Macros {
		if m.Step != milter.MS_MAIL {
			continue
		}
		fmt.Printf("  MACRO %s\n", m.String())
	}

	// Copy macro
	px.cli.Macros = srv.Macros

	// Send helo information
	action, err = px.cli.ExchangeMail(mail)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< MAIL %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnRCPT(srv *milter.Server, mail *milter.MsgMail)(*milter.Action, error) {
	var err error
	var action *milter.Action
	var m *milter.Macro

	// Dump received information
	fmt.Printf("> RCPT %s\n", mail.String())

	// Dump macros
	for _, m = range srv.Macros {
		if m.Step != milter.MS_RCPT {
			continue
		}
		fmt.Printf("  MACRO %s\n", m.String())
	}

	// Copy macro
	px.cli.Macros = srv.Macros

	// Send helo information
	action, err = px.cli.ExchangeRcpt(mail)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< RCPT %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnHEADER(srv *milter.Server, hdr *milter.MsgHeader)(*milter.Action, error) {
	var err error
	var action *milter.Action

	// Dump received information
	fmt.Printf("> HEADER %s\n", hdr.String())

	// Send helo information
	action, err = px.cli.ExchangeHeader(hdr)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< HEADER %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnEOH(srv *milter.Server)(*milter.Action, error) {
	var err error
	var action *milter.Action

	// Dump received information
	fmt.Printf("> EOH\n")

	// Send helo information
	action, err = px.cli.ExchangeEOH()
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< EOH %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnBODY(srv *milter.Server, body []byte)(*milter.Action, error) {
	var err error
	var action *milter.Action

	// Dump received information
	if len(body) > 20 {
		fmt.Printf("> BODY %q...\n", string(body[:20]))
	} else {
		fmt.Printf("> BODY %q\n", string(body))
	}

	// Send helo information
	action, err = px.cli.ExchangeBody(body)
	if err != nil {
		return nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< BODY %s\n", action.String())

	// Send action
	return action, nil
}

func (px *Proxy)OnBODYEOB(srv *milter.Server)([]*milter.Modification, *milter.Action, error) {
	var err error
	var action *milter.Action
	var modifications []*milter.Modification

	// Dump received information
	fmt.Printf("> BODYEOB\n")

	// Send helo information
	modifications, action, err = px.cli.ExchangeBodyEOB()
	if err != nil {
		return nil, nil, fmt.Errorf("Milter client error: %s", err.Error())
	}

	// Dump sent information
	fmt.Printf("< BODYEOB %s\n", action.String())

	// Send action
	return modifications, action, nil
}

func (px *Proxy)OnABORT(srv *milter.Server)(error) {
	var err error

	// Dump received information
	fmt.Printf("> ABORT\n")

	// Send abort information
	err = px.cli.ExchangeAbort()
	if err != nil {
		return fmt.Errorf("Milter client error: %s", err.Error())
	}

	return nil
}

func (px *Proxy)OnQUIT(srv *milter.Server)(error) {
	var err error

	// Dump received information
	fmt.Printf("> QUIT\n")

	// Send QUIT information
	err = px.cli.ExchangeQuit()
	if err != nil {
		return fmt.Errorf("Milter client error: %s", err.Error())
	}

	return nil
}

func (px *Proxy)OnERROR(srv *milter.Server, err error)() {
	if err != nil {
		fmt.Printf("> ERROR %q\n", err.Error())
	} else {
		fmt.Printf("> ERROR\n")
	}
}

func Example_exchangeProxy() {
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

	milter.Exchange(conn, &Proxy{})
	conn.Close()

	// Output: 
}
