ifneq (,$(wildcard .env))
    include .env
    export
endif

build:
	sh scripts/build.sh

start:
	sh scripts/start.sh

generate:
	sh scripts/generate.sh

lint:
	sh scripts/lint.sh
