GO=go
EXE=./dist/catelyn/bin/catelyn
LDFLAGS="-s -w"
BUILD_FLAGS=-x -v -ldflags $(LDFLAGS)
MAIN=./main.go
GO_SRC=$(wildcard ./catelyn/*.go ./catelyn/**/*.go)

all: $(EXE)

$(EXE): $(GO_SRC) $(MAIN)
	$(GO) build $(BUILD_FLAGS) -o $(EXE) $(MAIN)
