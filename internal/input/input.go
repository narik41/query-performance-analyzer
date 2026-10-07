package input

type ConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	DSN      string
}

type GetDatabaseInfoInput struct {
	IncludeTables bool `json:"include_tables,omitempty" jsonschema:"also list table names in the current database"`
}

type GetSchemaInput struct {
	Table string `json:"table,omitempty" jsonschema:"table name; omit to describe every table"`
}

type GetIndexesInput struct {
	Table string `json:"table,omitempty" jsonschema:"table name; omit to list indexes for every table"`
}

type ExplainInput struct {
	SQL string `json:"sql" jsonschema:"a single SELECT statement"`
}

type QueryStatsInput struct {
	SQL string `json:"sql" jsonschema:"a single SELECT statement"`
}
