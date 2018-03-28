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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List objects from database for reporting",
	Long:  "\nAvailable objects: 'trials', 'payments' ",
	Args:  cobra.ExactArgs(1),
	Run:   list,
}

func printActiveTrials() {
	systems := utils.ListActiveTrials()
	fmt.Printf("System\tUser\tSubscription_end\tEmail\t\n\n")
	for _, system := range systems {
		until := fmt.Sprintf("%d-%02d-%02d", system.Subscription.ValidUntil.Year(), system.Subscription.ValidUntil.Month(), system.Subscription.ValidUntil.Day())
		fmt.Printf("%-40v%-40v%-15v%s\n", system.UUID, system.CreatorID, until, system.Notification["emails"])
	}
}

func printTodayPayments() {
	var tax float64
	paymentList := utils.ListPayments(168)
	for _, payment := range paymentList {
		details := utils.GetPaypalPayment(payment.Payment)
		info := utils.GetBillingInfo(payment.CreatorID)

		fmt.Printf("===========================\n")
		fmt.Printf("Payment ID: %s\n", payment.Payment)
		fmt.Printf("Date (YYYY-MM-DD): %d-%02d-%02d\n", payment.Created.Year(), payment.Created.Month(), payment.Created.Day())
		fmt.Printf("Name: %s\n", info.Name)
		if info.Vat != "" {
			fmt.Printf("Vat: %s\n", info.Vat)
		}
		if info.Tax > 0 {
			fmt.Printf("Tax: %d\n", info.Tax)
		}
		fmt.Printf("Address: %s, City: %s, PostalCode: %s, Country: %s\n", info.Address, info.City, info.PostalCode, info.Country)

		if details.Tax > 0 {
			tax = details.Tax
		} else {
			tax = 0
		}
		fmt.Printf("Total: %.2f, Subtotal: %.2f, Tax: %.2f, Currency: %s, Description: %s\n", details.Total, details.Subtotal, tax, details.Currency, details.ItemDescription)
		fmt.Printf("===========================\n\n")
	}
}

func list(cmd *cobra.Command, args []string) {
	switch args[0] {
	case "trials":
		printActiveTrials()
	case "payments":
		printTodayPayments()
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
