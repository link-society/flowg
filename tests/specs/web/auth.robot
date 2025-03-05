*** Settings ***
Library    SeleniumLibrary
Variables  resources/vars.py
Resource   resources/auth.resource

*** Test Cases ***
Login With Valid Credentials
    Open Browser  ${BASE_URL}  ${BROWSER}
    Log as        username=root  password=root
    Close Browser


Login With Invalid Credentials
    Open Browser              ${BASE_URL}/web/login  ${BROWSER}
    Input Text                id=input:login.username  root
    Input Text                id=input:login.password  notroot
    Click Button              id=btn:login.submit
    Wait Until Page Contains  Invalid credentials      timeout=5s
    Close Browser


Logout
    Open Browser  ${BASE_URL}  ${BROWSER}
    Log as        username=root  password=root
    Logout
    Close Browser
