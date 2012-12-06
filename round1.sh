#!/bin/bash

./libclient/libclient a 1-1:1:2:PIT:SFO:10
./libclient/libclient q 1:10
./libclient/libclient q 0:0
./libclient/libclient b Bob@Com.com:1:1:1-1
./libclient/libclient b Bob@Com.com:1:1:1-1
./libclient/libclient b Bob@Com.com:1:1:1-1
./libclient/libclient b Bob@Com.com:1:1:1-1
./libclient/libclient b Bob@Com.com:1:1:1-1
./libclient/libclient c Bob@Com.com:5:1:1-1
./libclient/libclient b Bob@Com.com:2:1:1-1
./libclient/libclient c Bob@Com.com:1:1:1-1
./libclient/libclient c Bob@Com.com:1:1:1-0
./libclient/libclient c Alice@Com.com:1:1:1-1
./libclient/libclient d 1-1
