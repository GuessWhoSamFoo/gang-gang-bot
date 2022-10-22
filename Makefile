VERSION ?= v0.0.1

GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

.PHONY: version
version:
	@echo $(VERSION)


.PHONY: viz
viz:
	$(GOCMD) run ./internal/tools/tools.go
	dot -Tpng fsm_create.dot -o ./internal/tools/create.png
	dot -Tpng fsm_edit.dot -o ./internal/tools/edit.png

test:
	$(GOCMD) test -race -v ./...
