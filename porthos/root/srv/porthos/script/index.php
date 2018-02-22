<?php

include("lib.php");

list($system_id, $repo_version, $repo_name, $repo_arch) = array_slice(explode('/', $_SERVER['DOCUMENT_URI']), 1, 4);
$file_name = basename($_SERVER['DOCUMENT_URI']);
$file_path = substr($_SERVER['DOCUMENT_URI'], strpos($_SERVER['DOCUMENT_URI'], '/', 1));

$access = get_access_descriptor($system_id);

if($access['tier_id'] === FALSE) {
    exit_http(404);
} elseif ($access['tier_id'] === '') {
    exit_http(403);
} elseif($file_name == 'repomd.xml') {
    return_file(dirname($file_path) . '/repomd-t' . $access['tier_id'] . '.xml');
} else { 
    return_file($file_path);
}

