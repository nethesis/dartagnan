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

package main

import (
	"flag"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/methods"
	"github.com/nethesis/dartagnan/athos/middleware"
)

func main() {
	// read and init configuration
	ConfigFilePtr := flag.String("c", "/opt/dartagnan/athos/conf.json", "Path to configuration file")
	flag.Parse()
	configuration.Init(ConfigFilePtr)

	// init routers
	router := gin.Default()

	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = configuration.Config.Cors.Origins
	corsConf.AllowHeaders = configuration.Config.Cors.Headers
	corsConf.AllowMethods = configuration.Config.Cors.Methods
	router.Use(cors.New(corsConf))

	api := router.Group("/api")

	// protect API using SystemID middleware
	heartbeats := api.Group("/heartbeats").Use(middleware.AuthSystemID)
	{
		heartbeats.POST("/store", methods.SetHeartbeat)
	}
	alerts := api.Group("/alerts").Use(middleware.AuthSystemID)
	{
		alerts.POST("/store", methods.SetAlert)
	}
	inventories := api.Group("/inventories").Use(middleware.AuthSystemID)
	{
		inventories.POST("/store", methods.SetInventory)
	}

	// protect API using JWT middleware
	api.Use(middleware.AuthJWT)
	{
		heartbeats := api.Group("/heartbeats")
		{
			heartbeats.GET("/:system_id", methods.GetHeartbeat)
		}
		alerts := api.Group("/alerts")
		{
			alerts.GET("/:system_id", methods.GetAlerts)
			alerts.GET("/:system_id/histories", methods.GetAlertHistories)
		}
		inventories := api.Group("/inventories")
		{
			inventories.GET("/:system_id", methods.GetInventory)
			inventories.GET("/:system_id/histories", methods.GetInventoryHistories)
		}

		systems := api.Group("/systems")
		{
			systems.GET("", methods.GetSystems)
			systems.GET("/:system_id", methods.GetSystem)
			systems.POST("", methods.CreateSystem)
			systems.PUT("/:system_id", methods.UpdateSystem)
			systems.DELETE("/:system_id", methods.DeleteSystem)

			systems.POST("/:system_id/renewal", methods.RenewalPlan)

			systems.GET("/:system_id/upgrade_price", methods.UpgradePlanPrice)
			systems.POST("/:system_id/upgrade", methods.UpgradePlan)
		}

		plans := api.Group("/plans")
		{
			plans.GET("", methods.GetSubscriptionPlans)
		}
	}

	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})

	router.Run()

}
