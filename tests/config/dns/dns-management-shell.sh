#!/bin/sh

echo "$@"

echo "FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER=$FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER"

case "$1" in
"set")
	curl --location "$FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER/txt" \
		--header 'Content-Type: text/plain' \
		--data "$3:$4"
	;;
"del")
	curl --location --request DELETE "$FLOWG_CLUSTER_FORMATION_DNS_MANAGEMENT_SERVER/txt" \
		--header 'Content-Type: text/plain' \
		--data "$3$4"
	;;
*)
	echo "Unknown command '$1'"
	;;
esac
