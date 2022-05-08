#!/bin/sh

#curl -i -k --header "X-SMS-API-KEY:${1}" -F "requestFile=@get_filter.xml" "https://${2}/ipsProfileMgmt/getFilters"

#curl --trace-ascii tr.txt -v  -X POST -k --header "X-SMS-API-KEY: ${1}" --form name=@getFilters.xml https://${2}/ipsProfileMgmt/getFilters

curl -k --header "X-SMS-API-KEY: ${1}" "https://${2}/ipsProfileMgmt/distributeProfile?"