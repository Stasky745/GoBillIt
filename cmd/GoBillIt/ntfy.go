package main

import (
	"fmt"
	"strings"

	"github.com/Stasky745/GoBillIt/internal/email"
	"github.com/Stasky745/GoBillIt/internal/invoicegenerator"
	"github.com/Stasky745/GoBillIt/internal/ntfy"
	"github.com/Stasky745/go-libs/log"
)

var client ntfy.NtfyClient

var yesAction, noAction, cancelAction, acceptAction, addItemAction ntfy.Action

func createHTTPAction(label, url string, headers map[string]string) ntfy.Action {
	return ntfy.Action{
		Action:  ntfy.ActionHTTP,
		Label:   label,
		URL:     url,
		Headers: headers,
		Body:    label,
		Clear:   false,
	}
}

func getAuthHeaders(auth string) map[string]string {
	headers := map[string]string{}
	if auth != "" {
		headers = map[string]string{
			"Authorization": auth,
		}
	}
	return headers
}

func initializeActions() {
	headers := getAuthHeaders(client.Auth)
	headers["Priority"] = "1"

	yesAction = createHTTPAction(ACTION_LABEL_YES, client.PostURL, headers)
	noAction = createHTTPAction(ACTION_LABEL_NO, client.PostURL, headers)
	cancelAction = createHTTPAction(ACTION_LABEL_CANCEL, client.PostURL, headers)
	addItemAction = createHTTPAction(ACTION_LABEL_ADD_ITEMS, client.PostURL, headers)
}

func ntfyCheck(pdf string) (int, error) {
	var err error
	client, err = ntfy.Initialize(k.String("ntfy.server"), k.String("ntfy.topic"), k.String("ntfy.username"), k.String("ntfy.password"))
	if err != nil {
		return -1, err
	}
	TAGS = k.Strings("ntfy.tags")
	initializeActions()

	response, err := client.SendNotificationAndWaitForResponse(5, "New Invoice!", "Is it correct?", append(EMOJI_NEW_INVOICE, TAGS...), []ntfy.Action{
		yesAction,
		addItemAction,
		cancelAction,
	}, "", pdf)
	log.CheckErr(err, true, "can't send notification")

	switch response {
	case ACTION_LABEL_YES:
		return NTFY_CODE_DONE, nil
	case ACTION_LABEL_CANCEL:
		return NTFY_CODE_CANCEL, nil
	case ACTION_LABEL_ADD_ITEMS:
		extraItems := getExtras()
		err := k.Set(EXTRA_ITEMS, extraItems)
		log.CheckErr(err, false, "couldn't set extra items to koanf", "items", extraItems)
	default:
		log.Panic("received an unwanted response", "question", "New Invoice!", "response", response)
	}

	return NTFY_CODE_RECREATE, nil
}

func loadExtras(file string) []invoicegenerator.Item {
	extraItems := []invoicegenerator.Item{}
	if file != "" {
		err := loadYaml(file)
		if !log.CheckErr(err, false, "can't load ntfy extras file", "path", file) {
			err = k.Unmarshal("ntfy.extras.items", &extraItems)
			log.CheckErr(err, false, "can't unmarshal extra items", "path", file)
		}
	}
	return extraItems
}

func getExtraItem(items []invoicegenerator.Item, item string) invoicegenerator.Item {
	for _, i := range items {
		if i.Label == item {
			return i
		}
	}
	return invoicegenerator.Item{}
}

func getExtras() []invoicegenerator.Item {
	newItems := map[string]invoicegenerator.Item{}
	extraItems := []invoicegenerator.Item{}
	err := k.Unmarshal("ntfy.extras.items", &extraItems)
	if log.CheckErr(err, false, "can't unmarshal list of ntfy extra items", "items", k.Strings("ntfy.extras.items")) {
		extraItems = []invoicegenerator.Item{}
	}
	actions := []ntfy.Action{}

	title := "Add a new item!"
	body := "Choose which item to add. This will keep showing until you press done."

	headers := getAuthHeaders(client.Auth)
	headers["Priority"] = "1"

	for _, item := range extraItems {
		actions = append(actions, createHTTPAction(item.Label, client.PostURL, headers))
	}

	actions = append(actions, createHTTPAction(ACTION_LABEL_DONE, client.PostURL, headers))

	response := ""
	for response != ACTION_LABEL_DONE {
		extendedBody := body

		for _, val := range newItems {
			extendedBody += fmt.Sprintf("\n * %s (%.2f): %d", val.Label, val.Unit_cost, val.Quantity)
		}

		var err error
		response, err = client.SendNotificationAndWaitForResponse(5, title, extendedBody, append(EMOJI_EXTRA_ITEMS, TAGS...), actions, "", "")
		if log.CheckErr(err, false, "can't send new items notification") {
			return []invoicegenerator.Item{}
		} else if response != ACTION_LABEL_DONE {
			newItem := getExtraItem(extraItems, response)
			if val, ok := newItems[newItem.Label]; ok {
				val.Quantity += 1
				newItems[newItem.Label] = val
			} else {
				newItems[newItem.Label] = newItem
			}
		}
	}

	res := make([]invoicegenerator.Item, 0, len(newItems))
	for _, value := range newItems {
		res = append(res, value)
	}
	return res
}

func ntfyEmailCheck(e email.Email) (bool, error) {
	actions := []ntfy.Action{
		yesAction,
		noAction,
	}

	title := "Send Email?"
	body := printBody(e)

	res, err := client.SendNotificationAndWaitForResponse(5, title, body, append(EMOJI_EMAIL_CHECK, TAGS...), actions, "", "")
	if log.CheckErr(err, false, "can't send notification for email check") {
		return false, err
	}
	if res == ACTION_LABEL_YES {
		return true, nil
	}
	return false, nil
}

func printBody(e email.Email) string {
	var rawBodyParts []string
	rawBodyParts = append(rawBodyParts, "From: "+e.From)

	if len(e.To) > 0 {
		rawBodyParts = append(rawBodyParts, "To: "+strings.Join(e.To, ";"))
	}

	if len(e.Cc) > 0 {
		rawBodyParts = append(rawBodyParts, "CC: "+strings.Join(e.Cc, ";"))
	}

	if len(e.Bcc) > 0 {
		rawBodyParts = append(rawBodyParts, "BCC: "+strings.Join(e.Bcc, ";"))
	}

	if e.Subject != "" {
		rawBodyParts = append(rawBodyParts, "Subject: "+e.Subject)
	}

	if e.Body != "" {
		rawBodyParts = append(rawBodyParts, "Body: "+e.Body)
	}

	return strings.Join(rawBodyParts, ",\n")
}
