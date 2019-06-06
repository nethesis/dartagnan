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

if ( ! isset($_SERVER['PHP_AUTH_USER'])) {
    header('WWW-Authenticate: Basic realm="subscription"');
    header('HTTP/1.1 401 Unauthorized');
    echo "Provide system subscription credentials\n";
    exit;
}

// Disable the Content-Type header in PHP, so that nginx x-accel can add its own
ini_set('default_mimetype', FALSE);

// Mask any repo that does not belong to the site:
$repo = get_uri_part($_SERVER['DOCUMENT_URI'], 'repo');
if(! in_array($repo, $config['repositories'])) {
    exit_http(404);
}

$access = get_access_descriptor($_SERVER['PHP_AUTH_USER']);
$valid_credentials = $_SERVER['PHP_AUTH_PW'] === $access['secret'];
if($config['legacy_auth']) {
    $valid_credentials = $valid_credentials || $_SERVER['PHP_AUTH_USER'] ===  $_SERVER['PHP_AUTH_PW'];
}
if (! is_numeric($access['tier_id']) || ! $valid_credentials) {
    exit_http(403);
}

if($access['tier_id'] < 0) {
    $hash = 0;
    foreach(str_split($_SERVER['PHP_AUTH_USER']) as $c) {
        $hash += ord($c);
    }
    $hash = $hash % 256;
    if($hash < 12) { // 5%
        $tier_id = 0;
    } elseif($hash < 38) { // 15%
        $tier_id = 1;
    } elseif($hash < 76) { // 30%
        $tier_id = 2;
    } else { // 50%
        $tier_id = 3;
    }
} else {
    $tier_id = intval($access['tier_id']);
}

if(basename($_SERVER['DOCUMENT_URI']) == 'repomd.xml') {
    header('Cache-Control: private');
    error_log(sprintf('[NOTICE] %s: %s using tier %s%s on repo %s',
        $_SERVER['PORTHOS_SITE'],
        $_SERVER['PHP_AUTH_USER'],
        $tier_id,
        isset($hash) ? ' (automatic)' : '', $repo)
    );
}

return_file('/t' . $tier_id . $_SERVER['DOCUMENT_URI']);
