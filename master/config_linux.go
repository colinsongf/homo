// +build linux

package master

import (
	"time"

	"github.com/aiicy/aiicy/logger"
	"github.com/aiicy/aiicy/protocol/http"
	"github.com/aiicy/aiicy/sdk/aiicy-go/api"
)

// DBConf db config
type DBConf struct {
	Driver string
	Path   string
}

// Config master init config
type Config struct {
	Mode     string           `yaml:"mode" json:"mode" default:"docker" validate:"regexp=^(native|docker)$"`
	Server   http.ServerInfo  `yaml:"server" json:"server" default:"{\"address\":\"unix:///var/run/aiicy.sock\"}"`
	Database DBConf           `yaml:"database" json:"database" default:"{\"driver\":\"sqlite3\",\"path\":\"var/lib/aiicy/db\"}"`
	API      api.ServerConfig `yaml:"api" json:"api" default:"{\"address\":\"unix:///var/run/aiicy/api.sock\"}"`
	Logger   logger.LogInfo   `yaml:"logger" json:"logger" default:"{\"path\":\"var/log/aiicy/aiicy.log\"}"`
	OTALog   logger.LogInfo   `yaml:"otalog" json:"otalog" default:"{\"path\":\"var/db/aiicy/ota.log\",\"format\":\"json\"}"`
	Grace    time.Duration    `yaml:"grace" json:"grace" default:"30s"`
	SNFile   string           `yaml:"snfile" json:"snfile"`
	Docker   struct {
		APIVersion string `yaml:"api_version" json:"api_version" default:"1.38"`
	} `yaml:"docker" json:"docker"`
	// cache config file path
	File string
}
