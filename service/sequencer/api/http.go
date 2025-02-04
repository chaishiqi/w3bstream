package api

import (
	"log/slog"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/metrics"
	"github.com/iotexproject/w3bstream/service/sequencer/db"
	"github.com/iotexproject/w3bstream/task"
)

type RetrieveTask func(taskIDs []common.Hash) ([]*task.Task, error)

type CreateTaskReq struct {
	TaskID common.Hash `json:"taskID"                     binding:"required"`
}

type errResp struct {
	Error string `json:"error,omitempty"`
}

func newErrResp(err error) *errResp {
	return &errResp{Error: err.Error()}
}

type httpServer struct {
	engine   *gin.Engine
	db       *db.DB
	retrieve RetrieveTask
}

func (s *httpServer) createTask(c *gin.Context) {
	req := &CreateTaskReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		slog.Error("failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, newErrResp(errors.Wrap(err, "invalid request payload")))
		return
	}

	ts, err := s.retrieve([]common.Hash{req.TaskID})
	if err != nil {
		slog.Error("failed to retrieve task", "error", err)
		c.JSON(http.StatusBadRequest, newErrResp(errors.Wrap(err, "failed to retrieve task")))
		return
	}

	t := ts[0]
	metrics.NewTaskMtc.WithLabelValues(t.ProjectID.String()).Inc()

	if err := s.db.CreateTask(req.TaskID); err != nil {
		slog.Error("failed to create task", "error", err)
		c.JSON(http.StatusInternalServerError, newErrResp(errors.Wrap(err, "failed to create task")))
		return
	}

	slog.Info("get a new task", "project_id", t.ProjectID, "task_id", req.TaskID.String())
	c.Status(http.StatusOK)
}

// this func will block caller
func Run(address string, db *db.DB, retrieve RetrieveTask) error {
	s := &httpServer{
		retrieve: retrieve,
		db:       db,
		engine:   gin.Default(),
	}

	s.engine.POST("/task", s.createTask)
	metrics.RegisterMetrics(s.engine)

	if err := s.engine.Run(address); err != nil {
		slog.Error("failed to start http server", "address", address, "error", err)
		return errors.Wrap(err, "could not start http server; check if the address is in use or network is accessible")
	}
	return nil
}
