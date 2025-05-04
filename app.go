package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/rwpp/RzWeLook/pkg/saramax"
)

type App struct {
	web       *gin.Engine
	consumers []saramax.Consumer
	cron      *cron.Cron
}
