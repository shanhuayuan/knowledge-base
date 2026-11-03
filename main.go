package main

import (
	"app/global"
	"app/router"
)

func main() {
	global.InitConfig()
	global.InitMinio()

	r := router.SetupRouter()
	r.Run(global.OverallRuntime.Config.App.Port)
}
