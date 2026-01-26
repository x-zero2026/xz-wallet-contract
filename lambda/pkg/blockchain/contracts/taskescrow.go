// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// TaskEscrowMetaData contains all meta data concerning the TaskEscrow contract.
var TaskEscrowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"}],\"name\":\"ExecutorSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalPaid\",\"type\":\"uint256\"}],\"name\":\"MilestonePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"executorAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"creatorRefund\",\"type\":\"uint256\"}],\"name\":\"TaskCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TaskCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"executorAmount\",\"type\":\"uint256\"}],\"name\":\"cancelTask\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"createTask\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"getRemainingAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"getTask\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paidAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"cancelled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextTaskId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"payMilestone\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"}],\"name\":\"setExecutor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"tasks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paidAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"cancelled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TaskEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use TaskEscrowMetaData.ABI instead.
var TaskEscrowABI = TaskEscrowMetaData.ABI

// TaskEscrow is an auto generated Go binding around an Ethereum contract.
type TaskEscrow struct {
	TaskEscrowCaller     // Read-only binding to the contract
	TaskEscrowTransactor // Write-only binding to the contract
	TaskEscrowFilterer   // Log filterer for contract events
}

// TaskEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type TaskEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TaskEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TaskEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaskEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TaskEscrowSession struct {
	Contract     *TaskEscrow       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaskEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TaskEscrowCallerSession struct {
	Contract *TaskEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TaskEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TaskEscrowTransactorSession struct {
	Contract     *TaskEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TaskEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type TaskEscrowRaw struct {
	Contract *TaskEscrow // Generic contract binding to access the raw methods on
}

// TaskEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TaskEscrowCallerRaw struct {
	Contract *TaskEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// TaskEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TaskEscrowTransactorRaw struct {
	Contract *TaskEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTaskEscrow creates a new instance of TaskEscrow, bound to a specific deployed contract.
func NewTaskEscrow(address common.Address, backend bind.ContractBackend) (*TaskEscrow, error) {
	contract, err := bindTaskEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TaskEscrow{TaskEscrowCaller: TaskEscrowCaller{contract: contract}, TaskEscrowTransactor: TaskEscrowTransactor{contract: contract}, TaskEscrowFilterer: TaskEscrowFilterer{contract: contract}}, nil
}

// NewTaskEscrowCaller creates a new read-only instance of TaskEscrow, bound to a specific deployed contract.
func NewTaskEscrowCaller(address common.Address, caller bind.ContractCaller) (*TaskEscrowCaller, error) {
	contract, err := bindTaskEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowCaller{contract: contract}, nil
}

// NewTaskEscrowTransactor creates a new write-only instance of TaskEscrow, bound to a specific deployed contract.
func NewTaskEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*TaskEscrowTransactor, error) {
	contract, err := bindTaskEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowTransactor{contract: contract}, nil
}

// NewTaskEscrowFilterer creates a new log filterer instance of TaskEscrow, bound to a specific deployed contract.
func NewTaskEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*TaskEscrowFilterer, error) {
	contract, err := bindTaskEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowFilterer{contract: contract}, nil
}

// bindTaskEscrow binds a generic wrapper to an already deployed contract.
func bindTaskEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TaskEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaskEscrow *TaskEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaskEscrow.Contract.TaskEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaskEscrow *TaskEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaskEscrow.Contract.TaskEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaskEscrow *TaskEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaskEscrow.Contract.TaskEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaskEscrow *TaskEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaskEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaskEscrow *TaskEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaskEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaskEscrow *TaskEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaskEscrow.Contract.contract.Transact(opts, method, params...)
}

