.PHONY: all mk_build_dir streamer facer skuder deploy_bin

all: mk_build_dir streamer facer skuder

mk_build_dir:
	mkdir -p build

streamer:
	go build -o build/streamer cmd/streamer/*.go

facer:
	go build -o build/facer cmd/facer/*.go

skuder:
	go build -o build/skuder cmd/skuder/*.go

deploy_bin:
	sudo cp build/streamer /usr/sbin/streamer
	sudo cp build/facer /usr/sbin/facer
	sudo cp build/skuder /usr/sbin/skuder