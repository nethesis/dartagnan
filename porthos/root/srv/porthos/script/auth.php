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

// Mask any repo/version/arch that does not belong to the site:
if(! in_array($uri['repo'], $config['repositories'])
    || ! in_array($uri['version'], $config['versions'])
    || ! in_array($uri['arch'], $config['arches'])) {
    exit_http(404);
}

if(! isset($_SERVER['PHP_AUTH_USER']) || ! isset($_SERVER['PHP_AUTH_PW'])) {
    exit_http(403);
}

$access = get_access_descriptor($_SERVER['PHP_AUTH_USER']);
$valid_credentials = $_SERVER['PHP_AUTH_PW'] === $access['secret'];
if($config['legacy_auth']) {
    $valid_credentials = $valid_credentials || $_SERVER['PHP_AUTH_USER'] ===  $_SERVER['PHP_AUTH_PW'];
}
$has_access_disabled = ! is_numeric($access['tier_id']);

if($access['tier_id'] < 0) {
    $hash = 0;
    foreach(str_split($_SERVER['PHP_AUTH_USER']) as $c) {
        $hash += ord($c);
    }
    $hash = $hash % 256;
    if($hash < 13) { // 5%
        $tier_id = 0;
    } elseif($hash < 51) { // +15% = 20%
        $tier_id = 1;
    } elseif($hash < 128) { // +30% = 50%
        $tier_id = 2;
    } else { // +50% = 100%
        $tier_id = 3;
    }
    $tier_id += $config['tier_id_base'];
} else {
    $tier_id = intval($access['tier_id']);
}

if(basename($uri['rest']) == 'repomd.xml') {
    header('Cache-Control: private');
    application_log(json_encode(array(
        'porthos_site' => $_SERVER['PORTHOS_SITE'],
        'connection' => $_SERVER['CONNECTION'] ?: '',
        'system_id' => $_SERVER['PHP_AUTH_USER'],
        'remote_addr' => $_SERVER['REMOTE_ADDR'],
        'repo' => $uri['repo'],
        'version' => $uri['version'],
        'arch' => $uri['arch'],
        'tier_id' => $uri['prefix'] == 'autoupdate' ? $tier_id : NULL,
        'tier_auto' => isset($hash),
        'tls' => isset($_SERVER['HTTPS']),
        'auth_response' => ! $valid_credentials ? 'bad_credentials' : $has_access_disabled ? 'bad_access' : 'pass',
    )));
}

if ($has_access_disabled || ! $valid_credentials) {
    exit_http(403);
}

if($uri['prefix'] == 'autoupdate') {
    return_file('/T' . $tier_id . $uri['full_path']);
} else {
    return_file('/head' . $uri['full_path']);
}
