import socket
from datetime import datetime, timezone

_MONTHS = (
    "Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
)

# user.notice = facility(1) * 8 + severity(5), matching `logger`'s default.
_DEFAULT_PRI = 13


def _send(payload, host, port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        sock.sendto(payload.encode("utf-8"), (host, int(port)))
    finally:
        sock.close()


def send_syslog_rfc5424(message, tag="robotframework", host="localhost", port=5514):
    timestamp = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%fZ")
    hostname = socket.gethostname()
    payload = f"<{_DEFAULT_PRI}>1 {timestamp} {hostname} {tag} - - - {message}"
    _send(payload, host, port)


def send_syslog_rfc3164(message, tag="robotframework", host="localhost", port=5514):
    now = datetime.now()
    timestamp = f"{_MONTHS[now.month - 1]} {now.day:>2} {now:%H:%M:%S}"
    hostname = socket.gethostname()
    payload = f"<{_DEFAULT_PRI}>{timestamp} {hostname} {tag}: {message}"
    _send(payload, host, port)
