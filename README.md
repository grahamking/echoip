# echoip

Go server that echoes your IP address. Handy to check if you are on a VPN, or what your IP address looks like. Single UDP packet in each direction.

## Try it live

I run echoip as a service on `plebis.net` port 7777. You send it a UDP packet (containing anything), it responds with a UDP packet with your IP address as seen by remote host, and location as determined by MaxMind GeoIP Lite city database.  Feel free to use it in your scripts.

UDP (prefered):

> echo " " | nc -u plebis.net 7777

TCP (if you have to):

> echo " " | nc plebis.net 7777

TCP over TOR (assuming default tor port 9050):

> echo " " | nc -x 127.0.0.1:9050 plebis.net 7777

On Linux you can add `-q 1` to `nc` to auto-quit after a second.

## Run it yourself:

- Download [GeoLite2-City.mmdb](http://dev.maxmind.com/geoip/geoip2/geolite2/)
- Download and install [Go](http://golang.org)
- Get the geoip library: `go get github.com/oschwald/geoip2-golang`
- Build: `go build echoip.go`
- Run: `./echoip`

You may want to setup an `upstart` script. Here's my `/etc/init/echoip.conf`:

	author "Your Name <gloryglory@example.com>"
	description "Echo users IP address"

	start on static-network-up
	stop on shutdown

	console log
	respawn
	respawn limit 10 5

	setuid www-data
	setgid www-data
	chdir /usr/local/echoip

	exec /usr/local/echoip/echoip

This scripts assumes you have the `echoip` binary and `GeoLite2-City.mmdb` in directory `/usr/local/echoip`, owned by `www-data`.

### Misc

There's a million "what's my IP" websites ([duckduckgo](https://duckduckgo.com) has it right in the search engine), if you just want to casually know your IP address.

The by-committee way to do this is [STUN](https://en.wikipedia.org/wiki/Session_Traversal_Utilities_for_NAT) (echoip will suffer at least the same limitations as STUN servers).

`echoip` is intended to be something you hit from your bash prompt, or you tmux config, or even maybe a browser or gnome extension. Never be in the dark about how you face the world again!

Was there an easier way to do this? Something that already exists? Please let me know - maybe raise a github issue.

### License

Copyright 2014 Graham King

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This product includes GeoLite2 data created by MaxMind, available from
<a href="http://www.maxmind.com">http://www.maxmind.com</a>.
