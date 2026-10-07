package mcp_tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ToolSpec[In, Out any] interface {
	Name() string

	Description() string

	Spec() *mcp.Tool

	Handle(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, Out, error)
}
