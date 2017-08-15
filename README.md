# netdisk


baidu netdisk

pcs api for file download

### usage


config file

on unix  `/etc/disk.json` and use `mpv` player

on windows  `C:\Users\Default\disk.json`  and use `PotPlayerMini.exe` player


```
disk info

disk info /path/to/file

disk ls

disk cd /path

disk get /path/to/file

disk put /local/file

disk rm /path/to/file

disk wget /path/to/file

disk wget http://xxx

disk play /path/to/file

disk play http://xxx

disk hash /local/file

disk pwd

disk mv oldpath newpath

disk mkdir dirname

disk task list

disk task add / http://fileurl

disk task info taskId

disk task remove taskId

disk config set "token"

disk config setapp "app"

disk config get

disk config list

disk help
```

fast download and multithreading

send to stdout for gzip or xz decode play `disk play test.mp4.gz --stdout | gzip -d | mpv  -`



### speed and thunk control

`--fast/--slow   --fat/--thin` can be used for `disk wget file/url` and `disk play file/url`

for thread

default thread is 8
```
disk wget file/url --fast // up to 16 thread
disk wget file/url --slow // set to 4 thread
```

for thunk
default thunk is 2M
```
disk wget file/url --fat // set to 8M
disk wget file/url --thin // set to 256K
```




### request with cookie or refer or ua

`cookie,ua,refer` control can be used for `disk wget url` and `disk play url`

`--cookie "cookie string"`

`--refer "http refer string"`

`--ua "user agent string"`


`disk reverse` is a reverse proxy server like nginx but can work with upstream proxy

`--socks` use socks5 

`--proxy` use http/https proxy  

`-u https://backend --proxy http://your_http_proxy:6056` 
( or https://your_http_proxy:6056 , it doesn't matter)

backend is your proxy backend server (can be any url http/https or with uri)

proxy is your http_proxy https_proxy or

socks proxy `--socks 127.0.0.1:1080`

> if your proxy is only http_proxy proxy then you can only proxy http backend (no https backend) 


`disk proxy` give you a http_proxy https_proxy 

socks to http/https/socks

`-p` is the listen port

`--socks` is a upstream socks5 proxy  eg `x.x.x.x:6056`


### other flag


`disk put local_file -f`  fore rewrite remote file if conflicted

`disk info file -i` show file info and download link ,the link can be downloaded in multithread

use `GET` method to detect url `Content-Length` instead of `HEAD` method in case of some server declined

`disk wget http://xxx  --GET`

`disk wget file/url --debug` see debug log information

use  `--range:1230-123456` or `--range:45612-` to force get certain range content, it is supported  both `wget` and `play` action

you also can use
```
disk wget file --range:0-88000
disk wget file --range:88000-988000
disk wget file --range:988000-
```
thus will not break your file which just like `disk wget file`

```
https://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=fNThTaiSso4OtkgTsbtiFpyt&redirect_uri=oob&scope=netdisk
```

http://pan.plyz.net/Man.aspx