// GetRemainingAmount is a free data retrieval call binding the contract method 0xf6252ff2.
//
// Solidity: function getRemainingAmount(uint256 taskId) view returns(uint256)
func (_TaskEscrow *TaskEscrowCaller) GetRemainingAmount(opts *bind.CallOpts, taskId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "getRemainingAmount", taskId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingAmount is a free data retrieval call binding the contract method 0xf6252ff2.
//
// Solidity: function getRemainingAmount(uint256 taskId) view returns(uint256)
func (_TaskEscrow *TaskEscrowSession) GetRemainingAmount(taskId *big.Int) (*big.Int, error) {
	return _TaskEscrow.Contract.GetRemainingAmount(&_TaskEscrow.CallOpts, taskId)
}

// GetRemainingAmount is a free data retrieval call binding the contract method 0xf6252ff2.
//
// Solidity: function getRemainingAmount(uint256 taskId) view returns(uint256)
func (_TaskEscrow *TaskEscrowCallerSession) GetRemainingAmount(taskId *big.Int) (*big.Int, error) {
	return _TaskEscrow.Contract.GetRemainingAmount(&_TaskEscrow.CallOpts, taskId)
}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowCaller) GetTask(opts *bind.CallOpts, taskId *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "getTask", taskId)

	outstruct := new(struct {
		Creator     common.Address
		Executor    common.Address
		TotalAmount *big.Int
		PaidAmount  *big.Int
		Cancelled   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Executor = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TotalAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PaidAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Cancelled = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowSession) GetTask(taskId *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	return _TaskEscrow.Contract.GetTask(&_TaskEscrow.CallOpts, taskId)
}

// GetTask is a free data retrieval call binding the contract method 0x1d65e77e.
//
// Solidity: function getTask(uint256 taskId) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowCallerSession) GetTask(taskId *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	return _TaskEscrow.Contract.GetTask(&_TaskEscrow.CallOpts, taskId)
}

// NextTaskId is a free data retrieval call binding the contract method 0xfdc3d8d7.
//
// Solidity: function nextTaskId() view returns(uint256)
func (_TaskEscrow *TaskEscrowCaller) NextTaskId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "nextTaskId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextTaskId is a free data retrieval call binding the contract method 0xfdc3d8d7.
//
// Solidity: function nextTaskId() view returns(uint256)
func (_TaskEscrow *TaskEscrowSession) NextTaskId() (*big.Int, error) {
	return _TaskEscrow.Contract.NextTaskId(&_TaskEscrow.CallOpts)
}

