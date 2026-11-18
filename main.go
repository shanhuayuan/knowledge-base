package main

import (
	"app/global"
	"app/router"
)

func main() {
	_ = global.InitConfig()
	_ = global.InitMinio()
	_ = global.InitMysql()
	_ = global.InitMilvus()

	global.InitChunker()

	r := router.SetupRouter()
	//fmt.Printf("%v", global.OverallRuntime)
	_ = r.Run(global.OverallRuntime.Config.App.Port)
}
