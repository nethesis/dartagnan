/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dartagnan project.
 *
 * Dartagnan is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dartagnan is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dartagnan.  If not, see COPYING.
 *
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/nethesis/dartagnan/athos/cache"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "synckeys",
	Short: "synchronize keys",
	Long:  "Synchronize activation keys from database to Redis master instance.",
	Run:   sync,
}

func sync(*cobra.Command, []string) {
	_, errors := cache.BulkSetValidSystems()
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, err)
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
