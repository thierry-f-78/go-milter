// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

// This exemple shows usage of Send*/Receive* API in the same time for server and client.
package milter_test

import "fmt"
import "log"
import "net"

import "github.com/thierry-f-78/go-milter"

func proxy_send(conn net.Conn) {
	var srv *milter.Server
	var cli *milter.Client
	var msgType milter.MsgType
	var msg interface{}
	var err error
	var optNeg *milter.MsgOptNeg
	var connect *milter.MsgConnect
	var mail *milter.MsgMail
	var hdr *milter.MsgHeader
	var action *milter.Action
	var modification *milter.Modification
	var str string
	var data []byte
	var macros []*milter.Macro
	var macro *milter.Macro

	defer conn.Close()

	// New milter server
	srv = milter.ServerNew(conn)

	// Connect to server
	cli, err = milter.ClientNew("tcp", "127.0.0.1:7357", 10)
	if err != nil {
		log.Printf("Server side error: %s\n", err.Error())
		return
	}

	defer cli.Close()

	// Loop on messages
	for {

		// Next message
		msgType, msg, err = srv.ReceiveMessage()
		if err != nil {
			log.Printf("Client side error: %s", err.Error())
			return
		}

		// display information acording with message type
		switch msgType {

		case milter.SMFIC_OPTNEG:

			// Cast message according with protocol
			optNeg = msg.(*milter.MsgOptNeg)

			// Dump client request information
			fmt.Printf("> OPTNEG %s\n", optNeg.String())

			// Forward request to milter server
			err = cli.SendOptNeg(optNeg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect OptNeg
			if msgType != milter.SMFIC_OPTNEG {
				log.Printf("Server side error: Expect OptNeg message, got %s", msgType.String())
				return
			}
			optNeg = msg.(*milter.MsgOptNeg)

			// Dump server response information
			fmt.Printf("< OPTNEG %s\n", optNeg.String())

			// Send answer to client
			err = srv.SendOptNeg(optNeg)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_CONNECT:

			// Cast message according with protocol
			connect = msg.(*milter.MsgConnect)

			// Dump received information
			fmt.Printf("> CONNECT %s\n", connect.String())

			// Forward request to milter server
			err = cli.SendConnect(connect)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< CONNECT %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_HELO:

			// Cast message according with protocol
			str = msg.(string)

			// Dump received information
			fmt.Printf("> HELO %q\n", str)

			// Forward request to milter server
			err = cli.SendHelo(str)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< HELO %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_MAIL:

			// Cast message according with protocol
			mail = msg.(*milter.MsgMail)

			// Dump received information
			fmt.Printf("> MAIL %s\n", mail.String())

			// Forward request to milter server
			err = cli.SendMail(mail)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< MAIL %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_RCPT:

			// Cast message according with protocol
			mail = msg.(*milter.MsgMail)

			// Dump received information
			fmt.Printf("> RCPT %s\n", mail.String())

			// Forward request to milter server
			err = cli.SendRcpt(mail)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< RCPT %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_HEADER:

			// Cast message according with protocol
			hdr = msg.(*milter.MsgHeader)

			// Dump received information
			fmt.Printf("> HEADER %s\n", hdr.String())

			// Forward request to milter server
			err = cli.SendHeader(hdr)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< HEADER %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_EOH:

			// Dump received information
			fmt.Printf("> EOH\n")

			// Forward request to milter server
			err = cli.SendEOH()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< EOH %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_BODY:

			// Cast message according with protocol
			data = msg.([]byte)

			// Dump received information
			if len(data) > 20 {
				fmt.Printf("> BODY %q...\n", string(data[:20]))
			} else {
				fmt.Printf("> BODY %q\n", string(data))
			}

			// Forward request to milter server
			err = cli.SendBody(data)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Receive response from milter server
			msgType, msg, err = cli.ReceiveMessage()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Expect Action
			action, err = milter.AnswerToAction(msgType, msg)
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Dump server response information
			fmt.Printf("< BODY %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_BODYEOB:

			// Dump received information
			fmt.Printf("> BODYEOB\n")

			// Forward request to milter server
			err = cli.SendBodyEOB()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// Read all response until we got actio
			for {

				// Receive response from milter server
				msgType, msg, err = cli.ReceiveMessage()
				if err != nil {
					log.Printf("Server side error: %s", err.Error())
					return
				}

				// Expect modification or action
				modification, err = milter.AnswerToModification(msgType, msg)
				if err == nil {

					// Dump sent information
					fmt.Printf("< BODYEOB %s\n", modification.String())

					// Send answer to client
					err = srv.SendModification(modification)
					if err != nil {
						log.Printf("Client side error: %s", err.Error())
						return
					}

				} else {
					action, err = milter.AnswerToAction(msgType, msg)
					if err != nil {
						log.Printf("Server side error: %s", err.Error())
						return
					}
					break
				}
			}

			// Dump server response information
			fmt.Printf("< BODYEOB %s\n", action.String())

			// Send answer to client
			err = srv.SendAction(action)
			if err != nil {
				log.Printf("Client side error: %s", err.Error())
				return
			}

		case milter.SMFIC_QUIT:

			// Dump received information
			fmt.Printf("> QUIT\n")

			// Forward request to milter server
			err = cli.SendQuit()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

			// end of connexion
			return

		case milter.SMFIC_ABORT:

			// Dump received information
			fmt.Printf("> ABORT\n")

			// Forward request to milter server
			err = cli.SendAbort()
			if err != nil {
				log.Printf("Server side error: %s", err.Error())
				return
			}

		case milter.SMFIC_MACRO:

			// Cast message according with protocol
			macros = msg.([]*milter.Macro)

			// Dump received information
			for _, macro = range macros {
				fmt.Printf("> MACRO %s\n", macro.String())
			}

			// Forward request to milter server
			if len(macros) > 0 {
				err = cli.SendMacro(macros[0].Step, macros)
				if err != nil {
					log.Printf("Server side error: %s", err.Error())
					return
				}
			}
		}
	}
}

func Example_sendReceiveProxy() {
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
	proxy_send(conn)

	// Output: 
}
