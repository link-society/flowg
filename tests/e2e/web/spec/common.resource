*** Settings ***
Library   SeleniumLibrary
Library   RequestsLibrary

*** Variables ***
${BROWSER}   HeadlessFirefox
${BASE_URL}  http://localhost:5080

*** Keywords ***
API GET
    [Arguments]   ${path}  ${token}  ${expected_status}
    &{headers}=   Create Dictionary  Authorization=Bearer ${token}
    ${response}=  GET                ${BASE_URL}${path}  headers=&{headers}  expected_status=${expected_status}
    RETURN                           ${response}

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

Click Navbar Menu Item
    [Arguments]                    ${menu_locator}  ${item_locator}
    Click Element                  ${menu_locator}
    Wait Until Element Is Visible  ${item_locator}  timeout=5s
    Click Element                  ${item_locator}


Select From List
    [Arguments]                    ${list}  ${item}
    Element Should Be Visible      xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-left']//div[@id='label:generic.transfer-list.${item}']
    Click Element                  xpath=//div[@id='${list}']//input[@id='checkbox:generic.transfer-list.${item}']
    Click Element                  xpath=//div[@id='${list}']//button[@id='btn:generic.transfer-list.selected-right']
    Wait Until Element Is Visible  xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-right']//div[@id='label:generic.transfer-list.${item}']  timeout=5s
    Element Should Not Be Visible  xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-left']//div[@id='label:generic.transfer-list.${item}']


Wait Until Row Is Visible
    [Arguments]                    ${table}  ${row}
    Wait Until Element Is Visible  xpath=//div[@id='${table}']//div[@role='row'][@row-id='${row}']  timeout=5s


Remove Row
    [Arguments]                ${table}  ${row}
    Element Should Be Visible  xpath=//div[@id='${table}']//div[@role='row'][@row-id='${row}']
    Click Element              xpath=//div[@id='${table}']//div[@role='row'][@row-id='${row}']//button[@data-ref='btn:generic.tablerow.actions.delete']


Row Should Not Be Visible
    [Arguments]                    ${table}  ${row}
    Element Should Not Be Visible  xpath=//div[@id='${table}']//div[@role='row'][@row-id='${row}']
