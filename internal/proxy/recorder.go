package proxy

import (
	"net/http"
	"strings"
	"time"

	"halo/internal/model"
)

var sensitiveHeaders = map[string]bool{
	"Authorization":         true,
	"Proxy-Authorization":   true,
	"Cookie":                true,
	"Set-Cookie":            true,
}

func sanitizeHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if sensitiveHeaders[http.CanonicalHeaderKey(k)] {
			out[k] = "[redacted]"
			continue
		}
		out[k] = strings.Join(v, ", ")
	}
	return out
}

func buildRequestEvent(id string, r *http.Request, body []byte, sessionID string, ts time.Time) *model.Event {
	return &model.Event{
		ID:        id,
		Type:      model.EventRequest,
		Protocol:  model.ProtocolHTTP,
		Agent:     r.UserAgent(),
		SessionID: sessionID,
		Timestamp: ts,
		Method:    r.Method + " " + r.URL.Path,
		Request: &model.Payload{
			Headers: sanitizeHeaders(r.Header),
			Body:    body,
		},
		Risk: model.RiskNone,
	}
}

func buildResponseEvent(parentID string, h http.Header, status int, body []byte, sessionID string, ts time.Time) *model.Event {
	return &model.Event{
		ID:        model.NewID("evt"),
		Type:      model.EventResponse,
		Protocol:  model.ProtocolHTTP,
		SessionID: sessionID,
		Timestamp: ts,
		Response: &model.Payload{
			Headers: sanitizeHeaders(h),
			Body:    body,
			Status:  status,
		},
		ParentID: parentID,
		Risk:     model.RiskNone,
	}
}
