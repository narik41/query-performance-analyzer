package output

import "time"

// DatabaseInfo represents basic database information
type DatabaseInfo struct {
	Name    string
	Version string
}

// Table represents a database table
type Table struct {
	Name      string
	Schema    string
	Comment   string
	RowCount  int64
	SizeBytes int64
	Engine    string
	Charset   string
	Collation string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// Column represents a table column
type Column struct {
	Name            string
	Type            string
	Nullable        bool
	Key             string
	Default         *string
	Extra           string
	Comment         string
	OrdinalPosition int
	CharacterSet    string
	Collation       string
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	Name             string
	Column           string
	ReferencedTable  string
	ReferencedColumn string
	OnUpdate         string
	OnDelete         string
}

// TableSchema represents the complete schema of a table
type TableSchema struct {
	TableName   string
	Columns     []Column
	ForeignKeys []ForeignKey
	Indexes     []Index
}

// Index represents a table index
type Index struct {
	Name        string
	Table       string
	Unique      bool
	Type        string
	Columns     []string
	Cardinality int64
	Comment     string
	Visible     bool
}

// QueryPlan represents an EXPLAIN output
type QueryPlan struct {
	RawPlan       string
	Format        string // "json", "text", "tree"
	EstimatedCost float64
	EstimatedRows int64
}

// QueryAnalysis represents an EXPLAIN ANALYZE output
type QueryAnalysis struct {
	RawAnalysis     string
	ExecutionTimeMs float64
	ActualRows      int64
	Plan            *QueryPlan
}

// GetDatabaseInfoOutput represents the output for getting database info
type GetDatabaseInfoOutput struct {
	Database string   `json:"database"`
	Version  string   `json:"version"`
	Tables   []string `json:"tables,omitempty"`
}

// GetSchemaOutput represents the output for getting table schema
type GetSchemaOutput struct {
	Tables []TableSchema `json:"tables"`
}

// GetIndexesOutput represents the output for getting indexes
type GetIndexesOutput struct {
	Indexes []Index `json:"indexes"`
}

// ExplainOutput represents the output for explain queries
type ExplainOutput struct {
	Plan string `json:"plan"`
}

// QueryStatsOutput represents the output for getting query statistics
type QueryStatsOutput struct {
	SQL             string   `json:"sql"`
	ExecutionTimeMs int64    `json:"execution_time_ms"`
	RowsReturned    int64    `json:"rows_returned"`
	ExplainJSON     string   `json:"explain_json"`
	ExplainAnalyze  string   `json:"explain_analyze"`
	WarningCount    int      `json:"warning_count,omitempty"`
	Warnings        []string `json:"warnings,omitempty"`
}
