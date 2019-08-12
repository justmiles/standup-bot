VERSION=$$(git describe --tags $$(git rev-list --tags --max-count=1))

build:
	docker build . -t standup-bot

run:
	docker run -it -e SSH_AUTH -e SLACK_TOKEN standup-bot

gen:
	go get github.com/abice/go-enum
	go-enum --file plugins/terraform/types.go

dev:
	# install justrun with 
	# 	go get github.com/jmhodges/justrun

	# Grab env configs
	# eval "$(get-ssm-params -path /ops/standup-bot -output shell)"
	justrun -c 'go run main.go' -delay 10000ms main.go lib/plugins/** lib/slack/*  

push:
	# Push latest
	docker tag standup-bot:latest 965579072529.dkr.ecr.us-east-1.amazonaws.com/standup-bot:latest
	docker push 965579072529.dkr.ecr.us-east-1.amazonaws.com/standup-bot:latest

	# Push version
	docker tag standup-bot:latest 965579072529.dkr.ecr.us-east-1.amazonaws.com/standup-bot:$(VERSION)
	docker push 965579072529.dkr.ecr.us-east-1.amazonaws.com/standup-bot:$(VERSION)