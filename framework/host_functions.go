//go:build tinygo || (js && wasm)

package framework

// ==================== WES 宿主函数Go绑定库 ====================
//
// 🌟 **设计理念**：为WES合约提供统一的宿主函数访问接口
//
// 🎯 **核心特性**：
// - 封装所有WES宿主函数的底层调用
// - 提供类型安全的Go语言接口
// - 内置错误处理和参数验证
// - 支持UTXO操作、事件发出、环境查询等
// - 简化合约开发的复杂性
//

// ==================== 宿主函数原始声明 ====================

// 🔧 注意：TinyGo 0.31+ 要求 //go:wasmimport 函数必须是声明，不能有函数体
// 这些函数在WASM编译时会被链接到宿主函数
//
// 📋 版本兼容性：
// - TinyGo 0.30及以下：不兼容（需要函数体 { return 0 }）
// - TinyGo 0.31及以上：完全兼容（只需函数声明）
//
// 💡 如果您使用旧版本TinyGo，请升级到0.31+：
//   brew upgrade tinygo

// ABI 版本函数
//
//go:wasmimport env get_abi_version
func getABIVersion() uint32

// 基础环境函数
//
//go:wasmimport env get_caller
func getCaller(addrPtr uint32) uint32

//go:wasmimport env get_contract_address
func getContractAddress(addrPtr uint32) uint32

//go:wasmimport env set_return_data
func setReturnData(dataPtr uint32, dataLen uint32) uint32

//go:wasmimport env emit_event
func emitEvent(eventPtr uint32, eventLen uint32) uint32

//go:wasmimport env log_debug
func logDebug(messagePtr uint32, messageLen uint32) uint32

//go:wasmimport env get_contract_init_params
func getContractInitParams(bufPtr uint32, bufLen uint32) uint32

//go:wasmimport env get_timestamp
func getTimestamp() uint64

//go:wasmimport env get_block_height
func getBlockHeight() uint64

//go:wasmimport env get_block_hash
func getBlockHash(height uint64, hashPtr uint32) uint32

//go:wasmimport env get_merkle_root
func getMerkleRoot(height uint64, rootPtr uint32) uint32

//go:wasmimport env get_state_root
func getStateRoot(height uint64, rootPtr uint32) uint32

//go:wasmimport env get_miner_address
func getMinerAddress(height uint64, addrPtr uint32) uint32

//go:wasmimport env get_tx_hash
func getTxHash(hashPtr uint32) uint32

//go:wasmimport env get_tx_index
func getTxIndex() uint32

