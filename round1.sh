go run libclient/libclient.go a 1-1:1:2:PIT:SFO:10
go run libclient/libclient.go q 1:10
go run libclient/libclient.go q 0:0
go run libclient/libclient.go b Bob@Com.com:1:1:1-1
#go run libclient/libclient.go c Bob@Com.com:1:1:1-1
go run libclient/libclient.go b Bob@Com.com:2:1:1-1
#go run libclient/libclient.go c Bob@Com.com:1:1:1-1
go run libclient/libclient.go d 1-1

