package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// TODO: Download data from OpenRS2 (maybe as a one-time cli function that downloads everything at once)

type XTEA struct {
	Archive   int     `json:"archive"`
	Group     int     `json:"group"`
	NameHash  int     `json:"name_hash"`
	Name      string  `json:"name"`
	MapSquare int     `json:"mapsquare"`
	Key       []int32 `json:"key"`
}

var XTEAs []XTEA

func init() {
	content, err := os.ReadFile("/home/owner/Code/github.com/zsrv/rt5-server-go/data/xteas.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &XTEAs)
	if err != nil {
		panic(err)
	}
}

func GetXTEA(x int, z int) (XTEA, bool) {
	for _, v := range XTEAs {
		if v.MapSquare == x<<8|z {
			return v, true
		}
	}
	fmt.Printf("WARNING: BAD XTEA RETURNED for x %v, z %v\n", x, z)
	return XTEA{}, false
}

func GetGroup(archive uint8, group uint16) ([]byte, error) {
	path := fmt.Sprintf("data/cache/%d/%d.dat", archive, group)
	if _, err := os.Stat(path); err != nil {
		fmt.Printf("path %s does not exist!\n", path)
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("could not read file %s!\n", path)
		return nil, err
	}

	return file, nil
}
