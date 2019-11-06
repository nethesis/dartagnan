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

function application_log($message, $priority = LOG_INFO) {
    error_log($message);
    openlog('php', LOG_CONS | LOG_NDELAY, LOG_LOCAL3);
    syslog($priority, $message);
    closelog();
}

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

function get_snapshot_timestamp($snapshot_name) {
    $parts = array();
    if(!$snapshot_name || !preg_match('/^d(?<year>\d\d\d\d)(?<month>\d\d)(?<day>\d\d)$/', $snapshot_name, $parts)) {
        return time();
    }
    return mktime(0, 0, 0, $parts['month'], $parts['day'], $parts['year']);
}

function lookup_policy($version, $tier_id) {
    global $config;

    $policy_map = $config['autoupdate_policy'] ?: array();
    $policy_key = $version . '/' . $tier_id;
    $poldef_key = $version . '/*';

    if(isset($policy_map[$policy_key])) {
        $value = $policy_map[$policy_key];
    } elseif(isset($policy_map[$poldef_key])) {
        $value = $policy_map[$poldef_key];
    } else {
        $value = 'default';
    }

    return $value;
}

function lookup_fixed($path, $snapshot_start) {
    global $config;

    $root_path = $config['snapshots_dir'] ?: "/srv/porthos/webroot/";
    $snapshots = array_map('basename', glob($root_path . "d20*"));

    $start_key = array_search($snapshot_start, $snapshots);
    if( ! $start_key) {
        // the fixed snapshot does not exist: fall back to the empty repository.
        return 'empty';
    }
    $snapshots = array_slice($snapshots, $start_key);
    foreach($snapshots as $snapshot) {
        if(is_file($root_path . $snapshot . '/' . $path)) {
            return $snapshot;
        }
    }
    return 'head';
}

function lookup_monday($path, $tier_id) {
    global $config;

    $week_size = $config['week_size'] ?: 5;
    $root_path = $config['snapshots_dir'] ?: "/srv/porthos/webroot/";
    $snapshots = array_reverse(array_map('basename', glob($root_path . "d20*")));
    $last_snapshot_day_id = date('w', get_snapshot_timestamp($snapshots[0]));
    // $monday_offset formula:
    //     ($last_snapshot_day_id-1): rebase on Mondays
    //     ($last_snapshot_day_id > $tier_id ? 0 : $week_size): select current week Monday or previous one
    $monday_offset = ($last_snapshot_day_id-1) + ($last_snapshot_day_id > $tier_id ? 0 : $week_size);
    for($i = min($monday_offset, count($snapshots) - 1); $i >= 0; $i--) {
        if(is_file($root_path . $snapshots[$i] . '/' . $path)) {
            break;
        }
    }
    return $i < 0 ? 'head' : $snapshots[$i];
}

function lookup_snapshot($path, $version, $tier_id) {
    global $config;

    $policy_name = lookup_policy($version, $tier_id);
    if($policy_name == 'default') {
        return lookup_monday($path, $tier_id);
    } elseif(substr($policy_name, 0, 6) == 'fixed/') {
        return lookup_fixed($path, substr($policy_name, 6));
    } elseif($policy_name == 'empty') {
        return 'empty';
    } else {
        return 'head';
    }
}
