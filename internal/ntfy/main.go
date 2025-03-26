package ntfy

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Stasky745/go-libs/log"
)

type NtfyClient struct {
	Server  string
	Topic   string
	PostURL string
	GetURL  string
	Auth    string
}

func appendURL(baseURL, path string) (string, error) {
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	// Resolve the reference by appending the path
	// ResolveReference ensures no extra '/' is added
	fullURL := parsedURL.JoinPath(path)

	return fullURL.String(), nil
}

func Initialize(server, topic, username, password string) (NtfyClient, error) {
	postURL, err := appendURL(server, topic)
	if log.CheckErr(err, false, "can't create POST URL", "server", server, "topic", topic) {
		return NtfyClient{}, err
	}

	getURL, err := appendURL(postURL, "raw")
	if log.CheckErr(err, false, "can't create GET URL", "server", server, "topic", topic) {
		return NtfyClient{}, err
	}

	client := NtfyClient{
		Server:  server,
		Topic:   topic,
		PostURL: postURL,
		GetURL:  getURL,
		Auth:    "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
	}

	return client, nil
}

// Define a custom type for action types
type ActionType string

// Define valid action types as constants
const (
	ActionBroadcast ActionType = "broadcast"
	ActionView      ActionType = "view"
	ActionHTTP      ActionType = "http"
)

// Action struct defines an ntfy action
type Action struct {
	Action  ActionType        `json:"action"`            // Only allowed values
	Label   string            `json:"label"`             // Button text
	URL     string            `json:"url,omitempty"`     // For "view"
	Method  string            `json:"method,omitempty"`  // For "http"
	Headers map[string]string `json:"headers,omitempty"` // For "http"
	Body    string            `json:"body,omitempty"`    // For "http"
	Clear   bool              `json:"clear,omitempty"`
}

// Notification struct represents the notification to be sent
type Notification struct {
	Topic    string   `json:"topic"`
	Title    string   `json:"title,omitempty"`
	Message  string   `json:"message,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Priority int      `json:"priority"`
	Attach   string   `json:"attach,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Actions  []Action `json:"actions,omitempty"`
}

func CreateAction(actionType ActionType, label, url string) (Action, error) {
	// Validate action type (optional, but useful for runtime checks)
	switch actionType {
	case ActionBroadcast, ActionView, ActionHTTP:
		// Valid action
	default:
		err := fmt.Errorf("invalid action type: %s", actionType)
		log.Error("invalid action type.", "action_type", actionType)
		return Action{}, err
	}

	action := Action{
		Action: actionType,
		Label:  label,
		URL:    url,
	}

	return action, nil
}

// Convert a list of Actions to formatted string
func convertActionsToString(actions []Action) string {
	var result []string

	for _, action := range actions {
		parts := []string{string(action.Action), action.Label, action.URL}

		// Add method if applicable
		if action.Method != "" {
			parts = append(parts, fmt.Sprintf("method=%s", action.Method))
		}

		// Add headers if applicable
		for key, value := range action.Headers {
			parts = append(parts, fmt.Sprintf("headers.%s=%s", key, value))
		}

		// Add body if applicable
		if action.Body != "" {
			parts = append(parts, fmt.Sprintf("body=%s", action.Body))
		}

		if action.Clear {
			parts = append(parts, "clear=true")
		} else {
			parts = append(parts, "clear=false")
		}

		result = append(result, strings.Join(parts, ", "))
	}

	return strings.Join(result, "; ")
}

func (ntfy NtfyClient) SendNotification(priority int, title, message string, tags []string, actions []Action, attach, file string) error {
	var req *http.Request
	if file != "" {
		f, _ := os.Open(file)
		req, _ = http.NewRequest("POST", ntfy.PostURL, f)

		req.Header.Set("Message", message)
		req.Header.Set("Filename", filepath.Base(file))
	} else {
		req, _ = http.NewRequest("POST", ntfy.PostURL, strings.NewReader(message))
	}
	req.Header.Set("Title", title)
	req.Header.Set("Tags", strings.Join(tags, ","))
	req.Header.Set("Attach", attach)
	req.Header.Set("Priority", strconv.Itoa(priority))
	req.Header.Set("Actions", convertActionsToString(actions))
	req.Header.Set("Markdown", "yes")

	if ntfy.Auth != "" {
		req.Header.Set("Authorization", ntfy.Auth)
	}

	resp, err := http.DefaultClient.Do(req)
	log.CheckErr(err, false, "can't do request", "request", req, "response", resp)
	defer log.CheckErr(resp.Body.Close(), false, "can't close response body", "response", resp)

	return err
}

func (ntfy NtfyClient) listenForResponses() (string, error) {
	// Open HTTP connection to ntfy topic
	req, _ := http.NewRequest("GET", ntfy.GetURL, nil)
	if ntfy.Auth != "" {
		req.Header.Set("Authorization", ntfy.Auth)
	}

	resp, err := http.DefaultClient.Do(req)
	if log.CheckErr(err, false, "can't listen for response", "URL", ntfy.GetURL) {
		log.CheckErr(resp.Body.Close(), false, "can't close response body", "response", resp)
		return "", err
	}
	defer log.CheckErr(resp.Body.Close(), false, "can't close response body", "response", resp)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		message := scanner.Text()
		if message != "" {
			return scanner.Text(), nil
		}
	}
	return "", nil
}

func (ntfy NtfyClient) SendNotificationAndWaitForResponse(priority int, title, message string, tags []string, actions []Action, attach, filename string) (string, error) {
	if len(actions) == 0 {
		log.Errorf("can't send notification and wait for response without actions")
		return "", fmt.Errorf("can't send notification and wait for response without actions")
	}

	log.CheckErr(ntfy.SendNotification(priority, title, message, tags, actions, attach, filename), false, "can't send notification")
	return ntfy.listenForResponses()
}
