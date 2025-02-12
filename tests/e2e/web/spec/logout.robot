*** Settings ***
Library   SeleniumLibrary
Resource  common.resource

*** Test Cases ***
Logout
    Open Browser              ${BASE_URL}  ${BROWSER}
    Log as                    username=root  password=root
    Click Navbar Menu Item    id=menu:navbar.profile  id=link:navbar.profile.logout
    Wait Until Page Contains  Sign In  timeout=5s
    Close Browser
