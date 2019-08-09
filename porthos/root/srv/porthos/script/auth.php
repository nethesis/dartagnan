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
ini_set('date.timezone', $config['timezone']);

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

if($access['tier_id'] < 0) {
    $hash = 0;
    foreach(str_split($_SERVER['PHP_AUTH_USER']) as $c) {
        $hash += ord($c);
    }
    $hash = $hash % 256;
    if($hash < 26) { // 10%
        $tier_id = 0;
    } elseif($hash < 77) { // +20% = 30%
        $tier_id = 1;
    } else { // +70% = 100%
        $tier_id = 2;
    }
} else {
    $tier_id = $access['tier_id'];
}

$is_tier_request = is_numeric($tier_id) && $uri['prefix'] == 'autoupdate';
$real_full_path = resolve_version_symlinks($uri['full_path']);
if($is_tier_request && $valid_credentials) {
    // Seeking a snapshot is a time-consuming op. Ensure we have valid
    // credentials before running it!
    $snapshot = lookup_snapshot($real_full_path, $tier_id, $config['week_size']);
} else {
    $snapshot = 'head';
}

if(basename($uri['rest']) == 'repomd.xml') {
    // repomd.xml is the entry point of repository (meta)data: let's keep track
    // of every client access to repositories:
    application_log(json_encode(array(
        'porthos_site' => $_SERVER['PORTHOS_SITE'],
        'connection' => $_SERVER['CONNECTION'] ?: '',
        'system_id' => $_SERVER['PHP_AUTH_USER'],
        'remote_addr' => $_SERVER['REMOTE_ADDR'],
        'repo' => $uri['repo'],
        'version' => $uri['version'],
        'arch' => $uri['arch'],
        'tier_id' => $is_tier_request ? $tier_id : -1,
        'tier_auto' => isset($hash),
        'tls' => isset($_SERVER['HTTPS']),
        'auth_response' => ! $valid_credentials ? 'bad_credentials' : 'pass',
        'snapshot' => $snapshot,
    )));
}

if (! $valid_credentials) {
    // Exit here, after sending the application_log record for repomd.xml requests.
    exit_http(403);
}

header('Cache-Control: private');
return_file('/' . $snapshot . $real_full_path);
