.PHONY: all clean

OUTPUT= cctool

all: clean
	go build -o bin/${OUTPUT} main.go

clean:
	rm -f bin/${OUTPUT}