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

type System struct {
	ID        int       `db:"id" json:"id"`
	CreatorID string    `db:"creator_id" json:"creator_id"`
	UUID      string    `db:"uuid" json:"uuid"`
	Tags      string    `db:"tags" json:"tags"`
	PublicIP  string    `db:"public_ip" json:"public_ip"`
	Status    string    `db:"status" json:"status"`
	Created   time.Time `db:"created" json:"created"`

	Subscription   Subscription `json:"subscription"`
	SubscriptionID int          `db:"subscription_id" json:"-"`
}

type SystemJSON struct {
	Tags string `db:"tags" json:"tags"`
}
