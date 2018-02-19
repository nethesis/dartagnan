/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Icaro project.
 *
 * Icaro is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Icaro is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Icaro.  If not, see COPYING.
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
	subuscriptionPlan := utils.GetSubscriptionPlanByCode("trial")

	// create system
	system := models.System{
		CreatorID:   creatorID,
		UUID:        json.UUID,
		Hostname:    json.Hostname,
		Description: json.Description,
		PublicIP:    json.PublicIP,
		Status:      "active",
		Created:     time.Now().UTC(),
		Subscription: models.Subscription{
			UserID:             creatorID,
			ValidFrom:          time.Now().UTC(),
			ValidUntil:         time.Now().UTC().AddDate(0, 0, subuscriptionPlan.Period),
			Status:             "valid",
			Created:            time.Now().UTC(),
			SubscriptionPlanID: subuscriptionPlan.ID,
		},
	}

	// save new system
	db := database.Database()
	db.Create(&system)
	if err := db.Save(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not saved", "error": err.Error()})
		return
	}
	db.Close()

	if system.ID == 0 {
		c.JSON(http.StatusConflict, gin.H{"id": system.ID, "status": "system not added"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"id": system.ID, "status": "success"})
	}
}

func UpdateSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	var json models.SystemJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No system found!"})
		return
	}

	if len(json.Hostname) > 0 {
		system.Hostname = json.Hostname
	}
	if len(json.Description) > 0 {
		system.Description = json.Description
	}
	if len(json.PublicIP) > 0 {
		system.PublicIP = json.PublicIP
	}

	if json.SubscriptionPlanID != 0 && json.SubscriptionPlanID != system.Subscription.SubscriptionPlanID {
		// update subscription
		newSubscriptionPlan := utils.GetSubscriptionPlanById(json.SubscriptionPlanID)
		system.Subscription.SubscriptionPlanID = newSubscriptionPlan.ID
		system.Subscription.ValidFrom = time.Now().UTC()
		system.Subscription.ValidUntil = time.Now().UTC().AddDate(0, 0, newSubscriptionPlan.Period)
		system.Subscription.Status = "valid"
	}

	if err := db.Save(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not updated", "error": err.Error()})
		return
	}
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetSystems(c *gin.Context) {
	var systems []models.System
	creatorID := c.MustGet("authUser").(string)

	page := c.Query("page")
	limit := c.Query("limit")
	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Where("creator_id = ?", creatorID).Offset(offsets[0]).Limit(offsets[1]).Find(&systems)
	db.Close()

	if len(systems) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No systems found!"})
		return
	}

	c.JSON(http.StatusOK, systems)
}

func GetSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)
	db.Close()

	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No system found!"})
		return
	}

	c.JSON(http.StatusOK, system)
}

func DeleteSystem(c *gin.Context) {
	var system models.System
	creatorID := c.MustGet("authUser").(string)

	systemID := c.Param("system_id")

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)

	if system.ID == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "No system found!"})
		return
	}

	if err := db.Delete(&system).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "system not deleted", "error": err.Error()})
		return
	}
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
