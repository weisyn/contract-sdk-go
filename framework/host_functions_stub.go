//go:build !tinygo && !(js && wasm)

//nolint:golint,unused,U1000 // 该文件为非WASM环境提供存根实现，所有占位函数都是未使用的

package framework

// 该文件为非TinyGo/非WASM环境提供空实现，使得 go build ./... 能通过编译。
//
// 注意：这些类型定义与 contract_base.go 中的定义保持一致
// 但由于 build tag 的限制，需要在这里重新定义

// 基础类型定义（非WASM环境）
type Address [20]byte
type Hash [32]byte
type TokenID string
type Amount uint64

// ContractParams 合约参数（非WASM环境）
type ContractParams struct {
	data []byte
}

// NewContractParams 创建合约参数（非WASM环境）
func NewContractParams(data []byte) *ContractParams {
	return &ContractParams{data: data}
}

// Event 事件（非WASM环境）
type Event struct {
	Name string
	Data map[string]interface{}
}

// NewEvent 创建事件（非WASM环境）
func NewEvent(name string) *Event {
	return &Event{
		Name: name,
		Data: make(map[string]interface{}),
	}
}

// 错误码定义（非WASM环境）
const (
	SUCCESS                    = 0
	ERROR_INVALID_PARAMS       = 1
	ERROR_INSUFFICIENT_BALANCE = 2
	ERROR_UNAUTHORIZED         = 3
	ERROR_NOT_FOUND            = 4
	ERROR_ALREADY_EXISTS       = 5
	ERROR_EXECUTION_FAILED     = 6
	ERROR_INVALID_STATE        = 7
	ERROR_TIMEOUT              = 8
	ERROR_NOT_IMPLEMENTED      = 9
	ERROR_PERMISSION_DENIED    = 10
	ERROR_UNKNOWN              = 999
)

// 注意：这些实现仅用于宿主环境的编译占位，不会在合约WASM中使用。

// ABI 版本函数占位
func getABIVersion() uint32 { return 0x00010000 } // v1.0.0

// 基础环境函数占位
//
//nolint:unused // 这些是占位函数，用于非WASM环境的编译占位
func getCaller(addrPtr uint32) uint32                           { return 0 }
func getContractAddress(addrPtr uint32) uint32                  { return 0 }
func setReturnData(dataPtr uint32, dataLen uint32) uint32       { return SUCCESS }
func emitEvent(eventPtr uint32, eventLen uint32) uint32         { return SUCCESS }
func getContractInitParams(bufPtr uint32, bufLen uint32) uint32 { return 0 }
func getTimestamp() uint64                                      { return 0 }
func getBlockHeight() uint64                                    { return 0 }
func getBlockHash(height uint64, hashPtr uint32) uint32         { return SUCCESS }
func getMerkleRoot(height uint64, rootPtr uint32) uint32        { return SUCCESS }
func getStateRoot(height uint64, rootPtr uint32) uint32         { return SUCCESS }
func getMinerAddress(height uint64, addrPtr uint32) uint32      { return SUCCESS }
func getTxHash(hashPtr uint32) uint32                           { return SUCCESS }
func getTxIndex() uint32                                        { return 0 }

