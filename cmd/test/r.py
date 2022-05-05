import requests
import sys

headers = {
    'Content-Type': 'application/xml',
    'X-SMS-API-KEY': sys.argv[1],
}

url = f"https://{sys.argv[2]}/ipsProfileMgmt/getFilters"
xml='<getFilters><profile name="test"></profile><filter><number>51</number></filter></getFilters>'
print(url)
print(headers)
print(xml)
r = requests.post(url, headers=headers, data=xml, verify=False)
print(r)

