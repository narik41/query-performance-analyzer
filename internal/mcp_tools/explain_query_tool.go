package mcp_tools

import (
	"context"
	"database-query-analyzer/internal/input"
	"database-query-analyzer/internal/output"
	"database-query-analyzer/internal/service"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ExplainQueryTool struct {
	service service.QueryAnalyzerService
}

func NewExplainQueryTool(svc service.QueryAnalyzerService) ToolSpec[input.ExplainInput, output.ExplainOutput] {
	return &ExplainQueryTool{service: svc}
}

func (s *ExplainQueryTool) Name() string {
	return "explain_query"
}

func (s *ExplainQueryTool) Description() string {
	return "Shows the optimizer's planned execution (EXPLAIN FORMAT=JSON) for a SELECT. Does NOT run the query."
}

func (s *ExplainQueryTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *ExplainQueryTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.ExplainInput) (res *mcp.CallToolResult, out output.ExplainOutput, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s err=%v", s.Name(), err)
	}()
	out, err = s.service.ExplainQuery(ctx, in)
	return nil, out, err
}