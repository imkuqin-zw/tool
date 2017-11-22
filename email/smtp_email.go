package email

import (
	"net/smtp"
	"strings"
)

type Email struct {
	To			[]string
	Subject		string
	Body		string
	MailType	string
}

type MailWithDeal struct {
	Email	*Email
	Deal	*IDeal
}

type IDeal interface {
	Run()
	SetErr(err error)
}

type Deal struct {
	Err error
	IDeal
}

func (d *Deal) SetErr(err error) {
	d.Err = err
}

var (
	MailChannel = make(chan MailWithDeal, 50)
	DealChannel = make(chan IDeal, 50)

)

func SendMailLoop(addr, user, pass, from, nickname string) {
	for {
		emailWithDeal := <- MailChannel
		host := strings.Split(addr, ":")[0]
		auth := smtp.PlainAuth("", user, pass, host)
		contentType := getContentType(emailWithDeal.Email.MailType)
		msg := getMsg(emailWithDeal.Email.To, nickname, from, emailWithDeal.Email.Subject, contentType, emailWithDeal.Email.Body)
		err := smtp.SendMail(addr, auth, user, emailWithDeal.Email.To, msg)
		if emailWithDeal.Deal != nil {
			if err != nil {
				(*emailWithDeal.Deal).SetErr(err)
				DealChannel <- *emailWithDeal.Deal
			} else {
				DealChannel <- *emailWithDeal.Deal
			}
		}
	}
}

func DealLoop() {
	for {
		deal := <-DealChannel
		deal.Run()
	}
}

func SendMailAndDeal(addr, user, pass, from, nickname string) {
	go SendMailLoop(addr, user, pass, from, nickname)
	go DealLoop()
}

func getMsg(to []string, nickName, from, subject, contentType, body string) (msg []byte) {
	content := "to: " + strings.Join(to, ";") + "\r\n"
	content += "From: " + nickName + "<" + from + ">\r\n"
	content += "Subject: " + subject + "\r\n"
	content += contentType + "\r\n\r\n"
	content += body
	msg = []byte(content)
	return
}

func getContentType(mailType string) (contentType string) {
	switch mailType {
	case "html":
		contentType = "Content-Type: text/html; charset=UTF-8"
		break
	default:
		contentType = "Content-Type: text/plain; charset=UTF-8"
	}
	return
}
