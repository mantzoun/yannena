# Yannena

BINDIR = bin

PROJECT_TESTS = ./internal/socket \


.PHONY: all engine webserver clean test

all: engine webserver

engine: | $(BINDIR)
	go build -o $(BINDIR)/engine cmd/engine/main.go

webserver: | $(BINDIR)
	go build -o $(BINDIR)/webserver cmd/webserver/main.go

$(BINDIR):
	@mkdir -p $(BINDIR)

clean:
	@rm -rf $(BINDIR)

test:
	@go test $(PROJECT_TESTS)
