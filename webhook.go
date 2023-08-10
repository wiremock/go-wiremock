package wiremock

import (
	"encoding/json"
	"time"
)

type WebhookInterface interface {
	json.Marshaler
	WithName(name string) WebhookInterface
	ParseWebhook() map[string]interface{}
}

type Webhook struct {
	name       string
	parameters webhookParameters
}

type webhookParameters struct {
	method  string
	url     string
	body    string
	headers map[string]string
	delay   DelayInterface
}

func (w webhookParameters) MarshalJSON() ([]byte, error) {
	jsonMap := map[string]interface{}{
		"method":  w.method,
		"url":     w.url,
		"body":    w.body,
		"headers": w.headers,
		"delay":   w.delay.ParseDelay(),
	}

	return json.Marshal(jsonMap)
}

// WithName sets the name of the webhook and returns the webhook.
func (w Webhook) WithName(name string) WebhookInterface {
	w.name = name
	return w
}

// ParseWebhook returns a map representation of the webhook.
func (w Webhook) ParseWebhook() map[string]interface{} {
	return map[string]interface{}{
		"name":       w.name,
		"parameters": w.parameters,
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (w Webhook) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.ParseWebhook())
}

// WithMethod sets the HTTP method of the webhook.
func (w Webhook) WithMethod(method string) Webhook {
	w.parameters.method = method
	return w
}

// WithURL sets the URL of the webhook.
func (w Webhook) WithURL(url string) Webhook {
	w.parameters.url = url
	return w
}

// WithHeader sets a header of the webhook.
func (w Webhook) WithHeader(key string, value string) Webhook {
	if w.parameters.headers == nil {
		w.parameters.headers = make(map[string]string)
	}

	w.parameters.headers[key] = value

	return w
}

// WithBody sets the body of the webhook.
func (w Webhook) WithBody(body string) Webhook {
	w.parameters.body = body
	return w
}

// WithDelay sets the delay of the webhook.
func (w Webhook) WithDelay(delay DelayInterface) Webhook {
	w.parameters.delay = delay
	return w
}

// WithFixedDelay sets the fixed delay of the webhook.
func (w Webhook) WithFixedDelay(delay time.Duration) Webhook {
	w.parameters.delay = NewFixedDelay(delay)
	return w
}

// WithLogNormalRandomDelay sets the log normal delay of the webhook.
func (w Webhook) WithLogNormalRandomDelay(median time.Duration, sigma float64) Webhook {
	w.parameters.delay = NewLogNormalRandomDelay(median, sigma)
	return w
}

// WithUniformRandomDelay sets the uniform random delay of the webhook.
func (w Webhook) WithUniformRandomDelay(lower, upper time.Duration) Webhook {
	w.parameters.delay = NewUniformRandomDelay(lower, upper)
	return w
}

func NewWebhook() Webhook {
	return Webhook{}
}
