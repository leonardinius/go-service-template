package apiworker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/nats-io/nats.go"
)

func NewRequestFromMessage(msg *nats.Msg, subscribePath, handlerPath string) (*http.Request, error) {
	url := subjectToURL(msg.Subject, subscribePath, handlerPath)

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(msg.Data))

	req.Header.Set("Content-Type", "application/json")
	for k, v := range msg.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	return req, nil
}

func urlToSuscribeSubject(url string) string {
	subj := strings.Trim(url, "/")
	subj = strings.ReplaceAll(subj, "/", ".")
	return subj + ".>"
}

func subjectToURL(subj, subscribePath, handlerPath string) string {
	return handlerPath + subj[len(subscribePath)-1:]
}
