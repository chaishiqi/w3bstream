package project

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var (
	errEmptyConfig = errors.New("config is empty")
	errEmptyCode   = errors.New("code is empty")
)

// TODO: prefer protobuf for serialization and deserialization
type Project struct {
	DefaultVersion string    `json:"defaultVersion,omitempty"`
	Configs        []*Config `json:"config"`
}

type Attribute struct {
	Paused                bool
	RequestedProverAmount uint64
}

type SignedKey struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Config struct {
	Version            string      `json:"version"`
	VMTypeID           uint64      `json:"vmTypeID"`
	ProofType          string      `json:"proofType"` // liveness, movement
	SignedKeys         []SignedKey `json:"signedKeys"`
	SignatureAlgorithm string      `json:"signatureAlgorithm"`
	HashAlgorithm      string      `json:"hashAlgorithm"`
	Metadata           string      `json:"metadata,omitempty"`
	Code               string      `json:"code"`
}

func (p *Project) Config(version string) (*Config, error) {
	if len(version) == 0 {
		return p.DefaultConfig()
	}
	for _, c := range p.Configs {
		if c.Version == version {
			return c, nil
		}
	}
	return nil, errors.New("project config not exist")
}

func (p *Project) DefaultConfig() (*Config, error) {
	if len(p.DefaultVersion) > 0 {
		return p.Config(p.DefaultVersion)
	}
	return p.Configs[0], nil
}

func (p *Project) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Project) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return errors.Wrap(err, "failed to unmarshal project")
	}

	if len(p.Configs) == 0 {
		return errEmptyConfig
	}
	for _, c := range p.Configs {
		if err := c.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) validate() error {
	if len(c.Code) == 0 {
		return errEmptyCode
	}
	return nil
}
