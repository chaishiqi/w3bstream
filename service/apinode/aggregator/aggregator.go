package aggregator

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	adb "github.com/iotexproject/w3bstream/service/apinode/db"
	"github.com/iotexproject/w3bstream/service/sequencer/api"
	"github.com/iotexproject/w3bstream/smartcontracts/go/geonet"
	"github.com/pkg/errors"
)

func Run(db *adb.DB, sequencerAddr string, interval time.Duration, geoInstance *geonet.Geonet, startGeo bool, prv *ecdsa.PrivateKey, chainID *big.Int) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		ts, err := db.FetchAllTask()
		if err != nil {
			slog.Error("failed to fetch all tasks", "error", err)
			continue
		}
		if len(ts) == 0 {
			continue
		}
		if err := db.CreateTasks(ts); err != nil {
			slog.Error("failed to create tasks", "error", err)
			continue
		}
		time.Sleep(1 * time.Second) // after writing to clickhouse, reading immediately will not return the value.

		tids := make([]string, 0, len(ts))
		pt := make(map[string]*adb.Task)
		for _, t := range ts {
			tids = append(tids, t.TaskID)
			pt[t.ProjectID] = t
		}
		for pid, t := range pt {
			if pid == "942" && startGeo {
				tx, err := geoInstance.Tick(
					&bind.TransactOpts{
						From: crypto.PubkeyToAddress(prv.PublicKey),
						Signer: func(a common.Address, t *types.Transaction) (*types.Transaction, error) {
							return types.SignTx(t, types.NewLondonSigner(chainID), prv)
						},
					},
					common.HexToAddress(t.DeviceID),
				)
				if err != nil {
					slog.Error("failed to call geo contract", "error", err)
				}
				slog.Info("send tx to geo contract success", "hash", tx.Hash().String())
			}
			if err := notify(sequencerAddr, common.HexToHash(t.TaskID)); err != nil {
				slog.Error("failed to notify sequencer", "error", err)
				continue
			}
		}
		if err := db.DeleteTasks(tids); err != nil {
			slog.Error("failed to delete tasks at local", "error", err)
		}
	}
}

func notify(sequencerAddr string, taskID common.Hash) error {
	reqSequencer := &api.CreateTaskReq{TaskID: taskID}
	reqSequencerJ, err := json.Marshal(reqSequencer)
	if err != nil {
		return errors.Wrap(err, "failed to marshal sequencer request")
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/task", sequencerAddr), "application/json", bytes.NewBuffer(reqSequencerJ))
	if err != nil {
		return errors.Wrap(err, "failed to call sequencer service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			err = errors.New(string(body))
		}
		return errors.Wrap(err, "failed to call sequencer service")
	}
	return nil
}
