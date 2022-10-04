
devserver:
	./scripts/dev-server.sh

rebuild-server:
	./scripts/rebuild-server.sh

test:
	go test -v ./...

.PHONY: devserver rebuild-server
