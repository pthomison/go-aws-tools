dist_dir := $(CURDIR)/dist
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

build:
	go build -o $(dist_dir)/aws-tools

clean:
	rm -rf $(dist_dir)

builder:
	docker build . -t go-node:latest -f ./Dockerfile.builder	

docker-go: builder
	docker run \
	-it --rm \
	-v $(mkfile_dir):/go/src/github.com/pthomison/aws-tools \
	-v $(HOME)/.aws:/root/.aws \
	-w /go/src/github.com/pthomison/aws-tools \
	go-node:latest 
	
docker-build: builder
	docker run \
	-it --rm \
	-v $(mkfile_dir):/go/src/github.com/pthomison/aws-tools \
	-w /go/src/github.com/pthomison/aws-tools \
	go-node:latest \
	make build