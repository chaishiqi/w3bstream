// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package geonet

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// GeonetMetaData contains all meta data concerning the Geonet contract.
var GeonetMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_device\",\"type\":\"address\"}],\"name\":\"tick\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GeonetABI is the input ABI used to generate the binding from.
// Deprecated: Use GeonetMetaData.ABI instead.
var GeonetABI = GeonetMetaData.ABI

// Geonet is an auto generated Go binding around an Ethereum contract.
type Geonet struct {
	GeonetCaller     // Read-only binding to the contract
	GeonetTransactor // Write-only binding to the contract
	GeonetFilterer   // Log filterer for contract events
}

// GeonetCaller is an auto generated read-only Go binding around an Ethereum contract.
type GeonetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeonetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GeonetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeonetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GeonetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GeonetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GeonetSession struct {
	Contract     *Geonet           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GeonetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GeonetCallerSession struct {
	Contract *GeonetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// GeonetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GeonetTransactorSession struct {
	Contract     *GeonetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GeonetRaw is an auto generated low-level Go binding around an Ethereum contract.
type GeonetRaw struct {
	Contract *Geonet // Generic contract binding to access the raw methods on
}

// GeonetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GeonetCallerRaw struct {
	Contract *GeonetCaller // Generic read-only contract binding to access the raw methods on
}

// GeonetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GeonetTransactorRaw struct {
	Contract *GeonetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGeonet creates a new instance of Geonet, bound to a specific deployed contract.
func NewGeonet(address common.Address, backend bind.ContractBackend) (*Geonet, error) {
	contract, err := bindGeonet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Geonet{GeonetCaller: GeonetCaller{contract: contract}, GeonetTransactor: GeonetTransactor{contract: contract}, GeonetFilterer: GeonetFilterer{contract: contract}}, nil
}

// NewGeonetCaller creates a new read-only instance of Geonet, bound to a specific deployed contract.
func NewGeonetCaller(address common.Address, caller bind.ContractCaller) (*GeonetCaller, error) {
	contract, err := bindGeonet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GeonetCaller{contract: contract}, nil
}

// NewGeonetTransactor creates a new write-only instance of Geonet, bound to a specific deployed contract.
func NewGeonetTransactor(address common.Address, transactor bind.ContractTransactor) (*GeonetTransactor, error) {
	contract, err := bindGeonet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GeonetTransactor{contract: contract}, nil
}

// NewGeonetFilterer creates a new log filterer instance of Geonet, bound to a specific deployed contract.
func NewGeonetFilterer(address common.Address, filterer bind.ContractFilterer) (*GeonetFilterer, error) {
	contract, err := bindGeonet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GeonetFilterer{contract: contract}, nil
}

// bindGeonet binds a generic wrapper to an already deployed contract.
func bindGeonet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GeonetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Geonet *GeonetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Geonet.Contract.GeonetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Geonet *GeonetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Geonet.Contract.GeonetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Geonet *GeonetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Geonet.Contract.GeonetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Geonet *GeonetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Geonet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Geonet *GeonetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Geonet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Geonet *GeonetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Geonet.Contract.contract.Transact(opts, method, params...)
}

// Tick is a paid mutator transaction binding the contract method 0x80334736.
//
// Solidity: function tick(address _device) returns()
func (_Geonet *GeonetTransactor) Tick(opts *bind.TransactOpts, _device common.Address) (*types.Transaction, error) {
	return _Geonet.contract.Transact(opts, "tick", _device)
}

// Tick is a paid mutator transaction binding the contract method 0x80334736.
//
// Solidity: function tick(address _device) returns()
func (_Geonet *GeonetSession) Tick(_device common.Address) (*types.Transaction, error) {
	return _Geonet.Contract.Tick(&_Geonet.TransactOpts, _device)
}

// Tick is a paid mutator transaction binding the contract method 0x80334736.
//
// Solidity: function tick(address _device) returns()
func (_Geonet *GeonetTransactorSession) Tick(_device common.Address) (*types.Transaction, error) {
	return _Geonet.Contract.Tick(&_Geonet.TransactOpts, _device)
}
