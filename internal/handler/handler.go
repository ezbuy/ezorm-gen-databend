package handler

import (
	"context"

	"github.com/ezbuy/ezorm/v2/pkg/plugin"
)

type SchemaHandler interface {
	Handle(context.Context, plugin.Schema) error
}

type Printer interface {
	Print(ctx context.Context, dir string) error
}
