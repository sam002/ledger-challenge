*** Settings ***
Documentation     This .robot file is a suite
...
...               Keywords are imported from the resource file
Resource          ../resources/credentials.resource
Resource          keywords.resource
Suite Setup       Connect to Server
Test Teardown     Logout User
Suite Teardown    Disconnect

*** Test Cases ***
Login with correct login
    Given User ${email} exist in DB
    When Login by ${email}
    Then Success login
    And Token is correct