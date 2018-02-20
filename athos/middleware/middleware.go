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

package middleware

import (
	"fmt"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"github.com/logpacker/PayPal-Go-SDK"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/nethesis/dartagnan/athos/configuration"
)

func respondWithError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"message": message})
	c.Abort()
}

func AuthUUID(c *gin.Context) {
	// check UUID
	if true {
		c.Next()
	} else {
		respondWithError(http.StatusUnauthorized, "invalid UUID", c)
		return
	}
}

func AuthJWT(c *gin.Context) {
	// define api endpoint and audience
	AUTH0_DOMAIN := "https://" + configuration.Config.Auth0.Domain + "/"
	JWKS_URI := "https://" + configuration.Config.Auth0.Domain + "/.well-known/jwks.json"
	AUDIENCE := []string{configuration.Config.Auth0.IdentifierAPI}

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
	c, errSDK := paypalsdk.NewClient(configuration.Config.PayPal.ClientID, configuration.Config.PayPal.ClientSecret, paypalsdk.APIBaseSandBox)
	if errSDK != nil {
		fmt.Println(errSDK.Error())
	}
	_, err := c.GetAccessToken()

	payment, err := c.GetPayment(paymentID)
	if err != nil {
		fmt.Println(err.Error())
	}

	if payment.State == "approved" {
		if payment.Transactions[0].ItemList.Items[0].Name == planCode && payment.Transactions[0].ItemList.Items[0].SKU == uuid {
			return true
		}
		return false
	}
	return false
}
