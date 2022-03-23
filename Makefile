test:
	go clean -testcache './...'
	go generate './...'
	go test './...'