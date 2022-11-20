package repository

import (
	"context"
	"strings"
	sqlclient "weather/internal/sql-client"

	"github.com/uptrace/bun/schema"
)

var FusionSqlClient sqlclient.ISqlClientCon

func CreateTableCollate(client sqlclient.ISqlClientCon, ctx context.Context, table interface{}) error {
	query := client.GetDB().NewCreateTable().Model(table).IfNotExists()
	value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	queryStr := string(value) + " COLLATE utf8mb4_general_ci"
	_, err := client.GetDB().QueryContext(ctx, queryStr)
	return err
}

func CreateTable(client sqlclient.ISqlClientCon, ctx context.Context, table interface{}) error {
	query := client.GetDB().NewCreateTable().Model(table).IfNotExists()
	value, _ := query.AppendQuery(schema.NewFormatter(query.Dialect()), nil)
	queryStr := string(value)
	if client.GetDriver() == sqlclient.POSTGRESQL {
		queryStr = strings.ReplaceAll(queryStr, " char(36)", " uuid")
		queryStr = strings.ReplaceAll(queryStr, " timestamp", " timestamptz")
	}
	_, err := client.GetDB().QueryContext(ctx, queryStr)
	return err
}
