*** Settings ***
Library   DependencyLibrary
Library   SeleniumLibrary
Library   RequestsLibrary
Resource  common.resource

*** Test Cases ***
Create test role
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.profile  id=link:navbar.profile.admin
    Wait Until Page Contains       Roles
    Click Element                  id=btn:admin.roles.create
    Wait Until Element Is Visible  id=input:admin.roles.modal.name
    Input Text                     id=input:admin.roles.modal.name  test
    Select From List               list=field:admin.roles.modal.scopes  item=read_pipelines
    Click Element                  id=btn:admin.roles.modal.submit
    Wait Until Page Contains       Role created  timeout=5s
    Wait Until Row Is Visible      table=table:admin.roles  row=test
    Close Browser

Create test user
    Depends On Test                Create test role
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.profile  id=link:navbar.profile.admin
    Wait Until Page Contains       Users
    Click Element                  id=btn:admin.users.create
    Wait Until Element Is Visible  id=input:admin.users.modal.username
    Input Text                     id=input:admin.users.modal.username  test
    Input Text                     id=input:admin.users.modal.password  test
    Select From List               list=field:admin.users.modal.roles  item=test
    Click Element                  id=btn:admin.users.modal.submit
    Wait Until Page Contains       User created  timeout=5s
    Wait Until Row Is Visible      table=table:admin.users  row=test
    Close Browser

Log as test user
    Depends On Test  Create test user
    Open Browser     ${BASE_URL}  ${BROWSER}
    Log as           username=test  password=test
    Close Browser
