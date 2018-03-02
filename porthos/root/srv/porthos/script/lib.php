<?php

function return_file($file_path) {
    header('X-Accel-Redirect: ' . $file_path);
    exit(0);
}

function exit_http($code) {
    header('X-Accel-Redirect: /error/' . $code);
    exit(1);
}

function get_access_descriptor($system_id, $secret) {
    $redis = new Redis();
    if( ! $redis->connect($_SERVER['PORTHOS_REDIS'])) {
        error_log(sprintf('[ERROR] redis connect(%s) failed!', $_SERVER['PORTHOS_REDIS']));
        exit_http(503);
    };
    $descriptor = $redis->hMGet($system_id, array('tier_id', 'secret'));
    $redis->close();
    return $descriptor;
}
