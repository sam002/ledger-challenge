*** Settings ***
Documentation     This .robot file is a suite
...
...               Keywords are imported from the resource file
Resource          ../resources/credentials.resource
Resource          ../resources/accounts.resource
Suite Setup       Run keywords
...     Check running infra
...     AND       Connected to DB
#Test Teardown     Clean test data
Suite Teardown    Disconnect From Database

*** Test Cases ***
Create issuer account with correct user_id
    Given Random userID
    When Create issuer account by API
    Then Acoount should be created in DB

Create invocer account with correct user_id
    Given Random userID
    When Create investor account by API
    Then Acoount should be created in DB

Make deposit for new account
    Given Random userID
        AND Create issuer account by API
    When Make deposit with correct amount
    Then Acoount should be created in DB
        AND Balance should be increased

#Make deposit for account with positive balance
#    Given New account
#    When Make deposit with correct amount twice
#    Then Balance should be increased
#
#Make dublicate deposit for account with positive balance
#    Given Create new account is succesful
#    When Make double deposit with correct amount at some time
#    Then Balance should be increased
#
#Make widthrawal for new account
#
#Make widthrawal for account with positive balance
#
#Skip one of dublicate widthrawal for account

