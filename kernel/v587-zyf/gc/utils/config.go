package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

func LoadConf(fileName string, rawVal interface{}) error {
	pathName := path.Base(fileName)
	suffix := path.Ext(pathName)

	if len(suffix) <= 1 {
		// return errcode.ERR_FILE_NAME_INVALID
		return fmt.Errorf("err: filename invalid, %s", fileName)
	}
	viper.SetConfigType(suffix[1:])

	viper.SetConfigFile(fileName)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&rawVal)
	if err != nil {
		return err
	}

	return nil
}

func Load(configVal interface{}, path string, fileName ...string) error {
	var configFileName string
	if len(fileName) > 0 {
		configFileName = fileName[0]
	}
	if configFileName == "" {
		configFileName = getDefaultFileName(path)
	}

	// load from config file
	err := LoadConf(configFileName, configVal)
	if err != nil {
		return err
	}

	// load from env
	reflectT := reflect.TypeOf(configVal)
	reflectV := reflect.ValueOf(configVal)

	if reflectT.Kind() == reflect.Ptr {
		reflectT = reflectT.Elem()
	}
	if reflectV.Kind() == reflect.Ptr {
		reflectV = reflectV.Elem()
	}

	for i := 0; i < reflectT.NumField(); i++ {
		t := reflectT.Field(i)
		v := reflectV.Field(i)
		envKey := t.Tag.Get("env")
		if envKey == "" || envKey == "-" {
			continue
		}
		envValue := os.Getenv(envKey)
		if envValue == "" {
			continue
		}
		if !v.CanSet() {
			continue
		}

		if v.Kind() == reflect.String {
			v.SetString(envValue)
		} else if v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64 {
			tmpV, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				continue
			}
			v.SetInt(tmpV)
		} else if v.Kind() >= reflect.Uint && v.Kind() <= reflect.Uint64 {
			tmpV, err := strconv.ParseUint(envValue, 10, 64)
			if err != nil {
				continue
			}
			v.SetUint(tmpV)
		} else if v.Kind() == reflect.Bool {
			tmpV, err := strconv.ParseBool(envValue)
			if err != nil {
				continue
			}
			v.SetBool(tmpV)
		} else if v.Kind() == reflect.Float32 && v.Kind() == reflect.Float64 {
			tmpV, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				continue
			}
			v.SetFloat(tmpV)
		}
	}

	return nil
}

func getDefaultFileName(path ...string) string {
	configFileName := "./conf/config.yml"
	if len(path) > 0 && path[0] != "" {
		configFileName = path[0]
	}
	mode := os.Getenv("MODE")
	if mode != "" {
		idx := strings.LastIndex(configFileName, ".")
		if idx < 0 {
			return ""
		}
		configFileName = strings.Join([]string{configFileName[:idx+1], mode, configFileName[idx:]}, "")
	}

	return configFileName
}
