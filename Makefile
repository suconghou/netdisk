arm:
	CGO_ENABLED=0 GOARM=7 GOOS=linux GOARCH=arm  go build -v -o disk_arm -a -ldflags "-s -w" disk.go
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -v -o disk_linux -a -ldflags "-s -w" disk.go
windows32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386  go build -v -o disk_32.exe -a -ldflags "-s -w" disk.go  
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -v -o disk.exe -a -ldflags "-s -w" disk.go  
dev:
	go build disk.go
