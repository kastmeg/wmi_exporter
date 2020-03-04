export GOOS=windows

build:
	promu build -v

test:
	go test -v ./...

lint:
	golangci-lint -c .golangci.yaml run

fmt:
	gofmt -l -w -s .

crossbuild:
	# The prometheus/golang-builder image for promu crossbuild doesn't exist
	# on Windows, so for now, we'll just build twice
	GOARCH=amd64 promu build --prefix=output/amd64
	GOARCH=386   promu build --prefix=output/386

installer: build
	powershell.exe installer/build.ps1 -PathToExecutable wmi_exporter.exe -Arch amd64 -Version $(shell ./wmi_exporter.exe --version  2>&1 | grep -Po "version \d+\.\d+\.\d+" | sed 's/version //g')
	mv installer/Output/wmi_exporter-*.msi .
