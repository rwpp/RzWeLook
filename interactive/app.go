package main

import (
	"github.com/rwpp/RzWeLook/pkg/ginx"
	"github.com/rwpp/RzWeLook/pkg/grpcx"
	"github.com/rwpp/RzWeLook/pkg/saramax"
)

type App struct {
	server    *grpcx.Server
	consumers []saramax.Consumer
	webAdmin  *ginx.Server
}
