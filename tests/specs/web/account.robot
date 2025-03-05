*** Settings ***
Library    DependencyLibrary
Library    SeleniumLibrary

Variables  resources/vars.py

Resource   resources/auth.resource
Resource   resources/nav.resource
Resource   resources/api.resource
Resource   resources/components/table.resource

*** Test Cases ***
Create and Delete Personal Access Token
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.profile  id=link:navbar.profile.account
    Wait Until Page Contains       API Tokens  timeout=5s
    Click Element                  id=btn:account.tokens.create
    Wait Until Element Is Visible  id=input:account.tokens.modal.token  timeout=5s
    Element Should Be Visible      id=input:account.tokens.modal.token_uuid
    ${token}=       Get Value      id=input:account.tokens.modal.token
    ${token_uuid}=  Get Value      id=input:account.tokens.modal.token_uuid
    Click Element                  id=btn:account.tokens.modal.done
    Wait Until Page Contains       Token created  timeout=5s
    Element Should Not Be Visible  id=input:account.tokens.modal.token_uuid
    Wait Until Page Contains       ${token_uuid}  timeout=5s
    API GET                        path=/api/v1/auth/whoami  token=pat:${token}  expected_status=200
    Remove Row                     table=table:account.tokens  row=${token_uuid}
    Wait Until Page Contains       Token deleted  timeout=5s
    Row Should Not Be Visible      table=table:account.tokens  row=${token_uuid}
    API GET                        path=/api/v1/auth/whoami  token=pat:${token}  expected_status=401
    Close Browser

Change Password
    Open Browser              ${BASE_URL}  ${BROWSER}
    Log as                    username=root  password=root
    Click Navbar Menu Item    id=menu:navbar.profile  id=link:navbar.profile.account
    Wait Until Page Contains  Account Information  timeout=5s
    Input Text                id=input:account.settings.change-password.old  root
    Input Text                id=input:account.settings.change-password.new  rootroot
    Click Button              id=btn:account.settings.change-password.submit
    Wait Until Page Contains  Password changed  timeout=5s
    Logout
    Log as                    username=root  password=rootroot
    Wait Until Page Contains  Welcome  timeout=5s
    Close Browser

Restore Password
    Depends On Test           Change Password
    Open Browser              ${BASE_URL}  ${BROWSER}
    Log as                    username=root  password=rootroot
    Click Navbar Menu Item    id=menu:navbar.profile  id=link:navbar.profile.account
    Wait Until Page Contains  Account Information  timeout=5s
    Input Text                id=input:account.settings.change-password.old  rootroot
    Input Text                id=input:account.settings.change-password.new  root
    Click Button              id=btn:account.settings.change-password.submit
    Wait Until Page Contains  Password changed  timeout=5s
    Close Browser
