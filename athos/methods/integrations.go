/*
 * Copyright (C) 2024 Nethesis S.r.l.
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

package methods

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func CreateIntegration(c *gin.Context) {
	// define system model
	var system models.System

	// get integration
	integration := c.Param("integration")

	// get logged user
	creatorID := c.MustGet("authUser").(string)

	// search system
	db := database.Instance()
	db.Where("creator_id = ?", creatorID).First(&system)

	// check if there is at least one server
	if system.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no system found for this user"})
		return
	}

	// get email from system
	var email string
	switch x := system.Notification["emails"].(type) {
	case []interface{}:
		for _, e := range x {
			email = e.(string)
		}
	default:
	}

	// exec integration script
	out, err := exec.Command("/opt/dartagnan/"+integration+".sh", "/opt/dartagnan/config.json", email).Output()

	// parse json output
	var result map[string]interface{}
	_ = json.Unmarshal(out, &result)

	// check errors
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": integration + " integration not created",
			"data":    result,
		})
		return
	}

	// response 201
	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": integration + " integration created",
		"data":    result,
	})
}
