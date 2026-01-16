# Yannena

BINDIR = bin

PROJECT_TESTS = ./internal/socket \


.PHONY: all engine webserver clean test

all: engine webserver

engine: | $(BINDIR)
	go build -buildvcs=false -o $(BINDIR)/engine ./cmd/engine

webserver: | $(BINDIR)
	go build -buildvcs=false -o $(BINDIR)/webserver ./cmd/webserver

$(BINDIR):
	@mkdir -p $(BINDIR)

clean:
	@rm -rf $(BINDIR)

test:
	@go test -v $(PROJECT_TESTS)
