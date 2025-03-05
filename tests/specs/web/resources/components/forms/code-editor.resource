*** Settings ***
Library   SeleniumLibrary

*** Keywords ***
Input Code
    [Arguments]                       ${editor}  ${code}
    Wait Until Page Contains Element  css=div.monaco-editor  timeout=5s
    Execute Javascript                ${CURDIR}/monaco-input.js  ARGUMENTS  ${editor}  ${code}
