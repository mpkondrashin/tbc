#!/bin/sh

#curl -i -k --header "X-SMS-API-KEY:${1}" -F "requestFile=@get_filter.xml" "https://${2}/ipsProfileMgmt/getFilters"

curl -v  -X POST -k --header "X-SMS-API-KEY: ${1}" --form name=@getFilters.xml https://${2}/ipsProfileMgmt/getFilters
