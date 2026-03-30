FROM alpine:latest

RUN mkdir -p /run/openrc/ && touch /run/openrc/softlevel

RUN apk add --no-cache dnsmasq openrc mdevd-openrc socat
RUN rc-update add dnsmasq default

RUN sed -i '/getty/d' /etc/inittab

COPY tests/config/dns/management-service /etc/init.d/
RUN rc-update add management-service default

COPY tests/config/dns/dns-management-server.sh .
COPY tests/config/dns/dns-management-router.sh .

EXPOSE 53
EXPOSE 8080

CMD ["/sbin/init"]
