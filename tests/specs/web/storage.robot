*** Settings ***
Library    DependencyLibrary
Library    SeleniumLibrary

Variables  resources/vars.py

Resource   resources/auth.resource
Resource   resources/nav.resource

*** Test Cases ***
Create new stream
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.storage
    Wait Until Page Contains       No stream found  timeout=5s
    Click Element                  id=btn:streams.create
    Wait Until Element Is Visible  id=input:streams.modal.name
    Input Text                     id=input:streams.modal.name  test
    Click Button                   id=btn:streams.modal.save
    Wait Until Page Contains       Stream created  timeout=5s
    Close Browser

Configure stream
    Depends On Test                    Create new stream
    Open Browser                       ${BASE_URL}  ${BROWSER}
    Log as                             username=root  password=root
    Click Navbar Menu Item             id=menu:navbar.settings  id=link:navbar.settings.storage
    Wait Until Element Is Visible      id=label:streams.list-item.test
    Input Text                         id=input:editor.streams.retention-size  100
    Input Text                         id=input:editor.streams.retention-ttl  84600
    Input Text                         id=input:editor.streams.indexed-field.new.name  appname
    Click Button                       id=btn:editor.streams.indexed-field.new.add
    Wait Until Element Is Visible      id=input:editor.streams.indexed-field.item.appname.name
    Input Text                         id=input:editor.streams.indexed-field.new.name  tag
    Click Button                       id=btn:editor.streams.indexed-field.new.add
    Wait Until Element Is Visible      id=input:editor.streams.indexed-field.item.tag.name
    Click Button                       id=btn:editor.streams.indexed-field.item.tag.delete
    Wait Until Element Is Not Visible  id=input:editor.streams.indexed-field.item.tag
    Click Button                       id=btn:streams.save
    Wait Until Page Contains           Stream saved  timeout=5s
    Close Browser

Delete stream
    Depends On Test                Configure stream
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.storage
    Wait Until Element Is Visible  id=label:streams.list-item.test
    Click Element                  id=btn:streams.delete
    Wait Until Page Contains       No stream found  timeout=5s
    Close Browser
