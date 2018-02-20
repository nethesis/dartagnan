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

package models

import (
	"time"
)

type Alert struct {
	ID        int       `db:"id" json:"-"`
	AlertID   string    `db:"alert_id" json:"alert_id"`
	Priority  string    `db:"priority" json:"priority"`
	Note      string    `db:"note" json:"note"`
	Status    string    `db:"status" json:"status"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`

	System   System `json:"system"`
	SystemID int    `db:"system_id" json:"-"`
}

type AlertHistory struct {
	ID         int       `db:"id" json:"-"`
	AlertID    string    `db:"alert_id" json:"alert_id"`
	Priority   string    `db:"priority" json:"priority"`
	Resolution string    `db:"resolution" json:"resolution"`
	StatusFrom string    `db:"status_from" json:"status_from"`
	StatusTo   string    `db:"status_to" json:"status_to"`
	StartTime  time.Time `db:"start_time" json:"start_time"`
	EndTime    time.Time `db:"end_time" json:"end_time"`

	System   System `json:"system"`
	SystemID int    `db:"system_id" json:"-"`
}

type AlertJSON struct {
	SystemUUID string `db:"uuid" json:"uuid"`
	AlertID    string `db:"alert_id" json:"alert_id"`
	Status     string `db:"status" json:"status"`
}
