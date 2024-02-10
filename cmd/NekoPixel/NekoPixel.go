package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"github.com/wuhan005/NekoPixel/internal/conf"
	"github.com/wuhan005/NekoPixel/internal/db"
	"github.com/wuhan005/NekoPixel/internal/route"
)

func main() {
	port := flag.Int("port", 8080, "port to listen")
	flag.Parse()

	if err := conf.Init(); err != nil {
		logrus.WithError(err).Fatal("Failed to initialize config")
	}

	db, err := db.Init()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}

	f := route.New(db)
	f.Run(*port)
}
