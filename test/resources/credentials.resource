*** Variables ***
${URL_ACCOUNTS_SERVICE}=      %{URL=http://127.0.0.1:3000}
${URL_LEDGER_SERVICE}=      %{URL=http://127.0.0.1:3001}
${EMAIL}    %{USER=test@mail.com}
${PATH}     %{PATH=/}
${PG_DSN_ACCOUNT}   %{PG_DSN_ACCOUNT=host=172.17.0.1 user=postgres password=password dbname=ledger port=5432}
