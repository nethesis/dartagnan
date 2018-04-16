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
	"time"

	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

// List of trials active in the last 24 hours
func ListActiveTrials(hours int) []models.System {
	var systems []models.System
	var ret []models.System

	db := database.Database()
	db.Set("gorm:auto_preload", true)
	db.Preload("Subscription").Joins("JOIN heartbeats ON heartbeats.system_id = systems.id").Order("creator_id desc").Find(&systems)
	db.Set("gorm:auto_preload", false)

	for _, system := range systems {
		if system.Subscription.SubscriptionPlanID == 1 {
			var h models.Heartbeat
			db.Where("system_id = ?", system.ID).First(&h)
			duration := time.Since(h.Timestamp)
			if duration.Hours() < float64(hours) {
				ret = append(ret, system)
			}
		}
	}

	return ret
}
