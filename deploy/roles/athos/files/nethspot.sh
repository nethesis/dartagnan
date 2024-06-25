#!/bin/bash

# read data from input
config_path=$1
email=$2

# define variables
host=$(cat $config_path | jq -r .integrations.nethspot.host)
login_username=$(cat $config_path | jq -r .integrations.nethspot.username)
login_password=$(cat $config_path | jq -r .integrations.nethspot.password)

# get token login
token=$(curl -s -X POST $host/api/login -H 'Content-Type: application/json' --data '{"username": "'$login_username'", "password": "'$login_password'"}' | jq -r .token)

# generate account password
account_pass=($(cat /dev/urandom | head -c 10 | sha1sum))

# compose body
body=$(cat <<-END
{
    "uuid": "$(uuidgen)",
    "username": "$email",
    "name": " $email",
    "email": "$email",
    "type": "reseller",
    "password": "$account_pass",
    "hotspot_id": 0,
    "subscription_plan_id": 4
}
END
)

# create account
status=$(curl -s -X POST $host/api/accounts -H 'Content-Type: application/json' -H "Token: $token" --data "$body" | jq -r .status)

# check status
if [ "$status" = "success" ]; then
    JSON_STRING=$( jq -n -c \
                    --arg host "$host" \
                    --arg email "$email" \
                    --arg password "$account_pass" \
                    '{ host: $host, username: $email, password: $password }' )
    echo "$JSON_STRING"
else
    JSON_STRING=$( jq -n -c \
                    --arg status "$status" \
                    '{ error: $status }' )
    echo "$JSON_STRING"
    exit 1
fi