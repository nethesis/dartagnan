#!/bin/bash

# define main variables
PRODUCTS=(ns8 nsec)
LICENSES=(personal business)
ns8_SERVICES=()
nsec_SERVICES=(threat_shield hotspot)

# emulate associative arrays with individual variables
# license prices
LP_personal_ns8=120
LP_personal_nsec=47
LP_business_ns8=320
LP_business_nsec=297

# service prices
SP_threat_shield=48
SP_flashstart_lite=18
SP_hotspot=120
SP_nethvoice=280
SP_ldap_mattermost=24

# map names
MN_personal_ns8="Personal NethServer"
MN_business_ns8="Business NethServer"
MN_personal_nsec="Personal NethSecurity"
MN_business_nsec="Business NethSecurity"

# map descriptions
MD_personal_ns8="Personal NethServer"
MD_business_ns8="Business NethServer"
MD_personal_nsec="Personal NethSecurity"
MD_business_nsec="Business NethSecurity"

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

# collect rows for table, inserts and updates
table_rows=()
insert_rows=()
update_rows=()

# loop products
for p in "${PRODUCTS[@]}"; do

    # loop licenses
    for l in "${LICENSES[@]}"; do
        # get service length for this product
        length_var="${p}_n"
        service_length=${!length_var}

        # loop services
        for (( i = 1; i < (1 << service_length); i++ )); do
            # define local variables for service maps
            list=()
            service_total=0

            # combine services
            for (( j = 0; j < service_length; j++ )); do
                if (( (1 << j) & i )); then
                    # get service name from the product services list
                    service_list_var="${p}_services_list[$j]"
                    service="${!service_list_var}"

                    # add service to the map
                    list+=("${service}")

                    # get single service price
                    sp_var="SP_${service}"
                    service_price=${!sp_var}

                    # sum service price to total services price
                    service_total=$(( service_total + service_price ))
                fi
            done
            # get license price
            lp_var="LP_${l}_${p}"
            license_price=${!lp_var}

            # sum basic license price with services price
            total=$(( license_price + service_total ))

            # define license name and description
            mn_var="MN_${l}_${p}"
            name="${!mn_var}"
            md_var="MD_${l}_${p}"
            description="${!md_var}"

            # build services string
            services_str="$(IFS=,; echo "${list[*]}")"
            code="!${l}-${p}+${services_str}"

            # collect table row
            table_rows+=("$(printf "| %-4s | %-25s | %-7s | %-30s | %-14s | %-14s | %-10s |" \
                "$c" "$name" "$l-$p" "$services_str" "$license_price.00" "$service_total.00" "$total.00")")

            # collect insert query (column order: id, code, name, description, price, period, base_code, base_price)
            insert_rows+=("INSERT INTO subscription_plans VALUES ($c, '$code', '$name', '$description', $total.00, 365, '$l-$p', $license_price.00);")

            # collect update query (only updates rows where price or base_price changed)
            update_rows+=("UPDATE subscription_plans SET price = $total.00, base_price = $license_price.00 WHERE code = '$code' AND (price != $total.00 OR base_price != $license_price.00);")

            # increment insert id
            (( c++ ))
        done
    done
done

# print summary table
echo ""
echo "=== SUBSCRIPTION PLANS SUMMARY ==="
echo ""
printf "| %-4s | %-25s | %-7s | %-30s | %-14s | %-14s | %-10s |\n" \
    "ID" "Name" "Plan" "Services" "License Price" "Services Price" "Total"
printf "|%-6s|%-27s|%-9s|%-32s|%-16s|%-16s|%-12s|\n" \
    "------" "---------------------------" "---------" "--------------------------------" "----------------" "----------------" "------------"
for row in "${table_rows[@]}"; do
    echo "$row"
done
echo ""

# print insert queries (for new rows)
echo "=== SQL INSERT QUERIES (new rows) ==="
echo ""
for ins in "${insert_rows[@]}"; do
    echo "$ins"
done
echo ""

# print update queries (for existing rows with changed prices)
echo "=== SQL UPDATE QUERIES (existing rows, only if changed) ==="
echo ""
for upd in "${update_rows[@]}"; do
    echo "$upd"
done
