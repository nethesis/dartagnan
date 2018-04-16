/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
 *
 * WindMill is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * WindMill is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with WindMill.  If not, see COPYING.
 *
 */

package cache

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/mediocregopher/radix.v2/redis"

	"github.com/nethesis/dartagnan/athos/configuration"
	"github.com/nethesis/dartagnan/athos/database"
	"github.com/nethesis/dartagnan/athos/models"
)

func Cache() *redis.Client {

	client, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", configuration.Config.Redis.Host, configuration.Config.Redis.Port))
	if err != nil {
		panic(err.Error())
	}
	return client
}

func setRedisRecord(system models.System, client *redis.Client) (bool, string) {
	// Use 4 tiers
	tiers := uint32(4)

	now := time.Now()
	difference := system.Subscription.ValidUntil.Sub(now).Seconds()
	err := client.Cmd("HMSET", system.UUID, "tier_id", CalculateTier(system.UUID, tiers), "secret", system.Secret, "EX", int(difference)).Err
	if err != nil {
		return false, fmt.Sprintf("Can't save '%s' system inside Redis: %s", system.UUID, err)
	}
	// Set key expiration equal to Subscription.ValidUntil
	err = client.Cmd("EXPIRE", system.UUID, int(difference)).Err
	if err != nil {
		return false, fmt.Sprintf("Can't set expiration on '%s' system inside Redis: %s", system.UUID, err)
	}
	return true, ""
}

func deleteRedisRecord(uuid string, client *redis.Client) (bool, string) {
	err := client.Cmd("DEL", uuid).Err
	if err != nil {
		return false, fmt.Sprintf("Can't delete '%s' system inside Redis: %s", uuid, err)
	}
	return true, ""
}

func CalculateTier(uuid string, tiers uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(uuid))
	tier := h.Sum32() % tiers

	return tier
}

func SetValidSystem(system models.System) bool {
	client := Cache()
	ret, _ := setRedisRecord(system, client)
	if !ret {
		client.Close()
		return false
	}

	client.Close()
	return true
}

func DeleteValidSystem(system models.System) bool {
	client := Cache()
	ret, _ := deleteRedisRecord(system.UUID, client)
	if !ret {
		client.Close()
		return false
	}

	client.Close()
	return true
}

func BulkSetValidSystems() (bool, []string) {
	var errors []string
	systems := getValidSystems()

	client := Cache()

	err := client.Cmd("MULTI").Err
	if err != nil {
		client.Close()
		return false, []string{fmt.Sprintf("Can't start Redis transaction: %s", err)}
	}

	err = client.Cmd("FLUSHALL").Err
	if err != nil {
		client.Close()
		return false, []string{fmt.Sprintf("Can't FLUSH Redis keys: %s", err)}
	}

	for _, system := range systems {
		ret, err := setRedisRecord(system, client)
		if !ret {
			errors = append(errors, err)
		}
	}

	err = client.Cmd("EXEC").Err
	if err != nil {
		client.Close()
		return false, []string{fmt.Sprintf("Can't finish Redis transaction: %s", err)}
	}

	client.Close()

	if len(errors) > 0 {
		return false, errors
	} else {
		return true, errors
	}
}

func getValidSystems() []models.System {
	var systems []models.System

	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Joins("JOIN subscriptions ON systems.subscription_id = subscriptions.id").Where("valid_until > NOW()").Find(&systems)

	return systems
}
