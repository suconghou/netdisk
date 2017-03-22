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

```
https://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=fNThTaiSso4OtkgTsbtiFpyt&redirect_uri=oob&scope=netdisk
```

http://pan.plyz.net/Man.aspx


