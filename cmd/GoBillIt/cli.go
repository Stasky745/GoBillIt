package main

import (
	"os"

	"github.com/Stasky745/GoBillIt/internal/email"
	"github.com/Stasky745/GoBillIt/internal/invoicegenerator"
	"github.com/Stasky745/GoBillIt/internal/ntfy"
	"github.com/Stasky745/go-libs/log"
	"github.com/urfave/cli/v2"
)

var baseFlags = []cli.Flag{
	// base
	&cli.StringFlag{
		Name:  "config",
		Usage: "config file to load",
	},
	// // invoice
	// &cli.StringFlag{
	// 	Name:  "inv-apikey",
	// 	Usage: "invoice-generator.com API key",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-to",
	// 	Usage: "who the invoice is directed to",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-from",
	// 	Usage: "who is the invoice from",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-date",
	// 	Usage: "date the invoice is issued",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-duedate",
	// 	Usage: "duedate of the invoice",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-dateformat",
	// 	Usage: "format of the dates shown (see .env.example)",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-path",
	// 	Usage: "path to generate the invoice (can use templating)",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-number",
	// 	Usage: "invoice number (can use templating)",
	// },
	// &cli.StringFlag{
	// 	Name:  "inv-items-path",
	// 	Usage: "path to file containing items",
	// },
	// &cli.Float64Flag{
	// 	Name:  "inv-conversion",
	// 	Usage: "conversion rate",
	// },

	// // ApiLayer flags

	// &cli.BoolFlag{
	// 	Name:  "apilayer-enabled",
	// 	Usage: "use APILayer to manage conversions (overrides 'inv-conversion' if set)",
	// },
	// &cli.StringFlag{
	// 	Name:  "apilayer-apikey",
	// 	Usage: "APILayer API key",
	// },
	// &cli.StringFlag{
	// 	Name:  "apilayer-currency-base",
	// 	Usage: "Base currency",
	// },
	// &cli.StringFlag{
	// 	Name:  "apilayer-currency-new",
	// 	Usage: "New currency",
	// },
}

var app = &cli.App{
	Name:                 "GoBillIt",
	Version:              "v1.0.0",
	EnableBashCompletion: true,
	Flags:                baseFlags,
	Action: func(c *cli.Context) error {
		loadConfig(c)
		invoice, filename := createInvoice()
		invoice.CreatePDF(k.String("inv.apikey"), filename)

		ntfyEnabled := k.Bool("ntfy.enabled")
		ntfyCode := -1
		var err error
		for ntfyEnabled && NTFY_CODE_DONE != ntfyCode {
			ntfyCode, err = ntfyCheck(filename)
			if err != nil {
				ntfyEnabled = false
			}

			switch ntfyCode {
			case NTFY_CODE_RECREATE:
				extraItems := []invoicegenerator.Item{}
				err = k.Unmarshal(EXTRA_ITEMS, &extraItems)
				if !log.CheckErr(err, false, "can't unmarshall extra items") {
					invoice.Items = append(invoice.Items, extraItems...)
				}
				invoice.CreatePDF(k.String("inv.apikey"), filename)
			case NTFY_CODE_CANCEL:
				log.Info("User cancelled operation. Deleting file.")
				err := os.Remove(filename)
				log.CheckErr(err, false, "can't delete file", "file", filename)
				os.Exit(0)
			}
		}

		if k.Bool("email.enabled") {
			var e email.Email
			err := k.Unmarshal("email", &e)
			if log.CheckErr(err, false, "can't unmarshall email") {
				return err
			}

			e.Attachment = filename

			// use templating for subject and body
			e.Body = template(e.Body, map[string]string{})
			e.Subject = template(e.Subject, map[string]string{})

			sendEmail := true
			if k.Bool("email.ntfy.check") {
				sendEmail, err = ntfyEmailCheck(e)
				if log.CheckErr(err, false, "failed email check") {
					client.SendNotification(5, "ERROR Email Check", "There was an error sending notification for email check. Check logs for more info.", append(EMOJI_ERROR, TAGS...), []ntfy.Action{}, "", "")
					sendEmail = false
				}
			}

			if sendEmail {
				err := email.SendEmail(e)
				if log.CheckErr(err, false, "failed email check") {
					client.SendNotification(5, "ERROR Sending Email", "There was an error sending email. Check logs for more info.", append(EMOJI_ERROR, TAGS...), []ntfy.Action{}, "", "")
				} else {
					client.SendNotification(5, "Mission Complete!", "Email sent succesfully!", append(append(EMOJI_SUNGLASSES, EMOJI_MONEY_WITH_WINGS...), TAGS...), []ntfy.Action{}, "", "")
				}
			}
		}

		if ntfyEnabled {
			client.SendNotification(5, "Hooray!", "We are done!", append(EMOJI_TADA, TAGS...), []ntfy.Action{}, "", "")
		}

		return nil
	},
}
