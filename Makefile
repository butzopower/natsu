test:
	go clean -testcache './...'
	go test './...'