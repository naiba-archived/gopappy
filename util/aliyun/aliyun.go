/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package aliyun

import (
	"git.cm/naiba/gopappy"
	"github.com/naiba/com"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const api = "https://domainapi.aliyun.com/onsale/search?"

var Tags = map[string]string{
	"纯数字": "11",
	"纯字母": "12",
	"单拼":  "1207001",
	"双拼":  "1207002",
	"三拼":  "1207003",
	"杂米":  "13",
}

type aliyunResult struct {
	Code string
	Data struct {
		PageResult struct {
			CurrentPageNum int
			Data []struct {
				DomainName   string
				Price        string
				Introduction string
			}
		}
	}
}

func Domains(o gopappy.Option) (d []gopappy.Domain, err error) {
	d = make([]gopappy.Domain, 0)
	r := gorequest.New().Timeout(time.Second * 4)
	var res aliyunResult
	_, body, errs := r.Get(api + getURL(o)).
		Set("Referer", "https://mi.aliyun.com/").
		Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36").
		Set("X-Forward-For", com.RandomIP()).
		Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6").
		Set("Authority", "domainapi.aliyun.com").
		EndStruct(&res)
	if len(errs) != 0 {
		err = errs[0]
		return
	}
	if res.Code != "200" {
		err = errors.New("返回内容：" + string(body))
		return
	}
	for _, ali := range res.Data.PageResult.Data {
		var id gopappy.Domain
		id.Platform = 3
		id.Currency = "CNY"
		id.Price, _ = strconv.Atoi(strings.Replace(ali.Price, ",", "", -1))
		id.Description = ali.Introduction
		id.BuyURL = "https://mi.aliyun.com/detail/online.html?domainName=" + ali.DomainName
		id.Name = ali.DomainName
		d = append(d, id)
	}
	return
}

func getURL(o gopappy.Option) string {
	s := "productType=2&searchIntro=false&pageSize=20&fetchSearchTotal=true&token=tdomain-aliyun-com:" + com.RandomString(32)
	if o.Page > 0 {
		s += "&currentPage=" + strconv.Itoa(o.Page)
	}
	if len(o.TLDs) > 0 {
		s += "&suffix="
		for _, tld := range o.TLDs {
			s += gopappy.TLDs[tld] + "|"
		}
		s = s[:len(s)-1]
	}
	if len(o.Keyword) > 0 {
		s += "&keyWord=" + com.URLEncode(o.Keyword)
		if o.KwPos == 1 {
			s += "&keywordAsPrefix=true"
		} else if o.KwPos == 2 {
			s += "&keywordAsSuffix=true"
		}
	}
	if len(o.Exclude) > 0 {
		s += "&excludeKeyWord=" + o.Exclude
		if o.ExPos == 1 {
			s += "&exKeywordAsPrefix=true"
		} else if o.ExPos == 2 {
			s += "&exKeywordAsSuffix=true"
		}
	}
	{
		if o.MinPrice > 0 {
			s += "&minPrice=" + strconv.Itoa(o.MinPrice)
		}
		if o.MaxPrice > 0 {
			s += "&maxPrice=" + strconv.Itoa(o.MaxPrice)
		}
		if o.MinLength > 0 {
			s += "&minLength=" + strconv.Itoa(o.MinLength)
		}
		if o.MaxLength > 0 {
			s += "&maxLength=" + strconv.Itoa(o.MaxLength)
		}
	}
	if o.Order == 1 {
		s += "&sortBy=3"
		if o.Sort == 1 {
			s += "&sortType=2"
		} else {
			s += "&sortType=1"
		}
	}
	if o.Tag > 0 {
		s += "&constitute=" + Tags[gopappy.Tags[o.Tag]]
	}
	return s
}
