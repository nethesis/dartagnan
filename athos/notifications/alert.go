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

package notifications

import (
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
)


func AlertNotification(alert models.Alert, isNew bool) {
	if alert.System.Subscription.SubscriptionPlan.Code == "" {
		alert.System = utils.GetSystemById(alert.SystemID)
	}
	plan := alert.System.Subscription.SubscriptionPlan
	if plan.Code == "fiorentina" || plan.Code == "pizza" {
		switch x := alert.System.Notification["emails"].(type) {
		case []interface{}:
			for _, e := range x {
				MailNotification(e.(string), alert, isNew)
			}
		default:
		}
	}
}
