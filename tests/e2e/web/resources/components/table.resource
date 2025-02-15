*** Settings ***
Library    SeleniumLibrary

*** Keywords ***
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
