package handler

import (
	"context"
	"log"
	"net/url"
)

func HandleGitHubWebhook(_ context.Context, q url.Values, body string) (int, error) {
	log.Printf("q=%+v", q)
	log.Printf("body=%+v", body)
	return 200, nil
}
