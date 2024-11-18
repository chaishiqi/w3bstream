package monitor

import (
	"bytes"
	"context"
	"log/slog"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/metrics"
	"github.com/iotexproject/w3bstream/smartcontracts/go/minter"
	"github.com/iotexproject/w3bstream/smartcontracts/go/project"
	"github.com/iotexproject/w3bstream/smartcontracts/go/projectdevice"
	"github.com/iotexproject/w3bstream/smartcontracts/go/prover"
	"github.com/iotexproject/w3bstream/smartcontracts/go/taskmanager"
)

type (
	ScannedBlockNumber       func() (uint64, error)
	UpsertScannedBlockNumber func(number uint64) error
	AssignTask               func(taskID common.Hash, prover common.Address) error
	SettleTask               func(taskID, tx common.Hash) error
	UpsertProject            func(projectID uint64, uri string, hash common.Hash) error
	UpsertProver             func(addr common.Address) error
	UpsertProjectDevice      func(projectID uint64, address common.Address) error
	DeleteProjectDevice      func(projectID uint64, address common.Address) error
)

type Handler struct {
	ScannedBlockNumber
	UpsertScannedBlockNumber
	AssignTask
	SettleTask
	UpsertProject
	UpsertProver
	UpsertProjectDevice
	DeleteProjectDevice
}

type ContractAddr struct {
	Prover        common.Address
	Project       common.Address
	Minter        common.Address
	TaskManager   common.Address
	ProjectDevice common.Address
}

type contract struct {
	h                     *Handler
	addr                  *ContractAddr
	beginningBlockNumber  uint64
	listStepSize          uint64
	watchInterval         time.Duration
	client                *ethclient.Client
	minterInstance        *minter.Minter
	taskManagerInstance   *taskmanager.Taskmanager
	proverInstance        *prover.Prover
	projectInstance       *project.Project
	projectDeviceInstance *projectdevice.Projectdevice
}

var (
	taskAssignedTopic         = crypto.Keccak256Hash([]byte("TaskAssigned(uint256,bytes32,address,uint256)"))
	taskSettledTopic          = crypto.Keccak256Hash([]byte("TaskSettled(uint256,bytes32,address)"))
	projectConfigUpdatedTopic = crypto.Keccak256Hash([]byte("ProjectConfigUpdated(uint256,string,bytes32)"))
	proverSetTopic            = crypto.Keccak256Hash([]byte("BeneficiarySet(address,address)"))
	deviceApproveTopic        = crypto.Keccak256Hash([]byte("Approve(uint256,address)"))
	deviceUnapproveTopic      = crypto.Keccak256Hash([]byte("Unapprove(uint256,address)"))
)

var allTopic = []common.Hash{
	taskAssignedTopic,
	taskSettledTopic,
	projectConfigUpdatedTopic,
	proverSetTopic,
	deviceApproveTopic,
	deviceUnapproveTopic,
}

var emptyAddr = common.Address{}

func (a *ContractAddr) all() []common.Address {
	all := make([]common.Address, 0, 5)
	if !bytes.Equal(a.Minter[:], emptyAddr[:]) {
		all = append(all, a.Minter)
	}
	if !bytes.Equal(a.Project[:], emptyAddr[:]) {
		all = append(all, a.Project)
	}
	if !bytes.Equal(a.Prover[:], emptyAddr[:]) {
		all = append(all, a.Prover)
	}
	if !bytes.Equal(a.TaskManager[:], emptyAddr[:]) {
		all = append(all, a.TaskManager)
	}
	if !bytes.Equal(a.ProjectDevice[:], emptyAddr[:]) {
		all = append(all, a.ProjectDevice)
	}
	return all
}

