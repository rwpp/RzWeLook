package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rwpp/RzWeLook/internal/events"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
