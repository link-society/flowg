*** Settings ***
Library  SeleniumLibrary

*** Test Cases ***
Login With Valid Credentials
    Open Browser  http://localhost:5080/web/login  Chrome
    Input Text    id=username    root
    Input Text    id=password    root
    Click Button  id=submit
    Wait Until Page Contains  Welcome  timeout=5s
    Close Browser
