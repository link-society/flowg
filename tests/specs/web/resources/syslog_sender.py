"""Cross-platform syslog UDP sender for the web e2e tests.

The Robot suites originally shelled out to the util-linux ``logger`` command,
whose ``-n`` / ``-P`` / ``--rfc5424`` / ``--rfc3164`` flags do not exist on the
BSD ``logger`` shipped with macOS (its only options are ``[-is] [-f file]
[-p pri] [-t tag]``). Sending the datagram directly over UDP keeps the tests
working regardless of the host OS.
"""

import socket
from datetime import datetime, timezone

# Fixed English month abbreviations so the RFC 3164 timestamp never depends on
# the host locale (the go-syslog parser expects Go's ``time.Stamp`` layout,
# i.e. "Jan _2 15:04:05" with a space-padded day).
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
    # <PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID STRUCTURED-DATA MSG
    payload = f"<{_DEFAULT_PRI}>1 {timestamp} {hostname} {tag} - - - {message}"
    _send(payload, host, port)


def send_syslog_rfc3164(message, tag="robotframework", host="localhost", port=5514):
    now = datetime.now()
    timestamp = f"{_MONTHS[now.month - 1]} {now.day:>2} {now:%H:%M:%S}"
    hostname = socket.gethostname()
    # <PRI>TIMESTAMP HOSTNAME TAG: MSG
    payload = f"<{_DEFAULT_PRI}>{timestamp} {hostname} {tag}: {message}"
    _send(payload, host, port)
