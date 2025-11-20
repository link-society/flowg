*** Settings ***
Library    Process
Library    SeleniumLibrary
Library    RequestsLibrary

Variables  vars.py

Resource   api.resource

*** Keywords ***
Send log via API
    [Arguments]                      ${message}
    ${token}=    Get JSON Web Token  username=root  password=root
    ${payload}=  Evaluate            {"records": [{"appname": "robotframework", "message": "${message}"}]}
    API POST                         path=/api/v1/pipelines/default/logs/struct  token=${token}  body=${payload}  expected_status=200

Send log via Syslog (RFC 5424)
    [Arguments]              ${message}
    ${result}=  Run Process  logger -n localhost -P 5514 -t robotframework --rfc5424 "${message}"  shell=True


Send log via Syslog (RFC 3164)
    [Arguments]              ${message}
    ${result}=  Run Process  logger -n localhost -P 5514 -t robotframework --rfc3164 "${message}"  shell=True


Purge stream
    [Arguments]                    ${stream}
    Go To                          ${BASE_URL}/web/storage/${stream}
    Wait Until Element Is Visible  id=btn:streams.delete
    Click Element                  id=btn:streams.delete
    Sleep                          1s
