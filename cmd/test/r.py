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
files = {'file': ('getFilters.xml', xml)}
result = requests.post(url, files=files, headers=headers, verify=False)
print("=============")
print(result.text)
print("=============")
print(result)

"""
p = result.prepare()
print(p)
print(p.body)
print(p.body.decode('ascii'))
print(result.)

"""