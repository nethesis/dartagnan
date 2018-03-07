/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dartagnan project.
 *
 * Dartagnan is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dartagnan is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dartagnan.  If not, see COPYING.
 *
 */

package main

import (
	"os"
	"fmt"
	"flag"
        _ "github.com/jinzhu/gorm/dialects/mysql"

        "github.com/nethesis/dartagnan/athos/configuration"
        "github.com/nethesis/dartagnan/athos/cache"
)


func main() {
        // read and init configuration
        ConfigFilePtr := flag.String("c", "/opt/dartagnan/athos/conf.json", "Path to configuration file")
        flag.Parse()
        configuration.Init(ConfigFilePtr)

	ret, errors := cache.BulkSetValidSystems()
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, err)
	}

	if ret {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
