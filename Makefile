.PHONY: lint
lint:
	golangci-lint --exclude-use-default=false --out-format tab run ./...

.PHONY: build
build:
	$(eval REVISION=$(shell sh -c "git rev-parse --short HEAD" | awk '{print $$1}'))
	@echo ">>> Current commit hash $(REVISION)"
	@echo ">>> go build -o redirective_$(REVISION)"
	@go build -o redirective_$(REVISION)
