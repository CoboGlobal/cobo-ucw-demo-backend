//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"cobo-ucw-backend/internal/biz"
	"cobo-ucw-backend/internal/conf"
	"cobo-ucw-backend/internal/data"
	"cobo-ucw-backend/internal/middleware"
	"cobo-ucw-backend/internal/server"
	"cobo-ucw-backend/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.UCW_Auth, *conf.UCW, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, middleware.ProviderSet, biz.ProviderSet, data.ProviderSet, service.ProviderSet, biz.CronProviderSet, newApp))
}
