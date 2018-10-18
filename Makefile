# GO parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

DEPCMD=dep
DEPRUN=$(DEPCMD) ensure

all: clean build test
full: full_clean full_build

full_build: deps build 
deps:
	$(DEPRUN)
build: format
	$(GOBUILD) .
# Build static image with no outwards deps. 
build_static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix nocgo -ldflags '-w -extldflags "-static"' -v .
docker_build:
	docker build -t aa/gochromecast .

test:
	$(GOTEST) -v ./...
test_local:
	$(GOTEST) -v ./... -tags=local

full_clean: dep_clean clean
dep_clean: 
	rm -rf ./vendor
clean:
	$(GOCLEAN)

format:
	$(GOFMT) ./...

