package command

import (
	"github.com/urfave/cli/v2"
	"github.com/ytake/kfchc/config"
	"github.com/ytake/kfchc/pbdef"
	"google.golang.org/protobuf/encoding/protojson"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type GenerateServerConfig struct{}

func (gsc *GenerateServerConfig) stubServer() string {
	return "http://127.0.0.1:8083"
}

func (gsc *GenerateServerConfig) stubConnectors() []string {
	return []string{"replace_me", "replace_me"}
}

// Run call connectors status
func (gsc *GenerateServerConfig) Run(c *cli.Context) error {
	var cc []*pbdef.ConnectorConfig_Servers
	m := &protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}
	b, err := m.Marshal(&pbdef.ConnectorConfig{
		Servers: append(cc, &pbdef.ConnectorConfig_Servers{
			ConnectServer: gsc.stubServer(),
			Connectors:    gsc.stubConnectors(),
		})})
	if err != nil {
		return err
	}
	// protojson.Unmarshal(&pbdef.ConnectorConfig{})
	p := filepath.Join(c.String(config.FlagOutputPath), config.JsonConfigFileName)
	err = ioutil.WriteFile(p, b, fs.ModePerm)
	if err != nil {
		return err
	}
	perm32, _ := strconv.ParseUint("0644", 8, 32)
	err = os.Chmod(p, os.FileMode(perm32))
	if err != nil {
		return err
	}
	return nil
}
