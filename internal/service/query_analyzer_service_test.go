package service

import (
	"context"
	"database-query-analyzer/internal/input"
	"database-query-analyzer/internal/output"
	"database-query-analyzer/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryAnalyzerService_GetDatabaseInfo(t *testing.T) {
	ctx := context.Background()
	errInfo := errors.New("info failed")
	errTables := errors.New("tables failed")
	dbInfo := &output.DatabaseInfo{Name: "appdb", Version: "16.2"}

	tests := []struct {
		name    string
		in      input.GetDatabaseInfoInput
		setup   func(repo *mocks.DatabaseAdapter)
		want    output.GetDatabaseInfoOutput
		wantErr error
	}{
		{
			name: "repository error on database info returns empty output",
			in:   input.GetDatabaseInfoInput{IncludeTables: true},
			setup: func(repo *mocks.DatabaseAdapter) {
				repo.On("GetDatabaseInfo", ctx).Return(&output.DatabaseInfo{}, errInfo).Once()
			},
			want:    output.GetDatabaseInfoOutput{},
			wantErr: errInfo,
		},
		{
			name: "without IncludeTables does not fetch tables",
			in:   input.GetDatabaseInfoInput{IncludeTables: false},
			setup: func(repo *mocks.DatabaseAdapter) {
				repo.On("GetDatabaseInfo", ctx).Return(dbInfo, nil).Once()
			},
			want: output.GetDatabaseInfoOutput{Database: "appdb", Version: "16.2"},
		},
		{
			name: "with IncludeTables returns table names in order",
			in:   input.GetDatabaseInfoInput{IncludeTables: true},
			setup: func(repo *mocks.DatabaseAdapter) {
				repo.On("GetDatabaseInfo", ctx).Return(dbInfo, nil).Once()
				repo.On("GetTables", ctx).
					Return([]output.Table{{Name: "users"}, {Name: "roles"}}, nil).Once()
			},
			want: output.GetDatabaseInfoOutput{
				Database: "appdb",
				Version:  "16.2",
				Tables:   []string{"users", "roles"},
			},
		},
		{
			name: "with IncludeTables and no tables leaves Tables empty",
			in:   input.GetDatabaseInfoInput{IncludeTables: true},
			setup: func(repo *mocks.DatabaseAdapter) {
				repo.On("GetDatabaseInfo", ctx).Return(dbInfo, nil).Once()
				repo.On("GetTables", ctx).Return([]output.Table{}, nil).Once()
			},
			want: output.GetDatabaseInfoOutput{Database: "appdb", Version: "16.2"},
		},
		{
			name: "table fetch error returns partial output and the error",
			in:   input.GetDatabaseInfoInput{IncludeTables: true},
			setup: func(repo *mocks.DatabaseAdapter) {
				repo.On("GetDatabaseInfo", ctx).Return(dbInfo, nil).Once()
				repo.On("GetTables", ctx).Return(nil, errTables).Once()
			},
			want:    output.GetDatabaseInfoOutput{Database: "appdb", Version: "16.2"},
			wantErr: errTables,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.DatabaseAdapter)
			tt.setup(repo)
			svc := &QueryAnalyzerService{repository: repo}

			got, err := svc.GetDatabaseInfo(ctx, tt.in)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			repo.AssertExpectations(t)
		})
	}
}
