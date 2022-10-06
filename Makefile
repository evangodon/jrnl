
devbuild:
	./scripts/dev-build.sh

devserver:
	./scripts/dev-server.sh

rebuild-server:
	./scripts/rebuild-server.sh

test:
	gotest -v ./...

.PHONY: devserver rebuild-server test devbuild
