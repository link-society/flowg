*** Settings ***
Library    SeleniumLibrary

Variables  ../resources/vars.py

Resource   ../resources/auth.resource
Resource   ../resources/nav.resource
Resource   ../resources/logging.resource

*** Test Cases ***
View logs sent on default pipeline
    Open Browser              ${BASE_URL}  ${BROWSER}
    Log as                    username=root  password=root
    Click Element             id=link:navbar.streams
    Wait Until Page Contains  No stream found  timeout=5s
    Send log via API          hello world
    Reload Page
    Wait Until Page Contains  hello world  timeout=5s
    Page Should Contain       robotframework
    Purge stream              default
    Close Browser


View logs sent via Syslog (RFC 5424)
    Open Browser                    ${BASE_URL}  ${BROWSER}
    Log as                          username=root  password=root
    Click Element                   id=link:navbar.streams
    Wait Until Page Contains        No stream found  timeout=5s
    Send log via Syslog (RFC 5424)  hello world
    Reload Page
    Wait Until Page Contains        hello world  timeout=5s
    Page Should Contain             robotframework
    Purge stream                    default
    Close Browser


View logs sent via Syslog (RFC 3164)
    Open Browser                    ${BASE_URL}  ${BROWSER}
    Log as                          username=root  password=root
    Click Element                   id=link:navbar.streams
    Wait Until Page Contains        No stream found  timeout=5s
    Send log via Syslog (RFC 3164)  hello world
    Reload Page
    Wait Until Page Contains        hello world  timeout=5s
    Page Should Contain             robotframework
    Purge stream                    default
    Close Browser


Watch logs
    Send log via Syslog (RFC 3164)     hello world
    Open Browser                       ${BASE_URL}  Firefox
    Log as                             username=root  password=root
    Click Element                      id=link:navbar.streams
    Wait Until Page Contains           hello world  timeout=5s
    Click Element                      id=btn:streams.timewindow-selector.open
    Wait Until Element Is Visible      id=btn:streams.timewindow-selector.apply
    Click Element                      id=btn:streams.timewindow-selector.live
    Click Element                      id=btn:streams.timewindow-selector.apply
    Wait Until Element Is Not Visible  id=btn:streams.timewindow-selector.apply
    Send log via Syslog (RFC 3164)     john smith
    Wait Until Page Contains           john smith  timeout=5s
    Purge stream                       default
    Close Browser
