*** Settings ***
Documentation     This is a resource file, that can contain variables and keywords.
...               Keywords defined here can be used where this Keywords.resource in loaded.
Library           RequestsLibrary
Library           DatabaseLibrary
Library           FakerLibrary
Resource          credentials.resource
Resource          accounts.resource
Resource    accounts.resource
#Library           CustomLibrary.py

*** Variables ***
${ISSUER_UUID}
${INVESTOR_UUID}
${INVESTOR_BALANCE}
${INVOICE_UUID}
${INVOICE}
${BID}

*** Keywords ***
Check running ledger infra
    ${ping_url}=    Set Variable    ${URL_LEDGER_SERVICE}/ping
    ${ping_result}=        GET     ${ping_url}
    Log     ${ping_result}
    Status Should Be    200     ${ping_result}

Connected to ledger DB
    Connect To Database Using Custom Params    psycopg2    dsn="${PG_DSN_ACCOUNT}"

Random issuerId
    ${ISSUER_UUID}=   Uuid 4
    Set Test Variable   ${ISSUER_UUID}

Random investorId
    ${INVESTOR_UUID}=   Uuid 4
    Set Test Variable   ${INVESTOR_UUID}

Investor with account
    Random userID
    ${INVESTOR_UUID}=   Create investor account by API
    Set Test Variable    ${INVESTOR_UUID}

Create correct invoice by API
    ${create_invoice_url}=   Set Variable    ${URL_LEDGER_SERVICE}/create-invoice

    ${description}=     Set Variable    TEST INVOICE
    ${volume}=      Set Variable    ${{random.randint(1,100_000)}}
    ${price}=      Set Variable    ${{random.randint(1,100_000)}}.${{random.randint(0,9)}}${{random.randint(1,9)}}
    ${request_dict}=    Create Dictionary   issuer_id=${ISSUER_UUID}  description=${description}    volume=${volume}  unit_price=${price}

    ${res}=     POST    ${create_invoice_url}    json=${request_dict}
    Status Should Be    200     ${res}
    Should Be Equal As Strings    ${description}  ${res.json()}[description]
    Should Be Equal As Strings    ${volume}  ${res.json()}[volume]
    Should Be Equal As Strings    ${price}  ${res.json()}[unit_price]
    Should Be Equal As Strings    ${ISSUER_UUID}  ${res.json()}[issuer_id]

    ${INVOICE_UUID}=    Set Variable    ${res.json()}[id]
    Set Test Variable   ${INVOICE_UUID}
    ${INVOICE}=     Set Variable    ${res.json()}
    Set Test Variable   ${INVOICE}

Create correct bid by API
    ${create_bid_url}=   Set Variable    ${URL_LEDGER_SERVICE}/create-bid

    ${volume}=      Set Variable    ${{random.randint(1,${INVOICE}[volume])}}
    ${balance}=      Evaluate    2*$volume*${INVOICE}[unit_price]
    Make deposit with amount   ${balance}
    Set Test Variable    ${INVESTOR_BALANCE}

    ${request_dict}=    Create Dictionary   investor_id=${INVESTOR_UUID}    volume=${volume}    invoice_id=${INVOICE_UUID}

    ${res}=     POST    ${create_bid_url}    json=${request_dict}
    Status Should Be    200     ${res}
    Should Be Equal As Strings    ${volume}  ${res.json()}[volume]
    Should Be Equal As Strings    ${INVESTOR_UUID}  ${res.json()}[investor_id]

    ${BID}=    Set Variable    ${res.json()}
    Set Test Variable   ${BID}

Invoce should be created in DB
    ${select}=  Set Variable    SELECT id FROM invoice WHERE id='${INVOICE_UUID}' AND issuer_id='${ISSUER_UUID}'
    Check If Exists In Database    ${select}

Bid should be created in DB
    ${select}=  Set Variable    SELECT id FROM deal WHERE id='${BID}[id]' AND investor_id='${INVESTOR_UUID}' AND invoice_id='${INVOICE_UUID}'
    Check If Exists In Database     ${select}

Balance should be decrese
    ${exp_balance}=     Evaluate    ${INVESTOR_BALANCE} - ${BID}[volume] * ${INVOICE}[unit_price]
    ${select}=  Set Variable    SELECT id FROM accounts WHERE id='${ACCOUNT_UUID}' AND balance='${exp_balance}' AND equity='${INVESTOR_BALANCE}'
    Check If Exists In Database    ${select}