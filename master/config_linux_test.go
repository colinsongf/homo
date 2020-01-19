// +build linux

package master

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/countstarlight/homo/logger"
	"github.com/countstarlight/homo/protocol/http"
	"github.com/countstarlight/homo/sdk/homo-go"
	"github.com/countstarlight/homo/sdk/homo-go/api"
	"github.com/countstarlight/homo/utils"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		args []byte
	}{
		{
			name: "nil",
			args: nil,
		},
		{
			name: "empty",
			args: []byte{},
		},
		{
			name: "empty2",
			args: []byte(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			err := utils.UnmarshalYAML(tt.args, &cfg)
			assert.NoError(t, err)

			assert.Equal(t, "docker", cfg.Mode)

			if runtime.GOOS == "linux" {
				assert.Equal(t, "unix:///var/run/homo.sock", cfg.Server.Address)
			} else {
				assert.Equal(t, "tcp://127.0.0.1:50050", cfg.Server.Address)
			}
			assert.Equal(t, time.Duration(5*60*1000*1000000), cfg.Server.Timeout)

			assert.Equal(t, "var/log/homo/homo.log", cfg.Logger.Path)
			assert.Equal(t, "info", cfg.Logger.Level)
			assert.Equal(t, "text", cfg.Logger.Format)
			assert.Equal(t, 15, cfg.Logger.Age.Max)
			assert.Equal(t, 50, cfg.Logger.Size.Max)
			assert.Equal(t, 15, cfg.Logger.Backup.Max)

			assert.Equal(t, time.Duration(30*1000*1000000), cfg.Grace)
		})
	}
}

func TestConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	conf := Config{
		Server: http.ServerInfo{
			Address: "homo://127.0.0.1:51150",
		},
		API: api.ServerConfig{
			Address: "homo://127.0.0.1:51150",
		},
	}
	err = conf.Validate()
	assert.Error(t, err)
	assert.Equal(t, "only support unix domian socket or tcp socket", err.Error())

	filepath := path.Join(dir, "sn")
	sn := "homo"
	f, err := os.Create(filepath)
	assert.NoError(t, err)

	n, err := io.WriteString(f, sn)
	assert.NoError(t, err)
	assert.Len(t, sn, n)
	f.Sync()
	f.Close()

	conf = Config{
		Mode: "docker",
		Server: http.ServerInfo{
			Address: "unix:///tmp/run/homo.sock",
		},
		API: api.ServerConfig{
			Address: "unix:///tmp/run/homo/api.sock",
		},
		SNFile: filepath,
	}
	err = conf.Validate()
	assert.NoError(t, err)

	assert.Equal(t, "unix:///"+homo.DefaultSockFile, utils.GetEnv(homo.EnvKeyMasterAPIAddress))
	assert.Equal(t, "unix:///"+homo.DefaultGRPCSockFile, utils.GetEnv(homo.EnvKeyMasterGRPCAPIAddress))
	assert.Equal(t, "/tmp/run/homo.sock", utils.GetEnv(homo.EnvKeyMasterAPISocket))
	assert.Equal(t, "/tmp/run/homo/api.sock", utils.GetEnv(homo.EnvKeyMasterGRPCAPISocket))
	assert.Equal(t, sn, utils.GetEnv(homo.EnvKeyHostSN))
	assert.Equal(t, "v1", utils.GetEnv(homo.EnvKeyMasterAPIVersion))
	assert.Equal(t, runtime.GOOS, utils.GetEnv(homo.EnvKeyHostOS))
	assert.Equal(t, conf.Mode, utils.GetEnv(homo.EnvKeyServiceMode))

	conf = Config{
		Mode: "native",
		Server: http.ServerInfo{
			Address: "unix:///tmp/run/homo.sock",
		},
		API: api.ServerConfig{
			Address: "unix:///tmp/run/homo/api.sock",
		},
		SNFile: filepath,
	}
	err = conf.Validate()
	assert.NoError(t, err)

	assert.Equal(t, "unix://"+homo.DefaultSockFile, utils.GetEnv(homo.EnvKeyMasterAPIAddress))
	assert.Equal(t, "unix://"+homo.DefaultGRPCSockFile, utils.GetEnv(homo.EnvKeyMasterGRPCAPIAddress))
	assert.Equal(t, "/tmp/run/homo.sock", utils.GetEnv(homo.EnvKeyMasterAPISocket))
	assert.Equal(t, "/tmp/run/homo/api.sock", utils.GetEnv(homo.EnvKeyMasterGRPCAPISocket))
	assert.Equal(t, sn, utils.GetEnv(homo.EnvKeyHostSN))
	assert.Equal(t, "v1", utils.GetEnv(homo.EnvKeyMasterAPIVersion))
	assert.Equal(t, runtime.GOOS, utils.GetEnv(homo.EnvKeyHostOS))
	assert.Equal(t, conf.Mode, utils.GetEnv(homo.EnvKeyServiceMode))

	conf = Config{
		Mode: "docker",
		Server: http.ServerInfo{
			Address: "tcp://127.0.0.1:51150",
		},
		API: api.ServerConfig{
			Address: "tcp://127.0.0.1:51151",
		},
		SNFile: filepath,
	}
	err = conf.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "tcp://host.docker.internal:51150", utils.GetEnv(homo.EnvKeyMasterAPIAddress))
	assert.Equal(t, "host.docker.internal:51151", utils.GetEnv(homo.EnvKeyMasterGRPCAPIAddress))
	assert.Equal(t, sn, utils.GetEnv(homo.EnvKeyHostSN))
	assert.Equal(t, "v1", utils.GetEnv(homo.EnvKeyMasterAPIVersion))
	assert.Equal(t, runtime.GOOS, utils.GetEnv(homo.EnvKeyHostOS))
	assert.Equal(t, conf.Mode, utils.GetEnv(homo.EnvKeyServiceMode))

	conf = Config{
		Mode: "native",
		Server: http.ServerInfo{
			Address: "tcp://127.0.0.1:51150",
		},
		API: api.ServerConfig{
			Address: "tcp://127.0.0.1:51151",
		},
		SNFile: filepath,
	}
	err = conf.Validate()
	assert.NoError(t, err)
	assert.Equal(t, conf.Server.Address, utils.GetEnv(homo.EnvKeyMasterAPIAddress))
	assert.Equal(t, "127.0.0.1:51151", utils.GetEnv(homo.EnvKeyMasterGRPCAPIAddress))
	assert.Equal(t, sn, utils.GetEnv(homo.EnvKeyHostSN))
	assert.Equal(t, "v1", utils.GetEnv(homo.EnvKeyMasterAPIVersion))
	assert.Equal(t, runtime.GOOS, utils.GetEnv(homo.EnvKeyHostOS))
	assert.Equal(t, conf.Mode, utils.GetEnv(homo.EnvKeyServiceMode))
}

func TestOTALog(t *testing.T) {
	var cfg Config
	err := utils.UnmarshalYAML(nil, &cfg)
	assert.NoError(t, err)

	cfg.OTALog.Path = "testdata/ota.log"
	os.RemoveAll(cfg.OTALog.Path)
	defer os.RemoveAll(cfg.OTALog.Path)
	defer os.RemoveAll("testdata/ota.log.old")
	assert.False(t, utils.FileExists(cfg.OTALog.Path))
	logger.New(cfg.OTALog).With("step", "RECEIVE").With("trace", "xxxxxx").With("type", "APP").Info("receive update event")
	assert.True(t, utils.FileExists(cfg.OTALog.Path))
	os.Rename(cfg.OTALog.Path, "testdata/ota.log.old")
	assert.False(t, utils.FileExists(cfg.OTALog.Path))
	logger.New(cfg.OTALog).With("step", "SUCCESS").With("trace", "xxxxxx").With("type", "APP").Info("update application successfully")
	assert.True(t, utils.FileExists(cfg.OTALog.Path))
}