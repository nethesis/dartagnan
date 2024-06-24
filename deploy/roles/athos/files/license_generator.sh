#!/bin/bash

# define main variables
PRODUCTS=(ns8 nsec)
LICENSES=(personal business)
ns8_SERVICES=(nethvoice)
nsec_SERVICES=(threat_shield report hotspot)

# declare licenses and services price maps
declare -A LICENSES_PRICES
declare -A SERVICE_PRICES
declare -A MAP_NAMES
declare -A MAP_DESCRIPTIONS
LICENSES_PRICES=( ["personal-ns8"]=60 ["personal-nsec"]=60 ["business-ns8"]=250 ["business-nsec"]=250 )
SERVICE_PRICES=( ["threat_shield"]=5 ["report"]=6 ["hotspot"]=7 ["nethvoice"]=10)
MAP_NAMES=( ["personal-ns8"]="Personal (NS8)" ["business-ns8"]="Business (NS8)" ["personal-nsec"]="Personal (NSec)" ["business-nsec"]="Business (NSec)" )
MAP_DESCRIPTIONS=( ["personal-ns8"]="Personal NethServer 8" ["business-ns8"]="Business NethServer 8" ["personal-nsec"]="Personal NethSecurity 8" ["business-nsec"]="Business NethSecurity 8" )

# sort services by name
IFS=$'\n' ns8_services_list=($(sort -n <<<"${ns8_SERVICES[*]}"))
unset IFS
IFS=$'\n' nsec_services_list=($(sort -n <<<"${nsec_SERVICES[*]}"))
unset IFS

# get length of arrays
ns8_n=${#ns8_services_list[@]}
nsec_n=${#nsec_services_list[@]}

# define starting id for inserts
c=12

# loop products
for p in "${PRODUCTS[@]}"; do

    # loop licenses
    for l in "${LICENSES[@]}"; do
        # declare service length dynamic based on product
        declare -n service_length="${p}_n"

        # loop services
        for (( i = 1; i < (1 << $(echo "$service_length")); i++ )); do
            # define local variables for service maps
            list=()
            service_total=0

            # combine services
            for (( j = 0; j < $(echo "$service_length"); j++ )); do
                if (( (1 << j) & i )); then
                    # declare service list dynamic based on product
                    declare -n service="${p}_services_list[j]"

                    # add service to the map
                    list+=("${service}")

                    # get single service price
                    service_price=$(echo "${SERVICE_PRICES[${service}]}")

                    # sum service price to total services price
                    service_total=$(( $service_total + $service_price ))
                fi
            done
            # get license price
            license_price=$(echo "${LICENSES_PRICES[$l-$p]}")

            # sum basic license price with services price
            total=$(( $license_price + $service_total ))

            # define license name and description
            name=$(echo "${MAP_NAMES[$l-$p]}")
            description=$(echo "${MAP_DESCRIPTIONS[$l-$p]}")

            # print insert query
            (IFS=,; printf "INSERT INTO subscription_plans VALUES ($c, '!$l-$p+%s', '$l-$p', '$name', '$description', $total.00, $license_price.00, 365);\n"  "${list[*]}")

            # increment insert id
            (( c++ ))
        done
    done
done