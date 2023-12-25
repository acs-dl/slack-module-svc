c=run service

build:
	@go build -o bin/slack-module-svc

run: build
	@./bin/slack-module-svc $(c)