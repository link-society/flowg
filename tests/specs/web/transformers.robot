*** Settings ***
Library    DependencyLibrary
Library    SeleniumLibrary

Variables  resources/vars.py

Resource   resources/auth.resource
Resource   resources/nav.resource
Resource   resources/components/forms/kv-editor.resource
Resource   resources/components/forms/code-editor.resource

*** Test Cases ***
Create new transformer
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.transformers
    Wait Until Page Contains       No transformer found  timeout=5s
    Click Element                  id=btn:transformers.create
    Wait Until Element Is Visible  id=input:transformers.modal.name
    Input Text                     id=input:transformers.modal.name  test
    Click Button                   id=btn:transformers.modal.save
    Wait Until Page Contains       Transformer created  timeout=5s
    Close Browser

Configure transformer
    Depends On Test                Create new transformer
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.transformers
    Wait Until Element Is Visible  id=label:transformers.list-item.test
    Input Code                     editor=monaco:transformers.editor  code=.foo = "bar"
    Input Key/Value Pair           editor=kv:transformers.test.record  key=message  value=test
    Click Button                   id=btn:transformers.test.run
    ${result}=  Get Text           id=container:transformers.test.result
    ${payload}=  Evaluate          json.loads('''${result}''')
    Should Be Equal As Strings     ${payload['foo']}  bar
    Should Be Equal As Strings     ${payload['message']}  test
    Click Button                   id=btn:transformers.save
    Wait Until Page Contains       Transformer saved  timeout=5s
    Close Browser

Delete transformer
    Depends On Test                Configure transformer
    Open Browser                   ${BASE_URL}  ${BROWSER}
    Log as                         username=root  password=root
    Click Navbar Menu Item         id=menu:navbar.settings  id=link:navbar.settings.transformers
    Wait Until Element Is Visible  id=label:transformers.list-item.test
    Click Element                  id=btn:transformers.delete
    Wait Until Page Contains       Transformer deleted  timeout=5s
    Page Should Contain            No transformer found
    Close Browser
