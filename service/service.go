/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package service

import (
	"git.cm/naiba/gopappy"
)

var Stat *gopappy.Stat

func Init() {
	if Stat == nil {
		Stat = new(gopappy.Stat)
		Stat.Read()
	}
}
