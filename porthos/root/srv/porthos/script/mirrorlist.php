<?php

/*
 * Copyright (C) 2018 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with NethServer.  If not, see COPYING.
 */

$version = $_GET['nsversion'];
$arch = $_GET['arch'];
$repo = $_GET['repo'];
$use_tier = isset($_GET['usetier']) && ! in_array($_GET['usetier'], array('$YUM0', 'no', '0', ''));

$valid_version = in_array($version, array('7.4.1708', '7.5.1804', '7.6.1810'));
$valid_arch = in_array($arch, array('x86_64'));
$valid_repo = in_array($repo, array(
    'base',
    'updates',
    'extras',
    'epel',
    'centos-sclo-sclo',
    'centos-sclo-rh',
    'nethserver-base',
    'nethserver-updates',
));

header('Content-type: text/plain; charset=UTF-8');

if( ! $valid_arch || ! $valid_repo || ! $valid_version ) {
    header("HTTP/1.0 400 Bad request");
    error_log("[ERROR] invalid request: " . $_SERVER['REQUEST_URI']);
    echo "Invalid request\n";
    exit(1);
}

if($use_tier) {
    echo "https://m1.nethserver.com/autoupdate/${version}/${repo}/${arch}\n";
} else {
    echo "https://m1.nethserver.com/stable/${version}/${repo}/${arch}\n";
}
