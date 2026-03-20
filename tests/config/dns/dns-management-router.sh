#!/bin/sh
set -eu

CONF_FILE="/etc/dnsmasq.d/001-txts.conf"

read_line() {
	IFS= read -r line || return 1
	line=${line%"
"}
	printf '%s' "$line"
}

restart_server() {
	rc-service dnsmasq restart >/dev/null 2>&1

	return 0
}

add_dns_txt() {
	BODY="$(cat)"
	DOMAIN="${BODY%%:*}"
	RECORD="${BODY#*:}"

	echo "txt-record=$DOMAIN,\"$RECORD\"" >>"$CONF_FILE"

	restart_server

	return 0
}

del_dns_txt() {
	DOMAIN="$(cat)"

	sed -i "/txt-record=$DOMAIN/d" "$CONF_FILE"

	restart_server

	return 0
}

send_response() {
	status="$1"
	body="$2"

	printf 'HTTP/1.1 %s\r\n' "$status"
	printf 'Content-Type: text/plain\r\n'
	printf 'Content-Length: %s\r\n' "$(printf '%s' "$body" | wc -c | tr -d ' ')"
	printf 'Connection: close\r\n'
	printf '\r\n'
	printf '%s\r\n' "$body"

	return 0
}

request_line="$(read_line || true)"
[ -n "$request_line" ] || exit 0

# Example: POST /txt HTTP/1.1
method="$(echo "$request_line" | cut -d ' ' -f 1)"
path="$(echo "$request_line" | cut -d ' ' -f 2)"
# version="${3:-}"

content_length=0

while read -r line; do
	line="$(echo "$line" | tr -d '\r')"
	[ -z "$line" ] && break

	case "$line" in
	[Cc]ontent-[Ll]ength:*)
		content_length="$(printf '%s' "$line" | cut -d: -f2- | tr -d ' ')"
		;;

  *) ;;
	esac
done

body=""
if [ "$content_length" -gt 0 ] 2>/dev/null; then
	body="$(dd bs=1 count="$content_length" 2>/dev/null)"
fi

case "$method $path" in
"POST /txt")
	if printf '%s' "$body" | add_dns_txt; then
		send_response "200 OK" "OK\n"
	else
		send_response "500 Internal Server Error" "add_txt_record.sh failed\n"
	fi
	;;
"DELETE /txt")
	if printf '%s' "$body" | del_dns_txt; then
		send_response "200 OK" "OK\n"
	else
		send_response "500 Internal Server Error" "add_txt_record.sh failed\n"
	fi
	;;
"GET /healthcheck")
	send_response "200 OK" "OK\n"
	;;
*)
	send_response "404 Not Found" "Not Found\n"
	;;
esac
