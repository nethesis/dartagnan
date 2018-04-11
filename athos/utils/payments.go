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
 */

package utils

import (
	"fmt"
	"strconv"

	"github.com/nethesis/dartagnan/athos/database"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/models"
)

func parseFloatOrNegative(number string) float64 {
	if n, err := strconv.ParseFloat(number, 64); err == nil {
		return n
	}

	return -1
}

func GetPaypalPayment(paymentId string) models.PaypalPayment {
	var apiBase string
	var ret models.PaypalPayment
	if configuration.Config.PayPal.Sandbox {
		apiBase = paypalsdk.APIBaseSandBox
	} else {
		apiBase = paypalsdk.APIBaseLive
	}
	c, errSDK := paypalsdk.NewClient(configuration.Config.PayPal.ClientID, configuration.Config.PayPal.ClientSecret, apiBase)
	if errSDK != nil {
		fmt.Println(errSDK.Error())
	}
	_, err := c.GetAccessToken()

	payment, err := c.GetPayment(paymentId)
	if err != nil {
		fmt.Println(err.Error())
	}

	// our payments have always one transaction with one item
	t := payment.Transactions[0]
	i := t.ItemList.Items[0]

	ret.Item = i.SKU
	ret.ItemDescription = fmt.Sprintf("%s %s", i.Description, i.SKU)
	ret.Currency = t.Amount.Currency
	ret.Tax = parseFloatOrNegative(t.Amount.Details.Tax)
	ret.Total = parseFloatOrNegative(t.Amount.Total)
	ret.Subtotal = parseFloatOrNegative(t.Amount.Details.Subtotal)

	return ret
}

func ListPayments(sinceHours int) []models.Payment {
	var payments []models.Payment

	db := database.Database()
	db.Set("gorm:auto_preload", true)
	db.Where(fmt.Sprintf("(now() - created)::interval < '%d hour'::interval", sinceHours)).Find(&payments)

	return payments
}

func GetBillingInfo(user string) models.Billing {
	var info models.Billing
	db := database.Database()
	db.Where("creator_id = ?", user).First(&info)
	db.Close()

	return info
}
