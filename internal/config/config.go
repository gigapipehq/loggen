package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
)

type Config struct {
	URL           string            `yaml:"url" json:"url" validate:"required"`
	Headers       map[string]string `yaml:"headers" json:"headers"`
	APIKey        string            `yaml:"api_key" json:"api_key" validate:"required"`
	APISecret     string            `yaml:"api_secret" json:"api_secret" validate:"required"`
	Labels        map[string]string `yaml:"labels" json:"labels" validate:"required"`
	Rate          int               `yaml:"rate" json:"rate" validate:"required"`
	Timeout       Duration          `yaml:"timeout" json:"timeout" validate:"required"`
	LogConfig     LogConfig         `yaml:"log_config" json:"log_config" validate:"required"`
	EnableMetrics bool              `yaml:"enable_metrics" json:"enable_metrics"`
	Traces        TracesConfig      `yaml:"traces" json:"traces" validate:"required"`
}

type Duration time.Duration

func (s Duration) MarshalJSON() ([]byte, error) {
	d := time.Duration(s)
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func (s *Duration) UnmarshalJSON(data []byte) error {
	d, err := time.ParseDuration(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*s = Duration(d)
	return nil
}

func (s Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(s).String(), nil
}

func (s *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var d time.Duration
	if err := unmarshal(&d); err != nil {
		return err
	}
	*s = Duration(d)
	return nil
}

func (s Duration) String() string {
	return time.Duration(s).String()
}

func (s Duration) Seconds() float64 {
	return time.Duration(s).Seconds()
}

type LogConfig struct {
	Format    string            `yaml:"format" json:"format" validate:"oneof=logfmt json"`
	Structure map[string]string `yaml:"structure" json:"structure" validate:"required"`
}

type TracesConfig struct {
	Enabled  bool                  `yaml:"enabled" json:"enabled"`
	Defaults []attribute.Key       `yaml:"defaults" json:"defaults"`
	Custom   attributeKeyValueList `yaml:"custom" json:"custom"`
	Spans    []SpanStep            `json:"spans" json:"spans"`
}

type SpanStep struct {
	Kind       trace.SpanKind        `yaml:"kind" json:"kind" validate:"required"`
	Name       string                `yaml:"name" json:"name"`
	Duration   Duration              `yaml:"duration" json:"duration"`
	Attributes []SpanAttributeConfig `yaml:"attributes" json:"attributes"`
	Children   []SpanStep            `yaml:"children" json:"children"`
}

type SpanAttributeConfig struct {
	Name                string `yaml:"name" json:"name"`
	ValueType           string `yaml:"value_type" json:"value_type" validate:"required"`
	ResolveFake         string `yaml:"resolve_fake" json:"resolve_fake"`
	ResolveFromLogValue string `yaml:"resolve_from_log_value" json:"resolve_from_log_value"`
}

type logLine interface {
	ToLogFMT() string
	ToJSON() string
}

func init() {
	gofakeit.AddFuncLookup("url_path", gofakeit.Info{
		Display:     "URL path",
		Category:    "internet",
		Description: "Random URL path",
		Example:     "/users/turtle",
		Output:      "string",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u, _ := url.Parse(gofakeit.URL())
			return u.Path, nil
		},
	})
}

func (cfg *Config) GetHeaders() map[string]string {
	if cfg.Headers == nil {
		cfg.Headers = map[string]string{}
	}
	m := make(map[string]string)
	for k, v := range cfg.Headers {
		m[k] = v
	}
	m["X-API-Key"] = cfg.APIKey
	m["X-API-Secret"] = cfg.APISecret
	return m
}

type attributeKeyValueList []attribute.KeyValue

func (kvList *attributeKeyValueList) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	kvList = attributeKeyValueListFromMap(m)
	return nil
}

func (kvList *attributeKeyValueList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m map[string]interface{}
	if err := unmarshal(&m); err != nil {
		return err
	}
	kvList = attributeKeyValueListFromMap(m)
	return nil
}

func (kvList attributeKeyValueList) MarshalYAML() (interface{}, error) {
	m := map[string]interface{}{}
	for _, v := range kvList {
		m[string(v.Key)] = v.Value.AsInterface()
	}
	return m, nil
}

func GetLogLineMarshaller[T logLine](config LogConfig) func(T) string {
	if config.Format == "json" {
		return func(obj T) string {
			return obj.ToJSON()
		}
	}
	return func(obj T) string {
		return obj.ToLogFMT()
	}
}

var (
	BasePath       = fmt.Sprintf("%s/.loggen", os.Getenv("HOME"))
	ConfigFilename = fmt.Sprintf("%s/config.yaml", BasePath)
	c              = &Config{}
)

func Load() {
	f, err := os.Open(ConfigFilename)
	if err != nil {
		fmt.Println("Creating default config...")
		if err := os.MkdirAll(BasePath, os.ModePerm); err != nil {
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
		os.Exit(1)
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
		case "config.Duration":
			d, err := time.ParseDuration(value)
			if err != nil {
				return convertErr("config.Duration")
			}
			val = reflect.ValueOf(Duration(d))
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
		Timeout: Duration(30 * time.Second),
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
		Traces: TracesConfig{
			Enabled: true,
		},
	}
}

func writeConfig(c *Config) error {
	b, _ := yaml.Marshal(c)
	return os.WriteFile(ConfigFilename, b, os.ModePerm)
}

func attributeKeyValueListFromMap(m map[string]interface{}) *attributeKeyValueList {
	var kv []attribute.KeyValue
	for k, v := range m {
		kv = append(kv, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(fmt.Sprintf("%v", v)),
		})
	}
	kvl := attributeKeyValueList(kv)
	return &kvl
}
