Auto Allow Web: A web api for open iptables access permission
=============================================================

How to install
--------------
To get started, clone the source and build
```shell
git clone  
cd auto-allow
make
```

install to /usr/local/sbin  (need root permission)
```shell
sudo make install
```

Start
--------
###Usage
````
-allow string
IPs always allow(default:127.0.0.1) (default "127.0.0.1")
-k string
access key(must have)
-open string
ports always open(default:22) (default "22")
-p int
listen port(default:8080) (default 8080)
-protect string
special ports to open(default is all)
-t string
allow access duration(default:12h) (default "12h")
````

### Systemd service
copy script/auto-allow-web.service to /etc/systemd/system and replace password in that file
````shell
systemctl link auto-allow-web
systemctl daemon-reload
````

Web Access
----------------
|  url   |   |
|  ----  | ----  |
|http://host:8080/ |Simple Page|
|http://host:8080/add?key=password|add current ip to allow ip list|
|http://host:8080/add?key=password&ip=1.1.1.1|add 1.1.1.1 to allow ip list|
|http://host:8080/remove?key=password|remove current ip from allow ip list|
|http://host:8080/remove?key=password&ip=1.1.1.1|remove 1.1.1.1 from allow ip list|
|http://host:8080/get?key=password |List all allow ip|
