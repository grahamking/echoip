# sudo docker run --rm -p 7777:7777 -p 7777:7777/udp graham/echoip
FROM debian:stable
MAINTAINER Graham King <graham@gkgk.org>
RUN mkdir -p /usr/local/echoip && chown www-data:www-data /usr/local/echoip
COPY GeoLite2-City.mmdb /usr/local/echoip/
COPY echoip /usr/local/echoip/
RUN chown www-data:www-data /usr/local/echoip/GeoLite2-City.mmdb /usr/local/echoip/echoip
USER www-data
WORKDIR /usr/local/echoip
EXPOSE 7777
CMD ["/usr/local/echoip/echoip", "-i", "eth0"]
