go run libclient/libclient.go -n 40 q 1:10 &
go run libclient/libclient.go -n 40 b fucker@fuck.com:1:2:1-1:2-1 &
go run libclient/libclient.go -n 40 b fucker2@fuck.com:2:2:2-1:3-1 &
go run libclient/libclient.go -n 40 c fucker3@fuck.com:3:3:1-1:2-1:3-1 &

go run libclient/libclient.go -n 40 b shit@fuck.com:5:2:1-1:3-1 &
go run libclient/libclient.go -n 40 c shit@fuck.com:3:2:1-1:3-1 &
go run libclient/libclient.go -n 40 c shit@fuck.com:2:1:3-1 &
go run libclient/libclient.go -n 40 c shit@fuck.com:2:1:2-1 &

