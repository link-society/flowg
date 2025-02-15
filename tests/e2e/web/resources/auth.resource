*** Settings ***
Library    SeleniumLibrary
Resource   nav.resource
Variables  vars.py

*** Keywords ***
Log as
    [Arguments]               ${username}  ${password}
    Go To                     url=${BASE_URL}/web/login
    Input Text                id=input:login.username  ${username}
    Input Text                id=input:login.password  ${password}
    Click Element             id=btn:login.submit
    Wait Until Page Contains  Welcome                  timeout=5s


Logout
    Click Navbar Menu Item    id=menu:navbar.profile  id=link:navbar.profile.logout
    Wait Until Page Contains  Sign In  timeout=5s
