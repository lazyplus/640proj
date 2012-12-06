for ((i=0;i<10;i++))
do
	go run libclient/libclient.go q 1:20 &
	go run libclient/libclient.go q 0:0 &
done

