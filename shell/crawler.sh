#!/bin/bash
du -h --max-depth=1 $1
#每12小时一次
host="http://127.0.0.1:8888"
url_timing=$host"/baidu/getAllTags"
url_timing_detail=$host"/baidu/getDetail"
echo "getAllTags========>$url_timing"
echo "getDetail========>$url_timing_detail"
echo "---------------------------------------------------task_api_start"
curl -X GET -G --data-urlencode "code=xx"  -i  $url_timing -k
curl -X GET -G --data-urlencode "code=xx"  -i  $url_timing_detail -k
echo "---------------------------------------------------task_api_end"
