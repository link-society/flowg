*** Settings ***
Library   SeleniumLibrary
Resource  common.resource

*** Test Cases ***
Logout
    Open Browser  ${BASE_URL}  ${BROWSER}
    Log as        username=root  password=root
    Logout
    Close Browser
