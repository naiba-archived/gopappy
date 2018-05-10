/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gopappy

import (
	"encoding/json"
	"github.com/naiba/com"
	"io/ioutil"
)

type Stat struct {
	Search int64
	Domain int64
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
