#!/bin/bash

./libclient/libclient a 1-4:31:32:PIT:SFO:10
./libclient/libclient q 31:40
./libclient/libclient q 0:0
./libclient/libclient b Bob@Com.com:1:1:1-4
./libclient/libclient b Bob@Com.com:10:1:1-4
./libclient/libclient c Bob@Com.com:1:1:1-4
./libclient/libclient b Bob@Com.com:2:1:1-4
./libclient/libclient c Bob@Com.com:1:1:1-4
./libclient/libclient c Bob@Com.com:2:1:1-4
./libclient/libclient r 1-4:33:34:PIT:SFO:10
./libclient/libclient q 31:32
./libclient/libclient q 33:34
./libclient/libclient d 1-4
