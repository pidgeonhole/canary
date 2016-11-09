build: lib/runner.go main/main.go

build-windows: out/runner-python-windows-amd64.exe
	go build -o out/runner-python-windows-amd64.exe main/main.go

out/runner-python-windows-amd64.exe: build

out/runner-python-linux-amd64: build

build-linux: out/runner-python-linux-amd64
	export GOOS=linux
	go build -o out/runner-python-linux-amd64 main/main.go
