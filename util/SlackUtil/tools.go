package SlackUtil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Verify(w http.ResponseWriter, r *http.Request, body []byte, SLACK_VERIFY_TOKEN string) {
	var jsonbody map[string]interface{}
	err := json.Unmarshal(body, &jsonbody)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusBadRequest)
		return
	}

	token, ok := jsonbody["token"].(string)
	if !ok {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	eventType, ok := jsonbody["type"].(string)
	if !ok {
		http.Error(w, "Missing event type", http.StatusBadRequest)
		return
	}

	if token != SLACK_VERIFY_TOKEN {
		http.Error(w, "Invalid verification token", http.StatusUnauthorized)
		return
	}
	
	if eventType == "url_verification" {
		challenge, ok := jsonbody["challenge"].(string)
		if !ok {
			http.Error(w, "Missing challenge", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, challenge)
	}
	w.WriteHeader((http.StatusOK))
}