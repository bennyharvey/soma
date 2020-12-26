#!/usr/bin/env bash

for e in "/etc/soma.d/facer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    echo "starting facer@$fn"
    sudo systemctl start facer@$fn
  fi
done

for e in "/etc/soma.d/streamer"/* 
do
  if [ -f "$e" ];then
    bn=$(basename -- "$e")
    ext="${bn##*.}"
    fn="${bn%.*}"
    echo "starting streamer@$fn"
    sudo systemctl start streamer@$fn
  fi
done

echo "starting skuder"
sudo systemctl start skuder