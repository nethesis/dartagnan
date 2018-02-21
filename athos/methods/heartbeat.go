/*
 * Copyright (C) 2017 Nethesis S.r.l.
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

func heartbeatExists(SystemID int) (bool, models.Heartbeat) {
	var heartbeat models.Heartbeat
	db := database.Database()
	db.Where("system_id = ?", SystemID).First(&heartbeat)
	db.Close()

	if heartbeat.ID == 0 {
		return false, models.Heartbeat{}
	}

	return true, heartbeat
}

func SetHeartbeat(c *gin.Context) {
	var json models.HeartbeatJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	// get system from uuid
	system := utils.GetSystemFromUUID(json.SystemID)

	// check if heartbeat exists
	exists, heartbeat := heartbeatExists(system.ID)
	if exists {
		// update current heartbeat
		heartbeat.Timestamp = time.Now().UTC()

		// save current heartbeat
		db := database.Database()
		if err := db.Save(&heartbeat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "heartbeat not updated", "error": err.Error()})
			return
		}

		db.Close()
	} else {
		// create heartbeat
		heartbeat := models.Heartbeat{
			Timestamp: time.Now().UTC(),
			SystemID:  system.ID,
		}

		// save new heartbeat
		db := database.Database()
		if err := db.Save(&heartbeat).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "heartbeat not saved", "error": err.Error()})
			return
		}
		db.Close()
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetHeartbeat(c *gin.Context) {
	var heartbeat models.Heartbeat
	creatorID := c.MustGet("authUser").(string)
	systemID := c.Param("system_id")

	db := database.Database()
	db.Set("gorm:auto_preload", true).Preload("System", "creator_id = ?", creatorID).Where("system_id = ?", systemID).First(&heartbeat)
	db.Close()

	if heartbeat.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no heartbeat found!"})
		return
	}

	c.JSON(http.StatusOK, heartbeat)
}
