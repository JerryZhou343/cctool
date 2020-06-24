.PHONY: all clean

OUTPUT= cctool

all: clean
	go build -o bin/${OUTPUT} cmd/main.go

clean:
	rm -f bin/${OUTPUT}


windows:
	GOOS=windows GOARCH=amd64 go build -o bin/cctool.exe cmd/main.go