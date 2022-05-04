#!/bin/sh

curl -i -k --header "X-SMS-API-KEY:${1}" -F "requestFile=@get_filter.xml" "https://${2}/ipsProfileMgmt/getFilters"