# TBCheck - the most cautious handle of Trend Micro Tipping Point policy

Version 0.1
by Mikhail Kondrashin


This utility changes all filters with block action in given profile to specified ActionSet.

## Configuration

tbcheck.yaml configuration file:
```yaml
SMS:
  URL: https://<SMS address>
  APIKey: 123413441234-1234-1234-12341234
  SkipTLSVerify: false
Profile: <profile name>
Actionset: Permit / Notify / Trace
Distribution:
  Priority: high #low
  SegmentGroup: <segment group>
```

tbc searches for this file in following locations:

1. /etc/tbcheck/
2. $HOME/.tbcheck
3. Current folder

