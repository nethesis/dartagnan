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
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/satori/go.uuid"

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

func GetAlertHumanName(alertId string, locale string) string {
	i18n.MustLoadTranslationFile("templates/en-US-alert.json")
	T, _ := i18n.Tfunc(locale, locale, locale)
	parts := strings.Split(alertId, ":")

	switch parts[0] {
	/*
		system:heartbeat:link
		system:backup:failure
	*/
	case "system":
		if parts[1] == "heartbeat" && parts[2] == "link" {
			return T("Link failed")
		} else if parts[1] == "backup" && parts[2] == "failure" {
			return T("Backup failed")
		}
	/*
		load:load
	*/
	case "load":
		return T("System load")
	/*
		service:([a-z0-9-@]+):stopped
	*/
	case "service":
		return fmt.Sprintf("%s %s", strings.ToUpper(parts[1]), T(parts[2]))
	/*
		df:boot:percent_bytes:free
		df:root:percent_bytes:free
	*/
	case "df":
		if parts[1] == "boot" {
			return T("Boot partition")
		} else if parts[1] == "root" {
			return T("Root partition")
		}
	/*
		swap:percent:free
	*/
	case "swap":
		return T("SWAP partition")
	/*
		md:([a-z0-9-]+):md_disks:([a-z]+)
	*/
	case "md":
		return fmt.Sprintf("%s %s %s", T("RAID"), parts[1], T("failed"))
	/*
		wan:([a-z0-9-]+):down
	*/
	case "wan":
		return fmt.Sprintf("%s %s %s", T("WAN"), parts[1], T("down"))
	/*
		ping:ping:([a-z0-9-.]+)
		ping:ping_([a-z]+):([a-z0-9-.]+)
	*/
	case "ping":
		if strings.Index(parts[1], "_") > 0 {
			n := strings.Split(parts[1], "_")
			return fmt.Sprintf("%s %s (%s)", T("PING"), n[1], parts[2])
		} else {
			return fmt.Sprintf("%s %s", T("PING"), parts[2])
		}
	/*
		nut:ups:voltage:input
	*/
	case "nut":
		return T("UPS on battery")
	}
	return alertId
}

func GetAlertPriority(alertID string) string {
	parts := strings.Split(alertID, ":")
	switch parts[0] {
	case "system":
		return "HIGH"
	case "load":
		return "HIGH"
	case "service":
		return "AVERAGE"
	case "df":
		return "AVERAGE"
	case "swap":
		return "WARNING"
	case "md":
		return "HIGH"
	case "wan":
		return "WARNING"
	case "ping":
		return "AVERAGE"
	case "nut":
		return "HIGH"
	}

	return "AVERAGE"
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

func CheckSystemOwnership(systemID string, creatorID string) bool {
	var system models.System
	db := database.Database()
	db.Where("id = ? AND creator_id = ?", systemID, creatorID).First(&system)
	db.Close()

	if system.ID == 0 {
		return false
	}

	return true
}

func GetSystemById(systemID int) models.System {
	var system models.System
	db := database.Database()
	db.Preload("Subscription.SubscriptionPlan").Where("id = ?", systemID).First(&system)
	db.Close()

	return system
}

func CanAccessAlerts(plan models.SubscriptionPlan) bool {
	switch plan.ID {
	case 1: // Trial
		return true
	case 4: // Fiorentina
		return true
	case 5: // Pizza
		return true
	}
	return false
}
