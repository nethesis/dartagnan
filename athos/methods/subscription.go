/*
 * Copyright (C) 2017 Nethesis S.r.l.
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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package methods

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func GetSubscriptionPlans(c *gin.Context) {
	var subscriptionPlans []models.SubscriptionPlan

	db := database.Instance()
	db.Find(&subscriptionPlans)

	if len(subscriptionPlans) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no subscription plans found!"})
		return
	}

	c.JSON(http.StatusOK, subscriptionPlans)
}

func VolumeDiscountPrice(c *gin.Context) {
	var systems []models.System
	count := 0
	discount := 0

	creatorID := c.MustGet("authUser").(string)

	db := database.Instance()

	// count servers with payed subscriptions
	db.Preload("Subscription.SubscriptionPlan").Where("creator_id = ?", creatorID).Find(&systems)
	for _, system := range systems {
		if system.Subscription.SubscriptionPlanID > 1 && system.Subscription.ValidUntil.After(time.Now().UTC()) {
			count++
		}
	}

	// calculate volume discount
	if count > 5 && count <= 10 {
		discount = 15
	}

	if count > 10 && count <= 20 {
		discount = 20
	}

	if count > 20 {
		discount = 25
	}

	c.JSON(http.StatusOK, gin.H{"discount": discount, "count": count})
}
