<?php

/*
 * Copyright (C) 2019 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of Dartagnan.
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
 */

// this is just an example!

// categories (array):
//    this is the list of YUM comps categories that require an entitlement
$config['categories'] = array();

// base_urls (array):
//    base URLs of other porthos instances for load balancing
$config['base_urls'] = array("https://{$_SERVER['SERVER_NAME']}/");

// repositories (array):
//    list of valid repository names
$config['repositories'] = array('base', 'updates');

// versions (array):
//    list of valid version numbers
$config['versions'] = array('7.6.1810');

// arches (array):
//    list of valid architectures
$config['arches'] = array('x86_64');

// legacy_auth (boolean)
//     if TRUE, authenticate by user name only. The password must be equal to the user name.
//     if FALSE, check the user name and the password separately
$config['legacy_auth'] = FALSE;

// week_size (int)
//     number of days/snapshots per week. Ensure this value is consistent with
//     your crontab settings
$config['week_size'] = 5;

// timezone (string)
//     the PHP timezone for this application
$config['timezone'] = 'UTC';

// stop_autoupdate (array)
//     forcibly stop autoupdate contents distribution for the listed versions.
//     The /empty repository is served instead.
$config['stop_autoupdate'] = array();
