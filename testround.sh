#!/bin/bash

for ((i=0; i<10; i++))
do
    rm -f tmp$1
	./round$1.sh > tmp$1
	diff -s tmp$1 round$1.out
done
