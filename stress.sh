#!/bin/bash

./libclient/libclient -n 40 q 1:10 &

./libclient/libclient -n 40 b Alice@Com.com:1:2:1-1:2-1 &
./libclient/libclient -n 40 b Bob@Com.com:2:2:2-1:3-1 &
./libclient/libclient -n 40 b Copper@Com.com:3:2:1-1:3-1 &

./libclient/libclient -n 40 c Alice@Com.com:1:2:1-1:2-1 &
./libclient/libclient -n 40 c Bob@Com.com:2:2:2-1:3-1 &
./libclient/libclient -n 40 c Copper@Com.com:3:2:1-1:3-1 &
