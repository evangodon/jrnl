
devbuild:
	./scripts/dev-build.sh

devserver:
	./scripts/dev-server.sh

deployprod:
	./scripts/deploy-prod.sh

test:
	gotest -v ./...

.PHONY: devserver deployprod test devbuild
