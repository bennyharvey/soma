#!/usr/bin/env bash

for e in "/etc/soma.d/facer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    echo "stoping facer@$fn"
    sudo systemctl stop facer@$fn
  fi
done

for e in "/etc/soma.d/streamer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    echo "stoping streamer@$fn"
    sudo systemctl stop streamer@$fn
  fi
done

echo "stoping skuder"
sudo systemctl stop skuder