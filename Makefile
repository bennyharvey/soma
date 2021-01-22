.PHONY: all mk_build_dir streamer facer skuder deploy_bin help mk_conf_dir deploy_services
.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: install streamer facer skuder ## Build all binaries to "./deploy/build" dir
	
install: mk_conf_dir mk_build_dir ## Prepare dirs

rebuild: stop all deploy_bin start ## rebuild, deploy all binaries and restart server

restart: stop start ## restart all services

conf_unsafe: stop clear_deployed_conf deploy_conf start ## !! hard redeploy all conf files

start: ## start
	@bash deploy/scripts/start.sh

stop: ## stop
	@bash deploy/scripts/stop.sh

status: ## status
	@systemctl list-units --type=service | awk '/.*\.service/ {print $1}' | grep 'facer\|streamer\|skuder'

log: ##log
	journalctl -u skuder -u facer@test110 | less

logf: ##log follow
	@bash deploy/scripts/logf.sh

streamer: ## Build streamer binary
	go build -o deploy/build/streamer cmd/streamer/*.go

facer: ## Build facer binary
	go build -o deploy/build/facer cmd/facer/*.go

skuder: ## Build skuder binary
	go build -o deploy/build/skuder cmd/skuder/*.go

deploy: deploy_bin deploy_conf deploy_services ## Deploy all

deploy_bin: ## Copy built binaries to /usr/bin/* 
	@sudo cp deploy/build/streamer /usr/bin/streamer
	@sudo cp deploy/build/facer /usr/bin/facer
	@sudo cp deploy/build/skuder /usr/bin/skuder

deploy_conf: ## Copy files from deploy/configs/* to /etc/soma.d/* 
	@sudo mkdir -p /etc/soma.d/facer
	@sudo mkdir -p /etc/soma.d/streamer
	@sudo mkdir -p /etc/soma.d/skuder
	@sudo cp -R deploy/configs/facer/*.conf /etc/soma.d/facer/
	@sudo cp -R deploy/configs/streamer/*.conf /etc/soma.d/streamer/
	@sudo cp -R deploy/configs/skuder/*.conf /etc/soma.d/skuder/

deploy_services: ## Copy systemd service files to /etc/systemd/system/
	@sudo cp -R deploy/systemd/. /etc/systemd/system/

clear_deployed_conf: ## !! Delete all files in /etc/soma.d/ 
	@sudo rm -rvf /etc/soma.d/*

mk_build_dir:
	@mkdir -p deploy/build

mk_conf_dir:
	@mkdir -p deploy/configs/facer
	@mkdir -p deploy/configs/streamer
	@mkdir -p deploy/configs/skuder

mk_conf_deploy_dir:
	@sudo mkdir -p /etc/soma.d

secret:
	@echo "no"