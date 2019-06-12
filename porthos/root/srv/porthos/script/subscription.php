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

require_once("lib.php");
require_once("config-" . $_SERVER['PORTHOS_SITE'] . ".php");

$uri = parse_uri($_SERVER['DOCUMENT_URI']);

if ( isset($_SERVER['HTTPS']) && ! $uri['system_id'] && ! isset($_SERVER['PHP_AUTH_USER'])) {
    exit_basic_auth_required();
} else {
    if($uri['system_id']) {
        // override PHP authentication with system_id token:
        $_SERVER['PHP_AUTH_USER'] = $uri['system_id'];
        $_SERVER['PHP_AUTH_PW'] = $uri['system_id'];
    }
}

// Disable the Content-Type header in PHP, so that nginx x-accel can add its own
ini_set('default_mimetype', FALSE);

if(! isset($_SERVER['PHP_AUTH_USER']) || ! isset($_SERVER['PHP_AUTH_PW'])) {
    exit_http(403);
}

$access = get_access_descriptor($_SERVER['PHP_AUTH_USER']);
$valid_credentials = $_SERVER['PHP_AUTH_PW'] === $access['secret'];
if($config['legacy_auth']) {
    $valid_credentials = $valid_credentials || $_SERVER['PHP_AUTH_USER'] ===  $_SERVER['PHP_AUTH_PW'];
}
$has_access_disabled = ! is_numeric($access['tier_id']) || $access['icat'] === FALSE;
if ($has_access_disabled || ! $valid_credentials) {
    exit_http(403);
}

$include_categories = array_filter(explode(',', $access['icat']));
$exclude_categories = array_values(array_diff($config['categories'], $include_categories));

header('Content-type: application/json; charset=UTF-8');
echo json_encode(array(
    'fmt_version' => 1,
    'include_categories' => $include_categories,
    'exclude_categories' => $exclude_categories,
));
