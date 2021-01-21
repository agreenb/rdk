
format:
	gofmt -s -w .
	clang-format -i --style="{BasedOnStyle: Google, IndentWidth: 4}" `find samples utils -iname "*.cpp" -or -iname "*.h" -or -iname "*.ino"`
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run -v ./...

docker:
	docker build -f Dockerfile.fortest -t 'echolabs/robotcoretest:latest' .
	docker push 'echolabs/robotcoretest:latest'

minirover2: 
	go build -o minirover2 samples/minirover2/control.go samples/minirover2/util.go

python-macos:
	cp etc/darwin/python-2.7.pc /usr/local/lib/pkgconfig/
