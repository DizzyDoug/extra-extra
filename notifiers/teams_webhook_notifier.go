package notifiers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DizzyDoug/extra-extra/config"
)

// TeamsNotifier sends notifications via webhook to MS Teams.
type TeamsNotifier struct {
	webhookURL string
}

// NewTeamsNotifier is the constructor for a TeamsNotifier
func NewTeamsNotifier(config config.Config) TeamsNotifier {
	return TeamsNotifier{
		webhookURL: config.Teams.WebhookURL,
	}
}

// Publish sends a notification
func (n TeamsNotifier) Publish(payload []byte) {
	resp, err := http.Post(n.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Error while publishing to MS Teams: " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bz, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Error while publishing to MS Teams: %d_%s: %s\n", resp.StatusCode, resp.Status, string(bz))
	}
}
