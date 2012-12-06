#go run peer/peer.go -id 0 -name 1 2>&1 | tee log1_0 &
#P1=$!
#echo $P1
go run peer/peer.go -id 1 -name 1 2>&1 | tee log1_1 &
go run peer/peer.go -id 2 -name 1 2>&1 | tee log1_2 &

go run peer/peer.go -id 0 -name 2 2>&1 | tee log2_0 &
go run peer/peer.go -id 1 -name 2 2>&1 | tee log2_1 &
go run peer/peer.go -id 2 -name 2 2>&1 | tee log2_2 &

go run peer/peer.go -id 0 -name 3 2>&1 | tee log3_0 &
go run peer/peer.go -id 1 -name 3 2>&1 | tee log3_1 &
go run peer/peer.go -id 2 -name 3 2>&1 | tee log3_2 &

go run delegateserver/delegateserver.go -name 1 -port 12300 | tee log_d1 &
go run delegateserver/delegateserver.go -name 2 -port 12310 | tee log_d2 &
go run delegateserver/delegateserver.go -name 3 -port 12320 | tee log_d3 &

