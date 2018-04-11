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

	"github.com/nethesis/dartagnan/athos/utils"
	"github.com/spf13/cobra"
)

var trialHours int

var trialsCmd = &cobra.Command{
	Use:   "trials",
	Short: "List active trials",
	Long: `List all systems with an associated trial subscription and an heartbeat sent in the last X hours (default 24)

	Examples

	List trials active in the last 24 hours (default)
 	  foil payments -H 24

	`,
	Run: printActiveTrials,
}

func printActiveTrials(cmd *cobra.Command, args []string) {
	systems := utils.ListActiveTrials(trialHours)
	fmt.Printf("System\tUser\tSubscription_end\tEmail\t\n\n")
	for _, system := range systems {
		until := fmt.Sprintf("%d-%02d-%02d", system.Subscription.ValidUntil.Year(), system.Subscription.ValidUntil.Month(), system.Subscription.ValidUntil.Day())
		fmt.Printf("%-40v%-40v%-15v%s\n", system.UUID, system.CreatorID, until, system.Notification["emails"])
	}
}

func init() {
	trialsCmd.Flags().IntVarP(&trialHours, "hours", "H", 24, "Hours from last hearbeat")

	rootCmd.AddCommand(trialsCmd)
}
