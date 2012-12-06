#!/bin/bash

./libclient/libclient a 1-3:21:22:PIT:SFO:10
./libclient/libclient q 21:30
./libclient/libclient q 0:0
./libclient/libclient b Bob@Com.com:1:1:1-3
./libclient/libclient b Bob@Com.com:10:1:1-3
./libclient/libclient c Bob@Com.com:1:1:1-3
./libclient/libclient b Bob@Com.com:2:1:1-3
./libclient/libclient c Bob@Com.com:1:1:1-3
./libclient/libclient c Bob@Com.com:2:1:1-3
./libclient/libclient d 1-3
