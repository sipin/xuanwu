package models

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
)
