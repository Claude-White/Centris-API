package server

import (
	"net/http"
	"strings"
)

func SendNotification(title string, body string) {
	req, _ := http.NewRequest("POST", "https://ntfy.omniquark.me/centris-api",
		strings.NewReader(body))
	req.Header.Set("Title", title)
	http.DefaultClient.Do(req)
}
