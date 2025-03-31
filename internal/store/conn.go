package store

import (
	"context"

	"github.com/hexiaopi/rdm-toy/internal/model"
)

type ConnStore interface {
	List(ctx context.Context) ([]model.Conn, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, conn *model.Conn) error
	Update(ctx context.Context, conn *model.Conn) error
	Get(ctx context.Context, id string) (*model.Conn, error)
	Delete(ctx context.Context, id string) error
}
