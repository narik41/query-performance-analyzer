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

type GetSchemaSpecTool struct {
	service service.QueryAnalyzerService
}

func NewGetSchemaSpecTool(svc service.QueryAnalyzerService) ToolSpec[input.GetSchemaInput, output.GetSchemaOutput] {
	return &GetSchemaSpecTool{service: svc}
}

func (s *GetSchemaSpecTool) Name() string {
	return "get_schema"
}

func (s *GetSchemaSpecTool) Description() string {
	return "Get the schema spec"
}

func (s *GetSchemaSpecTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *GetSchemaSpecTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.GetSchemaInput) (res *mcp.CallToolResult, out output.GetSchemaOutput, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s err=%v", s.Name(), err)
	}()
	out, err = s.service.GetSchema(ctx, in)
	return nil, out, err
}
