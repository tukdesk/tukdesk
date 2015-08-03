package main

import (
	"log"

	"github.com/tukdesk/tukdesk/backend/apis"
	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/labstack/echo"
	emw "github.com/tukdesk/httputils/echomiddleware"
	"github.com/tukdesk/httputils/graceful"
	"github.com/tukdesk/mgoutils"
)

func main() {
	// config for dev
	cfg := &config.Config{
		Addr: "127.0.0.1:52081",
		Salt: "ivcnwHHpGZ",
		Database: config.DatabaseConfig{
			DBURL:  "127.0.0.1:27017",
			DBName: "tukdesk_dev",
		},
		File: config.FileConfig{
			Internal: config.InternalFileConfig{
				Dir: "./_file",
			},
		},
	}

	// init database storage
	stg, err := mgoutils.NewMgoPool(cfg.Database.DBURL, cfg.Database.DBName)
	if err != nil {
		log.Fatalln(err)
	}

	if err := helpers.InitWithStorage(stg); err != nil {
		log.Fatalln(err)
	}

	// init app
	app := echo.New()
	app.SetHTTPErrorHandler(emw.JSONErrHandlerForAPIError())
	app.Use(emw.RequestLogger())
	app.Use(emw.RequestTimer())
	app.Use(emw.JSONRecoverForAPIError())

	mux := app.Group("/apis/v1")
	apis.RegisterBaseModule(cfg, mux)
	apis.RegisterBrandModule(cfg, mux)
	apis.RegisterProfileModule(cfg, mux)
	apis.RegisterTicketsModule(cfg, mux)
	apis.RegisterUserModule(cfg, mux)

	log.Println("Server on ", cfg.Addr)
	if err := graceful.Serve(app, cfg.Addr); err != nil {
		log.Fatalln(err)
	}
}
