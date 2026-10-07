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

type ExplainAnalyzeTool struct {
	service service.QueryAnalyzerService
}

func NewExplainAnalyzeTool(svc service.QueryAnalyzerService) ToolSpec[input.ExplainInput, output.ExplainOutput] {
	return &ExplainAnalyzeTool{service: svc}
}

func (s *ExplainAnalyzeTool) Name() string {
	return "explain_analyze"
}

func (s *ExplainAnalyzeTool) Description() string {
	return "Runs a SELECT for real (EXPLAIN ANALYZE) and returns the plan with actual row counts and timings. Use explain_query first for expensive queries."
}

func (s *ExplainAnalyzeTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *ExplainAnalyzeTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.ExplainInput) (res *mcp.CallToolResult, out output.ExplainOutput, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s err=%v", s.Name(), err)
	}()
	out, err = s.service.ExplainAnalyze(ctx, in)
	return nil, out, err
}