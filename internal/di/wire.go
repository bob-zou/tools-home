//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package di

import (
	"tools-home/internal/dao"
	"tools-home/internal/server/http"
	"tools-home/internal/service"

	"github.com/google/wire"
)

//go:generate wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, NewApp))
}
