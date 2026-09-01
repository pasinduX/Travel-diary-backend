package integrations

import "net/http"

func NewHTTPClient() *http.Client {
	return &http.Client{}
}
