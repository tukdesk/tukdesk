package main

import (
	"log"

	"github.com/tukdesk/tukdesk/backend/apis"
	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/mgoutils"
)

func main() {
	// config for dev
	cfg := config.Config{
		Addr: "127.0.0.1:52081",
		Database: config.DatabaseConfig{
			DBURL:  "127.0.0.1:27017",
			DBName: "tukdesk_dev",
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
	app := gojimiddleware.NewApp()

	apis.RegisterBaseModule(cfg, app.Mux())
	apis.RegisterBrandModule(cfg, app.Mux())
	apis.RegisterProfileModule(cfg, app.Mux())

	app.Mux().Use(gojimiddleware.RequestLogger)
	app.Mux().Use(gojimiddleware.RequestTimer)
	app.Mux().Use(gojimiddleware.RecovererJson)

	if err := app.Run(cfg.Addr); err != nil {
		log.Fatalln(err)
	}
}
