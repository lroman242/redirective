.PHONY: lint
lint:
	golangci-lint --exclude-use-default=false --out-format tab run ./...

.PHONY: build
build:
	$(eval REVISION=$(shell sh -c "git rev-parse --short HEAD" | awk '{print $$1}'))
	@echo ">>> Current commit hash $(REVISION)"
	@echo ">>> go build -o redirective_$(REVISION)"
	@go build -o redirective_$(REVISION)

.PHONY: docker-scratch
docker-scratch:
	$(eval REVISION=$(shell sh -c "git rev-parse --short HEAD" | awk '{print $$1}'))
	@echo ">>> CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o redirective ."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o redirective .
	@echo "docker build -t redirective-$(REVISION) -f Dockerfile.scratch ."
	@docker build -t redirective-scratch -f Dockerfile.scratch .