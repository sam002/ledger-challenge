*** Settings ***
Documentation     This is a resource file, that can contain variables and keywords.
...               Keywords defined here can be used where this Keywords.resource in loaded.
Library           RequestsLibrary
Library           DatabaseLibrary
Library           FakerLibrary
Resource          credentials.resource
#Library           CustomLibrary.py

*** Variables ***
${USER_UUID}
${ACCOUNT_UUID}

*** Keywords ***
Check running infra
    ${ping_url}=    Set Variable    ${URL_ACCOUNTS_SERVICE}/ping
    ${ping_result}=        GET     ${ping_url}
    Log     ${ping_result}
    Status Should Be    200     ${ping_result}

Connected to DB
    Connect To Database Using Custom Params    psycopg2    dsn="${PG_DSN_ACCOUNT}"

Disconnect From DB
    Disconnect From Database

Random userID
    ${USER_UUID}=   Uuid 4
    Set Test Variable   ${USER_UUID}

Create issuer account by API
    ${create_issuer_acc_url}=   Set Variable    ${URL_ACCOUNTS_SERVICE}/create-issuer-account

    ${request_dict}=    Create Dictionary   user_id=${USER_UUID}
    ${request}=     Set Variable     ${request_dict}

    ${res}=     POST    ${create_issuer_acc_url}    json=${request_dict}
    Status Should Be    200     ${res}
    Should Be Equal As Strings    0  ${res.json()}[balance]
    Should Be Equal As Strings    0  ${res.json()}[equity]
    Should Be Equal As Strings    ${USER_UUID}  ${res.json()}[user_id]
#    Should Be Equal As Strings    ""  ${res.json()}[vat]

    ${ACCOUNT_UUID}=    Set Variable    ${res.json()}[id]
    Set Test Variable   ${ACCOUNT_UUID}

Create investor account by API
    ${create_issuer_acc_url}=   Set Variable    ${URL_ACCOUNTS_SERVICE}/create-investor-account

    ${request_dict}=    Create Dictionary   user_id=${USER_UUID}
    ${request}=     Set Variable     ${request_dict}

    ${res}=     POST    ${create_issuer_acc_url}    json=${request_dict}
    Status Should Be    200     ${res}
    Should Be Equal As Strings    0  ${res.json()}[balance]
    Should Be Equal As Strings    0  ${res.json()}[equity]
    Should Be Equal As Strings    ${USER_UUID}  ${res.json()}[user_id]
#    Should Be Equal As Strings    ""  ${res.json()}[company]

    ${ACCOUNT_UUID}=     Set Variable    ${res.json()}[id]
    Set Test Variable   ${ACCOUNT_UUID}

Acoount should be created in DB
    Check If Exists In Database    SELECT id FROM accounts WHERE id='${ACCOUNT_UUID}'

Make deposit with correct amount
    ${deposit_url}=   Set Variable    ${URL_ACCOUNTS_SERVICE}/deposit-account
    ${AMOUNT}=      Set Variable    ${{random.randint(1,100_000_000)}}.${{random.randint(0,9)}}${{random.randint(0,9)}}
    Set Test variable   ${AMOUNT}
    ${request_dict}=    Create Dictionary   account_id=${ACCOUNT_UUID}  amount=${amount}
    ${request}=     Set Variable     ${request_dict}

    ${res}=     POST    ${deposit_url}    json=${request_dict}
    Status Should Be    200     ${res}

Balance should be increased
    Check If Exists In Database    SELECT id FROM accounts WHERE id='${ACCOUNT_UUID}' AND balance='${AMOUNT}'


