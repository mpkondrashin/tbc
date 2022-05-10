# TBCheck - the most cautious handle of Trend Micro Tipping Point policy

Version 0.1
by Mikhail Kondrashin

[![License](https://img.shields.io/badge/License-Apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)


This utility changes all policies in given profile to specified ActionSet.

**Note:** cs-report does not initiate registry scans. It uses last scan result.

## Configuration

```yaml
 -url string
    	Smart Check/Container Security URL (i.e. https://10.1.1.10:31616)
 -user string
    	User name (default "administrator")  
 -password string
    	Password
 -ignore_tls
    	Ignore TLS Errors when connecting to Smart Check/Container Security
```

After execution cs-report generates report into filename of following format:
```
report_YYYYMMDD.html
```

## Report Example

![Report example](screen.png)