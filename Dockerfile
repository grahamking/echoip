# sudo docker run --rm -p 7777:7777 -p 7777:7777/udp graham/echoip
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /usr/local/echoip && chown www-data /usr/local/echoip
ADD http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz /usr/local/echoip/
RUN gunzip /usr/local/echoip/GeoLite2-City.mmdb.gz && chown www-data /usr/local/echoip/GeoLite2-City.mmdb
USER www-data
WORKDIR /usr/local/echoip
COPY echoip /usr/local/echoip/
EXPOSE 7777
CMD ["/usr/local/echoip/echoip", "-i", "eth0"]
