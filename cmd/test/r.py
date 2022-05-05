import requests
import sys

headers = {'content-type': 'application/xml'}
xml='<getFilters><profile name="test"></profile><filter><number>51</number></filter></getFilters>'
r = requests.post(f"https://{sys.argv[1]}/ipsProfileMgmt/getFilters", data=xml)
print(r.text)