// UTXO操作函数
//
//go:wasmimport env create_utxo_output
func createUTXOOutput(recipientPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32

// ⚠️ **已移除**：execute_utxo_transfer
// 原因：违背WES"无业务语义"架构原则
// 该函数包含业务语义（UTXO选择、找零计算），不应在HostABI层实现
// 请使用原语函数：append_asset_output (TxAddAssetOutput)
// 完整的转账逻辑应在SDK的helpers层实现（见 helpers/token/transfer.go）

//go:wasmimport env query_utxo_balance
func queryUTXOBalance(addressPtr uint32, tokenIDPtr uint32, tokenIDLen uint32) uint64

// 状态查询函数（可选）
//
//go:wasmimport env state_get
func stateGet(keyPtr uint32, keyLen uint32, valuePtr uint32, valueLen uint32) uint32

//go:wasmimport env state_get_from_chain
func stateGetFromChain(stateIDPtr uint32, stateIDLen uint32, valuePtr uint32, valueLen uint32, versionPtr uint32) uint32

// ⚠️ **已删除**：state_put 宿主函数声明
// 原因：违背WES架构原则，EUTXO模型无全局状态存储

// ⚠️ **已删除**：state_exists 宿主函数声明
// 原因：违背WES架构原则，EUTXO模型无全局状态存储

// 追加输出/高级UTXO/批量接口
//
//go:wasmimport env append_state_output
func appendStateOutput(stateIDPtr uint32, stateIDLen uint32, stateVersion uint64, execHashPtr uint32, publicInputsPtr uint32, publicInputsLen uint32, parentHashPtr uint32) uint32

//go:wasmimport env append_resource_output
func appendResourceOutput(resourcePtr uint32, resourceLen uint32, ownerPtr uint32, ownerLen uint32, lockingPtr uint32, lockingLen uint32) uint32

//go:wasmimport env create_asset_output_with_lock
func createAssetOutputWithLock(recipientPtr uint32, recipientLen uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32, lockingPtr uint32, lockingLen uint32) uint32

// ⚠️ **已移除**：execute_utxo_transfer_ex
// 原因：违背WES"无业务语义"架构原则
// 该函数包含业务语义（UTXO选择、找零计算），不应在HostABI层实现
// 请使用原语函数：create_asset_output_with_lock + append_tx_input
// 完整的转账逻辑应在SDK的helpers层实现（见 helpers/token/transfer.go）

//go:wasmimport env batch_create_outputs
func batchCreateOutputs(batchPtr uint32, batchLen uint32) uint32

// 内存管理函数
//
//go:wasmimport env malloc
func malloc(size uint32) uint32

// 地址编码转换函数（复用宿主 AddressManager）
//
//go:wasmimport env address_bytes_to_base58
func addressBytesToBase58(addrPtr uint32, resultPtr uint32, maxLen uint32) uint32

//go:wasmimport env address_base58_to_bytes
func addressBase58ToBytes(base58Ptr uint32, base58Len uint32, resultPtr uint32) uint32

// HostABI v1 新增原语
//
//go:wasmimport env get_chain_id
func getChainID(chainIDPtr uint32) uint32

//go:wasmimport env utxo_lookup
func utxoLookup(txIDPtr uint32, txIDLen uint32, index uint32, outputPtr uint32, outputSize uint32) uint32

//go:wasmimport env utxo_lookup_json
func utxoLookupJSON(txIDPtr uint32, txIDLen uint32, index uint32, outputPtr uint32, outputSize uint32) uint32

//go:wasmimport env utxo_exists
func utxoExists(txIDPtr uint32, txIDLen uint32, index uint32) uint32

//go:wasmimport env resource_lookup
func resourceLookup(contentHashPtr uint32, contentHashLen uint32, resourcePtr uint32, resourceSize uint32) uint32

//go:wasmimport env resource_lookup_json
func resourceLookupJSON(contentHashPtr uint32, contentHashLen uint32, resourcePtr uint32, resourceSize uint32) uint32

//go:wasmimport env resource_exists
func resourceExists(contentHashPtr uint32, contentHashLen uint32) uint32

//go:wasmimport env append_tx_input
func appendTxInput(txIDPtr uint32, txIDLen uint32, index uint32, isRefOnly uint32, proofPtr uint32, proofLen uint32) uint32

// ==================== 受控外部交互函数（ISPC创新）====================
//
// 🌟 **ISPC核心创新**：受控外部交互，替代传统预言机
//
// **ISPC 创新点**：
//   传统区块链是封闭系统，无法直接访问外部数据，需要"预言机"将外部数据喂入链上。
//   WES ISPC 通过"受控声明+佐证+验证"机制，让合约可以直接调用外部 API、查询数据库
//   或读取文件，无需传统预言机。这是 ISPC 的核心创新之一。
//
// **ISPC 工作原理**：
//   1. 声明外部状态预期（declareExternalState）：
//      - 告诉系统"我要调用这个外部数据源，预期得到这样的数据"
//      - 系统记录声明，生成 claimID
//   2. 提供验证佐证（provideEvidence）：
//      - 提供 API 数字签名、响应哈希、时间戳证明等密码学佐证
//      - 系统验证佐证的有效性
//   3. 运行时验证并记录到执行轨迹：
//      - ISPC 运行时验证佐证的有效性
//      - 外部调用被记录到执行轨迹
//   4. 查询已验证的外部状态数据（queryControlledState）：
//      - 返回验证后的外部数据
//   5. 生成 ZK 证明：
//      - 执行轨迹自动生成 ZK 证明（包含外部交互验证）
//   6. 验证节点验证证明：
//      - 其他节点验证证明，无需重复调用外部 API
//
// **与传统区块链的对比**：
//   传统区块链：
//     - 需要预言机服务调用外部 API
//     - 预言机将结果喂入链上
//     - 合约使用预言机提供的数据
//     - 问题：预言机是中心化瓶颈，需要支付费用，存在延迟
//
//   WES ISPC：
//     - 直接调用外部 API
//     - 单次调用，多点验证，自动生成 ZK 证明
//     - 无需传统预言机，直接获取外部数据
//     - 实时调用，无延迟
//
// **使用建议**：
//   - ✅ **推荐**：使用 `helpers/external` 模块的业务语义接口
//   - ⚠️ **不推荐**：直接使用这些底层 HostABI 函数（除非有特殊需求）
//
// **进一步了解**：
//   - [ISPC 快速开始指南](../docs/ISPC_QUICK_START.md)
//   - [ISPC vs 传统区块链对比](../docs/ISPC_VS_TRADITIONAL.md)
//   - [ISPC 最佳实践](../docs/ISPC_BEST_PRACTICES.md)
//
// ⚠️ **注意**：这些函数可能还在开发中，如果底层未实现，会返回错误

// host_declare_external_state 声明外部状态预期
//
// 🎯 **用途**：声明要调用的外部数据源和预期结果
//
// **ISPC 机制**：
//
//	这是 ISPC 受控外部交互的第一步。合约声明要调用的外部数据源
//	（API、数据库、文件等）和预期结果，系统记录声明并生成 claimID。
//
// **参数格式（JSON）**:
//
//	{
//	  "claim_type": "api_response|database_query|file_content",  // 声明类型
//	  "source": "API端点/数据库标识/文件标识",                      // 数据源标识
//	  "query_params": {...},                                      // 查询参数
//	  "timestamp": 1640995200,                                    // 时间戳（可选）
//	  "expected_response": {...}                                  // 预期响应（可选）
//	}
//
// **返回**：
//   - claimID: 声明ID（用于后续提供佐证和查询）
//   - error: 错误信息
//
// **示例**：
//
//	claim := &ExternalStateClaim{
//	    ClaimType:   "api_response",
//	    Source:     "https://api.example.com/price",
//	    QueryParams: map[string]interface{}{"symbol": "BTC"},
//	}
//	claimID, err := DeclareExternalState(claim)
//
//go:wasmimport env host_declare_external_state
func hostDeclareExternalState(claimPtr uint32, claimLen uint32, claimIDPtr uint32, claimIDSize uint32) uint32

// host_provide_evidence 提供验证佐证
//
// 🎯 **用途**：提供密码学验证佐证，证明外部数据的可信性
//
// **ISPC 机制**：
//
//	这是 ISPC 受控外部交互的第二步。合约提供密码学验证佐证
//	（API 数字签名、响应哈希、时间戳证明等），系统验证佐证的有效性。
//
// **参数格式（JSON）**:
//
//	{
//	  "claim_id": "...",          // 声明ID（从 declareExternalState 获取）
//	  "api_signature": "...",      // API 数字签名（从外部服务获取）
//	  "response_hash": "...",       // 响应数据哈希（从外部服务获取）
//	  "timestamp_proof": "...",    // 时间戳证明（可选）
//	  "data_integrity": "...",     // 数据完整性证明（可选）
//	  "attestation": "..."         // 其他证明（可选）
//	}
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	evidence := &Evidence{
//	    APISignature: apiSignature,  // API 数字签名（从外部服务获取）
//	    ResponseHash: responseHash,  // 响应数据哈希（从外部服务获取）
//	}
//	err := ProvideEvidence(claimID, evidence)
//
//go:wasmimport env host_provide_evidence
func hostProvideEvidence(claimIDPtr uint32, claimIDLen uint32, evidencePtr uint32, evidenceLen uint32) uint32

// host_query_controlled_state 查询受控外部状态
//
// 🎯 **用途**：查询已验证的外部状态数据
//
// **ISPC 机制**：
//
//	这是 ISPC 受控外部交互的第三步。合约查询已验证的外部状态数据。
//	只有在提供了有效的验证佐证后，才能查询到外部数据。
//
// **参数**：
//   - claimID: 声明ID（从 declareExternalState 获取）
//
// **返回**：
//   - data: 验证后的外部数据（JSON格式）
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	data, err := QueryControlledState(claimID)
//	if err != nil {
//	    return ERROR_EXECUTION_FAILED
//	}
//	// 使用data进行业务逻辑处理
//
//go:wasmimport env host_query_controlled_state
func hostQueryControlledState(claimIDPtr uint32, claimIDLen uint32, resultPtr uint32, resultSize uint32) uint32

// ==================== 封装的宿主函数接口 ====================

// ===== ABI 版本查询 =====

// GetABIVersion 获取引擎支持的 Host ABI 版本
//
// 🎯 **用途**: 合约启动时校验 ABI 版本兼容性
//
// **返回值**:
//   - version: 版本号（(major<<16)|(minor<<8)|patch）
//
// **示例**:
//
//	engineVersion := GetABIVersion()
//	expectedVersion := uint32(0x00010000) // v1.0.0
//	if (engineVersion >> 16) != (expectedVersion >> 16) {
//	    return ERROR_ABI_VERSION_MISMATCH
//	}
func GetABIVersion() uint32 {
	return getABIVersion()
}

// CheckABICompatibility 检查 ABI 版本兼容性
//
// 🎯 **用途**: 简化合约启动时的版本检查
//
// **参数**:
//   - expectedVersion: 合约编译时的 ABI 版本
//
// **返回值**:
//   - error: 兼容返回 nil，不兼容返回错误
func CheckABICompatibility(expectedVersion uint32) error {
	engineVersion := GetABIVersion()

	// 主版本号必须相同
	if (engineVersion >> 16) != (expectedVersion >> 16) {
		engineMajor := (engineVersion >> 16) & 0xFF
		engineMinor := (engineVersion >> 8) & 0xFF
		enginePatch := engineVersion & 0xFF
		expectedMajor := (expectedVersion >> 16) & 0xFF
		expectedMinor := (expectedVersion >> 8) & 0xFF
		expectedPatch := expectedVersion & 0xFF

		// 构造错误消息（不使用 fmt.Sprintf，因为 TinyGo 可能不支持）
		msg := "ABI major version mismatch: engine=" +
			Uint64ToString(uint64(engineMajor)) + "." +
			Uint64ToString(uint64(engineMinor)) + "." +
			Uint64ToString(uint64(enginePatch)) +
			", expected=" +
			Uint64ToString(uint64(expectedMajor)) + "." +
			Uint64ToString(uint64(expectedMinor)) + "." +
			Uint64ToString(uint64(expectedPatch))

		return NewContractError(ERROR_EXECUTION_FAILED, msg)
	}

	// 合约次版本号不能高于引擎
	engineMinor := (engineVersion >> 8) & 0xFF
	expectedMinor := (expectedVersion >> 8) & 0xFF
	if expectedMinor > engineMinor {
		engineMajor := (engineVersion >> 16) & 0xFF
		expectedMajor := (expectedVersion >> 16) & 0xFF

		// 构造错误消息
		msg := "ABI minor version too new: engine=" +
			Uint64ToString(uint64(engineMajor)) + "." +
			Uint64ToString(uint64(engineMinor)) +
			", expected=" +
			Uint64ToString(uint64(expectedMajor)) + "." +
			Uint64ToString(uint64(expectedMinor))

		return NewContractError(ERROR_EXECUTION_FAILED, msg)
	}

	return nil
}

// ===== 环境信息函数 =====

// GetCaller 获取合约调用者地址
//
// 🎯 **修复说明**：
//   - 严格校验宿主返回长度为 20 字节
//   - 防御性错误处理，避免使用损坏的地址数据
func GetCaller() Address {
	addr := malloc(20)
	if addr == 0 {
		return Address{}
	}

	// 🔧 关键修复：接收宿主返回的实际长度
	actualLen := getCaller(addr)

	// 严格校验返回长度必须为 20 字节
	if actualLen != 20 {
		// 返回零地址，避免使用非法数据
		return Address{}
	}

	return AddressFromBytes(GetBytes(addr, 20))
}

// GetContractAddress 获取当前合约地址
//
// 🎯 **修复说明**：
//   - 严格校验宿主返回长度为 20 字节
//   - 防御性错误处理，避免使用损坏的地址数据
func GetContractAddress() Address {
	addr := malloc(20)
	if addr == 0 {
		return Address{}
	}

	// 🔧 关键修复：接收宿主返回的实际长度
	actualLen := getContractAddress(addr)

	// 严格校验返回长度必须为 20 字节
	if actualLen != 20 {
		// 返回零地址，避免使用非法数据
		return Address{}
	}

	return AddressFromBytes(GetBytes(addr, 20))
}

// GetTimestamp 获取当前时间戳
func GetTimestamp() uint64 {
	return getTimestamp()
}

// GetBlockHeight 获取当前区块高度
func GetBlockHeight() uint64 {
	return getBlockHeight()
}

// GetBlockHash 获取指定高度的区块哈希
func GetBlockHash(height uint64) Hash {
	hashPtr := malloc(32)
	if hashPtr == 0 {
		return Hash{}
	}

	result := getBlockHash(height, hashPtr)
	if result != SUCCESS {
		return Hash{}
	}

	return HashFromBytes(GetBytes(hashPtr, 32))
}

// GetMerkleRoot 获取指定高度区块的交易Merkle根
//
// 🎯 **用途**：用于Merkle Proof验证、跨链桥、轻节点验证
func GetMerkleRoot(height uint64) Hash {
	rootPtr := malloc(32)
	if rootPtr == 0 {
		return Hash{}
	}

	result := getMerkleRoot(height, rootPtr)
	if result != SUCCESS {
		return Hash{}
	}

	return HashFromBytes(GetBytes(rootPtr, 32))
}

// GetStateRoot 获取指定高度区块的状态根
//
// 🎯 **用途**：用于状态证明、跨链验证、审计
func GetStateRoot(height uint64) Hash {
	rootPtr := malloc(32)
	if rootPtr == 0 {
		return Hash{}
	}

	result := getStateRoot(height, rootPtr)
	if result != SUCCESS {
		return Hash{}
	}

	return HashFromBytes(GetBytes(rootPtr, 32))
}

// GetMinerAddress 获取指定高度区块的矿工地址
//
// 🎯 **用途**：用于矿工奖励分配、治理权重计算
func GetMinerAddress(height uint64) Address {
	addrPtr := malloc(20)
	if addrPtr == 0 {
		return Address{}
	}

	result := getMinerAddress(height, addrPtr)
	if result != SUCCESS {
		return Address{}
	}

	return AddressFromBytes(GetBytes(addrPtr, 20))
}

// GetTxHash 获取当前执行交易的哈希
//
// 🎯 **用途**：交易唯一标识、幂等性检查、事件关联
func GetTxHash() Hash {
	hashPtr := malloc(32)
	if hashPtr == 0 {
		return Hash{}
	}

	result := getTxHash(hashPtr)
	if result != SUCCESS {
		return Hash{}
	}

	return HashFromBytes(GetBytes(hashPtr, 32))
}

// GetTxIndex 获取当前交易在区块内的索引
//
// 🎯 **用途**：区块内排序、状态快照
func GetTxIndex() uint32 {
	return getTxIndex()
}

// ===== 合约参数和返回值函数 =====

// GetContractParams 获取合约调用参数
func GetContractParams() *ContractParams {
	// 分配足够大的缓冲区
	bufSize := uint32(8192)
	buffer := malloc(bufSize)
	if buffer == 0 {
		return NewContractParams([]byte{})
	}

	actualLen := getContractInitParams(buffer, bufSize)
	if actualLen == 0 {
		return NewContractParams([]byte{})
	}

	data := GetBytes(buffer, actualLen)
	return NewContractParams(data)
}

// SetReturnData 设置合约返回数据
func SetReturnData(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	dataPtr, dataLen := AllocateBytes(data)
	if dataPtr == 0 {
		return NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate return data")
	}

	result := setReturnData(dataPtr, dataLen)
	if result != SUCCESS {
		return NewContractError(result, "failed to set return data")
	}

	return nil
}

// SetReturnString 设置字符串返回数据
func SetReturnString(s string) error {
	return SetReturnData([]byte(s))
}

// SetReturnJSON 设置JSON格式返回数据
func SetReturnJSON(obj interface{}) error {
	jsonStr := serializeToJSON(obj)
	if jsonStr == "" {
		return NewContractError(ERROR_INVALID_PARAMS, "unsupported return type")
	}
	return SetReturnString(jsonStr)
}

// serializeToJSON 递归序列化为 JSON 字符串
//
// 🎯 **修复说明**：
//   - 新增对 Amount (uint64 别名) 的显式支持
//   - 确保所有数值类型都能正确序列化
func serializeToJSON(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return `"` + escapeJSONString(v) + `"`
	case Amount:
		// 🔧 关键修复：显式支持 Amount 类型
		return Uint64ToString(uint64(v))
	case uint64:
		return Uint64ToString(v)
	case int64:
		if v < 0 {
			return "-" + Uint64ToString(uint64(-v))
		}
		return Uint64ToString(uint64(v))
	case int:
		return serializeToJSON(int64(v))
	case uint32:
		return Uint64ToString(uint64(v))
	case int32:
		return serializeToJSON(int64(v))
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	case map[string]interface{}:
		return serializeMapToJSON(v)
	case map[string]string:
		// 特化处理纯字符串 map
		result := make(map[string]interface{}, len(v))
		for k, val := range v {
			result[k] = val
		}
		return serializeMapToJSON(result)
	case map[string]uint64:
		// 特化处理纯数字 map
		result := make(map[string]interface{}, len(v))
		for k, val := range v {
			result[k] = val
		}
		return serializeMapToJSON(result)
	case []interface{}:
		return serializeArrayToJSON(v)
	case []string:
		// 特化处理字符串数组
		arr := make([]interface{}, len(v))
		for i, s := range v {
			arr[i] = s
		}
		return serializeArrayToJSON(arr)
	case []uint64:
		// 特化处理数字数组
		arr := make([]interface{}, len(v))
		for i, n := range v {
			arr[i] = n
		}
		return serializeArrayToJSON(arr)
	default:
		return ""
	}
}

// serializeMapToJSON 序列化 map 为 JSON 对象
func serializeMapToJSON(m map[string]interface{}) string {
	if len(m) == 0 {
		return "{}"
	}

	fields := make([]string, 0, len(m))
	for key, value := range m {
		valueJSON := serializeToJSON(value)
		if valueJSON != "" {
			fields = append(fields, `"`+escapeJSONString(key)+`":`+valueJSON)
		}
	}

	result := "{"
	for i, field := range fields {
		if i > 0 {
			result += ","
		}
		result += field
	}
	result += "}"
	return result
}

// serializeArrayToJSON 序列化数组为 JSON 数组
func serializeArrayToJSON(arr []interface{}) string {
	if len(arr) == 0 {
		return "[]"
	}

	result := "["
	for i, item := range arr {
		if i > 0 {
			result += ","
		}
		result += serializeToJSON(item)
	}
	result += "]"
	return result
}

// escapeJSONString 转义 JSON 字符串中的特殊字符
func escapeJSONString(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '"':
			result += `\"`
		case '\\':
			result += `\\`
		case '\n':
			result += `\n`
		case '\r':
			result += `\r`
		case '\t':
			result += `\t`
		default:
			result += string(c)
		}
	}
	return result
}

