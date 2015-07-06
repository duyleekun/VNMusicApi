package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["musicsrus/controllers:LinkController"] = append(beego.GlobalControllerRouter["musicsrus/controllers:LinkController"],
		beego.ControllerComments{
			"Inspect",
			`/inspect`,
			[]string{"get"},
			nil})

}
