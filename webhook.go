package wiremock

type WebhookDefinition struct {
	m map[string]any
}

func Webhook() WebhookDefinition {
	return WebhookDefinition{
		m: make(map[string]any),
	}
}

func (d WebhookDefinition) WithMethod(method string) WebhookDefinition {
	d.m["method"] = method
	return d
}

func (d WebhookDefinition) WithURL(url string) WebhookDefinition {
	d.m["url"] = url
	return d
}

func (d WebhookDefinition) WithHeader(key string, value string) WebhookDefinition {
	var headers map[string]string

	if headersAny, ok := d.m["headers"]; ok {
		headers = headersAny.(map[string]string)
	} else {
		headers = make(map[string]string)
		d.m["headers"] = headers
	}

	headers[key] = value

	return d
}

func (d WebhookDefinition) WithBody(body string) WebhookDefinition {
	d.m["body"] = body
	return d
}

func (d WebhookDefinition) ToMap() map[string]any {
	return d.m
}
