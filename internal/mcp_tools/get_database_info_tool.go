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

type GetDatabaseInfoTool struct {
	service service.QueryAnalyzerService
}

func NewGetDatabaseInfoTool(svc service.QueryAnalyzerService) ToolSpec[input.GetDatabaseInfoInput, output.GetDatabaseInfoOutput] {
	return &GetDatabaseInfoTool{service: svc}
}

func (s *GetDatabaseInfoTool) Name() string {
	return "get_database_info"
}
func (s *GetDatabaseInfoTool) Description() string {
	return "Get the database engine, version and list of schemas. Call this first to learn what the database contains."
}

func (s *GetDatabaseInfoTool) Spec() *mcp.Tool {
	return &mcp.Tool{
		Name:        s.Name(),
		Description: s.Description(),
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}
}

func (s *GetDatabaseInfoTool) Handle(ctx context.Context, req *mcp.CallToolRequest, in input.GetDatabaseInfoInput) (res *mcp.CallToolResult, out output.GetDatabaseInfoOutput, err error) {
 
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s panicked: %v", s.Name(), r)
		}
		log.Printf("tool=%s  err=%v", s.Name(), err)
	}()
	out, err = s.service.GetDatabaseInfo(ctx, in)
	return nil, out, err
}
