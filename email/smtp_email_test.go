package email

import (
	"fmt"
	"testing"
)

type QqCoinDeal struct {
	Deal
}

func (deal *QqCoinDeal) Run() {
	if deal.Err == nil {
		fmt.Println("发送成功")
	} else {
		fmt.Println(deal.Err.Error())
	}
}

func qqCoinNotify() {
	var addr= ""
	var user= ""
	var pass= ""
	var from= ""
	var nickname = ""
	SendMailAndDeal(addr, user, pass, from, nickname)
}

func checkQqCoinBalance() {
	var deal IDeal
	qqDeal := &QqCoinDeal{}
	deal = qqDeal
	mailWithDeal := MailWithDeal{
		Email: &Email{
			To: []string{"zhangwei@kuaifazs.com"},
			Subject: "test发送邮件",
			MailType: "html",
		},
		Deal: &deal,
	}
	body := "test"
	mailWithDeal.Email.Body = body
	MailChannel <- mailWithDeal
}

func Test_SendAndDealEmail(t *testing.T) {
	qqCoinNotify()
	checkQqCoinBalance()
	for true {

	}
}