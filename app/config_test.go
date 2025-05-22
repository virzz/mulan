package app_test

import (
	"os"
	"strings"
	"testing"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestUnmarshalConfig(t *testing.T) {
	type (
		Configer interface{}
		Test     struct {
			Name string `json:"name" yaml:"name"`
			Age  int    `json:"age" yaml:"age"`
		}
		Config struct {
			Test Test `json:"test" yaml:"test"`
		}
		WrapConfig struct {
			//lint:ignore SA5008 Ignore JSON option "squash"
			Config `json:",inline,squash" yaml:",inline"`
		}
	)
	var config Configer

	os.Setenv("MULAN_TEST_NAME", "mulan")
	os.Setenv("MULAN_TEST_AGE", "19")

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("test.name", "test", "test name")
	fs.Int("test.age", 0, "test age")
	viper.BindPFlags(fs)
	viper.SetEnvPrefix("Mulan")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	config = &WrapConfig{}
	err := viper.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) { dc.TagName = "json" })
	if err != nil {
		t.Fatal(err)
	}
	if config.(*WrapConfig).Test.Name != "mulan" {
		t.Fatal("test name is not mulan")
	}
	if config.(*WrapConfig).Test.Age != 19 {
		t.Fatal("test age is not 19")
	}
}
