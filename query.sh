#!/bin/bash

./libclient/libclient a 1-1:1:2:PIT:SFO:10
./libclient/libclient a 2-1:3:4:PIT:SFO:10
./libclient/libclient a 3-1:5:6:PIT:SFO:10

./libclient/libclient -n 5 b Alice@Com.com:1:2:1-1:2-1
./libclient/libclient -n 5 b Bob@Com.com:1:2:2-1:3-1
./libclient/libclient b Alice@Com.com:3:2:1-1:3-1

./libclient/libclient c Alice@Com.com:1:3:1-1:2-1:3-1
./libclient/libclient c Bob@Com.com:1:1:3-1
./libclient/libclient b Alice@Com.com:1:2:1-1:3-1

./libclient/libclient c Alice@Com.com:1:3:1-1:2-1:3-1
./libclient/libclient c Alice@Com.com:4:1:2-1
./libclient/libclient d 2-1
