go run peer/peer.go -id 0 -name 1 | tee log &
go run peer/peer.go -id 1 -name 1 &
go run peer/peer.go -id 2 -name 1 &

go run peer/peer.go -id 0 -name 2 &
go run peer/peer.go -id 1 -name 2 &
go run peer/peer.go -id 2 -name 2 &

go run peer/peer.go -id 0 -name 3 &
go run peer/peer.go -id 1 -name 3 &
go run peer/peer.go -id 2 -name 3 &

