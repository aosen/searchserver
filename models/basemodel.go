package models

import (
	"github.com/aosen/kernel"
)

type BaseModel struct {
	kernel.D
}

type S struct {
	Second     string
	SecondDesc string
	SecondPic  string
	SecondUrl  string
}

type F struct {
	First     string
	FirstDesc string
	FirstPic  string
	FirstUrl  string
	SecondMsg []S
}

func (self *BaseModel) Classifction() []F {
	sql_first := "SELECT * FROM firstclass"
	sql_second := "SELECT * FROM secondclass"
	ret_first, err := self.Query(sql_first)
	kernel.PutError(err)
	ret_second, err := self.Query(sql_second)
	kernel.PutError(err)
	ret := []F{}
	for i := 0; i < len(ret_first); i++ {
		secondmsglist := []S{}
		for j := 0; j < len(ret_second); j++ {
			if ret_second[j]["first_id"] == ret_first[i]["id"] {
				tmp2 := S{
					Second:     ret_second[j]["second"],
					SecondDesc: ret_second[j]["seconddesc"],
					SecondPic:  ret_second[j]["secondpic"],
					SecondUrl:  ret_second[j]["secondurl"],
				}
				secondmsglist = append(secondmsglist, tmp2)
			}
		}
		tmp1 := F{
			First:     ret_first[i]["first"],
			FirstDesc: ret_first[i]["firstdesc"],
			FirstPic:  ret_first[i]["firstpic"],
			FirstUrl:  ret_first[i]["firsturl"],
			SecondMsg: secondmsglist,
		}
		ret = append(ret, tmp1)
	}
	return ret
}

type Cl struct {
	FirstName   string
	FirstUrl    string
	SecondLeft  []SecondInfo
	SecondRight []SecondInfo
}

type SecondInfo struct {
	SecondName string
	SecondUrl  string
}

func (self *BaseModel) LoadClass() (c []Cl) {
	//加载分类
	ret := self.Classifction()
	for _, f := range ret {
		temp_first := Cl{
			FirstName: f.First,
			FirstUrl:  "/" + f.FirstUrl,
		}
		l := len(f.SecondMsg)
		var i int = 0
		secondleft := []SecondInfo{}
		for ; i < l/2; i++ {
			temp_second := SecondInfo{
				SecondName: f.SecondMsg[i].Second,
				SecondUrl:  "/" + f.FirstUrl + "/" + f.SecondMsg[i].SecondUrl,
			}
			secondleft = append(secondleft, temp_second)

		}
		temp_first.SecondLeft = secondleft
		secondright := []SecondInfo{}
		for ; i < l; i++ {
			temp_second := SecondInfo{
				SecondName: f.SecondMsg[i].Second,
				SecondUrl:  "/" + f.FirstUrl + "/" + f.SecondMsg[i].SecondUrl,
			}
			secondright = append(secondright, temp_second)

		}
		temp_first.SecondRight = secondright
		c = append(c, temp_first)
	}
	return
}

//商品信息
type Good struct {
	GoodId       string //商品ID
	First        string //一级分类
	Second       string //二级分类
	GoodPrice    string //金额
	GoodUrl      string //目的地址
	Title        string //标题
	Author       string //作者
	AuthorUrl    string //作者主页
	Pic          string //图片
	Support      string //支持点赞地址
	SupportCount string //点赞量
	SupportLevel string //支持等级
}
