.PHONY: build test race lint vet vulncheck catalog check smoke clean

BINARY := codinho

build:
	go build -o $(BINARY) ./cmd/codinho

test:
	go test ./... -race

race: test

lint:
	bash scripts/ci/check-format.sh

vet:
	go vet ./...

vulncheck:
	govulncheck ./...

catalog: build
	./$(BINARY) catalog validate

# check runs everything CI runs, in the same order, so a local failure
# never surprises a pull request (requirement R4, R5).
check:
	bash scripts/ci/validate.sh

smoke: build
	./scripts/smoke-install.sh ./$(BINARY)

clean:
	rm -f $(BINARY)
