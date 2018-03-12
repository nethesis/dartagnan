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

func alertExists(SystemID int, AlertID string) (bool, models.Alert) {
	var alert models.Alert
	db := database.Database()
	db.Where("alert_id = ? AND system_id = ?", AlertID, SystemID).First(&alert)
	db.Close()

	if alert.ID == 0 {
		return false, models.Alert{}
	}

	return true, alert
}

func SetAlert(c *gin.Context) {
	var json models.AlertJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	// get system from uuid
	system := utils.GetSystemFromUUID(json.SystemID)

	// check if alert exists
	exists, alert := alertExists(system.ID, json.AlertID)
	if exists {
		if alert.Status == json.Status {
			// reject
			c.JSON(http.StatusOK, gin.H{"status": "no update"})
			return
		}

		if json.Status == "OK" {
			// handle resolve
			var oldStatus = alert.Status

			// add to history with RESOLVED
			alertHistory := models.AlertHistory{
				AlertID:    alert.AlertID,
				Priority:   alert.Priority,
				Resolution: "RESOLVED",
				StatusFrom: oldStatus,
				StatusTo:   json.Status,
				StartTime:  alert.Timestamp,
				EndTime:    time.Now().UTC(),
				SystemID:   system.ID,
			}

			// save to history
			db := database.Database()
			if err := db.Save(&alertHistory).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "alert not moved to history", "error": err.Error()})
				return
			}

			// delete current alert
			if err := db.Delete(&alert).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "alert not deleted", "error": err.Error()})
				return
			}

		} else {
			// handle change state
			var oldStatus = alert.Status
			alert.Status = json.Status

			// save alert
			db := database.Database()
			if err := db.Save(&alert).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "alert not updated", "error": err.Error()})
				return
			}

			// add to history with CHANGE_STATUS
			alertHistory := models.AlertHistory{
				AlertID:    alert.AlertID,
				Priority:   alert.Priority,
				Resolution: "CHANGE_STATUS",
				StatusFrom: oldStatus,
				StatusTo:   json.Status,
				StartTime:  alert.Timestamp,
				EndTime:    time.Now().UTC(),
				SystemID:   system.ID,
			}

			if err := db.Save(&alertHistory).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "alert not moved to history", "error": err.Error()})
				return
			}

			db.Close()
		}
	} else {
		if json.Status == "INIT" {
			// reject
			c.JSON(http.StatusOK, gin.H{"status": "no update"})
			return
		}

		if json.Status == "OK" {
			// reject
			c.JSON(http.StatusOK, gin.H{"status": "no update"})
			return
		}

		// create alert
		alert := models.Alert{
			AlertID:   json.AlertID,
			Priority:  utils.GetAlertPriority(json.AlertID),
			Note:      "",
			Status:    json.Status,
			Timestamp: time.Now().UTC(),
			SystemID:  system.ID,
		}

		// save alert
		db := database.Database()
		if err := db.Save(&alert).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "alert not saved", "error": err.Error()})
			return
		}

		db.Close()
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpdateAlertNote(c *gin.Context) {
	var alert models.Alert
	creatorID := c.MustGet("authUser").(string)
	alertID := c.Param("alert_id")

	var json models.AlertNoteJSON
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	if !utils.CheckSystemOwnership(json.SystemID, creatorID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "this systems is not yours!"})
		return
	}

	db := database.Database()
	db.Where("id = ? AND system_id = ?", alertID, json.SystemID).First(&alert)

	if alert.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no alert found!"})
		return
	}

	alert.Note = json.Note
	db.Save(&alert)
	db.Close()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetAlerts(c *gin.Context) {
	var alerts []models.Alert
	creatorID := c.MustGet("authUser").(string)
	systemID := c.Param("system_id")

	page := c.Query("page")
	limit := c.Query("limit")
	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Set("gorm:auto_preload", true).Preload("System", "creator_id = ?", creatorID).Where("system_id = ?", systemID).Offset(offsets[0]).Limit(offsets[1]).Find(&alerts)
	db.Close()

        for i, alert := range alerts {
                alerts[i].NameI18n = utils.GetAlertHumanName(alert.AlertID, "en-US")
        }

	c.JSON(http.StatusOK, alerts)
}

func GetAlertHistories(c *gin.Context) {
	var alertHistories []models.AlertHistory
	creatorID := c.MustGet("authUser").(string)
	systemID := c.Param("system_id")

	page := c.Query("page")
	limit := c.Query("limit")
	offsets := utils.OffsetCalc(page, limit)

	db := database.Database()
	db.Set("gorm:auto_preload", true).Preload("System", "creator_id = ?", creatorID).Where("system_id = ?", systemID).Offset(offsets[0]).Limit(offsets[1]).Find(&alertHistories)
	db.Close()

	c.JSON(http.StatusOK, alertHistories)
}
