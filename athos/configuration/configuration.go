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

package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Notifications struct {
	PortalUrl   string `json:"portal_url"`
	HelpUrl     string `json:"help_url"`
	DocsUrl     string `json:"docs_url"`
	PortalTitle string `json:"portal_title"`
	Email       struct {
		From         string `json:"from"`
		SMTPHost     string `json:"smtp_host"`
		SMTPPort     int    `json:"smtp_port"`
		SMTPUser     string `json:"smtp_user"`
		SMTPPassword string `json:"smtp_password"`
	} `json:"email"`
}

type Configuration struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"database"`
	Redis struct {
		Host  string `json:"host"`
		Port  string `json:"port"`
	} `json:"redis"`
	Cors       struct {
		Headers []string `json:"headers"`
		Origins []string `json:"origins"`
		Methods []string `json:"methods"`
	} `json:"cors"`
	Auth0 struct {
		Domain        string `json:"domain"`
		IdentifierAPI string `json:"identifier_api"`
	} `json:"auth0"`
	PayPal struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"paypal"`
	Notifications Notifications `json:"notifications"`
}

var Config = Configuration{}

func Init(ConfigFilePtr *string) {
	// read configuration
	if _, err := os.Stat(*ConfigFilePtr); err == nil {
		file, _ := os.Open(*ConfigFilePtr)
		decoder := json.NewDecoder(file)
		// check errors or parse JSON
		err := decoder.Decode(&Config)
		if err != nil {
			fmt.Println("Configuration parsing error:", err)
		}
	}

	// read from ENV variables
	if os.Getenv("DB_USER") != "" {
		Config.Database.User = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		Config.Database.Password = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_HOST") != "" {
		Config.Database.Host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_PORT") != "" {
		Config.Database.Port = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		Config.Database.Name = os.Getenv("DB_NAME")
	}
	if os.Getenv("REDIS_HOST") != "" {
		Config.Redis.Host = os.Getenv("REDIS_HOST")
	}
	if os.Getenv("REDIS_PORT") != "" {
		Config.Redis.Port = os.Getenv("REDIS_PORT")
	}

	if os.Getenv("CORS_ORIGINS") != "" {
		Config.Cors.Origins = strings.Split(os.Getenv("CORS_ORIGINS"), " ")
	}

	if os.Getenv("AUTH0_DOMAIN") != "" {
		Config.Auth0.Domain = os.Getenv("AUTH0_DOMAIN")
	}
	if os.Getenv("AUTH0_IDENTIFIER_API") != "" {
		Config.Auth0.IdentifierAPI = os.Getenv("AUTH0_IDENTIFIER_API")
	}
}
