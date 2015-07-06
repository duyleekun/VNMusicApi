package main

import (
	_ "musicsrus/docs"
	_ "musicsrus/routers"

	"github.com/astaxie/beego"
	"os"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter,cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Made-By"},
		AllowCredentials: true,
	}))

	beego.Run(":"+os.Getenv("PORT"))
}
