package controllers
import (
	"github.com/astaxie/beego"
	"regexp"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"github.com/mikkyang/id3-go/v2"
	"github.com/jfbus/httprs"
	"github.com/mikkyang/id3-go/v1"
)



// Operations about object
type LinkController struct {
	beego.Controller
}

type Tag struct {
	title, artist, album, year, genre string
}

func OpenUrl(url string) (tag Tag) {
	response, _ := http.Get(url)
	ReadSeekerStream := httprs.NewHttpReadSeeker(response)
	if v2Tag := v2.ParseTag(ReadSeekerStream); v2Tag != nil {
		return Tag{
			title: v2Tag.Title(),
			album: v2Tag.Album(),
			artist: v2Tag.Artist(),
			year: v2Tag.Year(),
			genre: v2Tag.Genre()}
	} else if v1Tag := v1.ParseTag(ReadSeekerStream); v1Tag != nil {
		return Tag{
			title: v1Tag.Title(),
			album: v1Tag.Album(),
			artist: v1Tag.Artist(),
			year: v1Tag.Year(),
			genre: v1Tag.Genre()}
	}
	return Tag{
		title: "Unknown",
		album: "Unknown",
		artist: "Unknown",
		year: "Unknown",
		genre: "Unknown"}
}


// @Title Inspect
// @router /inspect [get]
func (o *LinkController) Inspect() {
	SourceLink := o.Ctx.Input.Query("url")
	Format := o.Ctx.Input.Query("format")

	var DownloadUrl string
	regex , _ := regexp.Compile("http:\\/\\/(.*)\\/bai-hat\\/.*[\\.\\/](.+)\\.html.*")
	for _,matches := range regex.FindAllStringSubmatch(SourceLink,-1) {
		song_id := matches[2]
		switch matches[1] {
		case "mp3.zing.vn":
			DownloadUrl = fmt.Sprintf("http://v3.mp3.zing.vn/download/vip/song/%s",song_id)
			break;
		case "www.nhaccuatui.com":
			response, _ := http.Get(fmt.Sprintf("http://www.nhaccuatui.com/download/song/%s",song_id))
			json, _ := simplejson.NewFromReader(response.Body)
			stream_url, _ := json.Get("data").Get("stream_url").String()
			DownloadUrl = stream_url
			break;
		}
	}

	switch Format {
	case "json":
		Info:=OpenUrl(DownloadUrl)

		o.Data["json"] = &map[string]interface{} {
			"info": map[string]interface{}{
				"title": Info.title,
				"album": Info.album,
				"artist": Info.artist,
				"year": Info.year,
				"genre": Info.genre},
			"stream_url": DownloadUrl }
		o.Ctx.Output.Header("Made-By","Le Duc Duy, Vietnam")
		o.ServeJson()
	default:
		o.Redirect(DownloadUrl,302)

	}
}