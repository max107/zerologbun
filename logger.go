package zerologbun

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

type QueryHook struct {
	logger zerolog.Logger
}

func NewQueryHook(logger zerolog.Logger) *QueryHook {
	return &QueryHook{logger}
}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	dur := time.Now().Sub(event.StartTime)

	l := h.logger.
		With().
		Str("type", queryOperation(event.Query)).
		Str("query", event.Query).
		Str("operation", eventOperation(event)).
		Str("duration", dur.String()).
		Logger()

	if event.Err == nil {
		l.Info().Msg("query")
	} else {
		l.Error().Err(event.Err).Msg("query")
	}
}

func eventOperation(event *bun.QueryEvent) string {
	switch event.IQuery.(type) {
	case *bun.SelectQuery:
		return "SELECT"
	case *bun.InsertQuery:
		return "INSERT"
	case *bun.UpdateQuery:
		return "UPDATE"
	case *bun.DeleteQuery:
		return "DELETE"
	case *bun.CreateTableQuery:
		return "CREATE TABLE"
	case *bun.DropTableQuery:
		return "DROP TABLE"
	}
	return queryOperation(event.Query)
}

func queryOperation(name string) string {
	if idx := strings.Index(name, " "); idx > 0 {
		name = name[:idx]
	}
	if len(name) > 16 {
		name = name[:16]
	}
	return name
}
