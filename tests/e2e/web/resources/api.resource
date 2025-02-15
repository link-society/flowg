*** Settings ***
Library    RequestsLibrary
Variables  vars.py

*** Keywords ***
API GET
    [Arguments]   ${path}  ${token}  ${expected_status}
    &{headers}=   Create Dictionary  Authorization=Bearer ${token}
    ${response}=  GET                ${BASE_URL}${path}  headers=&{headers}  expected_status=${expected_status}
    RETURN                           ${response}
