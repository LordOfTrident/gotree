BIN = ./bin
APP = $(BIN)/app
CMD = ./cmd/gotree

GO = go

compile:
	$(GO) build -o $(APP) $(CMD)

$(BIN):
	mkdir -p $(BIN)

run:
	$(GO) run $(CMD)

install:
	cp $(APP) /usr/bin/gotree

clean:
	rm -r $(BIN)/*

all:
	@echo compile, run, install, clean
