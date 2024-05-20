GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: 
	go install

test: build
	go test -cover ./rabbitmq

testacc: build
	scripts/testacc.sh

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: build test testacc vet
