package samplermiddlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"github.com/svix/svix-webhooks/go"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
)

type ReplicateMiddlewareWebhook struct {}

func NewReplicateMiddlewareWebhook() *ReplicateMiddlewareWebhook {
	return &ReplicateMiddlewareWebhook{}
}

func (r *ReplicateMiddlewareWebhook) VeriFyWebhook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Println("WEBHOOK RECEIVED")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Cannot read body", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		wh, err := svix.NewWebhook(config.Envs.ReplicateWebhookSecret)
		if err != nil {
			log.Printf("Webhook init error: %v", err)
			http.Error(w, "Webhook config error", http.StatusInternalServerError)
			return
		}
		headers := make(http.Header)
		headers.Set("webhook-id", r.Header.Get("webhook-id"))
		headers.Set("webhook-timestamp", r.Header.Get("webhook-timestamp"))
		headers.Set("webhook-signature", r.Header.Get("webhook-signature"))

		err = wh.Verify(body, headers)
		if err != nil {
			log.Printf("Webhook verification failed: %v", err)
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
		log.Println("Webhook verification successful")
		next.ServeHTTP(w, r)
	})
}