package apinode

import (
	"context"
	"crypto/ecdsa"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iotexproject/w3bstream/monitor"
	"github.com/iotexproject/w3bstream/service/apinode/aggregator"
	"github.com/iotexproject/w3bstream/service/apinode/api"
	"github.com/iotexproject/w3bstream/service/apinode/config"
	"github.com/iotexproject/w3bstream/service/apinode/db"
	"github.com/iotexproject/w3bstream/smartcontracts/go/geonet"
	"github.com/pkg/errors"
)

type APINode struct {
	cfg *config.Config
	db  *db.DB
	prv *ecdsa.PrivateKey
}

func NewAPINode(cfg *config.Config, db *db.DB, prv *ecdsa.PrivateKey) *APINode {
	return &APINode{
		cfg: cfg,
		db:  db,
		prv: prv,
	}
}

func (n *APINode) Start() error {
	client, err := ethclient.Dial(n.cfg.ChainEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to dial chain endpoint")
	}
	if err := monitor.Run(
		&monitor.Handler{
			ScannedBlockNumber:       n.db.ScannedBlockNumber,
			UpsertScannedBlockNumber: n.db.UpsertScannedBlockNumber,
			AssignTask:               n.db.UpsertAssignedTask,
			SettleTask:               n.db.UpsertSettledTask,
			UpsertProjectDevice:      n.db.UpsertProjectDevice,
		},
		&monitor.ContractAddr{
			TaskManager: common.HexToAddress(n.cfg.TaskManagerContractAddr),
			IoID:        common.HexToAddress(n.cfg.IoIDContractAddr),
		},
		n.cfg.BeginningBlockNumber,
		client,
	); err != nil {
		return errors.Wrap(err, "failed to run contract monitor")
	}

	geoInstance, err := geonet.NewGeonet(common.HexToAddress(n.cfg.GeoContractAddr), client)
	if err != nil {
		return errors.Wrap(err, "failed to new geo contract instance")
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get chain id")
	}

	go aggregator.Run(n.db, n.cfg.SequencerServiceEndpoint, time.Duration(n.cfg.TaskAggregatorIntervalSecond)*time.Second, geoInstance, n.cfg.StartGeo, n.prv, chainID)

	go func() {
		if err := api.Run(n.db, n.cfg.ServiceEndpoint, n.cfg.SequencerServiceEndpoint, n.cfg.ProverServiceEndpoint); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (n *APINode) Stop() error {
	return nil
}
