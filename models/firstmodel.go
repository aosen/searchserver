package models

import (
//"github.com/aosen/kernel"
)

type FirstModel struct {
	BaseModel
}

type second struct {
	Title string //二级分类的标题
	Url   string //二级分类的地址
	Id    int    //二级分类的ID
	Count int    //话题个数
}

type Fream struct {
	Id         int
	Title      string
	Desc       string
	BigPic     string
	SecondList []second
	Count      int //一级分类下的总发布量
}

func (self *FirstModel) LoadFrame() (c Fream) {
	c.Id = 1
	c.Title = "搞笑&奇葩"
	c.Desc = "各种搞笑，各种奇葩"
	c.BigPic = "http://i13.tietuku.com/0d2a7be856ebf374.jpg"
	c.Count = 9999999
	c.SecondList = []second{
		{
			Title: "搞笑",
			Url:   "/categorys/fun/",
			Id:    1,
			Count: 10000,
		},
		{
			Title: "奇葩",
			Url:   "/categorys/miracle/",
			Id:    2,
			Count: 8000,
		},
	}
	return
}

func (self *FirstModel) LoadData() (d []Good) {
	for i := 0; i < 48; i++ {
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
