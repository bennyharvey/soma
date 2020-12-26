#!/usr/bin/env bash

services=''

for e in "/etc/soma.d/facer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    services="$services -u facer@$fn"
    # sudo systemctl start facer@$fn
  fi
done

for e in "/etc/soma.d/streamer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    services="$services -u streamer@$fn"
    # sudo systemctl start streamer@$fn
  fi
done

services="$services -u skuder"

journalctl $services -f