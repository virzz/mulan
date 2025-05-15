package db

import (
	"strings"

	"github.com/spf13/pflag"
)

type GenParams struct {
	DSN               string   `yaml:"dsn"`               // consult[https://gorm.io/docs/connecting_to_the_database.html]"
	DB                string   `yaml:"db"`                // input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
	Tables            []string `yaml:"tables"`            // enter the required data table or leave it blank
	OnlyModel         bool     `yaml:"onlyModel"`         // only generate model
	OutPath           string   `yaml:"outPath"`           // specify a directory for output
	OutFile           string   `yaml:"outFile"`           // query code file name, default: gen.go
	WithUnitTest      bool     `yaml:"withUnitTest"`      // generate unit test for query code
	ModelPkgName      string   `yaml:"modelPkgName"`      // generated model code's package name
	FieldNullable     bool     `yaml:"fieldNullable"`     // generate with pointer when field is nullable
	FieldCoverable    bool     `yaml:"fieldCoverable"`    // generate with pointer when field has default value
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
	FieldSignable     bool     `yaml:"fieldSignable"`     // detect integer field's unsigned type, adjust generated data type
}

func (c *GenParams) Revise() *GenParams {
	if c == nil {
		return c
	}
	if c.DB == "" {
		c.DB = "mysql"
	}
	if c.OutPath == "" {
		c.OutPath = defaultQueryPath
	}
	if len(c.Tables) == 0 {
		return c
	}
	tableList := make([]string, 0, len(c.Tables))
	for _, tableName := range c.Tables {
		_tableName := strings.TrimSpace(tableName) // trim leading and trailing space in tableName
		if _tableName == "" {                      // skip empty tableName
			continue
		}
		tableList = append(tableList, _tableName)
	}
	c.Tables = tableList
	return c
}

type GenConfig struct {
	Version  string     `yaml:"version"`
	Database *GenParams `yaml:"database"`
}

func GenFlagSet() *pflag.FlagSet {
	fs := pflag.NewFlagSet("gen", pflag.ContinueOnError)

	fs.StringP("config", "c", "", "is path for gen.yml")

	fs.String("dsn", "", "consult[https://gorm.io/docs/connecting_to_the_database.html]")
	fs.String("db", "", "input mysql|postgres. consult[https://gorm.io/docs/connecting_to_the_database.html]")
	fs.String("tables", "", "enter the required data table or leave it blank")
	fs.Bool("onlyModel", false, "only generate models (without query file)")
	fs.String("outPath", "", "specify a directory for output")
	fs.String("outFile", "", "query code file name, default: gen.go")
	fs.Bool("withUnitTest", false, "generate unit test for query code")
	fs.String("modelPkgName", "", "generated model code's package name")
	fs.Bool("fieldNullable", false, "generate with pointer when field is nullable")
	fs.Bool("fieldCoverable", false, "generate with pointer when field has default value")
	fs.Bool("fieldWithIndexTag", false, "generate field with gorm index tag")
	fs.Bool("fieldWithTypeTag", false, "generate field with gorm column type tag")
	fs.Bool("fieldSignable", false, "detect integer field's unsigned type, adjust generated data type")
	return fs
}
