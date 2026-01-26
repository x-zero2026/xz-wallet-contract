.PHONY: help build clean deploy test generate-bindings

help:
	@echo "XZ Wallet Lambda Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  generate-bindings  - Generate Go bindings from contract ABIs"
	@echo "  build             - Build all Lambda functions (for SAM)"
	@echo "  clean             - Clean build artifacts"
	@echo "  deploy            - Deploy all Lambda functions to AWS"
	@echo "  test              - Run tests"

# Generate Go bindings from contract ABIs
generate-bindings:
	@echo "Generating Go bindings from contract ABIs..."
	@mkdir -p pkg/blockchain/contracts
	@if [ ! -f ../contracts/artifacts/contracts/XZToken.sol/XZToken.json ]; then \
		echo "Error: Contract artifacts not found. Run 'npm run compile' in contracts directory first."; \
		exit 1; \
	fi
	abigen --abi ../contracts/artifacts/contracts/XZToken.sol/XZToken.json \
		--pkg contracts \
		--type XZToken \
		--out pkg/blockchain/contracts/xztoken.go
	abigen --abi ../contracts/artifacts/contracts/TaskEscrow.sol/TaskEscrow.json \
		--pkg contracts \
		--type TaskEscrow \
		--out pkg/blockchain/contracts/taskescrow.go
	@echo "✅ Go bindings generated successfully"

# SAM build targets (called by sam build)
build-GetBalanceFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/get-balance/main.go

build-CreateTaskFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/create-task/main.go

build-ListTasksFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/list-tasks/main.go

build-GetTaskFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/get-task/main.go

build-BidTaskFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/bid-task/main.go

build-SelectBidderFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/select-bidder/main.go

build-ApproveWorkFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/approve-work/main.go

build-CancelTaskFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/cancel-task/main.go

build-SubmitWorkFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap ./cmd/submit-work/main.go

# Build all Lambda functions locally
build:
	@echo "Building Lambda functions..."
	@for dir in cmd/*; do \
		if [ -d "$$dir" ]; then \
			echo "Building $$dir..."; \
			cd $$dir && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go && cd ../..; \
		fi \
	done
	@echo "✅ All Lambda functions built"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@find cmd -name "bootstrap" -delete
	@find cmd -name "*.zip" -delete
	@rm -rf .aws-sam
	@echo "✅ Clean complete"

# Deploy all Lambda functions
deploy:
	@echo "Deploying Lambda functions with SAM..."
	sam build
	sam deploy

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