// NextTaskId is a free data retrieval call binding the contract method 0xfdc3d8d7.
//
// Solidity: function nextTaskId() view returns(uint256)
func (_TaskEscrow *TaskEscrowCallerSession) NextTaskId() (*big.Int, error) {
	return _TaskEscrow.Contract.NextTaskId(&_TaskEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaskEscrow *TaskEscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaskEscrow *TaskEscrowSession) Owner() (common.Address, error) {
	return _TaskEscrow.Contract.Owner(&_TaskEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaskEscrow *TaskEscrowCallerSession) Owner() (common.Address, error) {
	return _TaskEscrow.Contract.Owner(&_TaskEscrow.CallOpts)
}

// Tasks is a free data retrieval call binding the contract method 0x8d977672.
//
// Solidity: function tasks(uint256 ) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowCaller) Tasks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "tasks", arg0)

	outstruct := new(struct {
		Creator     common.Address
		Executor    common.Address
		TotalAmount *big.Int
		PaidAmount  *big.Int
		Cancelled   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Executor = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TotalAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PaidAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Cancelled = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Tasks is a free data retrieval call binding the contract method 0x8d977672.
//
// Solidity: function tasks(uint256 ) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowSession) Tasks(arg0 *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	return _TaskEscrow.Contract.Tasks(&_TaskEscrow.CallOpts, arg0)
}

// Tasks is a free data retrieval call binding the contract method 0x8d977672.
//
// Solidity: function tasks(uint256 ) view returns(address creator, address executor, uint256 totalAmount, uint256 paidAmount, bool cancelled)
func (_TaskEscrow *TaskEscrowCallerSession) Tasks(arg0 *big.Int) (struct {
	Creator     common.Address
	Executor    common.Address
	TotalAmount *big.Int
	PaidAmount  *big.Int
	Cancelled   bool
}, error) {
	return _TaskEscrow.Contract.Tasks(&_TaskEscrow.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_TaskEscrow *TaskEscrowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaskEscrow.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_TaskEscrow *TaskEscrowSession) Token() (common.Address, error) {
	return _TaskEscrow.Contract.Token(&_TaskEscrow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_TaskEscrow *TaskEscrowCallerSession) Token() (common.Address, error) {
	return _TaskEscrow.Contract.Token(&_TaskEscrow.CallOpts)
}

// CancelTask is a paid mutator transaction binding the contract method 0x1397e04a.
//
// Solidity: function cancelTask(uint256 taskId, uint256 executorAmount) returns()
func (_TaskEscrow *TaskEscrowTransactor) CancelTask(opts *bind.TransactOpts, taskId *big.Int, executorAmount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "cancelTask", taskId, executorAmount)
}

// CancelTask is a paid mutator transaction binding the contract method 0x1397e04a.
//
// Solidity: function cancelTask(uint256 taskId, uint256 executorAmount) returns()
func (_TaskEscrow *TaskEscrowSession) CancelTask(taskId *big.Int, executorAmount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.CancelTask(&_TaskEscrow.TransactOpts, taskId, executorAmount)
}

// CancelTask is a paid mutator transaction binding the contract method 0x1397e04a.
//
// Solidity: function cancelTask(uint256 taskId, uint256 executorAmount) returns()
func (_TaskEscrow *TaskEscrowTransactorSession) CancelTask(taskId *big.Int, executorAmount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.CancelTask(&_TaskEscrow.TransactOpts, taskId, executorAmount)
}

// CreateTask is a paid mutator transaction binding the contract method 0xf6ea2c5d.
//
// Solidity: function createTask(address creator, address executor, uint256 amount) returns(uint256)
func (_TaskEscrow *TaskEscrowTransactor) CreateTask(opts *bind.TransactOpts, creator common.Address, executor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "createTask", creator, executor, amount)
}

// CreateTask is a paid mutator transaction binding the contract method 0xf6ea2c5d.
//
// Solidity: function createTask(address creator, address executor, uint256 amount) returns(uint256)
func (_TaskEscrow *TaskEscrowSession) CreateTask(creator common.Address, executor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.CreateTask(&_TaskEscrow.TransactOpts, creator, executor, amount)
}

// CreateTask is a paid mutator transaction binding the contract method 0xf6ea2c5d.
//
// Solidity: function createTask(address creator, address executor, uint256 amount) returns(uint256)
func (_TaskEscrow *TaskEscrowTransactorSession) CreateTask(creator common.Address, executor common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.CreateTask(&_TaskEscrow.TransactOpts, creator, executor, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address to, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowTransactor) EmergencyWithdraw(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "emergencyWithdraw", to, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address to, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowSession) EmergencyWithdraw(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.EmergencyWithdraw(&_TaskEscrow.TransactOpts, to, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address to, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowTransactorSession) EmergencyWithdraw(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.EmergencyWithdraw(&_TaskEscrow.TransactOpts, to, amount)
}

// PayMilestone is a paid mutator transaction binding the contract method 0x80652153.
//
// Solidity: function payMilestone(uint256 taskId, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowTransactor) PayMilestone(opts *bind.TransactOpts, taskId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "payMilestone", taskId, amount)
}

// PayMilestone is a paid mutator transaction binding the contract method 0x80652153.
//
// Solidity: function payMilestone(uint256 taskId, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowSession) PayMilestone(taskId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.PayMilestone(&_TaskEscrow.TransactOpts, taskId, amount)
}

// PayMilestone is a paid mutator transaction binding the contract method 0x80652153.
//
// Solidity: function payMilestone(uint256 taskId, uint256 amount) returns()
func (_TaskEscrow *TaskEscrowTransactorSession) PayMilestone(taskId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _TaskEscrow.Contract.PayMilestone(&_TaskEscrow.TransactOpts, taskId, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaskEscrow *TaskEscrowTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaskEscrow *TaskEscrowSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaskEscrow.Contract.RenounceOwnership(&_TaskEscrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaskEscrow *TaskEscrowTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaskEscrow.Contract.RenounceOwnership(&_TaskEscrow.TransactOpts)
}

// SetExecutor is a paid mutator transaction binding the contract method 0xc37874cb.
//
// Solidity: function setExecutor(uint256 taskId, address executor) returns()
func (_TaskEscrow *TaskEscrowTransactor) SetExecutor(opts *bind.TransactOpts, taskId *big.Int, executor common.Address) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "setExecutor", taskId, executor)
}

// SetExecutor is a paid mutator transaction binding the contract method 0xc37874cb.
//
// Solidity: function setExecutor(uint256 taskId, address executor) returns()
func (_TaskEscrow *TaskEscrowSession) SetExecutor(taskId *big.Int, executor common.Address) (*types.Transaction, error) {
	return _TaskEscrow.Contract.SetExecutor(&_TaskEscrow.TransactOpts, taskId, executor)
}

// SetExecutor is a paid mutator transaction binding the contract method 0xc37874cb.
//
// Solidity: function setExecutor(uint256 taskId, address executor) returns()
func (_TaskEscrow *TaskEscrowTransactorSession) SetExecutor(taskId *big.Int, executor common.Address) (*types.Transaction, error) {
	return _TaskEscrow.Contract.SetExecutor(&_TaskEscrow.TransactOpts, taskId, executor)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaskEscrow *TaskEscrowTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TaskEscrow.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaskEscrow *TaskEscrowSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaskEscrow.Contract.TransferOwnership(&_TaskEscrow.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaskEscrow *TaskEscrowTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaskEscrow.Contract.TransferOwnership(&_TaskEscrow.TransactOpts, newOwner)
}

// TaskEscrowExecutorSetIterator is returned from FilterExecutorSet and is used to iterate over the raw logs and unpacked data for ExecutorSet events raised by the TaskEscrow contract.
type TaskEscrowExecutorSetIterator struct {
	Event *TaskEscrowExecutorSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskEscrowExecutorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskEscrowExecutorSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskEscrowExecutorSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskEscrowExecutorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskEscrowExecutorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskEscrowExecutorSet represents a ExecutorSet event raised by the TaskEscrow contract.
type TaskEscrowExecutorSet struct {
	TaskId   *big.Int
	Executor common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterExecutorSet is a free log retrieval operation binding the contract event 0x5a7f91ad5127433cb5cf0dad01983c2100378fe0ccfb87cf6c9f007816c84a04.
//
// Solidity: event ExecutorSet(uint256 indexed taskId, address indexed executor)
func (_TaskEscrow *TaskEscrowFilterer) FilterExecutorSet(opts *bind.FilterOpts, taskId []*big.Int, executor []common.Address) (*TaskEscrowExecutorSetIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.FilterLogs(opts, "ExecutorSet", taskIdRule, executorRule)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowExecutorSetIterator{contract: _TaskEscrow.contract, event: "ExecutorSet", logs: logs, sub: sub}, nil
}

// WatchExecutorSet is a free log subscription operation binding the contract event 0x5a7f91ad5127433cb5cf0dad01983c2100378fe0ccfb87cf6c9f007816c84a04.
//
// Solidity: event ExecutorSet(uint256 indexed taskId, address indexed executor)
func (_TaskEscrow *TaskEscrowFilterer) WatchExecutorSet(opts *bind.WatchOpts, sink chan<- *TaskEscrowExecutorSet, taskId []*big.Int, executor []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.WatchLogs(opts, "ExecutorSet", taskIdRule, executorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskEscrowExecutorSet)
				if err := _TaskEscrow.contract.UnpackLog(event, "ExecutorSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExecutorSet is a log parse operation binding the contract event 0x5a7f91ad5127433cb5cf0dad01983c2100378fe0ccfb87cf6c9f007816c84a04.
//
// Solidity: event ExecutorSet(uint256 indexed taskId, address indexed executor)
func (_TaskEscrow *TaskEscrowFilterer) ParseExecutorSet(log types.Log) (*TaskEscrowExecutorSet, error) {
	event := new(TaskEscrowExecutorSet)
	if err := _TaskEscrow.contract.UnpackLog(event, "ExecutorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskEscrowMilestonePaidIterator is returned from FilterMilestonePaid and is used to iterate over the raw logs and unpacked data for MilestonePaid events raised by the TaskEscrow contract.
type TaskEscrowMilestonePaidIterator struct {
	Event *TaskEscrowMilestonePaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskEscrowMilestonePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskEscrowMilestonePaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskEscrowMilestonePaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskEscrowMilestonePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskEscrowMilestonePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskEscrowMilestonePaid represents a MilestonePaid event raised by the TaskEscrow contract.
type TaskEscrowMilestonePaid struct {
	TaskId    *big.Int
	Executor  common.Address
	Amount    *big.Int
	TotalPaid *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMilestonePaid is a free log retrieval operation binding the contract event 0xa5c2138f4ee89547657e692c9d954668da150bf271d1e1382addcb9bb4233c37.
//
// Solidity: event MilestonePaid(uint256 indexed taskId, address indexed executor, uint256 amount, uint256 totalPaid)
func (_TaskEscrow *TaskEscrowFilterer) FilterMilestonePaid(opts *bind.FilterOpts, taskId []*big.Int, executor []common.Address) (*TaskEscrowMilestonePaidIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.FilterLogs(opts, "MilestonePaid", taskIdRule, executorRule)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowMilestonePaidIterator{contract: _TaskEscrow.contract, event: "MilestonePaid", logs: logs, sub: sub}, nil
}

// WatchMilestonePaid is a free log subscription operation binding the contract event 0xa5c2138f4ee89547657e692c9d954668da150bf271d1e1382addcb9bb4233c37.
//
// Solidity: event MilestonePaid(uint256 indexed taskId, address indexed executor, uint256 amount, uint256 totalPaid)
func (_TaskEscrow *TaskEscrowFilterer) WatchMilestonePaid(opts *bind.WatchOpts, sink chan<- *TaskEscrowMilestonePaid, taskId []*big.Int, executor []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.WatchLogs(opts, "MilestonePaid", taskIdRule, executorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskEscrowMilestonePaid)
				if err := _TaskEscrow.contract.UnpackLog(event, "MilestonePaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMilestonePaid is a log parse operation binding the contract event 0xa5c2138f4ee89547657e692c9d954668da150bf271d1e1382addcb9bb4233c37.
//
// Solidity: event MilestonePaid(uint256 indexed taskId, address indexed executor, uint256 amount, uint256 totalPaid)
func (_TaskEscrow *TaskEscrowFilterer) ParseMilestonePaid(log types.Log) (*TaskEscrowMilestonePaid, error) {
	event := new(TaskEscrowMilestonePaid)
	if err := _TaskEscrow.contract.UnpackLog(event, "MilestonePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskEscrowOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TaskEscrow contract.
type TaskEscrowOwnershipTransferredIterator struct {
	Event *TaskEscrowOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskEscrowOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskEscrowOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskEscrowOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskEscrowOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskEscrowOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskEscrowOwnershipTransferred represents a OwnershipTransferred event raised by the TaskEscrow contract.
type TaskEscrowOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaskEscrow *TaskEscrowFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaskEscrowOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaskEscrow.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowOwnershipTransferredIterator{contract: _TaskEscrow.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaskEscrow *TaskEscrowFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TaskEscrowOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaskEscrow.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskEscrowOwnershipTransferred)
				if err := _TaskEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaskEscrow *TaskEscrowFilterer) ParseOwnershipTransferred(log types.Log) (*TaskEscrowOwnershipTransferred, error) {
	event := new(TaskEscrowOwnershipTransferred)
	if err := _TaskEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskEscrowTaskCancelledIterator is returned from FilterTaskCancelled and is used to iterate over the raw logs and unpacked data for TaskCancelled events raised by the TaskEscrow contract.
type TaskEscrowTaskCancelledIterator struct {
	Event *TaskEscrowTaskCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskEscrowTaskCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskEscrowTaskCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskEscrowTaskCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskEscrowTaskCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskEscrowTaskCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskEscrowTaskCancelled represents a TaskCancelled event raised by the TaskEscrow contract.
type TaskEscrowTaskCancelled struct {
	TaskId         *big.Int
	ExecutorAmount *big.Int
	CreatorRefund  *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTaskCancelled is a free log retrieval operation binding the contract event 0x99a4b4cbd1e738b902f770df96bb758050edc0dcd3531a19ba389de9dd89553a.
//
// Solidity: event TaskCancelled(uint256 indexed taskId, uint256 executorAmount, uint256 creatorRefund)
func (_TaskEscrow *TaskEscrowFilterer) FilterTaskCancelled(opts *bind.FilterOpts, taskId []*big.Int) (*TaskEscrowTaskCancelledIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _TaskEscrow.contract.FilterLogs(opts, "TaskCancelled", taskIdRule)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowTaskCancelledIterator{contract: _TaskEscrow.contract, event: "TaskCancelled", logs: logs, sub: sub}, nil
}

// WatchTaskCancelled is a free log subscription operation binding the contract event 0x99a4b4cbd1e738b902f770df96bb758050edc0dcd3531a19ba389de9dd89553a.
//
// Solidity: event TaskCancelled(uint256 indexed taskId, uint256 executorAmount, uint256 creatorRefund)
func (_TaskEscrow *TaskEscrowFilterer) WatchTaskCancelled(opts *bind.WatchOpts, sink chan<- *TaskEscrowTaskCancelled, taskId []*big.Int) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}

	logs, sub, err := _TaskEscrow.contract.WatchLogs(opts, "TaskCancelled", taskIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskEscrowTaskCancelled)
				if err := _TaskEscrow.contract.UnpackLog(event, "TaskCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskCancelled is a log parse operation binding the contract event 0x99a4b4cbd1e738b902f770df96bb758050edc0dcd3531a19ba389de9dd89553a.
//
// Solidity: event TaskCancelled(uint256 indexed taskId, uint256 executorAmount, uint256 creatorRefund)
func (_TaskEscrow *TaskEscrowFilterer) ParseTaskCancelled(log types.Log) (*TaskEscrowTaskCancelled, error) {
	event := new(TaskEscrowTaskCancelled)
	if err := _TaskEscrow.contract.UnpackLog(event, "TaskCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaskEscrowTaskCreatedIterator is returned from FilterTaskCreated and is used to iterate over the raw logs and unpacked data for TaskCreated events raised by the TaskEscrow contract.
type TaskEscrowTaskCreatedIterator struct {
	Event *TaskEscrowTaskCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TaskEscrowTaskCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaskEscrowTaskCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TaskEscrowTaskCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TaskEscrowTaskCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaskEscrowTaskCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaskEscrowTaskCreated represents a TaskCreated event raised by the TaskEscrow contract.
type TaskEscrowTaskCreated struct {
	TaskId   *big.Int
	Creator  common.Address
	Executor common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTaskCreated is a free log retrieval operation binding the contract event 0xf3efa663e8763e3719e2bdc58b7fdc03d43b6fcec97b7bcf371e6a4ea8704488.
//
// Solidity: event TaskCreated(uint256 indexed taskId, address indexed creator, address indexed executor, uint256 amount)
func (_TaskEscrow *TaskEscrowFilterer) FilterTaskCreated(opts *bind.FilterOpts, taskId []*big.Int, creator []common.Address, executor []common.Address) (*TaskEscrowTaskCreatedIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.FilterLogs(opts, "TaskCreated", taskIdRule, creatorRule, executorRule)
	if err != nil {
		return nil, err
	}
	return &TaskEscrowTaskCreatedIterator{contract: _TaskEscrow.contract, event: "TaskCreated", logs: logs, sub: sub}, nil
}

// WatchTaskCreated is a free log subscription operation binding the contract event 0xf3efa663e8763e3719e2bdc58b7fdc03d43b6fcec97b7bcf371e6a4ea8704488.
//
// Solidity: event TaskCreated(uint256 indexed taskId, address indexed creator, address indexed executor, uint256 amount)
func (_TaskEscrow *TaskEscrowFilterer) WatchTaskCreated(opts *bind.WatchOpts, sink chan<- *TaskEscrowTaskCreated, taskId []*big.Int, creator []common.Address, executor []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}

	logs, sub, err := _TaskEscrow.contract.WatchLogs(opts, "TaskCreated", taskIdRule, creatorRule, executorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaskEscrowTaskCreated)
				if err := _TaskEscrow.contract.UnpackLog(event, "TaskCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskCreated is a log parse operation binding the contract event 0xf3efa663e8763e3719e2bdc58b7fdc03d43b6fcec97b7bcf371e6a4ea8704488.
//
// Solidity: event TaskCreated(uint256 indexed taskId, address indexed creator, address indexed executor, uint256 amount)
func (_TaskEscrow *TaskEscrowFilterer) ParseTaskCreated(log types.Log) (*TaskEscrowTaskCreated, error) {
	event := new(TaskEscrowTaskCreated)
	if err := _TaskEscrow.contract.UnpackLog(event, "TaskCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
