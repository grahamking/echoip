# sudo docker run --rm -p 7777:7777 -p 7777:7777/udp graham/echoip
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /usr/local/echoip && chown www-data /usr/local/echoip
USER www-data
WORKDIR /usr/local/echoip
COPY GeoLite2-City.mmdb /usr/local/echoip/
COPY echoip /usr/local/echoip/
CMD ["/usr/local/echoip/echoip", "-i", "eth0"]
