package main

const (
	APP_PREFIX      = "GBI_"
	CONFIG_FILE     = "/etc/gbi/config.yaml"
	CONFIG_DIR      = "/etc/gbi/conf.d"
	EXTRA_ITEMS     = "ntfy.extras.items"
	TEMPLATE_PREFIX = "var."

	// Action Labels
	ACTION_LABEL_YES       = "Yes"
	ACTION_LABEL_CANCEL    = "Cancel"
	ACTION_LABEL_ADD_ITEMS = "Add Items"
	ACTION_LABEL_NO        = "No"
	ACTION_LABEL_RESTART   = "Restart"
	ACTION_LABEL_DONE      = "Done"

	// ntfy Codes
	NTFY_CODE_DONE     = 0
	NTFY_CODE_RECREATE = 1
	NTFY_CODE_EMAIL    = 2
	NTFY_CODE_CANCEL   = 3

	// regex
	REGEX_TEMPLATING = `\{\{\s*([\w\.]+)(?:\s*\[\s*(\w+)\s*\])?\s*\}\}`
)

var (
	LIST_ENVS = []string{
		"GBI_EMAIL_TO",
		"GBI_EMAIL_CC",
		"GBI_EMAIL_BCC",
		"GBI_NTFY_TAGS",
		"GBI_NTFY_EXTRAS_ITEMS",
		"GBI_INV_ITEMS_LIST",
	}
	REQUIRED_FIELDS = []string{
		"inv.apikey", "inv.to", "inv.from",
	}
	TAGS                   = []string{}
	EMOJI_ERROR            = []string{"x"}
	EMOJI_NEW_INVOICE      = []string{"receipt"}
	EMOJI_EXTRA_ITEMS      = []string{"money_mouth_face"}
	EMOJI_EMAIL_CHECK      = []string{"email"}
	EMOJI_TADA             = []string{"tada"}
	EMOJI_ROCKET           = []string{"rocket"}
	EMOJI_MONEY_WITH_WINGS = []string{"money_with_wings"}
	EMOJI_SUNGLASSES       = []string{"sunglasses"}
)
