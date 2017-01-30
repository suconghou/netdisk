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

disk config set "token"

disk config get

disk config list

disk help
```

fast download and multithreading

send to stdout for gzip or xz decode play `disk play test.mp4.gz --stdout | gzip -d | mpv  -`

```
https://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=fNThTaiSso4OtkgTsbtiFpyt&redirect_uri=oob&scope=netdisk
```


