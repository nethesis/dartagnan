<?php

include("lib.php");

if ( ! isset($_SERVER['PHP_AUTH_USER'])) {
    header('WWW-Authenticate: Basic realm="subscription"');
    header('HTTP/1.1 401 Unauthorized');
    echo "Provide system subscription credentials\n";
    exit;
}

// Disable the Content-Type header in PHP, so that nginx x-accel can add its own
ini_set('default_mimetype', FALSE);

$access = get_access_descriptor($_SERVER['PHP_AUTH_USER'], $_SERVER['PHP_AUTH_PW']);

if($access['tier_id'] === FALSE) {
    exit_http(403);
} elseif ($access['tier_id'] === '' || $access['secret'] != $_SERVER['PHP_AUTH_PW']) {
    exit_http(403);
} else {
    if(basename($_SERVER['DOCUMENT_URI']) == 'repomd.xml') {
        header('Cache-Control: private');
    }
    return_file('/t' . $access['tier_id'] . $_SERVER['DOCUMENT_URI']);
}

