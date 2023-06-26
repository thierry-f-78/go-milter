[![GoDoc](https://pkg.go.dev/badge/github.com/thierry-f-78/go-milter)](https://pkg.go.dev/github.com/thierry-f-78/go-milter)

MILTER client/server library for the GO programming language.
-------------------------------------------------------------

This library implement milter protocol. It have no dependencies except standard
Go library. The library implement version 2 of the protocol. The next versions
doesn't provides documentation (to my knowledge). There are mainly Postfix
proposal. Note, according with original recommendations, the library accept
macros at any step, and this is compatible with postfix.

The library propose 3 ways to works:

- `Exchange*` function handle protocol I/O, call the right function for the
  right message and the API expect the right response for each message. This
  is the most simple way to use the library and it answer to the majority of
  use cases.
  
- `Send*`/ `Receive*` functions handleprotocol I/O, but the user choose the
  right answer to each request. This way allow a  lot of flexibility, offloading
  the user to code borring things.
  
- `Encode*` / `Decode*` function allow user to manager network connexion and
  I/O. It uses these function to encode/decode milter messages. This is the
  most complicated way to use the library, but it leaves a lot of freedom to
  the developer.

Check the documentation whoch contains exemples to understand how using the
library.

Milter protocol overview
------------------------

Milter is at the same time the protocole name avec the server name which offer
this protocole to communicate with it. Its a packet based binary protocol over
connecte protocol like TCP or Unix socket.

The connection could be closed at any time by the client or the server. If the
milter server closes connection before sending accept decision, the client
must apply its default behavior (accept, or reject).

Below, classic and complete communication between client and server using
library terminology.

| SMTP Client        | `*milter.Client` request | `*milter.Server` response | Notes                                                                              |
|--------------------|--------------------------|---------------------------|------------------------------------------------------------------------------------|
| `> TCP connection` | `ExchangeOptNeg()`       |                           | Send protocol and actions negotiation when connection is established               |
|                    |                          | `SendOptNeg()`            | Answer protocol and actions negociation                                            |
|                    | `ExchangeConnect()`      |                           | Send information relative to Client connection. Technically send also macro        |
| `< 220`            |                          | `SendAction()`            | Answer with continuation decision. Answer 220 or reject to the SMTP client         |
| `> HELO`           | `ExchangeHelo()`         |                           | Send information about HELO step. Technically send also macro                      |
| `< 250`            |                          | `SendAction()`            | Answer with continuation decision. Answer 250 or reject to the SMTP client         |
| `> MAIL FROM`      | `ExchangeMail()`         |                           | Send information about MAIL FROM to the milter. Technically send also macro        |
| `< 250`            |                          | `SendAction()`            | Answer with continuation decision. Answer 250 or reject to the SMTP client         |
| `> RCPT TO`        | `ExchangeRcpt()`         |                           | Send mulitple information about RCPT TO to the milter. Technically send also macro |
| `< 250`            |                          | `SendAction()`            | Answer with continuation decision. Answer 250 or reject to the SMTP client         |
| `> DATA`           | `ExchangeHeader()`       |                           | Receive multiple time information about each header.                               |
|                    |                          | `SendAction()`            | Answer with continuation decision. Send reject to the client if required           |
|                    | `ExchangeEOH()`          |                           | Receive information about last header was sent. No more headers to process.        |
|                    |                          | `SendAction()`            | Answer with continuation decision. Send reject to the client if required           |
|                    | `ExchangeBody()`         |                           | Receive multiple time body content.                                                |
|                    |                          | `SendAction()`            | Answer with continuation decision. Send reject to the client if required           |
|                    | `ExchangeBodyEOB()`      |                           | Receive information about body sent complete.                                      |
|                    |                          | `SendModification()`      | Sent mulitple modifications (header, body, recipients, ...                         |
| `< 250`            |                          | `SendAction()`            | Answer with continuation decision. Answer 250 or reject to the SMTP client         |
| `> QUIT`           | `ExchangeQuit()`         |                           | Close connexion between MTa and milter                                             |

Messages
--------

Milter client / server send messages which have four roles:

- `proto`: Message about milter protocol itself, like negociate capabilities or quit connexion.
- `info`: Information about current SMTP email exchange between client and MTA.
- `action`: Actions about message processing send by milter server to MTA client.
  The client should follow directives, but it chooses.
- `modification`: Email modifications ask by milter server to MTA client.
  The MTA client should follow directives, but it chooses.


| Message      | Origin          | Role         | description
|--------------|-----------------|--------------|------------
| `OPTNEG`     | client / server | proto        | Option negotiation.
| `MACRO`      | client          | info         | Define macros. Doesn't expect response.
| `CONNECT`    | client          | info         | SMTP connection information. Expect action.
| `HELO`       | client          | info         | `HELO` / `EHLO` name. Expect action.
| `MAIL`       | client          | info         | `MAIL FROM` information. Expect action.
| `RCPT`       | client          | info         | `RCPT TO` information. Expect action.
| `HEADER`     | client          | info         | Mail header. Expect action.
| `EOH`        | client          | info         | End of headers marker. Expect action.
| `BODY`       | client          | info         | Body chunk. Max size of 65535 bytes. Expect action.
| `BODYEOB`    | client          | info         | Final body chunk. Expect action.
| `ABORT`      | client          | proto        | Abort current filter checks. Resets internal state of milter program to before HELO, but keeps the connection open. Doesn't expect response.
| `QUIT`       | client          | proto        | Quit milter communication. Doesn't expect response.
| `ACCEPT`     | server          | action       | Accept message completely. This will skip to the end of the milter sequence, and recycle back to `MAIL` step. The MTA may, instead, close the connection at that point.
| `CONTINUE`   | server          | action       | Accept and keep processing. If issued at the end of the milter conversation, functions the same as ACCEPT.
| `DISCARD`    | server          | action       | Set discard flag for entire message. Note that message processing MAY continue afterwards, but the mail will not be delivered even if accepted with ACCEPT.
| `REJECT`     | server          | action       | Reject with a 5xx SMTP code.
| `TEMPFAIL`   | server          | action       | Reject with a 4xx SMTP code.
| `REPLYCODE`  | server          | action       | Send specific SMTP code and reply message.
| `ADDRCPT`    | server          | modification | Milter server ask to MTA to add recipient.
| `DELRCPT`    | server          | modification |  Milter server ask to MTA to remove recipient.
| `REPLBODY`   | server          | modification |  Milter server ask to MTA to replace body.
| `ADDHEADER`  | server          | modification |  Milter server ask to MTA to add header.
| `CHGHEADER`  | server          | modification |  Milter server ask to MTA to change header (remove it if the value is empty).
| `QUARANTINE` | server          | modification |  Milter server ask to MTA to put email in quarantine.

Macros
------

Macros are informations added by the MTA which is complementary to the SMTP
protocol exchange. With version 1 of the protocol, macros are sent on the
following steps. Postfix adds macros to other step.

| step      | macros
|-----------|-----------
| `CONNECT` | `_` `j` `{daemon_name}` `{if_name}` `{if_addr}`
| `HELO`    | `{tls_version}` `{cipher}` `{cipher_bits}` `{cert_subject}` `{cert_issuer}
| `MAIL`    | `i` `{auth_type}` `{auth_authen}` `{auth_ssf}` `{auth_author}` `{mail_mailer}` `{mail_host}` `{mail_addr}`
| `RCPT`    | `{rcpt_mailer}` `{rcpt_host}` `{rcpt_addr}`