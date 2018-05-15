/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package ename

import (
	"git.cm/naiba/gopappy"
	"github.com/PuerkitoBio/goquery"
	"github.com/naiba/com"
	"github.com/parnurzeal/gorequest"
	"log"
	"strconv"
	"strings"
	"time"
)

var TLDs = map[string]string{
	"com":    "1",
	"net":    "7",
	"org":    "8",
	"cn":     "2",
	"com.cn": "3",
	"cc":     "9",
	"me":     "15",
	"tv":     "16",
}

var Tags = map[string]string{
	"纯数字": "1",
	"纯字母": "7",
	"单拼":  "16",
	"双拼":  "17",
	"三拼":  "18",
	"杂米":  "13",
}

const api = "https://auction.ename.com/tao"

func Domains(o gopappy.Option) (d []gopappy.Domain, e error) {
	d = make([]gopappy.Domain, 0)

	appended := getURL(o)
	r := gorequest.New().Timeout(time.Second * 4)
	log.Println(api + appended)
	_, body, errs := r.Get(api+appended).
		Set("Referer", api+"/search"+appended).
		Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36").
		Set("X-Forward-For", com.RandomIP()).
		Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6").
		End()
	if len(errs) != 0 {
		e = errs[0]
		return
	}

	doc, e := goquery.NewDocumentFromReader(strings.NewReader(body))
	if e != nil {
		return
	}

	table := doc.Find("table.com_table tbody").First()
	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if len(tr.AttrOr("id", "")) > 0 {
			return
		}
		var id gopappy.Domain
		id.Platform = 2
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			switch j {
			case 1:
				a := td.Find("a").First()
				id.Name = strings.TrimSpace(a.Text())
				id.BuyURL = "https:" + a.AttrOr("href", "")
			case 2:
				id.Description = td.AttrOr("title", "")
			case 4:
				id.Price, _ = strconv.Atoi(strings.TrimSpace(td.Find("span.c_red").First().Text()))
				id.Currency = "CNY"
			}
		})
		d = append(d, id)
	})

	return d, nil
}

func getURL(o gopappy.Option) string {
	s := "?domain=2"
	if len(o.Keyword) > 0 {
		s += "&domainsld=" + com.URLEncode(o.Keyword)
		if o.KwPos == 1 {
			s += "&sldtypestart=1"
		} else if o.KwPos == 2 {
			s += "&sldtypeend=1"
		}
	}

	if o.Order == 1 {
		s += "&sort=2"
	}
	s += "&bidpricestart=" + strconv.Itoa(o.MinPrice)
	if o.MaxPrice > 0 {
		s += "&bidpriceend=" + strconv.Itoa(o.MaxPrice)
	}
	if len(o.Exclude) > 0 {
		s += "&skipword1=" + o.Exclude
		if o.ExPos == 1 {
			s += "&skipstart1=1"
		} else if o.ExPos == 2 {
			s += "&skipend1=1"
		}
	}
	if len(o.TLDs) > 0 {
		for _, tld := range o.TLDs {
			s += "&domaintld[]=" + TLDs[gopappy.TLDs[tld]]
		}
	}
	if o.MinLength > 0 {
		s += "&domainlenstart=" + strconv.Itoa(o.MinLength)
	}
	if o.MaxLength > 0 {
		s += "&domainlenend=" + strconv.Itoa(o.MaxLength)
	}
	if o.Tag > 0 {
		s += "&domaingroup=" + Tags[gopappy.Tags[o.Tag]]
	}
	if o.Page > 0 {
		s += "&page=" + strconv.Itoa(o.Page)
	}
	s += "&pageSize=20"
	return s
}
