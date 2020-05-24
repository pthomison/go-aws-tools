dist_dir := $(CURDIR)/dist
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

os := darwin
platform := amd64

package_name := go-aws-tools
package_path := github.com/pthomison/$(package_name)

build:
	GOOS=$(os) GOARCH=$(platform) go build -o $(dist_dir)/$(package_name)

clean:
	rm -rf $(dist_dir)

clean-images:
	docker rmi go-node:latest

builder:
	docker build . -t go-node:latest -f ./Dockerfile.builder	

delve-auth:
	dlv debug -- --profile blue-test auth --name integ-bastion --user ec2-user --pubkey /tmp/testing_key.pub

docker-go: builder
	docker run \
	-it --rm \
	-v $(mkfile_dir):/go/src/$(package_path) \
	-v $(HOME)/.aws:/root/.aws \
	-w /go/src/$(package_path) \
	go-node:latest 
	
docker-build: builder
	docker run \
	-it --rm \
	-v $(mkfile_dir):/go/src/$(package_path) \
	-w /go/src/$(package_path) \
	go-node:latest \
	make build

docker-delve-auth: builder
	docker run \
	-it --rm \
	-v ~/.aws:/root/.aws \
	-v $(mkfile_dir):/go/src/$(package_path) \
	-w /go/src/$(package_path) \
	go-node:latest \
	make delve-auth

# For hacky dev use
try-jump-name: clean docker-build 
	./dist/go-aws-tools --profile blue-test jump --name integ-bastion

try-jump-id: clean docker-build 
	./dist/go-aws-tools --profile blue-test jump --id i-0daf40ab8c0b5eb5a

try-jump-bastion: clean docker-build 
	./dist/go-aws-tools --profile blue-test jump --name integ-delivery-k8s-worker-default --bastion integ-bastion

try-jump-fail: clean docker-build 
	./dist/go-aws-tools --profile blue-test jump --name integ-bastion --id i-deadbeef

try-auth-name: clean docker-build
	./dist/go-aws-tools --profile blue-test auth --name integ-bastion --user ec2-user --pubkey /tmp/id_rsa.pub