*** Settings ***
Library    SeleniumLibrary

*** Keywords ***
Click Navbar Menu Item
    [Arguments]                    ${menu_locator}  ${item_locator}
    Click Element                  ${menu_locator}
    Wait Until Element Is Visible  ${item_locator}  timeout=5s
    Click Element                  ${item_locator}
