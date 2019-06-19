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

function return_file($file_path) {
    header('X-Accel-Redirect: ' . $file_path);
    exit(0);
}

function exit_basic_auth_required() {
    header('WWW-Authenticate: Basic realm="subscription"');
    header('HTTP/1.1 401 Unauthorized');
    echo "Provide system subscription credentials\n";
    exit(0);
}

function exit_http($code) {
    header('X-Accel-Redirect: /error/' . $code);
    exit(1);
}

function get_access_descriptor($system_id) {
    $redis = new Redis();
    if( ! $redis->connect($_SERVER['PORTHOS_REDIS'])) {
        error_log(sprintf('[ERROR] redis connect(%s) failed!', $_SERVER['PORTHOS_REDIS']));
        exit_http(503);
    };
    $descriptor = $redis->hMGet($system_id, array('tier_id', 'secret', 'icat'));
    $redis->close();
    return $descriptor;
}

function parse_uri($uri) {
    $matches = array();
    $parts = array(
        'uri' => $uri,
        'system_id' => NULL,
        'full_path' => $uri,
        'prefix' => NULL,
        'version' => NULL,
        'repo' => NULL,
        'arch' => NULL,
        'rest' => NULL,
    );
    preg_match('#^/(?P<prefix>autoupdate|stable)(?:/(?P<system_id>[\w-]{36,48}))?(?P<full_path>/(?P<version>[\d\.]+)/(?P<repo>[\w-]+)/(?P<arch>\w+)/(?P<rest>.*))$#', $uri, $matches);
    $parts = array_merge($parts, $matches);
    return $parts;
}
