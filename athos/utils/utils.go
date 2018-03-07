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
	"os"

	"hash/fnv"

	"github.com/satori/go.uuid"
	"github.com/mediocregopher/radix.v2/redis"

	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func GenerateUUID() string {
	u := uuid.Must(uuid.NewV4())
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

func GetSystemFromSecret(secret string) models.System {
	var system models.System
	db := database.Database()
	db.Where("secret = ?", secret).First(&system)
	db.Close()

	return system
}

func GetValidSystems() []models.System {
	var systems []models.System

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Joins("JOIN subscriptions ON systems.subscription_id = subscriptions.id").Where("valid_until > NOW()").Find(&systems)
	db.Close()

	return systems
}

func CalculateTier(uuid string, tiers uint32) uint32 {
        h := fnv.New32a()
        h.Write([]byte(uuid))
        tier := h.Sum32() % tiers

	return tier
}

func SetValidSystem(system models.System, client *redis.Client) (bool, string) {
	// Use 4 tiers
	tiers := uint32(4)

	now := time.Now()
	difference := system.Subscription.ValidUntil.Sub(now).Seconds()
	err := client.Cmd("HMSET", system.UUID, "tier_id", CalculateTier(system.UUID, tiers), "secret", system.Secret, "EX", int(difference)).Err
	if err != nil {
		return false, fmt.Sprintf("Can't save '%s' system inside Redis instance: %s", system.UUID, err)
	}
	// Set key expiration equal to Subscription.ValidUntil
	err = client.Cmd("EXPIRE", system.UUID, int(difference)).Err
	if err != nil {
		return false, fmt.Sprintf("Can't set expiratation on '%s' system inside Redis instance: %s", system.UUID, err)
	}
	return true, ""
}

func BulkSetValidSystems() (bool, []string) {
	var errors []string;
	systems := GetValidSystems()

	client, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", configuration.Config.RedisHost, configuration.Config.RedisPort))
	if err != nil {
		client.Close()
		return false, []string{fmt.Sprintf("Can't connect to Redis instance '%s:%s': %s", configuration.Config.RedisHost, configuration.Config.RedisPort, err)}
	}
	err = client.Cmd("FLUSHALL").Err
	if err != nil {
		client.Close()
		return false, []string{fmt.Sprintf("Can't FLUSH Redis instance: %s", err)}
	}

	for _, system := range systems {
		ret, err := SetValidSystem(system, client)
		if !ret {
			errors = append(errors, err)
		}
	}
	client.Close()

	if len(errors) > 0 {
		return false, errors
	} else {
		return true, errors
	}
}
