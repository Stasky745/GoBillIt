package main

import (
	"os"

	"github.com/Stasky745/GoBillIt/internal/email"
	"github.com/Stasky745/GoBillIt/internal/invoicegenerator"
	"github.com/Stasky745/GoBillIt/internal/ntfy"
	"github.com/Stasky745/GoBillIt/internal/utils"
	"github.com/Stasky745/go-libs/log"
	"github.com/urfave/cli/v2"
)

var version string

var baseFlags = []cli.Flag{
	// base
	&cli.StringFlag{
		Name:  "config",
		Usage: "config file to load",
	},
}

var app = &cli.App{
	Name:                 "GoBillIt",
	Version:              version,
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
					conversion := k.Float64("inv.conversion.value")
					if conversion > 0 {
						for i, item := range extraItems {
							extraItems[i].Unit_cost = utils.GetConvertedCost(item.Unit_cost, conversion)
						}
					}
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
			client.SendNotification(5, "Hooray!", "We are done!!", append(EMOJI_TADA, TAGS...), []ntfy.Action{}, "", "")
		}

		return nil
	},
}
