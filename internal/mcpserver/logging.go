package mcpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// loggingMiddleware logs one structured line per request to logger:
// correlation ID, tool (or method, for non tool-call requests), duration,
// status and response size — never the request or response payload
// itself (reliability-observability-compatibility requirement R3; never
// stdout, which mcp.StdioTransport owns exclusively for the protocol).
func loggingMiddleware(logger *slog.Logger) mcp.Middleware {
	var seq atomic.Uint64
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			correlationID := seq.Add(1)
			start := time.Now()

			result, err := next(ctx, method, req)

			name := method
			if ctr, ok := req.(*mcp.CallToolRequest); ok {
				name = ctr.Params.Name
			}
			status := "ok"
			if err != nil {
				status = "error"
			} else if ctr, ok := result.(*mcp.CallToolResult); ok && ctr.IsError {
				status = "tool_error"
			}
			size := 0
			if result != nil {
				if raw, marshalErr := json.Marshal(result); marshalErr == nil {
					size = len(raw)
				}
			}

			logger.Info("mcp_call",
				"correlation_id", correlationID,
				"tool", name,
				"duration_ms", time.Since(start).Milliseconds(),
				"status", status,
				"size_bytes", size,
			)
			return result, err
		}
	}
}
