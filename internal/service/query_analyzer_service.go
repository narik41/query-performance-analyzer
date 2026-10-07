package service

import (
	"context"
	"database-query-analyzer/internal/logger"
	"database-query-analyzer/internal/repository"
	"fmt"

	"database-query-analyzer/internal/input"
	"database-query-analyzer/internal/output"
)

type QueryAnalyzerService struct {
	repository repository.DatabaseAdapter
}

func NewQueryAnalyzerService(repository repository.DatabaseAdapter) QueryAnalyzerService {
	return QueryAnalyzerService{repository: repository}
}

func (s *QueryAnalyzerService) GetDatabaseInfo(ctx context.Context, in input.GetDatabaseInfoInput) (output.GetDatabaseInfoOutput, error) {
	info, err := s.repository.GetDatabaseInfo(ctx)
	if err != nil {
		return output.GetDatabaseInfoOutput{}, err
	}
	out := output.GetDatabaseInfoOutput{
		Database: info.Name,
		Version:  info.Version,
	}
	if in.IncludeTables {
		tables, err := s.repository.GetTables(ctx)
		if err != nil {
			return out, err
		}
		for _, t := range tables {
			out.Tables = append(out.Tables, t.Name)
		}
	}
	return out, nil
}

func (s *QueryAnalyzerService) GetSchema(ctx context.Context, in input.GetSchemaInput) (output.GetSchemaOutput, error) {
	if in.Table != "" {
		schema, err := s.repository.GetTableSchema(ctx, in.Table)
		if err != nil {
			return output.GetSchemaOutput{}, err
		}
		return output.GetSchemaOutput{Tables: []output.TableSchema{*schema}}, nil
	}

	tables, err := s.repository.GetTables(ctx)
	if err != nil {
		return output.GetSchemaOutput{}, err
	}

	var schemas []output.TableSchema
	for _, t := range tables {
		schema, err := s.repository.GetTableSchema(ctx, t.Name)
		if err != nil {
			return output.GetSchemaOutput{}, err
		}
		schemas = append(schemas, *schema)
	}
	return output.GetSchemaOutput{Tables: schemas}, nil
}

func (s *QueryAnalyzerService) GetIndexes(ctx context.Context, in input.GetIndexesInput) (output.GetIndexesOutput, error) {
	if in.Table != "" {
		indexes, err := s.repository.GetIndexes(ctx, in.Table)
		if err != nil {
			return output.GetIndexesOutput{}, err
		}
		return output.GetIndexesOutput{Indexes: indexes}, nil
	}

	tables, err := s.repository.GetTables(ctx)
	if err != nil {
		return output.GetIndexesOutput{}, err
	}

	var allIndexes []output.Index
	for _, t := range tables {
		indexes, err := s.repository.GetIndexes(ctx, t.Name)
		if err != nil {
			return output.GetIndexesOutput{}, err
		}
		allIndexes = append(allIndexes, indexes...)
	}
	return output.GetIndexesOutput{Indexes: allIndexes}, nil
}

func (s *QueryAnalyzerService) ExplainQuery(ctx context.Context, in input.ExplainInput) (output.ExplainOutput, error) {
	logger.Default().InfoWithCtxf(ctx, "Running explain cmd. %v", in.SQL)
	plan, err := s.repository.Explain(ctx, in.SQL)
	if err != nil {
		return output.ExplainOutput{}, err
	}
	return output.ExplainOutput{Plan: plan.RawPlan}, nil
}

func (s *QueryAnalyzerService) ExplainAnalyze(ctx context.Context, in input.ExplainInput) (output.ExplainOutput, error) {
	logger.Default().InfoWithCtxf(ctx, "Running explain analyze cmd. %v", in.SQL)
	analysis, err := s.repository.ExplainAnalyze(ctx, in.SQL)
	if err != nil {
		return output.ExplainOutput{}, err
	}
	return output.ExplainOutput{Plan: analysis.RawAnalysis}, nil
}

func (s *QueryAnalyzerService) GetQueryStats(ctx context.Context, in input.QueryStatsInput) (output.QueryStatsOutput, error) {
	logger.Default().InfoWithCtxf(ctx, "Get query stats. %v", in.SQL)
	plan, err := s.repository.Explain(ctx, in.SQL)
	if err != nil {
		return output.QueryStatsOutput{}, fmt.Errorf("explain: %w", err)
	}

	analysis, err := s.repository.ExplainAnalyze(ctx, in.SQL)
	explainAnalyze := ""
	if err != nil {
		explainAnalyze = "EXPLAIN ANALYZE not supported: " + err.Error()
	} else {
		explainAnalyze = analysis.RawAnalysis
	}

	rowCount, executionTimeMs, err := s.repository.ExecuteQuery(ctx, in.SQL)
	if err != nil {
		return output.QueryStatsOutput{}, fmt.Errorf("query execution: %w", err)
	}

	warnings, err := s.repository.GetWarnings(ctx)
	if err != nil {
		warnings = []string{}
	}

	return output.QueryStatsOutput{
		SQL:             in.SQL,
		ExecutionTimeMs: executionTimeMs,
		RowsReturned:    rowCount,
		ExplainJSON:     plan.RawPlan,
		ExplainAnalyze:  explainAnalyze,
		WarningCount:    len(warnings),
		Warnings:        warnings,
	}, nil
}
