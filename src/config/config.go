/*
  trader API Engine
*/

package config

import (
	"encoding/json"
	"io/ioutil"
)

var Config map[string]string

func get_config_path(file string) (filepath string) {
	//fmt.Println("config file:", Root + file)
	return Root + file
}

func load_config(file string) (config map[string]string, err error) {
	// Load 全局配置文件
	configFile := get_config_path(file)

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config = make(map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadConfig() (err error) {
	Config, err = load_config("/conf/config.json")
	if err != nil {
		panic("load config file failed")
		return
	}

	return
}
