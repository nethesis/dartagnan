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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
)

func BasicAuth(c *gin.Context) {
	// define model
	var system models.System

	// get http basic credentials
	uuid, secret, _ := c.Request.BasicAuth()

	// init db instance
	db := database.Instance()
	db.Preload("Subscription.SubscriptionPlan").Where("uuid = ? AND secret = ?", uuid, secret).First(&system)

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

	// check if system subscription is expired
	if time.Now().After(system.Subscription.ValidUntil) {
		// response 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "basic auth failed. system subscription is expired",
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
	// define model
	var system models.System

	// get http basic credentials
	uuid, secret, _ := c.Request.BasicAuth()

	// get service name
	service := c.Param("service")

	// init db instance
	db := database.Instance()
	db.Preload("Subscription.SubscriptionPlan").Where("uuid = ? AND secret = ?", uuid, secret).First(&system)

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

	// check if system subscription is expired
	if time.Now().After(system.Subscription.ValidUntil) {
		// response 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "basic auth failed. system subscription is expired",
			"data":    nil,
		})
		return
	}

	// extract service from subscription plan
	servicesParts := strings.Split(system.Subscription.SubscriptionPlan.Code, "+")

	// check if there are services
	if len(servicesParts) <= 1 {
		// response 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "basic auth failed. service not found in system subscription",
			"data":    nil,
		})
		return
	}

	// get single services
	services := strings.Split(servicesParts[1], ",")

	// check if service is in subscription plan
	if !utils.ContainsS(services, service) {
		// response 401
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "basic auth failed. service not found in system subscription",
			"data":    nil,
		})
		return
	}

	// response 200
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "basic auth for service ok",
		"data":    nil,
	})
}
