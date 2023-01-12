devserver:
	./scripts/dev-server.sh

deployprodserver:
	./scripts/deploy-prod-server.sh

test:
	gotest -v ./...

