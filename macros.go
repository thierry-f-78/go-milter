// Copyright (c) 2022 Thierry FOURNIER (tfournier@arpalert.org)

package milter

import "fmt"

type Macro struct {
	Step MacroStep
	Name string
	Value string
}

func (m *Macro)String()(string) {
	var msgType MsgType

	msgType = MsgType(m.Step)
	return fmt.Sprintf("step=%s, name=%s, value=%s", msgType.String(), qt(m.Name), qt(m.Value))
}

func macroAdd(macros *[]*Macro, step MacroStep, name string, value string)() {
	var m *Macro

	/* lookup for existing macro. Do not append macros if already exists. */
	for _, m = range *macros {
		if m.Name == name {
			return
		}
	}

	m = &Macro{}
	m.Step = step
	m.Name = name
	m.Value = value
	*macros = append(*macros, m)
}

func macroGet(macros []*Macro, name string)(MacroStep, string) {
	var m *Macro

	for _, m = range macros {
		if m.Name == name {
			return m.Step, m.Value
		}
	}
	return 0, ""
}

func macroDebug(macros []*Macro)() {
	var m *Macro

	for _, m = range macros {
		fmt.Printf("%s\n", m.String())
	}
}

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

func macroAdd__(macros *[]*Macro, value string)            { macroAdd(macros, MS_CONNECT, "_", value) }
func macroAdd_j(macros *[]*Macro, value string)            { macroAdd(macros, MS_CONNECT, "j", value) }
func macroAdd_daemon_name(macros *[]*Macro, value string)  { macroAdd(macros, MS_CONNECT, "{daemon_name}", value) }
func macroAdd_if_name(macros *[]*Macro, value string)      { macroAdd(macros, MS_CONNECT, "{if_name}", value) }
func macroAdd_if_addr(macros *[]*Macro, value string)      { macroAdd(macros, MS_CONNECT, "{if_addr}", value) }

func macroAdd_tls_version(macros *[]*Macro, value string)  { macroAdd(macros, MS_HELO,    "{tls_version}", value) }
func macroAdd_cipher(macros *[]*Macro, value string)       { macroAdd(macros, MS_HELO,    "{cipher}", value) }
func macroAdd_cipher_bits(macros *[]*Macro, value string)  { macroAdd(macros, MS_HELO,    "{cipher_bits}", value) }
func macroAdd_cert_subject(macros *[]*Macro, value string) { macroAdd(macros, MS_HELO,    "{cert_subject}", value) }
func macroAdd_cert_issuer(macros *[]*Macro, value string)  { macroAdd(macros, MS_HELO,    "{cert_issuer}", value) }

func macroAdd_i(macros *[]*Macro, value string)            { macroAdd(macros, MS_MAIL,    "i", value) }
func macroAdd_auth_type(macros *[]*Macro, value string)    { macroAdd(macros, MS_MAIL,    "{auth_type}", value) }
func macroAdd_auth_authen(macros *[]*Macro, value string)  { macroAdd(macros, MS_MAIL,    "{auth_authen}", value) }
func macroAdd_auth_ssf(macros *[]*Macro, value string)     { macroAdd(macros, MS_MAIL,    "{auth_ssf}", value) }
func macroAdd_auth_author(macros *[]*Macro, value string)  { macroAdd(macros, MS_MAIL,    "{auth_author}", value) }
func macroAdd_mail_mailer(macros *[]*Macro, value string)  { macroAdd(macros, MS_MAIL,    "{mail_mailer}", value) }
func macroAdd_mail_host(macros *[]*Macro, value string)    { macroAdd(macros, MS_MAIL,    "{mail_host}", value) }
func macroAdd_mail_addr(macros *[]*Macro, value string)    { macroAdd(macros, MS_MAIL,    "{mail_addr}", value) }

func macroAdd_rcpt_mailer(macros *[]*Macro, value string)  { macroAdd(macros, MS_RCPT,    "{rcpt_mailer}", value) }
func macroAdd_rcpt_host(macros *[]*Macro, value string)    { macroAdd(macros, MS_RCPT,    "{rcpt_host}", value) }
func macroAdd_rcpt_addr(macros *[]*Macro, value string)    { macroAdd(macros, MS_RCPT,    "{rcpt_addr}", value) }
