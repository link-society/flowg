*** Settings ***
Library  SeleniumLibrary

*** Variables ***
${BROWSER}  HeadlessFirefox

*** Test Cases ***
Login With Valid Credentials
    Open Browser  http://localhost:5080/web/login  ${BROWSER}
    Input Text    id=username    root
    Input Text    id=password    root
    Click Button  id=submit
    Wait Until Page Contains  Welcome  timeout=5s
    Close Browser

Login With Invalid Credentials
    Open Browser  http://localhost:5080/web/login  ${BROWSER}
    Input Text    id=username    root
    Input Text    id=password    notroot
    Click Button  id=submit
    Wait Until Page Contains  Invalid credentials  timeout=5s
    Close Browser
