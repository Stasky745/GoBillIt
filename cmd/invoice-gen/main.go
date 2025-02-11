package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Stasky745/invoice-generator/apilayer"
	gomail "github.com/Stasky745/invoice-generator/email"
	"github.com/Stasky745/invoice-generator/invoicegenerator"
	"github.com/Stasky745/invoice-generator/utils"
)

func main() {
	config := InitFlags()

	r := new(invoicegenerator.Invoice)

	r.From = config.invFrom
	r.To = config.invTo

	date := utils.GetLastDayCurrentMonth()

	// r.Date = utils.FormatDate(date)
	// r.DueDate = utils.FormatDate(utils.GetLastDayNextMonth())
	r.Date = utils.FormatDate(time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC))
	r.DueDate = utils.FormatDate(date)

	yyyymm := time.Now().Format("200601")

	if config.invNumberPrefix != "" {
		r.Number = config.invNumberPrefix + "-" + yyyymm + "-001"
	} else {
		r.Number = yyyymm + "-001"
	}

	items := make([]invoicegenerator.Item, 2)
	item := invoicegenerator.Item{}
	item.Name = utils.GetCurrentMonthName() + " Services"
	item.Quantity = 1

	amount, err := strconv.ParseFloat(config.invAmount, 64)
	if err != nil {
		panic(err)
	}

	notes := "Service Agreement amount: " + config.apiLayerBaseCurr + " " + utils.FormatFloatToAmount(amount) + "."

	var conversion float64
	if config.apiLayerBaseCurr != config.apiLayerNewCurr {

		if config.conversionRate == "" {
			conversion, err = apilayer.GetRate(config.apiLayerApiKey, config.apiLayerNewCurr, config.apiLayerBaseCurr)
			if err != nil {
				panic(err)
			}
		} else {
			conversion, err = strconv.ParseFloat(config.conversionRate, 64)
			if err != nil {
				panic(err)
			}
		}
		amount = utils.GetConvertedCost(amount, conversion)

		notes += "\nThe payment shall be made in " + config.apiLayerNewCurr + " based on the " + config.apiLayerBaseCurr + "-" + config.apiLayerNewCurr + " currency exchange rate for the last day of the service month.\nFx: " + strconv.FormatFloat(conversion, 'f', -1, 64)
	}

	notes += "\nBank account:  " + config.invBankAcc

	item.Unit_cost = amount

	items[0] = item

	item2 := invoicegenerator.Item{}
	item2.Name = "On-Call"
	item2.Quantity = 1
	item2.Unit_cost = utils.GetConvertedCost(210, conversion)

	items[1] = item2

	r.Items = items

	r.Notes = notes

	if config.apiLayerNewCurr != "" {
		r.Currency = config.apiLayerNewCurr
	} else {
		r.Currency = config.apiLayerBaseCurr
	}

	filename := "./" + config.pdfPrefix + "_" + yyyymm + "_" + r.Number + ".pdf"

	err = r.Create(config.invgenApiKey, filename)

	if err != nil {
		panic(err)
	}

	// sender := email.NewSender(smtpServer, smtpPort, smtpUsername, smtpPassword)
	// m := email.NewMessage(utils.GetCurrentMonthName()+" invoice for: "+companyName, "Hello!\n\nSee attached here the invoice for this month of "+utils.GetCurrentMonthName()+".\n\nCheers!")
	// m.To = []string{emailTo}
	// m.AttachFile(filename)

	// fmt.Println(sender.Send(m))

	email := gomail.Email{
		From:         config.emailFrom,
		To:           config.emailTo,
		Subject:      utils.GetCurrentMonthName() + " invoice for: " + config.companyName,
		Body:         "Hello!\n\nSee attached here the invoice for this month of " + utils.GetCurrentMonthName() + ".\n\nCheers!",
		SmtpServer:   config.smtpServer,
		SmtpPort:     config.smtpPort,
		SmtpUsername: config.smtpUsername,
		SmtpPassword: config.smtpPassword,
		Attachment:   filename,
	}

	fmt.Println(gomail.SendEmail(email))
}
