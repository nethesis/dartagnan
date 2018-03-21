/*
 * Copyright (C) 2018 Nethesis S.r.l.
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
 */

package methods

import (
	"net"
	"net/http"
	"github.com/gin-gonic/gin"
)

func ReverseLookup(c *gin.Context) {
	ip := c.Param("ip")

	names, err := net.LookupAddr(ip)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "reverse lookup failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, names)
}
