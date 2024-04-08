package server

import (
	"context"

	"github.com/hexiaopi/rdm-toy/internal/routers"
	"github.com/hexiaopi/rdm-toy/pkg/app"
	"github.com/hexiaopi/rdm-toy/pkg/server/http"

	log "github.com/sirupsen/logrus"
)

func Run() {
	router := routers.NewRouter()
	httpServer := http.NewServer(router,
		http.WithServerHost("0.0.0.0"),
		http.WithServerPort(8080),
	)
	if err := app.NewApp(app.WithServer(httpServer)).Run(context.Background()); err != nil {
		log.Errorf("start http server err:%v", err)
	}
}
