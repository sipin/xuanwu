package test

import (
	"encoding/json"
	"strconv"
)

var _ = strconv.Itoa
var _ = json.Marshal

var (
	UserStatus = struct {
		Banned, Offline, Online int32
	}{
		0, 1, 2,
	}

	UserStatusLabel = map[int32]string{
		0: "禁用",
		1: "离线",
		2: "在线",
	}
	UserStatusJSON      string
	UserStatusLabelJSON string
)

func init() {
	var ret []byte
	var err error
	{
		ret, _ = json.Marshal(UserStatus)
		UserStatusJSON = string(ret)
		tmp := make(map[string]string, len(UserStatusLabel))
		for k, v := range UserStatusLabel {
			tmp[strconv.Itoa(int(k))] = v
		}
		ret, err = json.Marshal(tmp)
		if err != nil {
			panic(err)
		}
		UserStatusLabelJSON = string(ret)
	}

}
