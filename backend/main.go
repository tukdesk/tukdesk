package main

import (
	"log"

	"github.com/tukdesk/tukdesk/backend/apis"
	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/mgoutils"
	"github.com/zenazn/goji/web"
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
	app := web.New()
	apis.RegisterBaseModule(cfg, app)
	apis.RegisterBrandModule(cfg, app)
	apis.RegisterProfileModule(cfg, app)
	apis.RegisterTicketsModule(cfg, app)
	apis.RegisterUserModule(cfg, app)
	apis.RegisterFocusModule(cfg, app)

	service := gojimiddleware.NewApp()
	service.Mux().Use(gojimiddleware.RequestLogger)
	service.Mux().Use(gojimiddleware.RequestTimer)
	service.Mux().Use(gojimiddleware.RecovererJson)

	gojimiddleware.RegisterSubroute("/apis/v1", service.Mux(), app)

	if err := service.Run(cfg.Addr); err != nil {
		log.Fatalln(err)
	}
}
