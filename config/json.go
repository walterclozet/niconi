package config

import (
	"elichika/utils"
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type AppConfigs struct {
	AppName  string   `json:"app_name"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	IsGlobal      bool   `json:"is_global"`
	MasterVersion string `json:"master_version"`
	StartUpKey    string `json:"startup_key"`
	CdnServer     string `json:"cdn_server"`
}

type LevelDbConfigs struct {
	DataPath string `json:"data_path"`
}

func DefaultConfigs() *AppConfigs {
	return &AppConfigs{
		AppName: "elichika",
		Settings: Settings{
			IsGlobal:      false,
			MasterVersion: "b66ec2295e9a00aa",
			StartUpKey:    "5f7IZY1QrAX0D49g",
			CdnServer:     "http://192.168.1.123/static",
		},
	}
}

func Load(p string) *AppConfigs {
	if !utils.PathExists(p) {
		_ = DefaultConfigs().Save(p)
	}
	c := AppConfigs{}
	err := json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	if err != nil {
		_ = os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		_ = DefaultConfigs().Save(p)
	}
	c = AppConfigs{}
	_ = json.Unmarshal([]byte(utils.ReadAllText(p)), &c)
	return &c
}

func (c *AppConfigs) Save(p string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	utils.WriteAllText(p, string(data))
	return nil
}
