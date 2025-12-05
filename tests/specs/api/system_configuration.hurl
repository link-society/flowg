###############################################################################
# DESCRIPTION: Test system configuration through API
###############################################################################

# -----------------------------------------------------------------------------
# TEST:
#  - Verify default configuration

GET http://localhost:5080/api/v1/system-configuration
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.configuration.syslog_allowed_origins" == null

# -----------------------------------------------------------------------------
# TEST:
#  - Test setting configuration

PUT http://localhost:5080/api/v1/system-configuration
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "syslog_allowed_origins": [
    "192.168.1.1"
  ]
}
HTTP 200
[Asserts]
jsonpath "$.success" == true

# -----------------------------------------------------------------------------
# TEST:
#  - Verify updated configuration

GET http://localhost:5080/api/v1/system-configuration
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.configuration.syslog_allowed_origins" count == 1
jsonpath "$.configuration.syslog_allowed_origins[0]" == "192.168.1.1"
