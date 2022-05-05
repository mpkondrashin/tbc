import requests
import sys
import io

headers = {
    'Content-Type': 'application/xml',
    'X-SMS-API-KEY': sys.argv[1],
}

url = f"https://{sys.argv[2]}/ipsProfileMgmt/getFilters"
xml='<getFilters><profile name="test"></profile><filter><number>51</number></filter></getFilters>'
print(url)
print(headers)
print(xml)
files = {'BackupFile': io.StringIO(xml)}
#result = requests.Request('POST', url, files=files, headers=headers)
result = requests.request('POST', url, files=files, headers=headers, verify=False)
print(result)
print("=============")
print(result.content)
"""
p = result.prepare()
print(p)
print(p.body)
print(p.body.decode('ascii'))
print(result.)

"""