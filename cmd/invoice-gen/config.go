package main

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	invgenApiKey    string
	invTo           string
	invFrom         string
	invNumberPrefix string
	invBankAcc      string
	invAmount       string

	companyName string

	smtpServer   string
	smtpPort     int
	smtpUsername string
	smtpPassword string

	emailFrom string
	emailTo   string

	apiLayerApiKey   string
	apiLayerBaseCurr string
	apiLayerNewCurr  string

	pdfPrefix string

	conversionRate string
}

func InitFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.invgenApiKey, "invgenApiKey", os.Getenv("INVGEN_API_KEY"), "API key for invoice-generator.com")
	flag.StringVar(&config.invTo, "invTo", os.Getenv("INV_TO"), "\"Bill To\" part of the invoice.")
	flag.StringVar(&config.invFrom, "invFrom", os.Getenv("INV_FROM"), "API key for invoice-generator.com")
	flag.StringVar(&config.invNumberPrefix, "invNumberPrefix", os.Getenv("INV_NUMBER_PREFIX"), "Invoice number.")
	flag.StringVar(&config.invBankAcc, "invBankAcc", os.Getenv("INV_BANK_ACCOUNT"), "The bank account it has to be paid to.")
	flag.StringVar(&config.invAmount, "invAmount", os.Getenv("INV_AMOUNT"), "The amount to bill.")

	flag.StringVar(&config.companyName, "companyName", os.Getenv("COMPANY_NAME"), "Name of the company being invoiced.")

	flag.StringVar(&config.smtpServer, "smtpServer", os.Getenv("SMTP_SERVER"), "SMTP Server to be used.")
	smtpPortEnv, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPortEnv = 1
	}
	flag.IntVar(&config.smtpPort, "smtpPort", smtpPortEnv, "SMTP Port to be used.")
	flag.StringVar(&config.smtpUsername, "smtpUsername", os.Getenv("SMTP_USERNAME"), "SMTP Username to be used.")
	flag.StringVar(&config.smtpPassword, "smtpPassword", os.Getenv("SMTP_PASSWORD"), "SMTP Password to be used.")

	flag.StringVar(&config.emailFrom, "emailFrom", os.Getenv("EMAIL_FROM"), "Where to send the invoice from.")
	flag.StringVar(&config.emailTo, "emailTo", os.Getenv("EMAIL_TO"), "Where to send the invoice to.")

	flag.StringVar(&config.apiLayerApiKey, "apiLayerApiKey", os.Getenv("APILAYER_API_KEY"), "API Key for apilayer.com")
	flag.StringVar(&config.apiLayerBaseCurr, "apiLayerBaseCurr", os.Getenv("APILAYER_BASE_CURRENCY"), "Base currency for conversion.")
	flag.StringVar(&config.apiLayerNewCurr, "apiLayerNewCurr", os.Getenv("APILAYER_NEW_CURRENCY"), "Currency to be converted to.")

	flag.StringVar(&config.pdfPrefix, "pdfPrefix", os.Getenv("PDF_PREFIX"), "Prefix of generated PDF.")

	flag.StringVar(&config.conversionRate, "conversionRate", os.Getenv("CONVERSION_RATE"), "Conversion rate to be applied.")

	flag.Parse()

	return config
}
