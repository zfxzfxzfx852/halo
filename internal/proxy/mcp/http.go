package mcp

import (
	"encoding/json"
	"errors"
	"time"

	"halo/internal/model"
)

// EventsFromRequest 解析 MCP streamable HTTP 请求体（单个或批量 JSON-RPC），
// 对 tools/call 生成事件；非 MCP 内容返回 nil。
func EventsFromRequest(body []byte, sessionID string, now time.Time) []*model.Event {
	reqs, err := decodeRequests(body)
	if err != nil {
		return nil
	}
	var events []*model.Event
	for _, req := range reqs {
		if req.Method != MethodToolsCall {
			continue
		}
		var params ToolCallParams
		_ = json.Unmarshal(req.Params, &params)
		events = append(events, &model.Event{
			ID:        model.NewID("evt"),
			Type:      model.EventToolCall,
			Protocol:  model.ProtocolMCPHTTP,
			SessionID: sessionID,
			Timestamp: now,
			Method:    req.Method,
			ToolCall: &model.ToolCall{
				ID:        string(req.ID),
				Name:      params.Name,
				Arguments: params.Arguments,
			},
		})
	}
	return events
}

func decodeRequests(body []byte) ([]Request, error) {
	var single Request
	if err := json.Unmarshal(body, &single); err == nil && single.JSONRPC == JSONRPCVersion {
		return []Request{single}, nil
	}
	var batch []Request
	if err := json.Unmarshal(body, &batch); err == nil && len(batch) > 0 {
		return batch, nil
	}
	return nil, errors.New("not a json-rpc payload")
}
