package apis

import (
	"github.com/labstack/echo"

	"github.com/tukdesk/tukdesk/backend/config"
)

func RegisterFileModule(cfg *config.Config, mux *echo.Group) error {
	group := mux.Group("/files")
	return newInternalFileModule(cfg, group)
}
