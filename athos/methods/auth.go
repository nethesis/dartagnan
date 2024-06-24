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
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func BasicAuth(c *gin.Context) {
	// define model
	var system models.System

	// get http basic credentials
	uuid, secret, _ := c.Request.BasicAuth()

	// init db instance
	db := database.Instance()
	db.Where("uuid = ? AND secret = ?", uuid, secret).First(&system)

	// check if system exists
	if system.ID == 0 {
		// response 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "basic auth failed. system not found",
			"data":    nil,
		})
		return
	}

	// response 200
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "basic auth ok",
		"data":    nil,
	})
}

func BasicAuthService(c *gin.Context) {
	// response 200
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "basic auth for service ok",
		"data":    nil,
	})
}
