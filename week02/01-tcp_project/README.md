# 多阶段构建tcp_server、tcp_client

> 构建背景：
> - 以 golang:1.22 作为构建镜像
> - 以 ubuntu:20.04 作为运行基础镜像
> - 在Dockerfile中应用 `ARG TARGET` 来构建时通过传参构建不同的二进制文件
> - 使用docker-compose构建并运行两个容器，同时指定共享网络空间，以支持基于localhost访问

## 构建镜像并运行容器

```bash
# docker-compose up --build -d
Building server
Step 1/16 : FROM golang:1.22 AS builder
 ---> 964d0e6869f9
Step 2/16 : WORKDIR /app
 ---> Using cache
 ---> 6105c2214b4a
Step 3/16 : COPY go.mod .
 ---> Using cache
 ---> c97866fe71fb
Step 4/16 : COPY go.sum .
 ---> Using cache
 ---> 38c0e23b21cc
Step 5/16 : RUN go mod download
 ---> Using cache
 ---> fe9e6f035367
Step 6/16 : COPY . .
 ---> 54d0a54595ed
Step 7/16 : ARG TARGET
 ---> Running in fdf240b87723
Removing intermediate container fdf240b87723
 ---> f268257b9867
Step 8/16 : RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tcp_app ./cmd/tcp_${TARGET}.go
 ---> Running in 237dc902cac2
Removing intermediate container 237dc902cac2
 ---> a7acc6e9b1e7
Step 9/16 : FROM ubuntu:20.04
 ---> 9df6d6105df2
Step 10/16 : RUN apt-get update &&     apt-get install -y --no-install-recommends ca-certificates &&     rm -rf /var/lib/apt/lists/* &&     addgroup --system appgroup &&     adduser --system --ingroup appgroup --no-create-home appuser &&     mkdir -p /opt &&     chown appuser:appgroup /opt
 ---> Running in 8a2b26e82dcb
Get:1 http://security.ubuntu.com/ubuntu focal-security InRelease [128 kB]
Get:2 http://archive.ubuntu.com/ubuntu focal InRelease [265 kB]
Get:3 http://security.ubuntu.com/ubuntu focal-security/universe amd64 Packages [1273 kB]
Get:4 http://archive.ubuntu.com/ubuntu focal-updates InRelease [128 kB]
Get:5 http://archive.ubuntu.com/ubuntu focal-backports InRelease [128 kB]
Get:6 http://security.ubuntu.com/ubuntu focal-security/multiverse amd64 Packages [30.9 kB]
Get:7 http://security.ubuntu.com/ubuntu focal-security/restricted amd64 Packages [4016 kB]
Get:8 http://archive.ubuntu.com/ubuntu focal/multiverse amd64 Packages [177 kB]
Get:9 http://archive.ubuntu.com/ubuntu focal/universe amd64 Packages [11.3 MB]
Get:10 http://security.ubuntu.com/ubuntu focal-security/main amd64 Packages [4008 kB]
Get:11 http://archive.ubuntu.com/ubuntu focal/main amd64 Packages [1275 kB]
Get:12 http://archive.ubuntu.com/ubuntu focal/restricted amd64 Packages [33.4 kB]
Get:13 http://archive.ubuntu.com/ubuntu focal-updates/multiverse amd64 Packages [33.5 kB]
Get:14 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 Packages [4478 kB]
Get:15 http://archive.ubuntu.com/ubuntu focal-updates/universe amd64 Packages [1559 kB]
Get:16 http://archive.ubuntu.com/ubuntu focal-updates/restricted amd64 Packages [4178 kB]
Get:17 http://archive.ubuntu.com/ubuntu focal-backports/main amd64 Packages [55.2 kB]
Get:18 http://archive.ubuntu.com/ubuntu focal-backports/universe amd64 Packages [28.6 kB]
Fetched 33.1 MB in 6s (5791 kB/s)
Reading package lists...
Reading package lists...
Building dependency tree...
Reading state information...
The following additional packages will be installed:
  libssl1.1 openssl
The following NEW packages will be installed:
  ca-certificates libssl1.1 openssl
0 upgraded, 3 newly installed, 0 to remove and 2 not upgraded.
Need to get 2096 kB of archives.
After this operation, 5817 kB of additional disk space will be used.
Get:1 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 libssl1.1 amd64 1.1.1f-1ubuntu2.23 [1323 kB]
Get:2 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 openssl amd64 1.1.1f-1ubuntu2.23 [621 kB]
Get:3 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 ca-certificates all 20230311ubuntu0.20.04.1 [152 kB]
debconf: delaying package configuration, since apt-utils is not installed
Fetched 2096 kB in 2s (949 kB/s)
Selecting previously unselected package libssl1.1:amd64.
(Reading database ... 4124 files and directories currently installed.)
Preparing to unpack .../libssl1.1_1.1.1f-1ubuntu2.23_amd64.deb ...
Unpacking libssl1.1:amd64 (1.1.1f-1ubuntu2.23) ...
Selecting previously unselected package openssl.
Preparing to unpack .../openssl_1.1.1f-1ubuntu2.23_amd64.deb ...
Unpacking openssl (1.1.1f-1ubuntu2.23) ...
Selecting previously unselected package ca-certificates.
Preparing to unpack .../ca-certificates_20230311ubuntu0.20.04.1_all.deb ...
Unpacking ca-certificates (20230311ubuntu0.20.04.1) ...
Setting up libssl1.1:amd64 (1.1.1f-1ubuntu2.23) ...
debconf: unable to initialize frontend: Dialog
debconf: (TERM is not set, so the dialog frontend is not usable.)
debconf: falling back to frontend: Readline
debconf: unable to initialize frontend: Readline
debconf: (Can't locate Term/ReadLine.pm in @INC (you may need to install the Term::ReadLine module) (@INC contains: /etc/perl /usr/local/lib/x86_64-linux-gnu/perl/5.30.0 /usr/local/share/perl/5.30.0 /usr/lib/x86_64-linux-gnu/perl5/5.30 /usr/share/perl5 /usr/lib/x86_64-linux-gnu/perl/5.30 /usr/share/perl/5.30 /usr/local/lib/site_perl /usr/lib/x86_64-linux-gnu/perl-base) at /usr/share/perl5/Debconf/FrontEnd/Readline.pm line 7.)
debconf: falling back to frontend: Teletype
Setting up openssl (1.1.1f-1ubuntu2.23) ...
Setting up ca-certificates (20230311ubuntu0.20.04.1) ...
debconf: unable to initialize frontend: Dialog
debconf: (TERM is not set, so the dialog frontend is not usable.)
debconf: falling back to frontend: Readline
debconf: unable to initialize frontend: Readline
debconf: (Can't locate Term/ReadLine.pm in @INC (you may need to install the Term::ReadLine module) (@INC contains: /etc/perl /usr/local/lib/x86_64-linux-gnu/perl/5.30.0 /usr/local/share/perl/5.30.0 /usr/lib/x86_64-linux-gnu/perl5/5.30 /usr/share/perl5 /usr/lib/x86_64-linux-gnu/perl/5.30 /usr/share/perl/5.30 /usr/local/lib/site_perl /usr/lib/x86_64-linux-gnu/perl-base) at /usr/share/perl5/Debconf/FrontEnd/Readline.pm line 7.)
debconf: falling back to frontend: Teletype
Updating certificates in /etc/ssl/certs...
137 added, 0 removed; done.
Processing triggers for libc-bin (2.31-0ubuntu9.16) ...
Processing triggers for ca-certificates (20230311ubuntu0.20.04.1) ...
Updating certificates in /etc/ssl/certs...
0 added, 0 removed; done.
Running hooks in /etc/ca-certificates/update.d...
done.
Adding group `appgroup' (GID 101) ...
Done.
Adding system user `appuser' (UID 101) ...
Adding new user `appuser' (UID 101) with group `appgroup' ...
Not creating home directory `/home/appuser'.
Removing intermediate container 8a2b26e82dcb
 ---> 91600332f172
Step 11/16 : WORKDIR /opt
 ---> Running in b6f836b660ac
Removing intermediate container b6f836b660ac
 ---> d2921b36fed7
Step 12/16 : ARG TARGET
 ---> Running in a340fba40b30
Removing intermediate container a340fba40b30
 ---> edccdd74ba2d
Step 13/16 : COPY --from=builder /tcp_app /opt/tcp_app
 ---> fb13ff39c725
Step 14/16 : RUN chown -R appuser:appgroup /opt &&     chmod +x /opt/tcp_app
 ---> Running in 6f2bcc327173
Removing intermediate container 6f2bcc327173
 ---> 0f811c74c97b
Step 15/16 : USER appuser
 ---> Running in 8e14c6354787
Removing intermediate container 8e14c6354787
 ---> 3f868ae6fa6a
Step 16/16 : CMD ["/opt/tcp_app"]
 ---> Running in 0bb62170d6da
Removing intermediate container 0bb62170d6da
 ---> 944fe0155e9f
Successfully built 944fe0155e9f
Successfully tagged 01-tcp_project_server:latest
Building client
Step 1/16 : FROM golang:1.22 AS builder
 ---> 964d0e6869f9
Step 2/16 : WORKDIR /app
 ---> Using cache
 ---> 6105c2214b4a
Step 3/16 : COPY go.mod .
 ---> Using cache
 ---> c97866fe71fb
Step 4/16 : COPY go.sum .
 ---> Using cache
 ---> 38c0e23b21cc
Step 5/16 : RUN go mod download
 ---> Using cache
 ---> fe9e6f035367
Step 6/16 : COPY . .
 ---> Using cache
 ---> 54d0a54595ed
Step 7/16 : ARG TARGET
 ---> Using cache
 ---> f268257b9867
Step 8/16 : RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tcp_app ./cmd/tcp_${TARGET}.go
 ---> Running in cb16a63174e7
Removing intermediate container cb16a63174e7
 ---> c90ffee68c81
Step 9/16 : FROM ubuntu:20.04
 ---> 9df6d6105df2
Step 10/16 : RUN apt-get update &&     apt-get install -y --no-install-recommends ca-certificates &&     rm -rf /var/lib/apt/lists/* &&     addgroup --system appgroup &&     adduser --system --ingroup appgroup --no-create-home appuser &&     mkdir -p /opt &&     chown appuser:appgroup /opt
 ---> Using cache
 ---> 91600332f172
Step 11/16 : WORKDIR /opt
 ---> Using cache
 ---> d2921b36fed7
Step 12/16 : ARG TARGET
 ---> Using cache
 ---> edccdd74ba2d
Step 13/16 : COPY --from=builder /tcp_app /opt/tcp_app
 ---> c29290678514
Step 14/16 : RUN chown -R appuser:appgroup /opt &&     chmod +x /opt/tcp_app
 ---> Running in 24a223d66e80
Removing intermediate container 24a223d66e80
 ---> b57b792d7b32
Step 15/16 : USER appuser
 ---> Running in cbb265720115
Removing intermediate container cbb265720115
 ---> e0d4043e01ca
Step 16/16 : CMD ["/opt/tcp_app"]
 ---> Running in 244d16b4610e
Removing intermediate container 244d16b4610e
 ---> 407e868cc70b
Successfully built 407e868cc70b
Successfully tagged 01-tcp_project_client:latest
Creating tcp_server ... done
Creating tcp_client ... done
```

## 查看容器日志输出

```bash
# docker-compose logs
Attaching to tcp_client, tcp_server
tcp_client | Message received.
tcp_client |
tcp_server | Listening on localhost:3333
tcp_server | Got:  Test Request
tcp_server |
```
