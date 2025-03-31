package service

import (
	"context"

	"github.com/hexiaopi/rdm-toy/internal/config"
	"github.com/hexiaopi/rdm-toy/internal/model"
)

func ListConn(ctx context.Context) ([]model.Conn, error) {
	result := make([]model.Conn, 0)
	for k, v := range config.AppEngine.Redis {
		v.Name = k
		result = append(result, v)
	}
	return result, nil
}
