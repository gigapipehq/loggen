package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	URL           string            `yaml:"url" json:"url" validate:"required"`
	APIKey        string            `yaml:"api_key" json:"api_key" validate:"required"`
	APISecret     string            `yaml:"api_secret" json:"api_secret" validate:"required"`
	Labels        map[string]string `yaml:"labels" json:"labels" validate:"required"`
	Rate          int               `yaml:"rate" json:"rate" validate:"required"`
	Timeout       time.Duration     `yaml:"timeout" json:"timeout" validate:"required"`
	LogConfig     LogConfig         `yaml:"log_config" json:"log_config" validate:"required"`
	EnableMetrics bool              `yaml:"enable_metrics" json:"enable_metrics"`
	EnableTraces  bool              `yaml:"enable_traces" json:"enable_traces"`
}

type LogConfig struct {
	Format    string            `yaml:"format" json:"format" validate:"oneof=logfmt json"`
	Structure map[string]string `yaml:"structure" json:"structure" validate:"required"`
}

type DetailedLogConfig struct {
	Format    string    `json:"format"`
	Structure []LogInfo `json:"structure"`
}

type LogInfo struct {
	Display     string           `json:"display"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Example     string           `json:"example"`
	Params      []gofakeit.Param `json:"params"`
}

func (lc *LogConfig) Detailed(categories ...string) *DetailedLogConfig {
	cfg := &DetailedLogConfig{Format: lc.Format, Structure: []LogInfo{}}
	for _, v := range lc.Structure {
		split := strings.Split(v, ":")
		if info := gofakeit.GetFuncLookup(split[0]); info != nil {
			if len(categories) > 0 {
				if !hasCategory(info.Category, categories) {
					continue
				}
			}

			if len(split) > 1 {
				params := strings.Split(split[1], ",")
				for i, p := range params {
					info.Params[i].Default = p
				}
			}
			cfg.Structure = append(cfg.Structure, LogInfo{
				Display:     info.Display,
				Category:    info.Category,
				Description: info.Description,
				Example:     info.Example,
				Params:      info.Params,
			})
		}
	}
	return cfg
}

var (
	basePath       = fmt.Sprintf("%s/.loggen", os.Getenv("HOME"))
	configFilename = fmt.Sprintf("%s/config.yaml", basePath)
	c              = &Config{}
)

func Load() {
	f, err := os.Open(configFilename)
	if err != nil {
		fmt.Println("Creating default config...")
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			log.Printf("unable create config file directory: %v", err)
			return
		}

		c = getDefaultConfig()
		_ = writeConfig(c)
		return
	}

	b, _ := io.ReadAll(f)
	if err := yaml.Unmarshal(b, c); err != nil {
		fmt.Println(err)
	}
	_ = f.Close()
}

func SettingNames() []string {
	var tags []string
	v := reflect.ValueOf(Config{})
	for i := 0; i < v.Type().NumField(); i++ {
		tags = append(tags, v.Type().Field(i).Tag.Get("yaml"))
	}
	return tags
}

func GetSettingValue(name string) any {
	v := reflect.ValueOf(*c)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if name == t.Field(i).Tag.Get("yaml") {
			return v.Field(i)
		}
	}
	return nil
}

func UpdateSettingValue(name string, value string) error {
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if name == typeField.Tag.Get("yaml") {
			if err := updateSettingValue(name, typeField.Name, value); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func Reset() error {
	return writeConfig(getDefaultConfig())
}

func updateSettingValue(tagName, name string, value string) error {
	convertErr := func(typeName string) error {
		return fmt.Errorf("unable to convert value %s for field %s to %s", value, tagName, typeName)
	}
	structValue := reflect.ValueOf(c).Elem()
	structFieldValue := structValue.FieldByName(name)
	structFieldType := structFieldValue.Type()

	val := reflect.ValueOf(value)
	switch structFieldType.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return convertErr(reflect.Int.String())
		}
		val = reflect.ValueOf(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return convertErr(reflect.Bool.String())
		}
		val = reflect.ValueOf(b)
	case reflect.Map:
		if strings.TrimSpace(value) == "" {
			val = reflect.ValueOf(map[string]string{})
			break
		}
		groups := strings.Split(value, ",")
		m := map[string]string{}
		for _, g := range groups {
			kvs := strings.Split(g, "=")
			if len(kvs) != 2 {
				return convertErr("string map")
			}
			m[strings.TrimSpace(kvs[0])] = strings.TrimSpace(kvs[1])
		}
		val = reflect.ValueOf(m)
	case reflect.Int64:
		switch structFieldType.String() {
		case "time.Duration":
			d, err := time.ParseDuration(value)
			if err != nil {
				return convertErr("time.Duration")
			}
			val = reflect.ValueOf(d)
		default:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return convertErr(reflect.Int.String())
			}
			val = reflect.ValueOf(i)
		}
	}
	structFieldValue.Set(val)
	return writeConfig(c)
}

func ValidArgSettingName(cmdName string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		if err := cobra.OnlyValidArgs(cmd, args[:1]); err != nil {
			return fmt.Errorf(
				"invalid argument \"%s\" for \"loggen config %s\". Check documentation for list of valid values",
				args[0],
				cmdName,
			)
		}
		return nil
	}
}

func Get() *Config {
	return c
}

func getDefaultConfig() *Config {
	return &Config{
		URL:     "https://qryn.gigapipe.com",
		Rate:    100,
		Timeout: 30 * time.Second,
		LogConfig: LogConfig{
			Format: "logfmt",
			Structure: map[string]string{
				"level":       "loglevel",
				"host":        "domainname",
				"method":      "httpmethod",
				"status_code": "httpstatuscodesimple",
				"bytes":       "number:0,300",
			},
		},
		EnableMetrics: true,
		EnableTraces:  true,
	}
}

func writeConfig(c *Config) error {
	b, _ := yaml.Marshal(c)
	return os.WriteFile(configFilename, b, os.ModePerm)
}

func hasCategory(category string, list []string) bool {
	for _, cat := range list {
		if cat == category {
			return true
		}
	}
	return false
}