// ===== 事件发出函数 =====

// EmitEvent 发出事件
func EmitEvent(event *Event) error {
	if event == nil {
		return NewContractError(ERROR_INVALID_PARAMS, "event cannot be nil")
	}

	eventJSON := event.ToJSON()
	eventPtr, eventLen := AllocateString(eventJSON)
	if eventPtr == 0 {
		return NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate event data")
	}

	result := emitEvent(eventPtr, eventLen)
	if result != SUCCESS {
		return NewContractError(result, "failed to emit event")
	}

	return nil
}

// EmitSimpleEvent 发出简单事件
func EmitSimpleEvent(name string, data map[string]string) error {
	event := NewEvent(name)
	for key, value := range data {
		event.AddStringField(key, value)
	}
	return EmitEvent(event)
}

// ===== UTXO操作函数 =====

// ⚠️ **已删除**：TransferUTXO 和 TransferUTXOWithLock
//
// **原因**：违背WES"无业务语义"架构原则，功能不完整（仅创建输出，不处理UTXO选择和找零）
//
// **替代方案**：
// 1. 使用 helpers/token/Transfer - 包含完整的转账逻辑（推荐）
// 2. 使用 framework.BeginTransaction().Transfer().Finalize() - 链式API，包含完整业务逻辑
// 3. 直接使用原语函数 create_utxo_output 或 create_asset_output_with_lock（仅创建输出）

