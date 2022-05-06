run:
	go run main.go

build:
	docker build -t standup-bot .

docker-run:
	docker run \
	-it \
	--env-file .env \
	-e AWS_REGION \
	-e AWS_ACCESS_KEY_ID \
	-e AWS_SECRET_ACCESS_KEY \
	-e AWS_SESSION_TOKEN \
	-t standup-bot

release:
	goreleaser --rm-dist