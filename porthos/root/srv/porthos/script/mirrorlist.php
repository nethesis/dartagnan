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

require_once("config-" . $_SERVER['PORTHOS_SITE'] . ".php");

$version = $_GET['nsversion'];
$arch = $_GET['arch'];
$repo = $_GET['repo'];
$system_id = isset($_GET['systemid']) ? $_GET['systemid'] : '';
$use_tier = isset($_GET['usetier']) && ! in_array($_GET['usetier'], array('$YUM0', 'no', '0', ''));

$valid_version = in_array($version, $config['versions']);
$valid_arch = in_array($arch, $config['arches']);
$valid_repo = in_array($repo, $config['repositories']);
$valid_system_id = ! $system_id || preg_match('/^[\w-]+$/', $system_id);

header('Content-type: text/plain; charset=UTF-8');

if( ! $valid_arch || ! $valid_repo || ! $valid_version || ! $valid_system_id ) {
    header("HTTP/1.0 400 Invalid request");
    echo "Invalid request\n";
    exit(1);
}

if($use_tier) {
    $path = "autoupdate/";
} else {
    $path = "stable/";
}

if($system_id) {
    $path .= $system_id . '/';
}

foreach($config['base_urls'] as $baseurl) {
    echo $baseurl . $path . "${version}/${repo}/${arch}\n";
}
