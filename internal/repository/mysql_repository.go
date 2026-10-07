package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"database-query-analyzer/internal/output"
)

type DatabaseAdapter interface {
	GetDatabaseInfo(ctx context.Context) (*output.DatabaseInfo, error)

	GetTables(ctx context.Context) ([]output.Table, error)

	GetTableSchema(ctx context.Context, table string) (*output.TableSchema, error)

	GetIndexes(ctx context.Context, table string) ([]output.Index, error)

	Explain(ctx context.Context, query string) (*output.QueryPlan, error)

	ExplainAnalyze(ctx context.Context, query string) (*output.QueryAnalysis, error)

	ExecuteQuery(ctx context.Context, query string) (int64, int64, error)

	GetWarnings(ctx context.Context) ([]string, error)
}

type MySQLAdapter struct {
	db *sql.DB
}

func NewMySQLAdapter(db *sql.DB) *MySQLAdapter {
	return &MySQLAdapter{
		db: db,
	}
}

func (a MySQLAdapter) GetDatabaseInfo(ctx context.Context) (*output.DatabaseInfo, error) {

	var name, version string
	if err := a.db.QueryRowContext(ctx, "SELECT DATABASE(), VERSION()").Scan(&name, &version); err != nil {
		return nil, fmt.Errorf("failed to get database info: %w", err)
	}

	return &output.DatabaseInfo{
		Name:    name,
		Version: version,
	}, nil
}

