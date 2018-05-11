/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package cn4

import (
	"git.cm/naiba/gopappy"
	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/com"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"strings"
)

const api = "https://4.cn/buynow/index"

var TLDs = map[string]string{
	"com":    "0",
	"net":    "1",
	"org":    "2",
	"cn":     "56",
	"com.cn": "560",
	"cc":     "47",
	"me":     "146",
	"tv":     "230",
}

var Tags = map[string]string{
	"纯数字": "type_1",
	"纯字母": "type_3",
	"单拼":  "79",
	"双拼":  "81",
	"三拼":  "82",
	"杂米":  "type_4",
}

func Domains(o gopappy.Option) (d []gopappy.Domain, err error) {
	r := gorequest.New()
	r.Header["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	r.Header["Referer"] = "https://4.cn/buynow"
	_, body, errs := r.Get(api + getURL(o)).End()
	if len(errs) != 0 {
		err = errs[0]
		return
	}

	d = make([]gopappy.Domain, 0)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return
	}

	table := doc.Find("table.grid tbody").First()
	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var id gopappy.Domain
		id.Platform = "1"
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 2:
				a := td.Find("a").First()
				id.BuyURL = "https://4.cn" + a.AttrOr("href", "")
				id.Name = strings.TrimSpace(a.Text())
			case 3:
				id.Description = strings.TrimSpace(td.Text())
			case 5:
				priceStr := strings.TrimSpace(td.Text())
				priceStr = strings.Replace(priceStr, ",", "", -1)
				if strings.HasPrefix(priceStr, "￥") {
					id.Currency = "CNY"
					id.Price, _ = strconv.Atoi(strings.TrimSpace(priceStr[strings.Index(priceStr, "￥")+3 : strings.Index(priceStr, "元")]))
				}
				if strings.HasPrefix(priceStr, "€") {
					id.Currency = "EUR"
					id.Price, _ = strconv.Atoi(strings.TrimSpace(priceStr[strings.Index(priceStr, "€")+3 : strings.Index(priceStr, "欧元")]))
				}
				if strings.HasPrefix(priceStr, "$") {
					id.Currency = "USD"
					id.Price, _ = strconv.Atoi(strings.TrimSpace(priceStr[strings.Index(priceStr, "$")+1 : strings.Index(priceStr, "美元")]))
				}
			}
		})
		d = append(d, id)
	})

	return
}

func getURL(o gopappy.Option) string {
	s := "/search/1"
	//排序
	if o.Order == 1 {
		s += "/so/price"
		if o.Sort == 1 {
			s += "/sb/asc"
		} else {
			s += "/sb/desc"
		}
	}
	//后缀
	if len(o.TLDs) > 0 {
		for _, tld := range o.TLDs {
			s += "/tlds/" + TLDs[gopappy.TLDs[tld]]
		}
	}
	//关键词
	if len(o.Keyword) > 0 {
		s += "/keyword/" + com.URLEncode(o.Keyword)
		if o.KwPos == 0 {
			s += "/kws2/1/kws/1"
		} else {
			s += "/kws2/1/kws/" + strconv.Itoa(o.KwPos)
		}
	}
	//排除
	if len(o.Exclude) > 0 {
		s += "/exclude/" + com.URLEncode(o.Exclude)
		if o.ExPos == 0 {
			s += "/ekws2/1/ekws/1"
		} else {
			s += "/ekws2/1/ekws/" + strconv.Itoa(o.ExPos)
		}
	}
	//分类
	if len(o.Tags) > 0 {
		for _, tag := range o.Tags {
			s += "/tags/" + Tags[gopappy.Tags[tag]]
		}
	}
	//价格
	if o.MaxPrice > 0 {
		s += "/pmax/" + strconv.Itoa(o.MaxPrice)
	}
	if o.MinPrice > 0 {
		s += "/pmin/" + strconv.Itoa(o.MinPrice)
	}
	//长度
	if o.MaxLength > 0 {
		s += "/lmax/" + strconv.Itoa(o.MaxLength)
	}
	if o.MinLength > 0 {
		s += "/lmin/" + strconv.Itoa(o.MinLength)
	}
	//分页
	if o.Page > 1 {
		s += "/page/" + strconv.Itoa(o.Page)
	}
	return s
}