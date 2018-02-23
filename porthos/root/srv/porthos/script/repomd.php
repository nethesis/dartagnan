<?php

include("lib.php");

$file_name = basename($_SERVER['DOCUMENT_URI']);
$access = get_access_descriptor($_SERVER['PHP_AUTH_USER'], $_SERVER['PHP_AUTH_PW']);

if($access['tier_id'] === FALSE) {
    exit_http(404);
} elseif ($access['tier_id'] === '' || $access['secret'] != $_SERVER['PHP_AUTH_PW']) {
    exit_http(403);
} elseif($file_name == 'repomd.xml') {
    header('Content-Type: text/xml');
    return_file(dirname($_SERVER['DOCUMENT_URI']) . '/repomd-t' . $access['tier_id'] . '.xml');
}