func (c *contract) processLogs(logs []types.Log) error {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber != logs[j].BlockNumber {
			return logs[i].BlockNumber < logs[j].BlockNumber
		}
		return logs[i].TxIndex < logs[j].TxIndex
	})

	for _, l := range logs {
		switch l.Topics[0] {
		case taskAssignedTopic:
			if c.taskManagerInstance == nil || c.h.AssignTask == nil {
				continue
			}
			e, err := c.taskManagerInstance.ParseTaskAssigned(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse task assigned event")
			}
			if err := c.h.AssignTask(e.TaskId, e.Prover); err != nil {
				return err
			}
			metrics.AssignedTaskMtc.WithLabelValues(e.ProjectId.String()).Inc()
		case taskSettledTopic:
			if c.taskManagerInstance == nil || c.h.SettleTask == nil {
				continue
			}
			e, err := c.taskManagerInstance.ParseTaskSettled(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse task settled event")
			}
			if err := c.h.SettleTask(e.TaskId, l.TxHash); err != nil {
				return err
			}
			metrics.SucceedTaskNumMtc.WithLabelValues(e.ProjectId.String()).Inc()
		case projectConfigUpdatedTopic:
			if c.projectInstance == nil || c.h.UpsertProject == nil {
				continue
			}
			e, err := c.projectInstance.ParseProjectConfigUpdated(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse project config updated event")
			}
			if err := c.h.UpsertProject(e.ProjectId.Uint64(), e.Uri, e.Hash); err != nil {
				return err
			}
		case proverSetTopic:
			if c.proverInstance == nil || c.h.UpsertProver == nil {
				continue
			}
			e, err := c.proverInstance.ParseBeneficiarySet(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse prover set event")
			}
			if err := c.h.UpsertProver(e.Prover); err != nil {
				return err
			}
		case deviceApproveTopic:
			if c.projectDeviceInstance == nil || c.h.UpsertProjectDevice == nil {
				continue
			}
			e, err := c.projectDeviceInstance.ParseApprove(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse device approve event")
			}
			if err := c.h.UpsertProjectDevice(e.ProjectId.Uint64(), e.Device); err != nil {
				return err
			}
		case deviceUnapproveTopic:
			if c.projectDeviceInstance == nil || c.h.DeleteProjectDevice == nil {
				continue
			}
			e, err := c.projectDeviceInstance.ParseUnapprove(l)
			if err != nil {
				return errors.Wrap(err, "failed to parse device unapprove event")
			}
			if err := c.h.DeleteProjectDevice(e.ProjectId.Uint64(), e.Device); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *contract) list() (uint64, error) {
	head := c.beginningBlockNumber
	h, err := c.h.ScannedBlockNumber()
	if err != nil {
		return 0, err
	}
	head = max(head, h)

	query := ethereum.FilterQuery{
		Addresses: c.addr.all(),
		Topics:    [][]common.Hash{allTopic},
	}
	ctx := context.Background()
	from := head + 1
	to := from
	for {
		header, err := c.client.HeaderByNumber(ctx, nil)
		if err != nil {
			return 0, errors.Wrap(err, "failed to retrieve latest block header")
		}
		currentHead := header.Number.Uint64()
		to = from + c.listStepSize
		if to > currentHead {
			to = currentHead
		}
		if from > to {
			break
		}
		slog.Debug("listing chain", "from", from, "to", to)
		query.FromBlock = new(big.Int).SetUint64(from)
		query.ToBlock = new(big.Int).SetUint64(to)
		logs, err := c.client.FilterLogs(ctx, query)
		if err != nil {
			return 0, errors.Wrap(err, "failed to filter contract logs")
		}
		if err := c.processLogs(logs); err != nil {
			return 0, err
		}
		if err := c.h.UpsertScannedBlockNumber(to); err != nil {
			return 0, err
		}
		metrics.SyncedBlockHeightMtc.Set(float64(to))
		from = to + 1
	}
	slog.Info("contract data synchronization completed", "current_height", to)
	return to, nil
}

func (c *contract) watch(listedBlockNumber uint64) {
	scannedBlockNumber := listedBlockNumber
	query := ethereum.FilterQuery{
		Addresses: c.addr.all(),
		Topics:    [][]common.Hash{allTopic},
	}
	ticker := time.NewTicker(c.watchInterval)

	go func() {
		for range ticker.C {
			target := scannedBlockNumber + 1
			currheight, err := c.client.BlockNumber(context.Background())
			if err != nil {
				slog.Error("failed to get current block number", "error", err)
				continue
			}
			if target > currheight {
				continue
			}

			query.FromBlock = new(big.Int).SetUint64(target)
			query.ToBlock = new(big.Int).SetUint64(target)
			logs, err := c.client.FilterLogs(context.Background(), query)
			if err != nil {
				if !strings.Contains(err.Error(), "start block > tip height") {
					slog.Error("failed to filter contract logs", "error", err)
				}
				continue
			}
			slog.Debug("listing chain", "from", target, "to", target)
			if err := c.processLogs(logs); err != nil {
				slog.Error("failed to process logs", "error", err)
				continue
			}
			if err := c.h.UpsertScannedBlockNumber(target); err != nil {
				slog.Error("failed to upsert scanned block number", "error", err)
				continue
			}
			scannedBlockNumber = target
		}
	}()
}

func Run(h *Handler, addr *ContractAddr, beginningBlockNumber uint64, chainEndpoint string) error {
	client, err := ethclient.Dial(chainEndpoint)
	if err != nil {
		return errors.Wrap(err, "failed to dial chain endpoint")
	}

	c := &contract{
		h:                    h,
		addr:                 addr,
		beginningBlockNumber: beginningBlockNumber,
		listStepSize:         500,
		watchInterval:        1 * time.Second,
		client:               client,
	}

	if !bytes.Equal(addr.Minter[:], emptyAddr[:]) {
		minterInstance, err := minter.NewMinter(addr.Minter, client)
		if err != nil {
			return errors.Wrap(err, "failed to new minter contract instance")
		}
		c.minterInstance = minterInstance
	}
	if !bytes.Equal(addr.TaskManager[:], emptyAddr[:]) {
		taskManagerInstance, err := taskmanager.NewTaskmanager(addr.TaskManager, client)
		if err != nil {
			return errors.Wrap(err, "failed to new task manager contract instance")
		}
		c.taskManagerInstance = taskManagerInstance
	}
	if !bytes.Equal(addr.Prover[:], emptyAddr[:]) {
		proverInstance, err := prover.NewProver(addr.Prover, client)
		if err != nil {
			return errors.Wrap(err, "failed to new prover contract instance")
		}
		c.proverInstance = proverInstance
	}
	if !bytes.Equal(addr.Project[:], emptyAddr[:]) {
		projectInstance, err := project.NewProject(addr.Project, client)
		if err != nil {
			return errors.Wrap(err, "failed to new project contract instance")
		}
		c.projectInstance = projectInstance
	}
	if !bytes.Equal(addr.ProjectDevice[:], emptyAddr[:]) {
		projectDeviceInstance, err := projectdevice.NewProjectdevice(addr.ProjectDevice, client)
		if err != nil {
			return errors.Wrap(err, "failed to new project device contract instance")
		}
		c.projectDeviceInstance = projectDeviceInstance
	}

	listedBlockNumber, err := c.list()
	if err != nil {
		return err
	}
	go c.watch(listedBlockNumber)

	return nil
}
