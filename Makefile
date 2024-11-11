# Setup the platform
setup:
	@echo "Generating env files..."
	cp sample.env .env
	cp sample.env test.env
	@echo ".env & test.env created. Now, update values in them"

# Migrate database
migrate:
	@echo "Running migrations..."
	sh scripts/migrate.sh

# Run pretest script
pretest:
	sh scripts/test-helper.sh

# Migrate Tests
migrate-test:
	@echo "Running migrations for tests..."
	sh scripts/migrate-tests.sh

# Run Tests
test-cover: migrate-test
	go test `go list ./... | grep -v cmd` -coverprofile=/tmp/coverage.out -coverpkg=./...
	go tool cover -html=/tmp/coverage.out

# Generate API documentation
doc:
	@echo "Generating swagger docs..."
	swag fmt --exclude ./internal/domain
	swag init --parseDependency --parseInternal -g internal/http/api/app_api.go -ot go,yaml -o internal/http/swagger

# Connect dependencies
wire:
	cd internal/dependency/ && wire && cd ../..

# Build the platform
build: migrate doc
	@echo "Building app-api..."
	sh scripts/build.sh

# Clean the platform
clean:
	@echo "Cleaning up..."
	rm ./bin/app || true
	go clean -testcache

# Stop the platform
stop:
	pkill app || true

# Start the platform
start: stop build
	nohup ./bin/app &

# Start the services on docker
start-docker:
	sh scripts/start-docker.sh

# Generate SDK
sdk-gen:
	openapi-generator-cli generate -i internal/http/swagger/swagger.yaml --generator-name typescript-axios --config ./config/codegen.json -o ../app-sdk-ts
	cp -r ./docs/sdk/templates/README.md ../app-sdk-ts/
	cd ./../app-sdk-ts; git add .; git commit -m '${app_SDK_VERSION}: Regenerated models & apis'; npm version ${app_SDK_VERSION}; rm -rf dist node_modules; npm install; npm run build; npm publish; git add .; git commit -m '${app_SDK_VERSION}'; git push origin main
