package project

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/w3bstream/util/filefetcher"
	"github.com/iotexproject/w3bstream/util/ipfs"
)

func TestProjectMeta_FetchProjectRawData_init(t *testing.T) {
	r := require.New(t)
	p := gomonkey.NewPatches()
	defer p.Reset()

	t.Run("InvalidUri", func(t *testing.T) {
		p = p.ApplyFuncReturn(url.Parse, nil, errors.New(t.Name()))

		_, err := (&filefetcher.Filedescriptor{}).FetchFile()
		r.ErrorContains(err, t.Name())
	})
}

func TestProjectMeta_FetchProjectFile_http(t *testing.T) {
	r := require.New(t)
	p := gomonkey.NewPatches()
	defer p.Reset()

	c := Project{}
	jc, err := json.Marshal(c)
	r.NoError(err)

	h := sha256.New()
	_, err = h.Write(jc)
	r.NoError(err)
	hash := h.Sum(nil)

	fd := &filefetcher.Filedescriptor{
		Uri:  "https://test.com/project_config",
		Hash: [32]byte(hash),
	}

	t.Run("FailedToGetHTTP", func(t *testing.T) {
		p = p.ApplyFuncReturn(http.Get, nil, errors.New(t.Name()))

		_, err := fd.FetchFile()
		r.ErrorContains(err, t.Name())
	})
	t.Run("FailedToIOReadAll", func(t *testing.T) {
		p = p.ApplyFuncReturn(http.Get, &http.Response{
			Body: io.NopCloser(bytes.NewReader(jc)),
		}, nil)
		p = p.ApplyFuncReturn(io.ReadAll, nil, errors.New(t.Name()))

		_, err := fd.FetchFile()
		r.ErrorContains(err, t.Name())
	})
	t.Run("HashMismatch", func(t *testing.T) {
		p = p.ApplyFuncReturn(io.ReadAll, jc, nil)

		nfd := *fd
		nfd.Hash = [32]byte{}
		_, err := nfd.FetchFile()
		r.ErrorContains(err, "failed to validate project file hash")
	})
	t.Run("Success", func(t *testing.T) {
		_, err := fd.FetchFile()
		r.NoError(err)
	})
}

func TestProjectMeta_FetchProjectFile_ipfs(t *testing.T) {
	r := require.New(t)
	p := gomonkey.NewPatches()
	defer p.Reset()

	fd := &filefetcher.Filedescriptor{
		Uri: "ipfs://test.com/123",
	}
	t.Run("FailedToGetIPFS", func(t *testing.T) {
		p = p.ApplyMethodReturn(&ipfs.IPFS{}, "Cat", nil, errors.New(t.Name()))

		_, err := fd.FetchFile()
		r.ErrorContains(err, t.Name())
	})
}

func TestProjectMeta_FetchProjectFile_default(t *testing.T) {
	r := require.New(t)
	p := gomonkey.NewPatches()
	defer p.Reset()

	fd := &filefetcher.Filedescriptor{
		Uri: "test.com/123",
	}

	t.Run("FailedToGetIPFS", func(t *testing.T) {
		_, err := fd.FetchFile()
		r.ErrorContains(err, "invalid project file uri")
	})
}

func TestProject_GetConfig(t *testing.T) {
	r := require.New(t)

	conf := &Config{
		Version: "0.1",
	}
	project := &Project{
		Configs: []*Config{conf},
	}

	t.Run("Success", func(t *testing.T) {
		c, err := project.Config("0.1")
		r.NoError(err)
		r.Equal(conf, c)
	})

	t.Run("NotExist", func(t *testing.T) {
		_, err := project.Config("0.3")
		r.ErrorContains(err, "project config not exist")
	})
}

func TestConfig_Validate(t *testing.T) {
	r := require.New(t)

	config := &Config{
		VMTypeID: 1,
		Code:     "testCode",
	}

	t.Run("EmptyCode", func(t *testing.T) {
		c := *config
		c.Code = ""
		err := c.validate()
		r.EqualError(err, errEmptyCode.Error())
	})

	t.Run("Success", func(t *testing.T) {
		err := config.validate()
		r.NoError(err)
	})
}
