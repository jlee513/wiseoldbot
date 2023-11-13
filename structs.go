package main

type hallOfFameInfo struct {
	Data struct {
		BossName string         `json:"skill"`
		Players  map[int]Player `json:"players"`
	} `json:"data"`
}

type Player struct {
	Username string `json:"username"`
	Kc       int    `json:"xp"`
}

type hallOfFameRequestInfo struct {
	Bosses   map[string]string
	DiscChan string
}

type collectionLogInfo struct {
	CollectionLog struct {
		Uniques int `json:"uniqueObtained"`
	} `json:"collectionLog"`
}
