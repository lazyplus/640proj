package inputFormat

import {
	"strconv"
	"os"
	"bufio"
	"strings"
}

type input struct {
	num_airline		int
	addr_airline	map[string] string
	coordaddr		string
}

func (in *input) read_config_file (path string) error {
	in.addr_airline = make(map[string] string)
	userFile := path
	fin, err := os.Open(userFile)
	defer fin.Close()
	if err != nil {
		return err
	}
	in.num_airline = 0
	cur := 0
	rf := bufio.NewReader(fin)
	for{
		s, err2 := rf.ReadString('\n')
		if err2 != nil {
			return err2
		}
		if in.num_airline == 0 {
			in.num_airline , _ = strconv.Atoi(s)
		}else if cur < in.num_airline{
			ss := strings.Split(s,"\t")
			in.addr_airline[ss[0]] = ss[1]
			cur ++
		}else{
			in.coordaddr = s
		}
	}
	return nil
}
