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

function get_uri_part($uri, $part_name) {
    $parts = array_combine(
        array('version', 'repo', 'arch'),
        array_slice(explode('/', $uri), 1, 3)
    );
    return isset($parts[$part_name]) ? $parts[$part_name] : FALSE;
}
