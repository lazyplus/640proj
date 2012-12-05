go run peer/peer.go -id 0 -name 1 2>&1 | tee log1_0 &
go run peer/peer.go -id 1 -name 1 2>&1 | tee log1_1 &
go run peer/peer.go -id 2 -name 1 2>&1 | tee log1_2 &

go run peer/peer.go -id 0 -name 2 &
go run peer/peer.go -id 1 -name 2 &
go run peer/peer.go -id 2 -name 2 &

go run peer/peer.go -id 0 -name 3 &
go run peer/peer.go -id 1 -name 3 &
go run peer/peer.go -id 2 -name 3 &

go run delegateserver/delegateserver.go -name 1 -port 12300 &
go run delegateserver/delegateserver.go -name 2 -port 12310 &
go run delegateserver/delegateserver.go -name 3 -port 12320 &

