testconfig:
	./scripts/testconfig.sh

test:
	go test ./...

build:
	./scripts/build.sh

run:
	./bin/twonicornd -c ./config.yml
