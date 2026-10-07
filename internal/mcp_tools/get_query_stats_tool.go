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

type GetQueryStatsTool struct {
	service service.QueryAnalyzerService
}

func NewGetQueryStatsTool(svc service.QueryAnalyzerService) ToolSpec[input.QueryStatsInput, output.QueryStatsOutput] {
	return &GetQueryStatsTool{service: svc}
}

func (s *GetQueryStatsTool) Name() string {
	return "get_query_stats"
}

func (s *GetQueryStatsTool) Description() string {
	return "Executes a SELECT and returns execution time, row count, EXPLAIN FORMAT=JSON, EXPLAIN ANALYZE, and warnings in one call."
}

func (s *GetQueryStatsTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *GetQueryStatsTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.QueryStatsInput) (res *mcp.CallToolResult, out output.QueryStatsOutput, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s err=%v", s.Name(), err)
	}()
	out, err = s.service.GetQueryStats(ctx, in)
	return nil, out, err
}