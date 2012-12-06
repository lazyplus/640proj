for ((i=0; i<10; i++))
do
	./round$1.sh > tmp$1
	diff tmp$1 round$1.out
done

