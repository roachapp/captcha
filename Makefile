.PHONY: all
.DEFAULT_TARGET := all

all:
	CGO_ENABLED=0 go build -o bin/captcha capexample/main.go &&\
	docker build -t captcha:latest . 
