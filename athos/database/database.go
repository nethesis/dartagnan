/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Windmill project.
 *
 * WindMill is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * WindMill is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with WindMill.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/nethesis/dartagnan/athos/configuration"
)

func Database() *gorm.DB {
	db, err := gorm.Open("postgres", "sslmode=disable dbname="+configuration.Config.Database.Name+" host="+configuration.Config.Database.Host+" port="+configuration.Config.Database.Port+" user="+configuration.Config.Database.User+" password="+configuration.Config.Database.Password)
	if err != nil {
		panic(err.Error())
	}
	return db
}
