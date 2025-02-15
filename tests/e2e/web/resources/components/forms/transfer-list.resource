*** Settings ***
Library    SeleniumLibrary

*** Keywords ***
Select From List
    [Arguments]                    ${list}  ${item}
    Element Should Be Visible      xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-left']//div[@id='label:generic.transfer-list.${item}']
    Click Element                  xpath=//div[@id='${list}']//input[@id='checkbox:generic.transfer-list.${item}']
    Click Element                  xpath=//div[@id='${list}']//button[@id='btn:generic.transfer-list.selected-right']
    Wait Until Element Is Visible  xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-right']//div[@id='label:generic.transfer-list.${item}']  timeout=5s
    Element Should Not Be Visible  xpath=//div[@id='${list}']//div[@data-ref='container:generic.transfer-list.items-left']//div[@id='label:generic.transfer-list.${item}']
