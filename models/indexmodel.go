package models

import (
	"github.com/aosen/kernel"
)

type IndexModel struct {
	BaseModel
}

type Trending struct {
	Title    string
	Desc     string
	BigPic   string
	BigUrl   string
	PicOne   string
	OneUrl   string
	PicTwo   string
	TwoUrl   string
	PicThree string
	ThreeUrl string
}

func (self *IndexModel) LoadTrending() (t []Trending) {
	sql := "SELECT * FROM trending"
	ret, err := self.Query(sql)
	kernel.PutError(err)
	for _, v := range ret {
		t = append(t, Trending{
			Title:    v["title"],
			Desc:     v["desc"],
			BigPic:   v["bigpic"],
			BigUrl:   v["bigurl"],
			PicOne:   v["picone"],
			OneUrl:   v["oneurl"],
			PicTwo:   v["pictwo"],
			TwoUrl:   v["twourl"],
			PicThree: v["picthree"],
			ThreeUrl: v["threeurl"],
		})
	}
	return
}

func (self *IndexModel) LoadData() (d []Good) {
	for i := 0; i < 8; i++ {
		d = append(d, Good{
			GoodId:       string(i),
			First:        "搞笑 & 奇葩",
			Second:       "搞笑",
			GoodPrice:    "5",
			GoodUrl:      "/",
			Title:        "I will fuck fuck fuck fuck you!",
			Author:       "fuck fuck",
			AuthorUrl:    "/",
			Pic:          "http://i11.tietuku.com/a6512ad4adb04b7d.jpg",
			Support:      "/join",
			SupportCount: "800",
			SupportLevel: "10",
		})
	}
	return
}
