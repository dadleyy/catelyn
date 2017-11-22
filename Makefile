GO=go
EXE=./dist/catelyn/bin/catelyn
LDFLAGS="-s -w"
BUILD_FLAGS=-x -v -ldflags $(LDFLAGS)
MAIN=./main.go
SRC_DIR=./catelyn
GO_SRC=$(wildcard $(SRC_DIR)/*.go $(SRC_DIR)/**/*.go)
GLIDE=glide
LINT=golint
TEST_FLAGS=-v -covermode=atomic
LINT_FLAGS=-set_exit_status
VET=$(GO) vet
MISSPELL=misspell

all: $(EXE)

$(EXE): $(GO_SRC) $(MAIN) vendor
	$(GO) build $(BUILD_FLAGS) -o $(EXE) $(MAIN)


test: vendor
	$(MISSPELL) -error $(GO_SRC)
	$(LINT) $(LINT_FLAGS) $(SRC_DIR)/...
	$(GO) vet $(SRC_DIR)/...
	$(GO) test $(TEST_FLAGS) $(SRC_DIR)/...

vendor:
	$(GO) get -v -u github.com/golang/lint/golint
	$(GO) get -v -u github.com/client9/misspell/cmd/misspell
	$(GO) get -u -v github.com/Masterminds/glide
	$(GLIDE) install
