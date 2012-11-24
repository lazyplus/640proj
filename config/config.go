package config

import (
    "strconv"
    "os"
    "bufio"
    "strings"
)

type Config struct {
    NumAirline     int
    AirlineAddr    map[string] string
    CoordAddr       string
}

func ReadConfigFile (path string) (*Config, error) {
    in := &Config{}

    in.AirlineAddr = make(map[string] string)
    userFile := path
    fin, err := os.Open(userFile)
    defer fin.Close()

    if err != nil {
        return nil, err
    }

    in.NumAirline = -1
    cur := 0
    rf := bufio.NewReader(fin)
    for{
        s, err2 := rf.ReadString('\n')
        s = strings.TrimSpace(s)
        if err2 != nil {
            return nil, err2
        }
        if in.NumAirline == -1 {
            in.NumAirline , _ = strconv.Atoi(s)
        }else if cur < in.NumAirline{
            ss := strings.Split(s,"\t")
            in.AirlineAddr[ss[0]] = ss[1]
            cur ++
        }else{
            in.CoordAddr = s
            break
        }
    }
    return in, nil
}
