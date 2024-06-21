/*
 * Copyright (C) 2017 Nethesis S.r.l.
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
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nethesis/PayPal-Go-SDK"
	auth0 "github.com/nethesis/go-auth0"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
	"github.com/nethesis/dartagnan/athos/utils"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func GetSecret(c *gin.Context) string {
	/* Header format:
	Authorization: token <TOKEN>
	*/
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) > 1 {
		return authHeader[1]
	} else {
		return ""
	}
}

func AuthSecret(c *gin.Context) {
	secret := GetSecret(c)
	if secret == "" {
		respondWithError(http.StatusUnauthorized, "invalid Secret", c)
		return
	}
	if utils.GetSystemFromSecret(secret).ID != 0 {
		c.Next()
	} else {
		respondWithError(http.StatusUnauthorized, "invalid Secret", c)
		return
	}
}

func AuthJWT(c *gin.Context) {
	// define api endpoint and audience
	AUTH0_DOMAIN := "https://" + configuration.Config.Auth0.Domain + "/"
	JWKS_URI := "https://" + configuration.Config.Auth0.Domain + "/.well-known/jwks.json"
	AUDIENCE := []string{configuration.Config.Auth0.Audience}

	// create client configuration instance to check jwt
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: JWKS_URI})
	configuration := auth0.NewConfiguration(client, AUDIENCE, AUTH0_DOMAIN, jose.RS256)
	validator := auth0.NewValidator(configuration)

	// check jwt validation
	token, err := validator.ValidateRequest(c.Request)
	if err != nil {
		respondWithError(http.StatusUnauthorized, "missing or invalid token", c)
		return
	} else {
		// extract claims from token
		claims := map[string]interface{}{}
		err := validator.Claims(c.Request, token, &claims)
		if err != nil {
			respondWithError(http.StatusUnauthorized, "claims extraction failed", c)
			return
		}

		// set current user
		c.Set("authUser", claims["sub"])
		c.Next()
	}
}

func PaymentCheck(paymentID string, planCode string, uuid string) bool {
	var apiBase string
	if configuration.Config.PayPal.Sandbox {
		apiBase = paypalsdk.APIBaseSandBox
	} else {
		apiBase = paypalsdk.APIBaseLive
	}
	c, errSDK := paypalsdk.NewClient(configuration.Config.PayPal.ClientID, configuration.Config.PayPal.ClientSecret, apiBase)
	if errSDK != nil {
		fmt.Println(errSDK.Error())
	}
	_, err := c.GetAccessToken()

	payment, err := c.GetPayment(paymentID)
	if err != nil {
		fmt.Println(err.Error())
	}

	SavePaymentDetails(paymentID, uuid)

	if payment.State == "approved" {
		if payment.Transactions[0].ItemList.Items[0].Name == planCode && payment.Transactions[0].ItemList.Items[0].SKU == uuid {
			return true
		}
		return false
	}
	return false
}

func SavePaymentDetails(paymentID string, systemUUID string) {
	var system models.System

	db := database.Instance()
	db.Set("gorm:auto_preload", false).Where("uuid = ?", systemUUID).First(&system)
	payment := models.Payment{CreatorID: system.CreatorID, Payment: paymentID, SystemID: system.ID, Created: time.Now().UTC()}
	db.Create(&payment)
}
