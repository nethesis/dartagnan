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

package utils

import (
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/satori/go.uuid"

	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func GenerateUUID() string {
	u := uuid.NewV4()
	return u.String()
}

func GenerateSecret(uuid string) string {
	h := sha256.New()
	h.Write([]byte(time.Now().UTC().String() + uuid))
	secret := fmt.Sprintf("%x", h.Sum(nil))
	return secret
}

func OffsetCalc(page string, limit string) [2]int {
	var resLimit = 0
	var resOffset = 0

	if len(page) == 0 {
		page = "1"
	}
	if len(limit) == 0 {
		limit = "-1"
	}

	limitInt, errLimit := strconv.Atoi(limit)
	if errLimit != nil {
		fmt.Println(errLimit.Error())
	}

	pageInt, errPage := strconv.Atoi(page)
	if errPage != nil {
		fmt.Println(errPage.Error())
	}

	resLimit = limitInt
	resOffset = (pageInt - 1) * resLimit

	result := [2]int{resOffset, resLimit}
	return result
}

func Round(val float64, roundOn float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	var round float64
	if val > 0 {
		if div >= roundOn {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if div >= roundOn {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}
	return round / pow
}

func Contains(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func GetAlertPriority(alertID string) string {
	return "HIGH"
}

func GetSubscriptionPlanByCode(code string) models.SubscriptionPlan {
	var subscriptionPlan models.SubscriptionPlan
	db := database.Database()
	db.Where("code = ?", code).First(&subscriptionPlan)
	db.Close()

	return subscriptionPlan
}

func GetSubscriptionPlanById(id int) models.SubscriptionPlan {
	var subscriptionPlan models.SubscriptionPlan
	db := database.Database()
	db.Where("id = ?", id).First(&subscriptionPlan)
	db.Close()

	return subscriptionPlan
}

func GetSubscription(id int) models.Subscription {
	var subscription models.Subscription
	db := database.Database()
	db.Where("id = ?", id).First(&subscription)
	db.Close()

	return subscription
}

func GetSystemFromUUID(uuid string) models.System {
	var system models.System
	db := database.Database()
	db.Where("uuid = ?", uuid).First(&system)
	db.Close()

	return system
}
