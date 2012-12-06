#!/bin/bash

./libclient/libclient a 1-2:11:12:PIT:SFO:10
./libclient/libclient q 11:20
./libclient/libclient q 0:0
./libclient/libclient b Bob@Com.com:1:1:1-2
./libclient/libclient c Bob@Com.com:1:1:1-2
./libclient/libclient b Bob@Com.com:2:1:1-2
./libclient/libclient c Bob@Com.com:1:1:1-2
./libclient/libclient d 1-2
