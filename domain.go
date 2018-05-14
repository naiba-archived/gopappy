/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gopappy

var Platform = map[int]string{
	1: "金名网[4.cn]",
	2: "易名[ename.com]",
	3: "阿里云",
}
var TLDs = map[int]string{
	1: "com",
	2: "net",
	3: "org",
	4: "cn",
	5: "com.cn",
	6: "cc",
	8: "me",
	7: "tv",
}

var Tags = map[int]string{
	1: "纯字母",
	2: "纯数字",
	3: "单拼",
	4: "双拼",
	5: "三拼",
	6: "杂米",
}

var Position = map[int]string{
	1: "开头",
	2: "结尾",
}

var Order = map[int]string{
	1: "价格",
}

var Sort = map[int]string{
	1: "正序",
	2: "倒叙",
}

type Domain struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"-"`
	PriceString string `json:"price"`
	Currency    string `json:"currency"`
	BuyURL      string `json:"buy_url"`
	Platform    int    `json:"platform"`
}
type SortDomain []Domain

func (sd SortDomain) Len() int {
	return len(sd)
}
func (sd SortDomain) Swap(i, j int) { sd[i], sd[j] = sd[j], sd[i] }
func (sd SortDomain) Less(i, j int) bool {
	return sd[i].Price < sd[j].Price
}

type Option struct {
	Page  int // 分页
	Order int //排序
	Sort  int //正序 倒叙

	Keyword string //关键词
	KwPos   int    //关键词位置
	Exclude string //排除
	ExPos   int    //排除词位置
	Tag     int    //分类
	TLDs    []int  //后缀

	MinPrice  int //最低价格
	MaxPrice  int //最高价格
	MinLength int //最短长度
	MaxLength int //最长长度
}
