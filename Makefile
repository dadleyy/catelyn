GO=go
EXE=./dist/catelyn/bin/catelyn
LDFLAGS="-s -w"
BUILD_FLAGS=-x -v -ldflags $(LDFLAGS)
MAIN=./main.go
GO_SRC=$(wildcard ./catelyn/*.go ./catelyn/**/*.go)
GLIDE=glide
LINT=golint
LINT_FLAGS=-set_exit_status
VET=$(GO) vet
MISSPELL=misspell

all: $(EXE)

$(EXE): $(GO_SRC) $(MAIN) vendor
	$(GO) build $(BUILD_FLAGS) -o $(EXE) $(MAIN)


test: vendor
	$(MISSPELL) -error $(GO_SRC)
	$(LINT) $(LINT_FLAGS) ./catelyn/...
	$(GO) vet ./catelyn/...

vendor:
	$(GO) get -v -u github.com/golang/lint/golint
	$(GO) get -v -u github.com/client9/misspell/cmd/misspell
	$(GO) get -u -v github.com/Masterminds/glide
	$(GLIDE) install
