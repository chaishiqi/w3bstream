package config

import (
	"os"
	"strconv"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/iotexproject/w3bstream/util/env"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestConfig_Init(t *testing.T) {
	r := require.New(t)

	t.Run("UseEnvConfig", func(t *testing.T) {
		os.Clearenv()
		expected := Config{
			ChainEndpoint:           "http://iotex.chainendpoint.io",
			BootNodeMultiAddr:       "/dns4/a.b.com/tcp/1000/ipfs/123123123",
			IoTeXChainID:            100,
			TaskProcessingBandwidth: 20,
			ProverContractAddr:      "0x",
			DatasourceDSN:           "111",
		}

		_ = os.Setenv("CHAIN_ENDPOINT", expected.ChainEndpoint)
		_ = os.Setenv("TASK_PROCESSING_BANDWIDTH", strconv.Itoa(expected.TaskProcessingBandwidth))
		_ = os.Setenv("DATASOURCE_DSN", expected.DatasourceDSN)
		_ = os.Setenv("BOOTNODE_MULTIADDR", expected.BootNodeMultiAddr)
		_ = os.Setenv("IOTEX_CHAINID", strconv.Itoa(expected.IoTeXChainID))
		_ = os.Setenv("PROVER_CONTRACT_ADDRESS", expected.ProverContractAddr)
		_ = os.Setenv("LOCAL_DB_DIRECTORY", expected.LocalDBDir)

		c := &Config{}
		r.Nil(c.init())
		r.Equal(*c, expected)
	})

	t.Run("CatchPanicCausedByEmptyRequiredEnvVar", func(t *testing.T) {
		os.Clearenv()

		c := &Config{}
		defer func() {
			r.NotNil(recover())
		}()
		_ = c.init()
	})

	t.Run("FailedToParseEnv", func(t *testing.T) {
		p := NewPatches()
		defer p.Reset()

		p.ApplyFuncReturn(env.ParseEnv, errors.New(t.Name()))

		c := &Config{}
		err := c.init()
		r.ErrorContains(err, t.Name())
	})
}
