/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package gopappy

import "github.com/naiba/com"

var Fx *com.Fixer

func init() {
	var err error
	Fx, err = com.NewFixer()
	if err != nil {
		panic(err)
	}
}
