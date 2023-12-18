
.PHONY: build
build:
	go build -o bin/releasepost .

.PHONY: test
test:
	./bin/releasepost --config config.example.yaml

clean:
	rm -Rf dist