// UTXO操作函数占位
//
//nolint:unused // 这些是占位函数，用于非WASM环境的编译占位
func createUTXOOutput(recipientPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32 {
	return SUCCESS
}
func executeUTXOTransfer(fromPtr uint32, toPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32 {
	return SUCCESS
}
func queryUTXOBalance(addressPtr uint32, tokenIDPtr uint32, tokenIDLen uint32) uint64 { return 0 }

// 状态查询占位
//
//nolint:unused // 这些是占位函数，用于非WASM环境的编译占位
func stateGet(keyPtr uint32, keyLen uint32, valuePtr uint32, valueLen uint32) uint32 { return SUCCESS }

// ⚠️ **已删除**：statePut 和 stateExists 宿主函数声明
// 原因：违背WES架构原则，EUTXO模型无全局状态存储

// 内存管理占位
//
//nolint:unused // 这是占位函数，用于非WASM环境的编译占位
func malloc(size uint32) uint32 { return 1 }

// 地址编码转换函数占位
//
//nolint:unused // 这些是占位函数，用于非WASM环境的编译占位
func addressBytesToBase58(addrPtr uint32, resultPtr uint32, maxLen uint32) uint32      { return 0 }
func addressBase58ToBytes(base58Ptr uint32, base58Len uint32, resultPtr uint32) uint32 { return 0 }

// ==================== 导出封装函数（宿主占位实现） ====================

// GetABIVersion 获取 ABI 版本（占位实现）
func GetABIVersion() uint32 { return getABIVersion() }

// CheckABICompatibility 检查 ABI 兼容性（占位实现）
func CheckABICompatibility(expectedVersion uint32) error { return nil }

// GetCaller 获取合约调用者地址（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func GetCaller() Address { return Address{} }

// GetContractAddress 获取当前合约地址（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func GetContractAddress() Address { return Address{} }

// GetTimestamp 获取当前时间戳（占位实现）
func GetTimestamp() uint64 { return 0 }

// GetBlockHeight 获取当前区块高度（占位实现）
func GetBlockHeight() uint64 { return 0 }

// GetBlockHash 获取指定高度的区块哈希（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func GetBlockHash(height uint64) Hash { return Hash{} }

// GetContractParams 获取合约调用参数（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func GetContractParams() *ContractParams { return NewContractParams([]byte{}) }

// SetReturnData 设置返回数据（占位实现）
func SetReturnData(data []byte) error { return nil }

// SetReturnString 设置字符串返回数据（占位实现）
func SetReturnString(s string) error { return nil }

// SetReturnJSON 设置JSON返回数据（占位实现）
func SetReturnJSON(obj interface{}) error { return nil }

// EmitEvent 发出事件（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func EmitEvent(event *Event) error { return nil }

// EmitSimpleEvent 发出简单事件（占位实现）
func EmitSimpleEvent(name string, data map[string]string) error { return nil }

// CreateUTXO 创建UTXO输出（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func CreateUTXO(recipient Address, amount Amount, tokenID TokenID) error { return nil }

// TransferUTXO 执行UTXO转移（占位实现）
//nolint:golint // 类型定义在文件前面，linter误报
// ⚠️ **已删除**：TransferUTXO
// 原因：违背WES"无业务语义"架构原则

// QueryBalance 查询UTXO余额（占位实现）
//
//nolint:golint // 类型定义在文件前面，linter误报
func QueryBalance(address Address, tokenID TokenID) Amount { return 0 }

// GetState 获取状态数据（占位实现）
func GetState(key string) ([]byte, error) { return []byte{}, nil }

// GetStateFromChain 从链上查询历史状态（占位实现）
func GetStateFromChain(stateID []byte) ([]byte, uint64, error) {
	return []byte{}, 0, nil
}

// GetStateVersion 获取状态的当前版本号（占位实现）
func GetStateVersion(stateID []byte) (uint64, error) {
	return 0, nil
}

// IncrementStateVersion 递增状态版本号（占位实现）
func IncrementStateVersion(stateID []byte) (uint64, error) {
	return 1, nil
}

// ⚠️ **已删除**：PutState 和 StateExists
// 原因：违背WES架构原则，EUTXO模型无全局状态存储

// Malloc 分配内存（占位实现）
func Malloc(size uint32) uint32 { return malloc(size) }

// AppendStateOutputSimple 追加状态输出（占位实现）
func AppendStateOutputSimple(stateID []byte, version uint64, execHash []byte, parentHash []byte) (uint32, error) {
	return 0, nil
}

// AppendStateOutput 追加状态输出（占位实现）
func AppendStateOutput(stateID []byte, version uint64, execHash []byte, zkProof []byte, parentHash []byte) (uint32, error) {
	return 0, nil
}

// AppendResourceOutput 追加资源输出（占位实现）
func AppendResourceOutput(resourceBytes []byte, owner Address, lockingBytes []byte) (uint32, error) {
	return 0, nil
}

// BatchCreateOutputsSimple 批量创建资产输出（占位实现）
func BatchCreateOutputsSimple(items []struct {
	Recipient []byte
	Amount    uint64
	TokenID   []byte
}) (uint32, error) {
	return 0, nil
}

// StateGet 状态只读查询（占位实现）
func StateGet(key []byte) ([]byte, error) {
	return nil, nil
}

// CreateAssetOutputWithLock 创建带锁定条件的资产输出（占位实现）
func CreateAssetOutputWithLock(recipient []byte, amount uint64, tokenID []byte, lockingConditionsBytes []byte) (uint32, error) {
	return 0, nil
}

// ⚠️ **已删除**：ExecuteUTXOTransferEx
// 原因：违背WES"无业务语义"架构原则
// 替代：使用 TransactionBuilder.Transfer() 或 helpers/token/Transfer

// GetMerkleRoot 获取指定高度区块的交易Merkle根（占位实现）
func GetMerkleRoot(height uint64) Hash { return Hash{} }

// GetStateRoot 获取指定高度区块的状态根（占位实现）
func GetStateRoot(height uint64) Hash { return Hash{} }

// GetMinerAddress 获取指定高度区块的矿工地址（占位实现）
func GetMinerAddress(height uint64) Address { return Address{} }

// GetTxHash 获取当前执行交易的哈希（占位实现）
func GetTxHash() Hash { return Hash{} }

// GetTxIndex 获取当前交易在区块内的索引（占位实现）
func GetTxIndex() uint32 { return 0 }
