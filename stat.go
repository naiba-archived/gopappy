/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gopappy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/naiba/com"
)

var StatInstance *Stat

func init() {
	if StatInstance == nil {
		StatInstance = new(Stat)
		StatInstance.Read()
	}
	go func() {
		time.Sleep(time.Minute * 10)
		StatInstance.Save()
		log.Println("stat auto save", StatInstance)
	}()
}

type Stat struct {
	Search int64 `json:"search"`
	Domain int64 `json:"domain"`
}

const filePath = "resource/data/stat.json"

func (s *Stat) Read() {
	bStat, err := ioutil.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(bStat, s)
		com.PanicIfNotNil(err)
	}
}

func (s *Stat) Save() {
	bStat, err := json.Marshal(s)
	com.PanicIfNotNil(err)
	err = ioutil.WriteFile(filePath, bStat, 0655)
	com.PanicIfNotNil(err)
}
