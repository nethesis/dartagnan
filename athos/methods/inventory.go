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

func inventoryExists(SystemID int) (bool, models.Inventory) {
	var inventory models.Inventory
	db := database.Database()
	db.Where("system_id = ?", SystemID).First(&inventory)
	db.Close()

	if inventory.ID == 0 {
		return false, models.Inventory{}
	}

	return true, inventory
}

func SetInventory(c *gin.Context) {
	var json models.InventoryJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	// get system from uuid
	system := utils.GetSystemFromUUID(json.SystemUUID)

	// check if inventory exists
	exists, inventory := inventoryExists(system.ID)
	if exists {
		// add to history
		inventoryHistory := models.InventoryHistory{
			Data:      inventory.Data,
			Timestamp: time.Now().UTC(),
			SystemID:  inventory.SystemID,
		}

		// save to inventory history
		db := database.Database()
		if err := db.Save(&inventoryHistory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "inventory not saved in history", "error": err.Error()})
			return
		}

		// update current inventory
		inventory.Data = json.Data
		inventory.Timestamp = time.Now().UTC()

		// save current inventory
		if err := db.Save(&inventory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "inventory not updated", "error": err.Error()})
			return
		}

		db.Close()
	} else {
		// create inventory
		inventory := models.Inventory{
			Data:      json.Data,
			Timestamp: time.Now().UTC(),
			SystemID:  system.ID,
		}

		// save new inventory
		db := database.Database()
		if err := db.Save(&inventory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "inventory not saved", "error": err.Error()})
			return
		}
		db.Close()
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetInventory(c *gin.Context) {
	var inventory models.Inventory
	creatorID := c.MustGet("authUser").(string)
	systemID := c.Param("system_id")

	db := database.Database()
	db.Set("gorm:auto_preload", true).Preload("System", "creator_id = ?", creatorID).Where("system_id = ?", systemID).First(&inventory)
	db.Close()

	if inventory.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no inventory found!"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func GetInventoryHistories(c *gin.Context) {
	var inventoryHistories []models.InventoryHistory
	creatorID := c.MustGet("authUser").(string)
	systemID := c.Param("system_id")

	page := c.Query("page")
	limit := c.Query("limit")
	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Set("gorm:auto_preload", true).Preload("System", "creator_id = ?", creatorID).Where("system_id = ?", systemID).Offset(offsets[0]).Limit(offsets[1]).Find(&inventoryHistories)
	db.Close()

	if len(inventoryHistories) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no inventory histories found!"})
		return
	}

	c.JSON(http.StatusOK, inventoryHistories)
}
