# Shadowsocks24

Shadowsocks24 is an advanced yet easy-to-use shadowsocks server and manager built on top of
the [outline-ss-server](https://github.com/Jigsaw-Code/outline-ss-server).
It exposes web-based dashboard panels and HTTP REST APIs to manage users, servers, and used traffic.

Features:
* Web-based dashboard panels to rule them all.
* Server replication to share keys on multiple servers.
* Exposing a single port for all keys (users).
* Exposing Outline SSCONF links with random load-balancing feature across available servers.
* Exposing subscription links to share all available servers with users.
* Quota (traffic limitation) management
* Public URL for users to see remaining traffic and shadowsocks keys.
* Single and multi-server friendly.
* Bridge server friendly.

# Documentation

## Installation

Clone this repository into the server:

```shell
git clone https://github.com/internet24/shadowsocks24.git
```

Pull Docker images and run:

```shell
cd shadowsocks24
docker compose pull
docker compose up -d
```

Surf the panel in your browser:

```
http://YOUR-SERVER-IP
```

The default admin username and password are:
* Username: `admin`
* Passwrod: `password`

Define keys (users) in the user tab and share the generate keys with them.
