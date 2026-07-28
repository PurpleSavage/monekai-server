package paymentsmiddlewares

import (
	"bytes"
	"context"
	"io"
	"net/http"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
)
type contextKey string

const WebhookBodyKey contextKey = "paddle_webhook_body"


type PaymentWebhookVerifier struct {
	verifier *paddle.WebhookVerifier
}

func NewPaymentWebhookVerifier()*PaymentWebhookVerifier{
	return &PaymentWebhookVerifier{
		verifier: paddle.NewWebhookVerifier(
			config.Envs.PaddleWebhookSecret,
		),
	}
}

func (p *PaymentWebhookVerifier) Verify(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ok, err := p.verifier.Verify(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !ok {
			http.Error(w, "invalid paddle signature", http.StatusForbidden)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read body", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		ctx := context.WithValue(
			r.Context(),
			WebhookBodyKey,
			body,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}