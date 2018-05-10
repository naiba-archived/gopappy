package com

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

//Fixer api.fixer.io
type Fixer struct {
	Base  string
	Date  string
	Rates map[string]float64
}

var rateInstance *Fixer

//NewFixer 获取汇率对象
func NewFixer() (*Fixer, error) {
	if rateInstance != nil {
		return rateInstance, nil
	}
	req := gorequest.New()
	_, body, err := req.Get("https://api.fixer.io/latest").End()
	if err != nil {
		return nil, err[0]
	}
	var rate Fixer
	errs := json.Unmarshal([]byte(body), &rate)
	if errs != nil {
		return nil, err[0]
	}
	rateInstance = &rate
	return rateInstance, nil
}

//Convert 获取汇率
func (m *Fixer) Convert(src string, dist string, num float64) float64 {
	rate, has := m.Rates[dist]
	rate1, has1 := m.Rates[src]
	if has && (has1 || src == m.Base) {
		if src == m.Base {
			return num * rate
		}
		return num / rate1 * rate
	}
	return num
}