// QueryBalance 查询UTXO余额
//
// 参数：
//   - address: 要查询的地址
//   - tokenID: 代币ID（空字符串""表示查询原生币）
func QueryBalance(address Address, tokenID TokenID) Amount {
	addressPtr, _ := AllocateBytes(address.ToBytes())
	if addressPtr == 0 {
		return 0
	}

	// tokenID可以为空（查询原生币），所以tokenIDPtr=0是合法的
	var tokenIDPtr, tokenIDLen uint32
	if tokenID != "" {
		tokenIDPtr, tokenIDLen = AllocateString(string(tokenID))
		if tokenIDPtr == 0 {
			// 分配失败
			return 0
		}
	}
	// 如果tokenID为空，tokenIDPtr=0, tokenIDLen=0，宿主函数会理解为查询原生币

	balance := queryUTXOBalance(addressPtr, tokenIDPtr, tokenIDLen)
	return Amount(balance)
}

// ===== 状态查询函数（可选，仅限只读操作）=====

// GetState 获取状态数据（只读）
func GetState(key string) ([]byte, error) {
	keyPtr, keyLen := AllocateString(key)
	if keyPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate key")
	}

	// 分配返回值缓冲区
	maxValueSize := uint32(4096)
	valuePtr := malloc(maxValueSize)
	if valuePtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate value buffer")
	}

	result := stateGet(keyPtr, keyLen, valuePtr, maxValueSize)
	if result != SUCCESS {
		return nil, NewContractError(result, "failed to get state")
	}

	// 简化实现：假设实际长度存储在特定位置
	// 实际实现中需要根据具体的宿主函数规范来处理
	value := GetBytes(valuePtr, maxValueSize)
	return value, nil
}

