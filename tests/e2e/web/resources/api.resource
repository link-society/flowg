*** Settings ***
Library    RequestsLibrary
Variables  vars.py

*** Keywords ***
Get JSON Web Token
    [Arguments]                 ${username}  ${password}
    ${payload}=   Evaluate      {"username": "${username}", "password": "${password}"}
    ${response}=  POST          ${BASE_URL}/api/v1/auth/login  json=${payload}  expected_status=200
    ${token}=     Set Variable  ${response.json()["token"]}
    RETURN                      ${token}


API GET
    [Arguments]                      ${path}  ${token}  ${expected_status}
    &{headers}=   Create Dictionary  Authorization=Bearer ${token}
    ${response}=  GET                ${BASE_URL}${path}  headers=&{headers}  expected_status=${expected_status}
    RETURN                           ${response}


API POST
    [Arguments]                      ${path}  ${token}  ${body}  ${expected_status}
    &{headers}=   Create Dictionary  Authorization=Bearer ${token}
    ${response}=  POST               ${BASE_URL}${path}  json=${body}  headers=&{headers}  expected_status=${expected_status}
    RETURN                           ${response}
