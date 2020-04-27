# netdisk


baidu netdisk use pcs api

and with some useful tools

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

disk help
```

fast download and multithreading


### other flag

`disk put local_file -f`  fore rewrite remote file if conflicted

`disk info file --link` show file info and download link ,the link can be downloaded in multithread


```
https://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=fNThTaiSso4OtkgTsbtiFpyt&redirect_uri=oob&scope=netdisk
```


## Static File Server

`disk serve` start a static file server 

`disk serve -h` see help

`-p` set the listen port

`-d` set the document root

> directory list is enabled by default

## Wget Download

`disk wget http://url`

like wget but only for http/https 

It is multithreading and with awesome features

### http header control

`--cookie "cookie string"`

`--refer "http refer string"`

`--ua "user agent string"`

### range control

use `--range:1230-123456` or `--range:1230-` to force get certain range content

you also can use
```
disk wget url --range:0-88000
disk wget url --range:88000-988000
disk wget url --range:988000-
```
thus will not break your file which just like `disk wget url`

### speed control

`--fast/--slow` `--fat/--thin` can be used for speed control

for thread 

default thread is 8
```
disk wget url --fast // up to 16 thread
disk wget url --slow // set to 4 thread
```

for thunk

default thunk is 2097152 (2MB)
```
disk wget url --fat // set to 8MB
disk wget url --thin // set to 256KB
```

## Play Url Video

`disk play url`

is like wget download but have another two features

### It call video player automatically

It calls player once download > 2%

on unix use `mpv` player

on windows use `PotPlayerMini.exe` player

those commands should be called directly or 
it will failed silently

### It can write data to stdout rather than file

use `--stdout` to write data to stdout ranther than file

for example 

use another player to play

`disk play url --stdout | ffplay -i -`

write to stdout for gzip or xz decode play 

`disk play url --stdout | gzip -d | mpv  -`

play exist file but the remain data download to stdout to play

`(cat a.flv && disk play /test/a.flv --stdout) | mpv -`


## Proxy 

`disk proxy` 

`disk reverse`

### Reverse Proxy

`disk reverse` start a reverse proxy server

it is like nginx reverse proxy , but can work with upstream proxy

`disk reverse -h` see help

`-u` is your reverse proxy url aka proxy_pass url

> it can be any url http/https or with uri

`-p` is the server port, default 8123

#### Reverse Proxy With Upstream Proxy

`-proxy` to use an upstream http(s) proxy

> `-proxy http://your_http_proxy:6056` or `-proxy https://your_http_proxy:6056` 

`-socks` to use an upstream socks5 proxy

> `-socks 127.0.0.1:1080`

**`-socks` is used if both proxy are configured**

> if your proxy is only http_proxy proxy then you can only proxy http backend (no https backend) 



### Forward Proxy

`disk proxy` start a http/https/socks5 proxy server 

`disk proxy -h` see help

`-p` is the proxy server port, default 8123

#### http(s) proxy

`disk proxy` start the server can be used as a http proxy and https proxy server

#### socks5 proxy

`disk proxy` start the server can be used as a socks5 proxy server

#### socks5 to http(s) proxy

use `-socks` to set an upstream socks5 proxy

which all the proxy request(http/https/socks5) will pass to 

## Network


`disk network` test download speed for given url 

`echo http://xxx.com/filepath  | tdisk network -s 0 -t 10`  
`-s 0` for unlimit filesize

test for ip with one host 

`disk network -host xxx.com -path /filepath`

`-path` begin with `/`