// GetStateFromChain 从链上查询历史状态
//
// 🎯 **用途**：查询链上已确认交易中的StateOutput，返回匹配stateID的最新状态值和版本号
//
// **参数**：
//   - stateID: 状态ID（字节数组）
//
// **返回**：
//   - value: 状态值（executionResultHash）
//   - version: 状态版本号
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 查询链上已确认的交易，查找包含匹配stateID的StateOutput
//   - 返回版本号最高的状态值
//   - 如果状态不存在，返回错误
//
// **示例**：
//
//	stateID := []byte("balance_user123")
//	value, version, err := framework.GetStateFromChain(stateID)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// 使用 value 和 version
func GetStateFromChain(stateID []byte) ([]byte, uint64, error) {
	// 验证参数
	if len(stateID) == 0 {
		return nil, 0, NewContractError(ERROR_INVALID_PARAMS, "stateID cannot be empty")
	}

	// 分配内存
	stateIDPtr, stateIDLen := AllocateBytes(stateID)
	if stateIDPtr == 0 {
		return nil, 0, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate stateID")
	}

	// 分配返回值缓冲区
	maxValueSize := uint32(4096)
	valuePtr := malloc(maxValueSize)
	if valuePtr == 0 {
		return nil, 0, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate value buffer")
	}

	// 分配版本号缓冲区（8字节uint64）
	versionPtr := malloc(8)
	if versionPtr == 0 {
		return nil, 0, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate version buffer")
	}

	// 调用宿主函数
	result := stateGetFromChain(stateIDPtr, stateIDLen, valuePtr, maxValueSize, versionPtr)
	if result != SUCCESS {
		return nil, 0, NewContractError(result, "failed to get state from chain")
	}

	// 读取状态值
	value := GetBytes(valuePtr, maxValueSize)
	// 移除尾部的零字节
	value = trimTrailingZeros(value)

	// 读取版本号（8字节uint64）
	versionBytes := GetBytes(versionPtr, 8)
	version := uint64(versionBytes[0])<<56 | uint64(versionBytes[1])<<48 | uint64(versionBytes[2])<<40 | uint64(versionBytes[3])<<32 |
		uint64(versionBytes[4])<<24 | uint64(versionBytes[5])<<16 | uint64(versionBytes[6])<<8 | uint64(versionBytes[7])

	return value, version, nil
}

