*** Settings ***
Library    DependencyLibrary
Library    SeleniumLibrary

Variables  resources/vars.py

Resource   resources/auth.resource
Resource   resources/nav.resource
Resource   resources/components/forms/kv-editor.resource
Resource   resources/components/forms/transfer-list.resource

*** Test Cases ***
Create new forwarder
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.forwarders
    Wait Until Page Contains       No forwarder found  timeout=5s
    Click Element                  id=btn:forwarders.create
    Wait Until Element Is Visible  id=input:forwarder.modal.name
    Input Text                     id=input:forwarder.modal.name  test
    Element Should Be Visible      id=container:editor.forwarders.http
    Input Text                     id=input:editor.forwarders.http.webhook_url  http://test.com
    Click Button                   id=btn:forwarder.modal.save
    Wait Until Page Contains       Forwarder created  timeout=5s
    Close Browser

Configure forwarder
    Depends On Test                       Create new forwarder
    Open Browser                          ${BASE_URL}  ${BROWSER}
    Log as                                username=root  password=root
    Click Navbar Menu Item                id=menu:navbar.settings  id=link:navbar.settings.forwarders
    Wait Until Element Is Visible         id=label:forwarders.list-item.test
    Element Should Be Visible             id=container:editor.forwarders.http
    Input Text                            id=input:editor.forwarders.http.webhook_url  http://test.com
    Input Key/Value Pair                  editor=field:editor.forwarders.http.headers  key=Test-Header  value=Test-Value
    Wait Until Key/Value Pair Is Visible  editor=field:editor.forwarders.http.headers  item=test-header
    Click Button                          id=btn:forwarders.save
    Wait Until Page Contains              Forwarder saved  timeout=5s
    Remove Key/Value Pair                 editor=field:editor.forwarders.http.headers  item=test-header
    Click Button                          id=btn:forwarders.save
    Wait Until Page Contains              Forwarder saved  timeout=5s
    Close Browser

Delete forwarder
    Depends On Test                Configure forwarder
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.forwarders
    Wait Until Element Is Visible  id=label:forwarders.list-item.test
    Click Element                  id=btn:forwarders.delete
    Wait Until Page Contains       Forwarder deleted  timeout=5s
    Page Should Contain            No forwarder found
    Close Browser