func (a MySQLAdapter) GetTables(ctx context.Context) ([]output.Table, error) {

	query := `
		SELECT
			table_name,
			table_schema,
			table_comment,
			table_rows,
			data_length + index_length as total_size,
			engine,
			table_collation,
			create_time,
			update_time
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		ORDER BY table_name
	`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []output.Table
	for rows.Next() {
		var t output.Table
		var collation string
		var rowCount, sizeBytes sql.NullInt64
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&t.Name, &t.Schema, &t.Comment, &rowCount, &sizeBytes, &t.Engine, &collation, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}
		t.RowCount = rowCount.Int64
		t.SizeBytes = sizeBytes.Int64
		t.Collation = collation
		if createdAt.Valid {
			t.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.Time
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (a MySQLAdapter) GetTableSchema(ctx context.Context, table string) (*output.TableSchema, error) {

	colQuery := `
		SELECT
			column_name,
			column_type,
			is_nullable,
			column_key,
			column_default,
			extra,
			column_comment,
			ordinal_position,
			character_set_name,
			collation_name
		FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY ordinal_position
	`

	colRows, err := a.db.QueryContext(ctx, colQuery, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer colRows.Close()

	var columns []output.Column
	for colRows.Next() {
		var c output.Column
		var nullable, key, extra, comment, charset, collation string
		var defaultVal sql.NullString
		if err := colRows.Scan(&c.Name, &c.Type, &nullable, &key, &defaultVal, &extra, &comment, &c.OrdinalPosition, &charset, &collation); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		c.Nullable = nullable == "YES"
		c.Key = key
		if defaultVal.Valid {
			c.Default = &defaultVal.String
		}
		c.Extra = extra
		c.Comment = comment
		c.CharacterSet = charset
		c.Collation = collation
		columns = append(columns, c)
	}
	if err := colRows.Err(); err != nil {
		return nil, fmt.Errorf("column iteration error: %w", err)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("table %q not found", table)
	}

	fkQuery := `
		SELECT
			rc.constraint_name,
			kcu.column_name,
			kcu.referenced_table_name,
			kcu.referenced_column_name,
			rc.update_rule,
			rc.delete_rule
		FROM information_schema.referential_constraints rc
		JOIN information_schema.key_column_usage kcu
			ON rc.constraint_name = kcu.constraint_name
			AND rc.table_schema = kcu.table_schema
		WHERE rc.table_schema = DATABASE()
		  AND rc.table_name = ?
		ORDER BY rc.constraint_name, kcu.ordinal_position
	`

	fkRows, err := a.db.QueryContext(ctx, fkQuery, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query foreign keys: %w", err)
	}
	defer fkRows.Close()

	var foreignKeys []output.ForeignKey
	for fkRows.Next() {
		var fk output.ForeignKey
		if err := fkRows.Scan(&fk.Name, &fk.Column, &fk.ReferencedTable, &fk.ReferencedColumn, &fk.OnUpdate, &fk.OnDelete); err != nil {
			return nil, fmt.Errorf("failed to scan foreign key: %w", err)
		}
		foreignKeys = append(foreignKeys, fk)
	}

	indexes, err := a.GetIndexes(ctx, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexes: %w", err)
	}

	return &output.TableSchema{
		TableName:   table,
		Columns:     columns,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
	}, nil
}

func (a MySQLAdapter) GetIndexes(ctx context.Context, table string) ([]output.Index, error) {

	query := `
		SELECT
			index_name,
			non_unique,
			index_type,
			column_name,
			cardinality,
			index_comment,
			is_visible
		FROM information_schema.statistics
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY index_name, seq_in_index
	`

	rows, err := a.db.QueryContext(ctx, query, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer rows.Close()

	type idxKey struct{ name string }
	indexMap := make(map[string]*output.Index)

	for rows.Next() {
		var name string
		var nonUnique int
		var idxType, column, comment, visible string
		var cardinality sql.NullInt64

		if err := rows.Scan(&name, &nonUnique, &idxType, &column, &cardinality, &comment, &visible); err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}

		idx, exists := indexMap[name]
		if !exists {
			idx = &output.Index{
				Name:    name,
				Table:   table,
				Unique:  nonUnique == 0,
				Type:    idxType,
				Columns: []string{},
				Comment: comment,
				Visible: visible == "YES",
			}
			indexMap[name] = idx
		}
		idx.Columns = append(idx.Columns, column)
		if cardinality.Valid {
			idx.Cardinality = cardinality.Int64
		}
	}

	var indexes []output.Index
	for _, idx := range indexMap {
		indexes = append(indexes, *idx)
	}

	return indexes, rows.Err()
}

func (a MySQLAdapter) Explain(ctx context.Context, query string) (*output.QueryPlan, error) {
	if err := ValidateReadOnlyQuery(query); err != nil {
		return nil, err
	}

	explainQuery := "EXPLAIN FORMAT=JSON " + query
	var planText string
	err := a.db.QueryRowContext(ctx, explainQuery).Scan(&planText)
	if err != nil {
		return nil, fmt.Errorf("explain failed: %w", err)
	}

	plan := &output.QueryPlan{
		RawPlan: planText,
		Format:  "json",
	}

	return plan, nil
}

func (a MySQLAdapter) ExplainAnalyze(ctx context.Context, query string) (*output.QueryAnalysis, error) {
	if err := ValidateReadOnlyQuery(query); err != nil {
		return nil, err
	}

	supportsAnalyze, prefix := a.supportsExplainAnalyze(ctx)
	if !supportsAnalyze {
		return nil, fmt.Errorf("EXPLAIN ANALYZE not supported on this server version")
	}

	explainQuery := prefix + " " + query

	start := time.Now()
	rows, err := a.db.QueryContext(ctx, explainQuery)
	if err != nil {
		return nil, fmt.Errorf("explain analyze failed: %w", err)
	}
	defer rows.Close()

	var analysisText strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, err
		}
		analysisText.WriteString(line)
		analysisText.WriteString("\n")
	}

	executionTime := float64(time.Since(start).Microseconds()) / 1000.0

	analysis := &output.QueryAnalysis{
		RawAnalysis:     analysisText.String(),
		ExecutionTimeMs: executionTime,
	}

	return analysis, rows.Err()
}

func (a MySQLAdapter) supportsExplainAnalyze(ctx context.Context) (bool, string) {
	var versionStr string
	err := a.db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&versionStr)
	if err != nil {
		return false, "EXPLAIN ANALYZE"
	}

	isMariaDB := strings.Contains(strings.ToLower(versionStr), "mariadb")
	versionStr = strings.ToLower(versionStr)

	var major, minor, patch int
	fmt.Sscanf(versionStr, "%d.%d.%d", &major, &minor, &patch)

	if isMariaDB {
		if major > 10 || (major == 10 && minor >= 11) {
			return true, "EXPLAIN ANALYZE"
		}
		return false, ""
	}

	if major > 8 || (major == 8 && minor > 0) || (major == 8 && minor == 0 && patch >= 18) {
		return true, "EXPLAIN ANALYZE"
	}

	return false, ""
}

func (a MySQLAdapter) ExecuteQuery(ctx context.Context, query string) (int64, int64, error) {
	if err := ValidateReadOnlyQuery(query); err != nil {
		return 0, 0, err
	}

	if a.db == nil {
		return 0, 0, fmt.Errorf("not connected")
	}

	start := time.Now()
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var rowCount int64
	for rows.Next() {
		rowCount++
	}
	if err := rows.Err(); err != nil {
		return 0, 0, err
	}

	executionTimeMs := int64(float64(time.Since(start).Microseconds()) / 1000.0)
	return rowCount, executionTimeMs, nil
}

func (a MySQLAdapter) GetWarnings(ctx context.Context) ([]string, error) {
	if a.db == nil {
		return nil, fmt.Errorf("not connected")
	}

	warnRows, err := a.db.QueryContext(ctx, "SHOW WARNINGS")
	if err != nil {
		return nil, err
	}
	defer warnRows.Close()

	var warnings []string
	for warnRows.Next() {
		var level, code, message string
		if err := warnRows.Scan(&level, &code, &message); err == nil {
			warnings = append(warnings, "["+level+"] "+code+": "+message)
		}
	}
	return warnings, warnRows.Err()
}
