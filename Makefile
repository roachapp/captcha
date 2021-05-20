.PHONY: all test
.DEFAULT_TARGET := all

all: test
	CGO_ENABLED=0 go build -o bin/captcha capexample/main.go &&\
	docker build -t captcha:latest . 

# postgres related stuff
# read https://hub.docker.com/_/postgres for insights
PGPW=captcha
PGUSR=captcha
PGDB=captcha
PGPORT=5432
DATABASE_URL=postgres://$(PGUSR):$(PGPW)@localhost:$(PGPORT)/$(PGDB)
PSQL_IMAGE=postgres:12.7-alpine
PSQL_RUNOPTS=\
	--rm\
	--name $(PGDB)-postgres\
	-d\
	-e POSTGRES_USER=$(PGUSR)\
	-e POSTGRES_PASSWORD=$(PGPW)\
	-e POSTGRES_DB=$(PGDB)\
	-p $(PGPORT):5432
PSQL_RUNARGS=\
	-c 'listen_addresses="*"'

# bind address already in use error: https://stackoverflow.com/questions/38249434/docker-postgres-failed-to-bind-tcp-0-0-0-05432-address-already-in-use
test:
	docker run $(PSQL_RUNOPTS) $(PSQL_IMAGE) $(PSQL_RUNARGS) &&\
	sleep 3 &&\
	psql -f /home/kashim/go/src/github.com/roachapp/db/postgres/captchas_schema.sql $(DATABASE_URL) &&\
	DATABASE_URL=$(DATABASE_URL) go test ./...
	docker kill $(PGDB)-postgres

kill:
	docker kill $(PGDB)-postgres
