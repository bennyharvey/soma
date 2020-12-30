#!/usr/bin/env bash

INPUT="$(dirname `pwd`)/configs/_import/streams.csv"
OLDIFS=$IFS
IFS=','
[ ! -f $INPUT ] && { echo "$INPUT file not found"; exit 99; }
# while read field1 field2 field3
while read f1 || [ -n "$f1" ]
do
	# echo "field1 : $field1"
	# echo "field2 : $field2"
	# echo "field3 : $field3"
	echo "f1 : $f1"
done < $INPUT
IFS=$OLDIFS