curl -H "Content-Type: application/json" -d \
'{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjbxQ","sequence":1,"event":"app_start","param_int":0,"param_str":"some text"}
{"client_time":"2020-12-01 23:59:00","device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E","device_os":"iOS 13.5.1","session":"ybuRi8mAUypxjAbxQ","sequence":123,"event":"app_start","param_int":0,"param_str":"some text"}' \
-X POST http://localhost:80/ 