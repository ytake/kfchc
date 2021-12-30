package command

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/ytake/kfchc/client"
	"github.com/ytake/kfchc/config"
	"github.com/ytake/kfchc/log"
	"github.com/ytake/kfchc/pbdef"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"os"
)

// HealthCheckHandle see https://docs.confluent.io/platform/current/connect/references/restapi.html#get--connectors-(string-name)-status
type HealthCheckHandle struct {
	Logger log.Logger
}

// Run call connectors status
func (cs *HealthCheckHandle) Run(c *cli.Context) error {
	p := c.String(config.FlagJsonConfigPath)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	req := &pbdef.ConnectorConfig{}
	err = protojson.Unmarshal(b, req)
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, v := range req.Servers {
		kc, err := client.NewKafkaConnect(v.ConnectServer, cs.Logger)
		if err != nil {
			cs.Logger.Error("kafka connect error", err)
			return err
		}
		for _, con := range v.Connectors {
			ccs := &client.CurrentStatus{ConnectorName: con}
			a := kc.Get(ccs, ctx)
			res := <-a
			if res.Err != nil {
				return res.Err
			}
			if res.ConnectorNotFound.ErrorCode != 0 {
				cs.Logger.Error(res.ConnectorNotFound.Message)
			}
			ts := res.ConnectorStatus.Tasks
			if len(ts) != 0 {
				for _, v := range ts {
					if v.IsFailed() {
						cs.Logger.Error(fmt.Sprintf("connector error %s, %s", con, v.Trace))
					}
				}
			}
		}
	}
	return nil
}
