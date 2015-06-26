package main

import (
	"log"

	"github.com/tukdesk/tukdesk/backend/apis"
	"github.com/tukdesk/tukdesk/backend/config"
	"github.com/tukdesk/tukdesk/backend/models/helpers"

	"github.com/tukdesk/httputils/gojimiddleware"
	"github.com/tukdesk/httputils/graceful"
	"github.com/tukdesk/httputils/jsonutils"
	"github.com/tukdesk/mgoutils"
	"github.com/zenazn/goji/web"
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
		Attachment: config.AttachmentConfig{
			Internal: config.InternalAttachmentConfig{
				Dir: "./_attachment",
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
	app := web.New()
	app.NotFound(jsonutils.NotFoundHandler)
	apis.RegisterBaseModule(cfg, app)
	apis.RegisterBrandModule(cfg, app)
	apis.RegisterProfileModule(cfg, app)
	apis.RegisterTicketsModule(cfg, app)
	apis.RegisterUserModule(cfg, app)
	apis.RegisterFocusModule(cfg, app)

	if _, err := apis.RegisterAttachmentModule(cfg, app); err != nil {
		log.Fatalln(err)
	}

	server := web.New()

	server.Use(gojimiddleware.RequestLogger)
	server.Use(gojimiddleware.RequestTimer)
	server.Use(gojimiddleware.RecovererJson)

	gojimiddleware.RegisterSubroute("/apis/v1", server, app)

	log.Println("Server on ", cfg.Addr)
	if err := graceful.Serve(server, cfg.Addr); err != nil {
		log.Fatalln(err)
	}
}
