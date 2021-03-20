all: build

.PHONY: build

build:
	mkdir -p build/{linux,solaris}
	GOOS=linux go build -o build/linux/ipmi_exporter main.go
	GOOS=solaris go build -o build/solaris/ipmi_exporter main.go

clean:
	rm -r build
