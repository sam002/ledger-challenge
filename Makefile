run:
	docker build -t ledger-accounts -f build/accounts.Dockerfile .
	docker run -p 3000:3000 --env-file='${PWD}/configs/accounts.env' -it --rm ledger-accounts
acceptance-tests:
	docker build -t ledger-tests -f build/tests.Dockerfile .
	docker run --rm -v ${PWD}/test:/test ledger-tests bash -c "robot --outputdir /test/output /test/"