// trimTrailingZeros 移除尾部的零字节
func trimTrailingZeros(data []byte) []byte {
	// 从后往前查找第一个非零字节
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return []byte{} // 全部是零
}

// ⚠️ **已删除**：PutState 和 StateExists
//
// **原因**：
// - PutState: 违背WES架构原则，EUTXO模型无全局状态存储
// - StateExists: 违背WES架构原则，EUTXO模型无全局状态存储
//
// **替代方案**：
// - 使用 AppendStateOutput 在交易草稿中显式记录状态
// - 使用 StateOutput 在交易中传递状态

// ===== 内存管理函数 =====

// Malloc 分配内存
func Malloc(size uint32) uint32 {
	return malloc(size)
}

// ==================== 高级封装函数 ====================

// ===== 合约标准接口辅助 =====

// StandardInitialize 标准合约初始化辅助
func StandardInitialize(contract *ContractBase, customInit func(*ContractParams) error) error {
	params := GetContractParams()

	// 执行自定义初始化逻辑
	if customInit != nil {
		if err := customInit(params); err != nil {
			return err
		}
	}

	// 发出初始化事件
	event := NewEvent("Initialize")
	event.AddStringField("contract_name", contract.Name)
	event.AddStringField("version", contract.Version)
	event.AddAddressField("contract_address", GetContractAddress())
	event.AddUint64Field("timestamp", GetTimestamp())

	return EmitEvent(event)
}

