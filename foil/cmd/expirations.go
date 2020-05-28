/*
 * Copyright (C) 2020 Nethesis S.r.l.
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

	"github.com/nethesis/dartagnan/athos/notifications"
	"github.com/nethesis/dartagnan/athos/utils"
	"github.com/spf13/cobra"
)

var period string
var sendMail bool
var quietMode bool

var expirationsCmd = &cobra.Command{
	Use:   "expirations",
	Short: "List expiration systems by period (default expire in 1 day)",
	Long: `List all systems that expire in a specified period (default expire in 1 day)

	 Examples

	 List expiration systems by period (default expire in 1 day)
		foil expirations -P 1d

	 List expiration systems by period and send expiration mail to systems owners (default expire in 1 day)
		foil expirations -P 1d -M


	 Supported periods:
	 - 1d (1 day)
	 - 1w (1 week)
	 - 2w (2 weeks) (only for trials)

	 `,
	Run: getExpireSystems,
}

func getExpireSystems(cmd *cobra.Command, args []string) {
	if !utils.ContainsS([]string{"1d", "1w", "2w"}, period) {
		fmt.Println(`Unsupported period.

Supported periods:
- 1d (1 day)
- 1w (1 week)
- 2w (2 weeks) (only for trials)
		`)
		os.Exit(1)
	}

	systems := utils.ListExpirationSystems(period)

	// prints results header
	if !quietMode {
		fmt.Printf("System\t\t\t\t\tUser\t\t\t\t\tSubscription\tSubscription End\tEmail\t\n\n")
	}

	for _, system := range systems {
		if system.Subscription.SubscriptionPlan.Code != "trial" && period == "2w" {
			// skip not trial in 2 weeks
			continue
		}
		// print results as command output
		if !quietMode {
			until := fmt.Sprintf("%d-%02d-%02d", system.Subscription.ValidUntil.Year(), system.Subscription.ValidUntil.Month(), system.Subscription.ValidUntil.Day())
			fmt.Printf("%-40v%-40v%-16v%-15v\t\t%s\n", system.UUID, system.CreatorID, system.Subscription.SubscriptionPlan.Code, until, system.Notification["emails"])
		}

		// send mail to owners
		if sendMail {
			switch x := system.Notification["emails"].(type) {
			case []interface{}:
				for _, e := range x {
					notifications.MailNotificationExpire(e.(string), system, period)
				}
			default:
			}
		}

	}
}

func init() {
	expirationsCmd.Flags().StringVarP(&period, "period", "P", "1d", "Period of expiration")
	expirationsCmd.Flags().BoolVarP(&sendMail, "send-mail", "M", false, "Send expiration mail to system owners")
	expirationsCmd.Flags().BoolVarP(&quietMode, "quiet-mode", "Q", false, "Do not print results (Used with -M flag)")

	rootCmd.AddCommand(expirationsCmd)
}
