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
package utils

import (
	"fmt"

	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

// List expiration systems in a specific period
func ListExpirationSystems(period string) []models.System {
	var systems []models.System

	db := database.Instance()
	db.Preload("Subscription.SubscriptionPlan").Joins("JOIN subscriptions ON systems.subscription_id = subscriptions.id JOIN subscription_plans ON subscription_plans.id = subscriptions.subscription_plan_id").Where(fmt.Sprintf("date(subscriptions.valid_until)::date = date(date('now') + interval '%s')::date", period)).Find(&systems)

	return systems
}