// StandardGetMetadata 标准元数据获取辅助
func StandardGetMetadata(contract *ContractBase) error {
	metadata := contract.BuildMetadataJSON()
	return SetReturnString(metadata)
}

// StandardGetVersion 标准版本获取辅助
func StandardGetVersion(contract *ContractBase) error {
	return SetReturnString(contract.Version)
}

// ===== 代币合约辅助函数 =====

// TokenTransfer 代币转账辅助
//
// ⚠️ **已更新**：使用 helpers/token/Transfer 实现完整转账逻辑
func TokenTransfer(tokenID TokenID, to Address, amount Amount) error {
	caller := GetCaller()

	// 使用 helpers/token/Transfer 实现完整转账逻辑（包含UTXO选择、找零计算）
	// 注意：需要导入 "github.com/weisyn/contract-sdk-go/helpers/token"
	// 这里使用 TransactionBuilder 作为替代实现
	success, _, errCode := BeginTransaction().
		Transfer(caller, to, tokenID, amount).
		Finalize()

	if !success {
		return NewContractError(errCode, "transfer failed")
	}

	// 发出转账事件
	event := NewEvent("Transfer")
	event.AddAddressField("from", caller)
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))

	return EmitEvent(event)
}

// TokenMint 代币铸造辅助
//
// 🎯 **用途**：在合约代码中铸造新代币
//
// **参数**：
//   - tokenID: 代币ID
//   - to: 接收者地址
//   - amount: 铸造数量
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 权限控制和总量控制是业务逻辑，需要在合约代码中实现
//   - ✅ **推荐**：在实际开发中，应使用 `helpers/token.Mint()` 等业务语义接口
//
// **示例**：
//
//	err := framework.TokenMint(framework.TokenID("my_token"), recipient, framework.Amount(1000))
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func TokenMint(tokenID TokenID, to Address, amount Amount) error {
	// 参数验证
	zeroAddr := Address{}
	if to == zeroAddr {
		return NewContractError(ERROR_INVALID_PARAMS, "to address cannot be zero")
	}
	if tokenID == "" {
		return NewContractError(ERROR_INVALID_PARAMS, "tokenID cannot be empty")
	}
	if amount == 0 {
		return NewContractError(ERROR_INVALID_PARAMS, "amount must be greater than 0")
	}

	// 使用framework层的交易构建API进行铸造
	success, _, errCode := BeginTransaction().
		AddAssetOutput(to, tokenID, amount).
		Finalize()

	if !success {
		return NewContractError(errCode, "mint failed")
	}

	// 发出铸造事件
	caller := GetCaller()
	event := NewEvent("Mint")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("minter", caller)
	return EmitEvent(event)
}

