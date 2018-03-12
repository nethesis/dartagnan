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
	"encoding/json"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
)


type Notification map[string][]string

func AlertNotification(alert models.Alert, isNew bool) {
	if (models.System{}) == alert.System {
		alert.System = utils.GetSystemById(alert.SystemID)
	}
	plan := alert.System.Subscription.SubscriptionPlan
	if plan.Code == "fiorentina" || plan.Code == "pizza" {
		// find all configured email addresses for this server
		var t Notification
		j := alert.System.Notification
		json.Unmarshal([]byte(j), &t)
		for _, email := range t["emails"] {
			MailNotification(email, alert, isNew)
		}
	}
}
