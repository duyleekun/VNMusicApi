package controllers
import (
	"github.com/astaxie/beego"
	"regexp"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
)


// Operations about object
type GetLinkController struct {
	beego.Controller
}


// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *GetLinkController) Get() {
	SourceLink := o.Ctx.Input.Query("url")


	regex , _ := regexp.Compile("http:\\/\\/(.*)\\/bai-hat\\/.*[\\.\\/](.+)\\.html.*")
	for _,matches := range regex.FindAllStringSubmatch(SourceLink,-1) {
		song_id := matches[2]
		switch matches[1] {
		case "mp3.zing.vn":
			o.Data["json"] = &map[string]string {"url": (fmt.Sprintf("http://v3.mp3.zing.vn/download/vip/song/%s",song_id))}
			break;
		case "www.nhaccuatui.com":
			response, _ := http.Get(fmt.Sprintf("http://www.nhaccuatui.com/download/song/%s",song_id))
			json, _ := simplejson.NewFromReader(response.Body)
			stream_url, _ := json.Get("data").Get("stream_url").String()
			o.Data["json"] = &map[string]string {"url": stream_url}
			break;
		}
	}

	o.ServeJson()
}