// TokenGetBalance 代币余额查询辅助
func TokenGetBalance(address Address, tokenID TokenID) error {
	balance := QueryBalance(address, tokenID)

	result := map[string]interface{}{
		"address":  address.ToString(),
		"token_id": string(tokenID),
		"balance":  uint64(balance),
	}

	return SetReturnJSON(result)
}

// ===== NFT合约辅助函数 =====

// NFTMint NFT铸造辅助
//
// 🎯 **用途**：在合约代码中铸造NFT
//
// **参数**：
//   - tokenID: NFT代币ID
//   - to: 接收者地址
//   - metadata: NFT元数据（可选）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - NFT铸造数量固定为1（NFT是唯一的）
//   - 检查NFT是否已存在（避免重复铸造）
//   - 权限控制是业务逻辑，需要在合约代码中实现
//   - ✅ **推荐**：在实际开发中，应使用 `helpers/token.Mint()` 等业务语义接口
//
// **示例**：
//
//	metadata := map[string]string{
//	    "name": "My NFT",
//	    "description": "A unique NFT",
//	}
//	err := framework.NFTMint(framework.TokenID("nft_001"), recipient, metadata)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func NFTMint(tokenID TokenID, to Address, metadata map[string]string) error {
	// 检查NFT是否已存在
	existingBalance := QueryBalance(to, tokenID)
	if existingBalance > 0 {
		return NewContractError(ERROR_ALREADY_EXISTS, "NFT already exists")
	}

	// 参数验证
	zeroAddr := Address{}
	if to == zeroAddr {
		return NewContractError(ERROR_INVALID_PARAMS, "to address cannot be zero")
	}
	if tokenID == "" {
		return NewContractError(ERROR_INVALID_PARAMS, "tokenID cannot be empty")
	}

	// 使用framework层的交易构建API进行NFT铸造（数量固定为1）
	success, _, errCode := BeginTransaction().
		AddAssetOutput(to, tokenID, Amount(1)).
		Finalize()

	if !success {
		return NewContractError(errCode, "NFT mint failed")
	}

	// 发出NFT铸造事件
	caller := GetCaller()
	event := NewEvent("NFTMint")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", 1)
	event.AddAddressField("minter", caller)

	// 添加元数据字段（如果提供）
	for key, value := range metadata {
		event.AddStringField("metadata_"+key, value)
	}

	return EmitEvent(event)
}

// NFTTransfer NFT转移辅助
//
// ⚠️ **已更新**：使用 TransactionBuilder 实现完整转账逻辑
func NFTTransfer(tokenID TokenID, from, to Address) error {
	// 检查所有权
	balance := QueryBalance(from, tokenID)
	if balance == 0 {
		return NewContractError(ERROR_NOT_FOUND, "NFT not found or not owned")
	}

	// 使用 TransactionBuilder 实现完整转账逻辑（包含UTXO选择、找零计算）
	success, _, errCode := BeginTransaction().
		Transfer(from, to, tokenID, 1).
		Finalize()

	if !success {
		return NewContractError(errCode, "NFT transfer failed")
	}

	// 发出转移事件
	event := NewEvent("NFTTransfer")
	event.AddStringField("token_id", string(tokenID))
	event.AddAddressField("from", from)
	event.AddAddressField("to", to)

	return EmitEvent(event)
}

// ===== 工具函数 =====

// ValidateAddress 验证地址格式
func ValidateAddress(addr Address) error {
	// 简单验证：检查是否为零地址
	zeroAddr := Address{}
	if addr == zeroAddr {
		return NewContractError(ERROR_INVALID_PARAMS, "invalid zero address")
	}
	return nil
}

// ValidateAmount 验证金额
func ValidateAmount(amount Amount) error {
	if amount == 0 {
		return NewContractError(ERROR_INVALID_PARAMS, "invalid zero amount")
	}
	return nil
}

// ValidateTokenID 验证代币ID
func ValidateTokenID(tokenID TokenID) error {
	if len(string(tokenID)) == 0 {
		return NewContractError(ERROR_INVALID_PARAMS, "invalid empty token ID")
	}
	return nil
}
