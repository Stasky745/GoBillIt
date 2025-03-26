package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/Stasky745/GoBillIt/internal/apilayer"
	"github.com/Stasky745/GoBillIt/internal/invoicegenerator"
	"github.com/Stasky745/GoBillIt/internal/utils"
	"github.com/Stasky745/go-libs/log"
	"gopkg.in/yaml.v2"
)

func createInvoice() (*invoicegenerator.Invoice, string) {
	var err error

	r := new(invoicegenerator.Invoice)

	r.From = k.String("inv.from")
	r.To = k.String("inv.to")

	date := time.Now()

	if d := k.String("inv.date"); d != "" {
		date, err = utils.GetDate(d)
		log.CheckErr(err, true, "can't parse date", "date", d)
	}

	dueDate := date.AddDate(0, 1, 0)

	if dd := k.String("inv.duedate"); dd != "" {
		dueDate, err = utils.GetDate(dd)
		log.CheckErr(err, true, "can't parse duedate", "date", dd)
	}

	dateLayout := "m D, YYYY"
	if dl := k.String("inv.dateformat"); dl != "" {
		dateLayout = dl
	}
	r.Date = utils.FormatDate(date, dateLayout)
	r.DueDate = utils.FormatDate(dueDate, dateLayout)

	// Deal with seq template
	filename := ""
	seq := fmt.Sprintf("%0*d", 3, 1)
	seqExists, seqWidthString := templateGetKeyParams(k.String("inv.path"), "seq")
	if seqExists {
		seqWidth := 3
		if seqWidthString != "" {
			seqWidth, err = strconv.Atoi(seqWidthString)
			log.CheckErr(err, true, "seq for invoice number/path is not a number", "sequence parameter", seqWidthString)
		}
		seqInt := 1
		seq = fmt.Sprintf("%0*d", seqWidth, seqInt)
		filename = template(k.String("inv.path"), map[string]string{"seq": seq})
		for utils.FileExists(filename) {
			seqInt += 1
			seq = fmt.Sprintf("%0*d", seqWidth, seqInt)
			filename = template(k.String("inv.path"), map[string]string{"seq": seq})
		}
	}

	invNumber := template(k.String("inv.number"), map[string]string{"seq": seq})

	if filename == "" {
		filename = fmt.Sprintf("./%s.pdf", invNumber)
	}

	r.Number = template(invNumber, map[string]string{})

	// Fetch the items
	itemsFile := k.String("inv.items.path")
	items := []invoicegenerator.Item{}
	itemsFromFile := []invoicegenerator.Item{}
	err = k.Unmarshal("inv.items.list", &items)
	if log.CheckErr(err, false, "can't unmarshal list of items", "items", k.Strings("inv.items.list")) {
		items = []invoicegenerator.Item{}
	}

	if itemsFile != "" {
		content, err := os.ReadFile(itemsFile)
		if !log.CheckErr(err, false, "can't get content from file", "file", itemsFile) {
			err := yaml.Unmarshal(content, &itemsFromFile)
			log.CheckErr(err, false, "can't unmarshal content from file", "file", itemsFile)
		}
	}
	items = append(items, itemsFromFile...)
	applyTemplateToItems(items)

	// Conversion if necessary
	conversion := k.Float64("inv.conversion.value")
	if k.Bool("apilayer.enabled") &&
		k.String("apilayer.apikey") != "" &&
		k.String("apilayer.currency.base") != "" &&
		k.String("apilayer.currency.new") != "" {
		api := k.String("apilayer.apikey")
		baseCurrency := k.String("apilayer.currency.base")
		newCurrency := k.String("apilayer.currency.new")
		conv, err := apilayer.GetRate(api, baseCurrency, newCurrency)
		if !log.CheckErr(err, false, "couldn't get conversion from API Layer, will use default value", "base currency", baseCurrency, "new currency", newCurrency, "default conversion", conversion) {
			conversion = conv
		}
	}

	// Set conversion to be the largest between new value and min conversion set
	conversion = math.Max(conversion, k.Float64("inv.conversion.min"))

	// Set new conversion value into koanf
	err = k.Set("inv.conversion.value", conversion)
	log.CheckErr(err, false, "can't set conversion into koanf", "conversion", conversion)

	// Convert the amounts
	if 0 != conversion {
		for i, item := range items {
			newItem := item
			newItem.Unit_cost = utils.GetConvertedCost(item.Unit_cost, conversion)
			items[i] = newItem
		}
	}

	notes := template(k.String("inv.notes"), map[string]string{"seq": seq})

	r.Items = items
	r.Notes = notes

	if k.String("apilayer.currency.new") != "" {
		r.Currency = k.String("apilayer.currency.new")
	} else {
		r.Currency = k.String("apilayer.currency.base")
	}

	return r, filename
}

func applyTemplateToItems(items []invoicegenerator.Item) []invoicegenerator.Item {
	for i, item := range items {
		item.Name = template(item.Name, map[string]string{})
		item.Description = template(item.Description, map[string]string{})
		items[i] = item
	}

	return items
}
