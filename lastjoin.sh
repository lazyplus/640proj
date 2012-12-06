#!/bin/bash

echo Compiling Binaries

cd libclient
go build libclient.go
cd ..

cd peer
go build peer.go
cd ..

cd delegateserver
go build delegateserver.go
cd ..

cd coordserver
go build coordserver.go
cd ..

echo Starting Servers...

./peer/peer -id 1 -name 1 2>&1 > log1_1 &
PEER1_PID=$!
./peer/peer -id 2 -name 1 2>&1 > log1_2 &

./peer/peer -id 0 -name 2 2>&1 > log2_0 &
./peer/peer -id 1 -name 2 2>&1 > log2_1 &
./peer/peer -id 2 -name 2 2>&1 > log2_2 &

./peer/peer -id 0 -name 3 2>&1 > log3_0 &
./peer/peer -id 1 -name 3 2>&1 > log3_1 &
./peer/peer -id 2 -name 3 2>&1 > log3_2 &

./delegateserver/delegateserver -name 1 -port 12300 2>&1 > log_d1 &
./delegateserver/delegateserver -name 2 -port 12310 2>&1 > log_d2 &
./delegateserver/delegateserver -name 3 -port 12320 2>&1 > log_d3 &

./coordserver/coordserver -port 12400 2>&1 > log_c &

sleep 5
echo Starting Clients

./testround.sh 1 &
T1=$!
./testround.sh 2 &
T2=$!
./testround.sh 3 &
T3=$!
./testround.sh 4 &
T4=$!

sleep 5
echo Starting Late Server

./peer/peer -id 0 -name 1 2>&1 > log1_0 &
PEER_PID=$!

sleep 8

echo Killing One Peer
kill $PEER_PID
wait $PEER_PID
echo Killed

sleep 8

echo Starting Late Server Again

./peer/peer -id 0 -name 1 2>&1 > log1_0 &
PEER_PID=$!

sleep 8

echo Killing Another Server
kill $PEER1_PID
wait $PEER1_PID
echo Killed

wait $T1
wait $T2
wait $T3
wait $T4

killall peer
killall coordserver
killall delegateserver

echo Test Finished!
