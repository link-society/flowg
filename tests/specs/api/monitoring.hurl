###############################################################################
# DESCRIPTION: Test monitoring endpoints
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Healthcheck endpoint is available
#  - Prometheus Exporter endpoint is available

GET http://localhost:9113/health
HTTP 200

GET http://localhost:9113/metrics
HTTP 200
