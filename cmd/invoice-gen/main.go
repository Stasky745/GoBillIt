package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Stasky745/invoice-generator/apilayer"
	gomail "github.com/Stasky745/invoice-generator/email"
	"github.com/Stasky745/invoice-generator/invoicegenerator"
	"github.com/Stasky745/invoice-generator/utils"
)

var (
	invgenApiKey    string
	invTo           string
	invFrom         string
	invNumberPrefix string
	invBankAcc      string
	invAmount       string

	companyName string

	emailFrom string
	emailTo   string

	smtpServer   string
	smtpPort     int
	smtpUsername string
	smtpPassword string

	apiLayerApiKey   string
	apiLayerBaseCurr string
	apiLayerNewCurr  string

	pdfPrefix string

	conversionRate string
)

func getArgs() {
	flag.StringVar(&invgenApiKey, "invgenApiKey", os.Getenv("INVGEN_API_KEY"), "API key for invoice-generator.com")
	flag.StringVar(&invTo, "invTo", os.Getenv("INV_TO"), "\"Bill To\" part of the invoice.")
	flag.StringVar(&invFrom, "invFrom", os.Getenv("INV_FROM"), "API key for invoice-generator.com")
	flag.StringVar(&invNumberPrefix, "invNumberPrefix", os.Getenv("INV_NUMBER_PREFIX"), "Invoice number.")
	flag.StringVar(&invBankAcc, "invBankAcc", os.Getenv("INV_BANK_ACCOUNT"), "The bank account it has to be paid to.")
	flag.StringVar(&invAmount, "invAmount", os.Getenv("INV_AMOUNT"), "The amount to bill.")

	flag.StringVar(&companyName, "companyName", os.Getenv("COMPANY_NAME"), "Name of the company being invoiced.")

	flag.StringVar(&smtpServer, "smtpServer", os.Getenv("SMTP_SERVER"), "SMTP Server to be used.")
	smtpPortEnv, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPortEnv = 1
	}
	flag.IntVar(&smtpPort, "smtpPort", smtpPortEnv, "SMTP Port to be used.")
	flag.StringVar(&smtpUsername, "smtpUsername", os.Getenv("SMTP_USERNAME"), "SMTP Username to be used.")
	flag.StringVar(&smtpPassword, "smtpPassword", os.Getenv("SMTP_PASSWORD"), "SMTP Password to be used.")

	flag.StringVar(&emailFrom, "emailFrom", os.Getenv("EMAIL_FROM"), "Where to send the invoice from.")
	flag.StringVar(&emailTo, "emailTo", os.Getenv("EMAIL_TO"), "Where to send the invoice to.")

	flag.StringVar(&apiLayerApiKey, "apiLayerApiKey", os.Getenv("APILAYER_API_KEY"), "API Key for apilayer.com")
	flag.StringVar(&apiLayerBaseCurr, "apiLayerBaseCurr", os.Getenv("APILAYER_BASE_CURRENCY"), "Base currency for conversion.")
	flag.StringVar(&apiLayerNewCurr, "apiLayerNewCurr", os.Getenv("APILAYER_NEW_CURRENCY"), "Currency to be converted to.")

	flag.StringVar(&pdfPrefix, "pdfPrefix", os.Getenv("PDF_PREFIX"), "Prefix of generated PDF.")

	flag.StringVar(&conversionRate, "conversionRate", os.Getenv("CONVERSION_RATE"), "Conversion rate to be applied.")

	flag.Parse()
}

func main() {
	getArgs()

	r := new(invoicegenerator.Invoice)

	r.From = invFrom
	r.To = invTo

	date := utils.GetLastDayCurrentMonth()

	// r.Date = utils.FormatDate(date)
	// r.DueDate = utils.FormatDate(utils.GetLastDayNextMonth())
	r.Date = utils.FormatDate(time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC))
	r.DueDate = utils.FormatDate(date)

	yyyymm := time.Now().Format("200601")

	if invNumberPrefix != "" {
		r.Number = invNumberPrefix + "-" + yyyymm + "-001"
	} else {
		r.Number = yyyymm + "-001"
	}

	items := make([]invoicegenerator.Item, 2)
	item := invoicegenerator.Item{}
	item.Name = utils.GetCurrentMonthName() + " Services"
	item.Quantity = 1

	amount, err := strconv.ParseFloat(invAmount, 64)
	if err != nil {
		panic(err)
	}

	notes := "Service Agreement amount: " + apiLayerBaseCurr + " " + utils.FormatFloatToAmount(amount) + "."

	var conversion float64
	if apiLayerBaseCurr != apiLayerNewCurr {

		if conversionRate == "" {
			conversion, err = apilayer.GetRate(apiLayerApiKey, apiLayerNewCurr, apiLayerBaseCurr)
			if err != nil {
				panic(err)
			}
		} else {
			conversion, err = strconv.ParseFloat(conversionRate, 64)
			if err != nil {
				panic(err)
			}
		}
		amount = utils.GetConvertedCost(amount, conversion)

		notes += "\nThe payment shall be made in " + apiLayerNewCurr + " based on the " + apiLayerBaseCurr + "-" + apiLayerNewCurr + " currency exchange rate for the last day of the service month.\nFx: " + strconv.FormatFloat(conversion, 'f', -1, 64)
	}

	notes += "\nBank account:  " + invBankAcc

	item.Unit_cost = amount

	items[0] = item

	item2 := invoicegenerator.Item{}
	item2.Name = "On-Call"
	item2.Quantity = 1
	item2.Unit_cost = utils.GetConvertedCost(210, conversion)

	items[1] = item2

	r.Items = items

	r.Notes = notes

	if apiLayerNewCurr != "" {
		r.Currency = apiLayerNewCurr
	} else {
		r.Currency = apiLayerBaseCurr
	}

	filename := "./" + pdfPrefix + "_" + yyyymm + "_" + r.Number + ".pdf"

	err = r.Create(invgenApiKey, filename)

	if err != nil {
		panic(err)
	}

	// sender := email.NewSender(smtpServer, smtpPort, smtpUsername, smtpPassword)
	// m := email.NewMessage(utils.GetCurrentMonthName()+" invoice for: "+companyName, "Hello!\n\nSee attached here the invoice for this month of "+utils.GetCurrentMonthName()+".\n\nCheers!")
	// m.To = []string{emailTo}
	// m.AttachFile(filename)

	// fmt.Println(sender.Send(m))

	email := gomail.Email{
		From:         emailFrom,
		To:           emailTo,
		Subject:      utils.GetCurrentMonthName() + " invoice for: " + companyName,
		Body:         "Hello!\n\nSee attached here the invoice for this month of " + utils.GetCurrentMonthName() + ".\n\nCheers!",
		SmtpServer:   smtpServer,
		SmtpPort:     smtpPort,
		SmtpUsername: smtpUsername,
		SmtpPassword: smtpPassword,
		Attachment:   filename,
	}

	fmt.Println(gomail.SendEmail(email))
}
