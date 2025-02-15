*** Settings ***
Library   SeleniumLibrary

*** Keywords ***
Input Key/Value Pair
    [Arguments]   ${editor}  ${key}  ${value}
    Input Text    xpath=//div[@id='${editor}']//div[@data-ref='input:generic.kv-editor.new.key']//input  ${key}
    Input Text    xpath=//div[@id='${editor}']//div[@data-ref='input:generic.kv-editor.new.value']//input  ${value}
    Click Button  xpath=//div[@id='${editor}']//button[@data-ref='btn:generic.kv-editor.new.submit']


Remove Key/Value Pair
    [Arguments]                        ${editor}  ${item}
    Element Should Be Visible          xpath=//div[@id='${editor}']//div[@data-ref='entry:generic.kv-editor.item.${item}']
    Click Button                       xpath=//div[@id='${editor}']//div[@data-ref='entry:generic.kv-editor.item.${item}']//button[@data-ref='btn:generic.kv-editor.item.delete']
    Wait Until Element Is Not Visible  xpath=//div[@id='${editor}']//div[@data-ref='entry:generic.kv-editor.item.${item}']


Wait Until Key/Value Pair Is Visible
    [Arguments]                    ${editor}  ${item}
    Wait Until Element Is Visible  xpath=//div[@id='${editor}']//div[@data-ref='entry:generic.kv-editor.item.${item}']  timeout=5s
