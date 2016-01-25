all: build

build:
		./bin/build.sh

deploy:
		docker push harlow/nsqd-discovery
