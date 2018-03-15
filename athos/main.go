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
	if configuration.Config.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = configuration.Config.Cors.Origins
	corsConf.AllowHeaders = configuration.Config.Cors.Headers
	corsConf.AllowMethods = configuration.Config.Cors.Methods
	router.Use(cors.New(corsConf))

	api := router.Group("/api")

	// protect API using SystemID middleware
	machine := api.Group("/machine")
	machine.Use(middleware.AuthSecret)
	{
		heartbeats := machine.Group("/heartbeats")
		{
			heartbeats.POST("/store", methods.SetHeartbeat)
		}
		alerts := machine.Group("/alerts")
		{
			alerts.POST("/store", methods.SetAlert)
		}
		inventories := machine.Group("/inventories")
		{
			inventories.POST("/store", methods.SetInventory)
		}
		info := machine.Group("/info")
		{
			info.GET("/", methods.GetSystemBySecret)
		}
	}

	// protect API using JWT middleware
	ui := api.Group("/ui")
	ui.Use(middleware.AuthJWT)
	{
		heartbeats := ui.Group("/heartbeats")
		{
			heartbeats.GET("/:system_id", methods.GetHeartbeat)
		}
		alerts := ui.Group("/alerts")
		{
			alerts.GET("", methods.GetAllAlerts)
			alerts.GET("/:system_id", methods.GetAlerts)
			alerts.GET("/:system_id/histories", methods.GetAlertHistories)
			alerts.PUT("/:alert_id", methods.UpdateAlertNote)
		}
		inventories := ui.Group("/inventories")
		{
			inventories.GET("/:system_id", methods.GetInventory)
			inventories.GET("/:system_id/histories", methods.GetInventoryHistories)
		}

		systems := ui.Group("/systems")
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

		plans := ui.Group("/plans")
		{
			plans.GET("", methods.GetSubscriptionPlans)
		}

		billings := ui.Group("/billings")
		{
			billings.GET("", methods.GetBilling)
			billings.POST("", methods.CreateBilling)
			billings.PUT("", methods.UpdateBilling)
		}

		taxes := ui.Group("/taxes")
		{
			taxes.GET("", methods.GetTaxes)
		}

		utils := ui.Group("/utils")
		{
			utils.GET("/reverse_lookup/:ip", methods.ReverseLookup)
		}
	}
	// handle missing endpoint
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "API not found"})
	})

	router.Run()

}
