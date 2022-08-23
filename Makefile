
deps-update:
	go get -u ./...
	go get -t -u ./...
deps-check:
	go list -u -m all
test-all:
	go test all