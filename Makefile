test:
	go clean -testcache './...'
	go generate './...'
	go test './...'

clean:
	-rm -r ./cli/build

build-cli: clean
	-mkdir -p ./cli/build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./cli/build/natsu.linux-amd64  .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ./cli/build/natsu.darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o ./cli/build/natsu.darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ./cli/build/natsu.windows-amd64.exe .
	cd ./cli/build && find . -name 'natsu*' | xargs -I{} tar czf {}.tar.gz {}
	cd ./cli/build && shasum -a 256 * > sha256sum.txt
	cat ./cli/build/sha256sum.txt

# example: make release V=0.0.0
release:
	git tag v$(V)
	@read -p "Press enter to confirm and push to origin ..." && git push origin v$(V)
