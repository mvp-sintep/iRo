package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Query - выполнение запросов
func (o *Pool) Query(sql string, args ...any) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(o.context, time.Duration(time.Duration(o.cfg.Timeout)*time.Millisecond))
	defer cancel()
	return o.pgx.Query(ctx, sql)
}
