package main

import (
	"context"
	"database-query-analyzer/api/controller"
	"database-query-analyzer/internal/config"
	"database-query-analyzer/internal/db"
	"database-query-analyzer/internal/logger"
	"database-query-analyzer/internal/mcp_tools"
	"database-query-analyzer/internal/repository"
	"database-query-analyzer/internal/service"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {

	ctx := context.Background()
	mysqlConnection := connectionMYSQL(ctx)
	defer func(mysqlConnection *db.MySQLConnectionManager) {
		err := mysqlConnection.Close()
		if err != nil {
			logger.Default().Errorf("failed to close mysql connection manager: %v", err)
		}
	}(mysqlConnection)
	mysqlRepository := repository.NewMySQLAdapter(mysqlConnection.DB())

	queryAnalyzerService := service.NewQueryAnalyzerService(mysqlRepository)

	// Create all MCP tools
	getDatabaseInfoTool := mcp_tools.NewGetDatabaseInfoTool(queryAnalyzerService)
	getSchemaInfoTool := mcp_tools.NewGetSchemaSpecTool(queryAnalyzerService)
	getIndexesTool := mcp_tools.NewGetIndexesTool(queryAnalyzerService)
	explainQueryTool := mcp_tools.NewExplainQueryTool(queryAnalyzerService)
	explainAnalyzeTool := mcp_tools.NewExplainAnalyzeTool(queryAnalyzerService)
	getQueryStatsTool := mcp_tools.NewGetQueryStatsTool(queryAnalyzerService)

	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "database-query-analyzer",
			Version: "v1.0.0",
		}, nil,
	)

	// Register all tools
	mcp.AddTool(mcpServer, getDatabaseInfoTool.Spec(), getDatabaseInfoTool.Handle)
	mcp.AddTool(mcpServer, getSchemaInfoTool.Spec(), getSchemaInfoTool.Handle)
	mcp.AddTool(mcpServer, getIndexesTool.Spec(), getIndexesTool.Handle)
	mcp.AddTool(mcpServer, explainQueryTool.Spec(), explainQueryTool.Handle)
	mcp.AddTool(mcpServer, explainAnalyzeTool.Spec(), explainAnalyzeTool.Handle)
	mcp.AddTool(mcpServer, getQueryStatsTool.Spec(), getQueryStatsTool.Handle)

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server {
			return mcpServer
		},
		&mcp.StreamableHTTPOptions{
			Stateless:    true,
			JSONResponse: true,
		},
	)

	mux := http.NewServeMux()
	healthHandler := controller.NewHcController()
	mux.HandleFunc("/healthz", healthHandler.Healthz)
	mux.HandleFunc("/readyz", healthHandler.Readyz)
	mux.Handle("/mcp", controller.LoggingMiddleware(handler))

	port := config.Default().GetInt("server.port")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	logger.Default().Infof("starting http server on port %d", port)
	if err := httpServer.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

func connectionMYSQL(ctx context.Context) *db.MySQLConnectionManager {
	mysqlConnection := db.NewMySQLConnectionManager()
	err := mysqlConnection.Connect(ctx, db.ConnectionConfig{
		User:     config.Default().GetString("db.mysql.user"),
		Password: config.Default().GetString("db.mysql.password"),
		Host:     config.Default().GetString("db.mysql.host"),
		Port:     config.Default().GetInt("db.mysql.port"),
		Database: config.Default().GetString("db.mysql.database"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return mysqlConnection
}
