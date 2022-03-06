BINARY = tools-home

all: $(BINARY)

$(BINARY):
	go build -ldflags "-s -w" -o bin/$@ ./cmd

clean:
	@rm -rf ./bin/*

.PHONY: clean all $(BINARY)