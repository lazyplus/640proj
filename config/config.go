package config

import (
    "strconv"
    "os"
    "bufio"
    "strings"
    "fmt"
)

type AirlineConfig struct {
    Name string
    DelegateHostPort string
    NumPeers int
    PeersHostPort []string
    UDPPort []int
}

type Config struct {
    NumAirline int
    Airlines map[string] *AirlineConfig
    CoordHostPort string
}

func ReadConfigFile (path string) (*Config, error) {
    conf := &Config{}
    conf.Airlines = make(map[string] *AirlineConfig)

    fin, err := os.Open(path)
    defer fin.Close()

    if err != nil {
        return nil, err
    }

    var line string
    rf := bufio.NewReader(fin)
    line, err = rf.ReadString('\n')
    conf.NumAirline, _ = strconv.Atoi(strings.TrimSpace(line))
    conf.Airlines = make(map[string] *AirlineConfig)

    for i:=0; i<conf.NumAirline; i++ {
        line, err = rf.ReadString('\n')
        parts := strings.Split(line, " ")
        ac := &AirlineConfig{}
        ac.Name = parts[0]
        ac.NumPeers, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
        ac.PeersHostPort = make([]string, ac.NumPeers)
        ac.UDPPort = make([]int, ac.NumPeers)
        line, err = rf.ReadString('\n')
        ac.DelegateHostPort = strings.TrimSpace(line)
        for j:=0; j<ac.NumPeers; j++ {
            line, err = rf.ReadString('\n')
            line = strings.TrimSpace(line)
            ss := strings.Split(line, " ")
            ac.PeersHostPort[j] = ss[0]
            ac.UDPPort[j], _ = strconv.Atoi(ss[1])
        }
        conf.Airlines[ac.Name] = ac
    }

    line, err = rf.ReadString('\n')
    conf.CoordHostPort = line

    return conf, nil
}

func DumpConfig(conf *Config) {
    fmt.Println(conf.NumAirline)
    for _, value := range(conf.Airlines) {
        fmt.Println(value.Name + " " + strconv.FormatInt(int64(value.NumPeers), 10))
        fmt.Println(value.DelegateHostPort)
        for i:=0; i<value.NumPeers; i++ {
            fmt.Println(value.PeersHostPort[i] + " " + strconv.FormatInt(int64(value.UDPPort[i]), 10))
        }
    }
    fmt.Println(conf.CoordHostPort)
}
