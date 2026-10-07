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

type GetIndexesTool struct {
	service service.QueryAnalyzerService
}

func NewGetIndexesTool(svc service.QueryAnalyzerService) ToolSpec[input.GetIndexesInput, output.GetIndexesOutput] {
	return &GetIndexesTool{service: svc}
}

func (s *GetIndexesTool) Name() string {
	return "get_indexes"
}

func (s *GetIndexesTool) Description() string {
	return "List indexes (name, uniqueness, type, ordered columns, cardinality) for one table or all tables."
}

func (s *GetIndexesTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *GetIndexesTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.GetIndexesInput) (res *mcp.CallToolResult, out output.GetIndexesOutput, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s err=%v", s.Name(), err)
	}()
	out, err = s.service.GetIndexes(ctx, in)
	return nil, out, err
}