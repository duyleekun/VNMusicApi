package main

import (
	_ "musicsrus/docs"
	_ "musicsrus/routers"

	"github.com/astaxie/beego"
	"os"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run(":"+os.Getenv("PORT"))
}
