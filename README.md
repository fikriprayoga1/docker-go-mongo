# Documentation
This is repository of fikriprayoga1/go-mongo in docker hub. The function of this repository is create image of fikriprayoga1/go-mongo in docker hub, so you can learn how to create image and run the docker system image easily. This repository is part of https://hub.docker.com/repository/docker/fikriprayoga1/go-mongo project

## Best Practice
Use this command in your CMD(Windows) or Terminal(Mac or Linux). Before start the command, you must be in the directory of the folder you pulled.

### Step 1. Create fikriprayoga1/go-mongo image
```
docker build -t fikriprayoga1/go-mongo .
```

### Step 2. Pull mongodb image
```
docker pull mongo
```

### Step 3. Run docker compose
```
docker-compose up -d
```

### Step 4. Test in Postman app
You can use Postman app to test server with this configuration
```
Request Method = POST
Address = localhost:8080/requestAccess
Body Type = raw
Body Format = JSON
Body = {"Series":1, "SerialNumber":"4GKL3", "UUID":"upToYou"}
```

## Serial Number List
This part show the serial number you can try. Rules in golang script is every different UUID has limited for request access 4 times

### Series 1
```
4GKL3
16IQS
H1AP0
WVXLD
AC77B
XENRC
XU29Y
MXQ7I
AJSSC
EFTKR
```

### Series 2
```
JJK75
H0A21
TPT7P
U7XXH
N3HT0
E6KY3
HQD7V
1U4U5
TEPV7
7416G
```

### Series 3
```
EI4H8
RDI38
CQV5I
NRK1O
QSNUW
MDUZJ
132AW
IB86R
IP75M
P3AR8
```

### Series 4
```
ZK8IP
Q6L4K
OPX8H
2TD1D
MZIHT
YDJFH
AXNXK
GXAWE
UV17Z
BN3E7
```

### Series 5
```
TOJST
ATO26
SFKZ0
0MM9T
RF3Q5
8SH6I
TPZQM
GYGNB
8JMTT
SEI0S
```

### Series 6
```
ONG3L
L1UOS
STUHH
BU99Q
FR48V
B7TO9
8C5UD
NOWZY
9CZR7
3RG3J
```

## Utility Command
This part is command line to help you something

### Show running container list
```
docker container ls
```

### Show image list
```
docker images
```

### Show  container list
```
docker container ls -a
```

### Show  container log for debugging
```
docker logs go-mongo
```

### Show network list
```
docker network ls
```
