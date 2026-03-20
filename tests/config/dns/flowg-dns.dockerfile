FROM linksociety/flowg:localdev

RUN apk add --no-cache curl

COPY tests/config/dns/dns-management-shell.sh /

ENV FLOWG_CLUSTER_FORMATION_STRATEGY="dns"
ENV FLOWG_CLUSTER_FORMATION_DNS_SCRIPT="/dns-management-shell.sh"
ENV FLOWG_CLUSTER_FORMATION_DNS_DOMAIN="example.com"