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
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/dartagnan/athos/cache"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/middleware"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
)

func CreateSystem(c *gin.Context) {
	creatorID := c.MustGet("authUser").(string)

	var json models.SystemJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	// get subscription plan
	subscriptionPlan := utils.GetSubscriptionPlanByCode("trial")

	// create system
	uuid := utils.GenerateUUID()
	secret := utils.GenerateSecret(uuid)
	system := models.System{
		CreatorID:    creatorID,
		UUID:         uuid,
		Secret:       secret,
		Tags:         "trial",
		PublicIP:     "",
		Status:       "active",
		Notification: json.Notification,
		Created:      time.Now().UTC(),
		Subscription: models.Subscription{
			UserID:             creatorID,
			ValidFrom:          time.Now().UTC(),
			ValidUntil:         time.Now().UTC().AddDate(0, 0, subscriptionPlan.Period),
			Status:             "valid",
			Created:            time.Now().UTC(),
			SubscriptionPlanID: subscriptionPlan.ID,
		},
	}

	// save new system
	db := database.Instance()
	db.Create(&system)
	if err := db.Save(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not saved", "error": err.Error()})
		return
	}

	if res := cache.SetValidSystem(system); !res {
		// Soft fail, chache can be restored later
		fmt.Println("[ERROR]: can't save %s inside the cache", system.UUID)
	}

	if system.ID == 0 {
		c.JSON(http.StatusConflict, gin.H{"status": "system not added"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"uuid": system.UUID, "secret": system.Secret, "status": "success"})
	}
}

func UpdateSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	var json models.SystemJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found!"})
		return
	}

	if len(json.Tags) > 0 {
		system.Tags = json.Tags
	}

	if err := db.Save(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not updated", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getStatus(id int) string {
	var heartbeat models.Heartbeat
	db := database.Instance()
	db.Where("system_id = ?", id).First(&heartbeat)

	if heartbeat.ID == 0 {
		return "no_comm"
	}

	sanity := heartbeat.Timestamp.Add(time.Minute * 30)
	if time.Now().After(sanity) {
		return "no_active"
	} else {
		return "active"
	}
}

func getAlertsNumber(system models.System) int {
	if !utils.CanAccessAlerts(system.Subscription.SubscriptionPlan) {
		return -1
	}
	type Result struct {
		Count int
	}

	var result Result
	db := database.Instance()
	db.Raw("SELECT COUNT(*) as count FROM alerts WHERE system_id = ?", system.ID).Scan(&result)

	return result.Count
}

func GetSystems(c *gin.Context) {
	var systems []models.System
	creatorID := c.MustGet("authUser").(string)

	page := c.Query("page")
	limit := c.Query("limit")
	offsets := utils.OffsetCalc(page, limit)

	db := database.Instance()
	db.Select("systems.*, inventories.data->'networking'->>'fqdn' AS hostname").Preload("Subscription.SubscriptionPlan").Joins("LEFT JOIN inventories ON systems.id = inventories.system_id").Where("creator_id = ?", creatorID).Offset(offsets[0]).Limit(offsets[1]).Find(&systems)

	if len(systems) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no systems found!"})
		return
	}

	for i, system := range systems {
		systems[i].Status = getStatus(system.ID)
		systems[i].Alerts = getAlertsNumber(system)
	}

	c.JSON(http.StatusOK, systems)
}

func GetSystemBySecret(c *gin.Context) {
	var system models.System
	sentSecret := middleware.GetSecret(c)

	db := database.Instance()
	db.Where("secret = ?", sentSecret).First(&system)

	db.Preload("Subscription.SubscriptionPlan").Where("id = ? ", system.ID).First(&system)

	system.Status = getStatus(system.ID)
	system.Alerts = getAlertsNumber(system)
	system.Secret = ""
	c.JSON(http.StatusOK, system)
}

func GetSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	db := database.Instance()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found!"})
		return
	}

	system.Status = getStatus(system.ID)
	system.Alerts = getAlertsNumber(system)
	c.JSON(http.StatusOK, system)
}

func DeleteSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	db := database.Instance()
	db.Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found!"})
		return
	}

	if err := db.Delete(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not deleted", "error": err.Error()})
		return
	}
	if res := cache.DeleteValidSystem(system); !res {
		// Soft fail, chache can be restored later
		fmt.Println("[ERROR]: can't delete %s from cache", system.UUID)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func RenewalPlan(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	var json models.SubscriptionRenewalJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Preload("Subscription").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found!"})
		return
	}

	// check payment
	if middleware.PaymentCheck(json.PaymentID, system.Subscription.SubscriptionPlan.Code, system.UUID) {
		// update subscription
		system.Subscription.ValidFrom = time.Now().UTC()
		system.Subscription.ValidUntil = system.Subscription.ValidUntil.AddDate(0, 0, system.Subscription.SubscriptionPlan.Period)
		system.Subscription.Status = "valid"

		// update system info
		if err := db.Save(&system).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "system subscription plan not renewed", "error": err.Error()})
			return
		}
		if res := cache.SetValidSystem(system); !res {
			// Soft fail, chache can be restored later
			fmt.Println("[ERROR]: can't save %s inside the cache", system.UUID)
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "no payment related to this plan for this server found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpgradePlanPrice(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")
	plan := c.Query("plan")

	newSubuscriptionPlan := utils.GetSubscriptionPlanByCode(plan)

	db := database.Instance()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	// calculate discount upgrade
	daysDiff := system.Subscription.ValidUntil.Sub(time.Now().UTC())
	discount := (daysDiff.Hours() / 24) * system.Subscription.SubscriptionPlan.Price / float64(system.Subscription.SubscriptionPlan.Period)
	dicscountPerc := discount * 100 / newSubuscriptionPlan.Price
	finalPrice := newSubuscriptionPlan.Price - discount

	c.JSON(http.StatusOK, gin.H{"discount": dicscountPerc, "full_price": newSubuscriptionPlan.Price, "price": utils.Round(finalPrice, 0.5, 2), "name": newSubuscriptionPlan.Code})
}

func UpgradePlan(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	var json models.SubscriptionUpgradeJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Instance()
	db.Preload("Subscription").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found!"})
		return
	}

	if json.SubscriptionPlanID != 0 && json.SubscriptionPlanID != system.Subscription.SubscriptionPlanID {
		// get subscription using id
		newSubscriptionPlan := utils.GetSubscriptionPlanById(json.SubscriptionPlanID)

		// check payment
		if middleware.PaymentCheck(json.PaymentID, newSubscriptionPlan.Code, system.UUID) {
			// update subscription
			system.Subscription.SubscriptionPlanID = newSubscriptionPlan.ID
			system.Subscription.ValidFrom = time.Now().UTC()
			system.Subscription.ValidUntil = time.Now().UTC().AddDate(0, 0, newSubscriptionPlan.Period)
			system.Subscription.Status = "valid"

			// update system info
			if err := db.Save(&system).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "system subscription plan not updated", "error": err.Error()})
				return
			}
			if res := cache.SetValidSystem(system); !res {
				// Soft fail, chache can be restored later
				fmt.Println("[ERROR]: can't save %s inside the cache", system.UUID)
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "no payment related to this plan for this server found"})
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{"message": "this plan is already associated with this server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
