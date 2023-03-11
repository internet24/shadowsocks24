# Shadowsocks24

Shadowsocks24 is an advanced yet easy-to-use shadowsocks server and manager built on top of
the [outline-ss-server](https://github.com/Jigsaw-Code/outline-ss-server).
It exposes web-based dashboard panels and HTTP REST APIs to manage users, servers, and data traffic.

Features:
* Web-based dashboard panels to rule them all.
* Server replication to share keys (users) on multiple servers.
* Exposing a single port for all keys (users).
* Exposing Outline SSCONF links with load-balancing feature across available servers.
* Exposing subscription links to share all available servers with users.
* Quota (traffic limitation) management
* Public URL for users to see remaining traffic and shadowsocks keys.
* Single and multi-server friendly.
* Bridge server friendly.

# Documentation

## Requirements

To run the Shadowsocks24 on a server, you need the following requirements.
* A Linux distribution (Other operating systems are not supported)
* Docker and Docker Compose
* Git client
* 1 GB of RAM
* 1 core of CPU

## Installation

It's so easy to install Shadowsocks24 on your Linux server.
Follow these steps:

1. Clone this repository: `git clone https://github.com/internet24/shadowsocks24.git`
1. Run docker compose in the `shadowsocks24` directory: `docker compose up -d`
1. Surf `http://YOUR-SERVER-IP` in your web browser to see dashboard.
1. Login using `admin` as the username and `password` as the password.
1. Define keys (users) in the user tab and share the generate keys with them.
