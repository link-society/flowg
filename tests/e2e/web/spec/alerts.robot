*** Settings ***
Library   DependencyLibrary
Library   SeleniumLibrary
Resource  common.resource

*** Test Cases ***
Create new alert
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.alerts
    Wait Until Page Contains       No alert found  timeout=5s
    Click Element                  id=btn:alerts.create
    Wait Until Element Is Visible  id=input:alerts.modal.name
    Input Text                     id=input:alerts.modal.name  test
    Click Button                   id=btn:alerts.modal.save
    Wait Until Page Contains       Alert created  timeout=5s
    Close Browser

Configure alert
    Depends On Test                       Create new alert
    Open Browser                          ${BASE_URL}  ${BROWSER}
    Log as                                username=root  password=root
    Click Navbar Menu Item                id=menu:navbar.settings  id=link:navbar.settings.alerts
    Wait Until Element Is Visible         id=label:alerts.list-item.test
    Element Should Be Visible             id=container:editor.alerts
    Input Text                            id=input:editor.alerts.webhook_url  http://test.com
    Input Key/Value Pair                  editor=field:editor.alerts.headers  key=Test-Header  value=Test-Value
    Wait Until Key/Value Pair Is Visible  editor=field:editor.alerts.headers  item=test-header
    Click Button                          id=btn:alerts.save
    Wait Until Page Contains              Alert saved  timeout=5s
    Remove Key/Value Pair                 editor=field:editor.alerts.headers  item=test-header
    Click Button                          id=btn:alerts.save
    Wait Until Page Contains              Alert saved  timeout=5s
    Close Browser

Delete alert
    Depends On Test                Configure alert
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.alerts
    Wait Until Element Is Visible  id=label:alerts.list-item.test
    Click Element                  id=btn:alerts.delete
    Wait Until Page Contains       Alert deleted  timeout=5s
    Page Should Contain            No alert found
    Close Browser
