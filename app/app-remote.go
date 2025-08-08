//go:build remote
// +build remote

package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pkg/errors"
	slogzap "github.com/samber/slog-zap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	defaultPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkhbTnK9POOy317u6ovxE
UqFT5FUPaTSSmAa0gDepG7B1SpDHmpsarJlf//doy9A4bqysxQ8Fu1njtxXU861s
J5lxS1p72UreuZoTbV+mnQFzeIqbPDiqQruzqws+hnKAVHdDcjy6NPvUH1na4bNf
snuVM/9FNik4bmd1bv362Oelhmj8jvx+sllf2L9/5H8/i35sW8oo811IE+cA+jow
BqvNT3/ayjtHlrYmnTOxGHv7H+j0JQ/yz2/ap7PWdIfspqGJZSV9iPKKfKfw37KF
H19ekMDgL248Y4PiK5BqWD4jY1hQfMsQf2ZVs2g6gGNQPAYLMiAXA4ngNurA4Kz+
xQIDAQAB
-----END PUBLIC KEY-----`
	defaultRemoteEndpoint = "config.app.virzz.com"
)

type Remote struct {
	project        string
	publicKey      string
	remoteEndpoint string
	remoteConfig   bool
	secretKey      []byte
}

func (app *App) EnableRemote(project string, publicKey ...string) error {
	app.rootCmd.Flags().String("remote-type", "json", "Remote config type")
	app.rootCmd.Flags().String("remote-endpoint", "", "Remote config endpoint")
	app.remote.project = project
	app.remote.remoteConfig = true
	if len(publicKey) > 0 && publicKey[0] != "" {
		app.remote.publicKey = publicKey[0]
	} else {
		app.remote.publicKey = defaultPublicKey
	}
	return nil
}

func (app *App) preInitRemote() {
	if app.remote.remoteConfig {
		app.remote.secretKey = make([]byte, 32)
		io.ReadFull(rand.Reader, app.remote.secretKey)
		block, _ := pem.Decode([]byte(app.remote.publicKey))
		if block == nil || block.Type != "PUBLIC KEY" {
			return errors.New("Failed to decode PEM block containing public key")
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		if key, ok := pub.(*rsa.PublicKey); ok {
			data, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, app.remote.secretKey, nil)
			viper.RemoteConfig = &RemoteProvider{EncryptSecret: data, logger: app.log.Named("remote")}
			viper.SupportedRemoteProviders = append(viper.SupportedRemoteProviders, "virzz")
			return nil
		}
		return errors.New("not an RSA public key")
	}
}
func (app *App) postInitRemote() {
	configLoaded := false
	if app.remote.remoteConfig {
		remoteEndpoint, _ := cmd.Flags().GetString("remote-endpoint")
		if remoteEndpoint == "" {
			remoteEndpoint = app.remote.remoteEndpoint
		}
		if remoteEndpoint == "" {
			remoteEndpoint = defaultRemoteEndpoint
		}
		key := fmt.Sprintf("/%s/%s/%s/%s", app.remote.project, app.ID, app.Version, instance)
		err = viper.AddSecureRemoteProvider("virzz", remoteEndpoint, key, string(app.remote.secretKey))
		if err != nil {
			app.log.Warn("Failed to add remote config provider", zap.Error(err))
		} else {
			err = viper.ReadRemoteConfig()
			if err != nil {
				app.log.Warn("Failed to load remote config", zap.Error(err))
			} else {
				configLoaded = true
			}
		}
	}
	if !configLoaded {
		if err = viper.ReadInConfig(); err != nil {
			app.log.Warn("Failed to read in config", zap.Error(err))
			viper.SetConfigType("yaml")
			if err = viper.ReadInConfig(); err != nil {
				app.log.Warn("Failed to read in config", zap.Error(err))
			}
		}
	}
}

func (app *App) ExecuteE(ctx context.Context) error {
	app.preInitRemote()
	app.rootCmd.PreRunE = func(cmd *cobra.Command, args []string) (err error) {
		fs := cmd.Flags()
		verbose, _ := fs.GetCount("verbose")
		instance, _ := fs.GetString("instance")
		configPath, _ := fs.GetString("config")
		if verbose > 1 {
			viper.SetOptions(
				viper.WithLogger(
					slog.New(
						slogzap.Option{Level: slog.LevelDebug, Logger: app.log.Named("viper")}.
							NewZapHandler(),
					),
				),
			)
		}
		if configPath != "" {
			viper.SetConfigFile(configPath)
		} else {
			viper.SetConfigType("json")
			viper.AddConfigPath(".")
			viper.SetConfigName("config_" + instance)
		}
		app.postInitRemote()
		if app.conf != nil {
			err = viper.Unmarshal(app.conf, func(dc *mapstructure.DecoderConfig) { dc.TagName = "json" })
			if err != nil {
				app.log.Error("Failed to unmarshal register config", zap.Error(err))
				return err
			}
		}
		if verbose >= 2 {
			viper.Debug()
		}
		return app.preRunE()
	}

	app.injectVersionCmd()

	if !disableConfigCmd {
		app.injectConfigCmd()
	}
	if app.validate != nil {
		app.injectValidateCmd()
	}
	return app.rootCmd.ExecuteContext(ctx)
}
