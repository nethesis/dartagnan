/*
 * Copyright (C) 2018 Nethesis S.r.l.
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
 */

package methods

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

//TODO: this probably should be rewritten after legal advice
func GetTaxPercentage(country string, vat string) int {
	var tax models.Tax
	if vat != "" {
		return 0
	}

	db := database.Database()
	db.Where("country = ?", country).First(&tax)
	defer db.Close()
	return tax.Percentage
}

func GetBilling(c *gin.Context) {
	var billing models.Billing
	creatorID := c.MustGet("authUser").(string)

	db := database.Database()
	db.Where("creator_id = ?", creatorID).First(&billing)
	defer db.Close()

	if billing.ID == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "no billing information found!"})
		return
	}

	billing.Tax = GetTaxPercentage(billing.Country, billing.Vat)

	c.JSON(http.StatusOK, billing)
}

func CreateBilling(c *gin.Context) {
	var json models.BillingJSON
	creatorID := c.MustGet("authUser").(string)

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	billing := models.Billing{
		CreatorID:    creatorID,
		Name:         json.Name,
		Address:      json.Address,
		Country:      json.Country,
		City:         json.City,
		PostalCode:   json.PostalCode,
		Vat:          json.Vat,
	}

	db := database.Database()
	defer db.Close()
	if err := db.Create(&billing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "billing not saved", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UpdateBilling(c *gin.Context) {
	var json models.BillingJSON
	var billing models.Billing
	creatorID := c.MustGet("authUser").(string)

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "request fields malformed", "error": err.Error()})
		return
	}

	db := database.Database()
	db.Where("creator_id = ?", creatorID).First(&billing)
	defer db.Close()

	if billing.ID == 0 {
		db.Close()
		c.JSON(http.StatusNotFound, gin.H{"message": "no billing found!"})
		return
	}

	billing.CreatorID = creatorID
	billing.Name = json.Name
	billing.Address = json.Address
	billing.Country = json.Country
	billing.City = json.City
	billing.PostalCode = json.PostalCode
	billing.Vat = json.Vat

	if err := db.Save(&billing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "billing not saved", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

