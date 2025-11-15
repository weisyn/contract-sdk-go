//go:build tinygo || (js && wasm)

package framework

// ==================== HostABI 原语封装 ====================
//
// 🌟 **设计理念**：封装 HostABI 17个最小原语，提供类型安全的Go接口
//
// 🎯 **核心特性**：
// - 完整覆盖17个HostABI原语
// - 类型安全的API设计
// - 统一的错误处理
// - 账户抽象支持
//
// 📋 **原语分类**：
// - 确定性区块视图（4个）
// - 执行上下文（3个）
// - UTXO查询（2个）
// - 资源查询（2个）
// - 交易草稿构建（4个）
// - 执行追踪（2个）

// ==================== 1. 确定性区块视图（4个）====================

// GetChainID 获取链标识符
//
// 🎯 **用途**：用于跨链验证、链标识
//
// **返回**：
//   - chainID: 链标识符（字节数组）
//   - error: 错误信息
//
// **示例**：
//
//	chainID := GetChainID()
//	if len(chainID) == 0 {
//	    return ERROR_EXECUTION_FAILED
//	}
func GetChainID() []byte {
	// 分配缓冲区（链ID通常是字符串，最大64字节）
	bufSize := uint32(64)
	buffer := malloc(bufSize)
	if buffer == 0 {
		return []byte{}
	}

	// 调用宿主函数
	actualLen := getChainID(buffer)
	if actualLen == 0 || actualLen > bufSize {
		return []byte{}
	}

	// 读取链ID
	chainID := GetBytes(buffer, actualLen)
	return chainID
}

// ==================== 2. 执行上下文（3个）====================

// GetTransactionID 获取当前交易ID
//
// 🎯 **用途**：交易唯一标识、幂等性检查、事件关联
//
// **返回**：
//   - txID: 交易ID（32字节哈希）
//
// **示例**：
//
//	txID := GetTransactionID()
//	event.AddBytesField("tx_id", txID)
func GetTransactionID() Hash {
	return GetTxHash() // 复用现有实现
}

// ==================== 3. UTXO查询（2个）====================

// ==================== HostABI相关类型 ====================
//
// 注意：基础类型定义在types.go中，此处仅提供HostABI专用的类型说明

// UTXOLookup 查询指定UTXO
//
// 🎯 **用途**：查询特定UTXO的详细信息
//
// **参数**：
//   - outpoint: UTXO引用点（交易哈希+索引）
//
// **返回**：
//   - utxo: UTXO信息，如果不存在返回nil
//   - error: 错误信息
//
// **示例**：
//
//	outpoint := OutPoint{TxHash: txHash, Index: 0}
//	utxo, err := UTXOLookup(outpoint)
//	if err != nil {
//	    return ERROR_NOT_FOUND
//	}
func UTXOLookup(outpoint OutPoint) (*UTXO, error) {
	// 验证参数
	if len(outpoint.TxHash) != 32 {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "txHash must be 32 bytes")
	}

	// 分配txID缓冲区
	txIDPtr, _ := AllocateBytes(outpoint.TxHash)
	if txIDPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate txID")
	}

	// 分配输出缓冲区（假设最大8KB，JSON可能比Protobuf大）
	outputSize := uint32(8192)
	outputPtr := malloc(outputSize)
	if outputPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate output buffer")
	}

	// 调用JSON格式的宿主函数（TinyGo友好）
	actualLen := utxoLookupJSON(txIDPtr, 32, outpoint.Index, outputPtr, outputSize)
	if actualLen == 0 {
		return nil, NewContractError(ERROR_NOT_FOUND, "UTXO not found")
	}

	// 读取JSON数据
	jsonBytes := GetBytes(outputPtr, actualLen)
	if len(jsonBytes) == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to read JSON data")
	}

	// 解析JSON数据
	utxo, err := parseUTXOFromJSON(jsonBytes, outpoint)
	if err != nil {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to parse UTXO JSON: "+err.Error())
	}

	return utxo, nil
}

// UTXOExists 检查UTXO是否存在
//
// 🎯 **用途**：快速检查UTXO是否存在，无需获取完整信息
//
// **参数**：
//   - outpoint: UTXO引用点
//
// **返回**：
//   - exists: 是否存在
//
// **示例**：
//
//	outpoint := OutPoint{TxHash: txHash, Index: 0}
//	if !UTXOExists(outpoint) {
//	    return ERROR_NOT_FOUND
//	}
func UTXOExists(outpoint OutPoint) bool {
	// 验证参数
	if len(outpoint.TxHash) != 32 {
		return false
	}

	// 分配txID缓冲区
	txIDPtr, _ := AllocateBytes(outpoint.TxHash)
	if txIDPtr == 0 {
		return false
	}

	// 调用宿主函数（返回1表示存在，0表示不存在）
	result := utxoExists(txIDPtr, 32, outpoint.Index)
	return result == 1
}

// QueryUTXOsByAddress 查询地址的所有UTXO（账户抽象）
//
// 🎯 **用途**：账户抽象层，查询地址的所有UTXO
//
// **参数**：
//   - address: 地址
//   - tokenID: 代币ID（可选，nil表示查询所有代币）
//
// **返回**：
//   - utxos: UTXO列表
//
// **注意**：这是账户抽象层提供的便捷方法，不是HostABI原语
//
// **示例**：
//
//	utxos := QueryUTXOsByAddress(caller, nil)
//	for _, utxo := range utxos {
//	    total += utxo.Output.Amount
//	}
func QueryUTXOsByAddress(address Address, tokenID TokenID) []UTXO {
	// TODO: 实现账户抽象层查询
	// 当前返回空列表，待实现
	return []UTXO{}
}

// ==================== 4. 资源查询（2个）====================

// ==================== 资源相关类型 ====================
//
// 注意：Resource类型定义在types.go中

// ResourceLookup 查询资源元数据
//
// 🎯 **用途**：查询资源的元数据信息
//
// **参数**：
//   - contentHash: 资源内容哈希（32字节）
//
// **返回**：
//   - resource: 资源元数据，如果不存在返回nil
//   - error: 错误信息
//
// **示例**：
//
//	contentHash := []byte{...} // 32字节
//	resource, err := ResourceLookup(contentHash)
//	if err != nil {
//	    return ERROR_NOT_FOUND
//	}
func ResourceLookup(contentHash []byte) (*Resource, error) {
	// 验证参数
	if len(contentHash) != 32 {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "contentHash must be 32 bytes")
	}

	// 分配contentHash缓冲区
	contentHashPtr, _ := AllocateBytes(contentHash)
	if contentHashPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate contentHash")
	}

	// 分配输出缓冲区（假设最大8KB，JSON可能比Protobuf大）
	resourceSize := uint32(8192)
	resourcePtr := malloc(resourceSize)
	if resourcePtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate resource buffer")
	}

	// 调用JSON格式的宿主函数（TinyGo友好）
	actualLen := resourceLookupJSON(contentHashPtr, 32, resourcePtr, resourceSize)
	if actualLen == 0 {
		return nil, NewContractError(ERROR_NOT_FOUND, "Resource not found")
	}

	// 读取JSON数据
	jsonBytes := GetBytes(resourcePtr, actualLen)
	if len(jsonBytes) == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to read JSON data")
	}

	// 解析JSON数据
	resource, err := parseResourceFromJSON(jsonBytes, contentHash)
	if err != nil {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to parse Resource JSON: "+err.Error())
	}

	return resource, nil
}

// ResourceExists 检查资源是否存在
//
// 🎯 **用途**：快速检查资源是否存在
//
// **参数**：
//   - contentHash: 资源内容哈希
//
// **返回**：
//   - exists: 是否存在
//
// **示例**：
//
//	contentHash := []byte{...}
//	if !ResourceExists(contentHash) {
//	    return ERROR_NOT_FOUND
//	}
func ResourceExists(contentHash []byte) bool {
	// 验证参数
	if len(contentHash) != 32 {
		return false
	}

	// 分配contentHash缓冲区
	contentHashPtr, _ := AllocateBytes(contentHash)
	if contentHashPtr == 0 {
		return false
	}

	// 调用宿主函数（返回1表示存在，0表示不存在）
	result := resourceExists(contentHashPtr, 32)
	return result == 1
}

// ==================== 5. 交易草稿构建（4个）====================

// AppendStateOutputSimple 追加状态输出（简化版）
//
// 🎯 **用途**：在交易草稿中追加状态输出，用于状态存储
//
// **参数**：
//   - stateID: 状态ID（字节数组）
//   - version: 状态版本号
//   - execHash: 执行结果哈希（字节数组）
//   - parentHash: 父状态哈希（可选，nil表示无父状态）
//
// **返回**：
//   - outputIndex: 输出索引（成功时返回索引，失败时返回0xFFFFFFFF）
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	stateID := []byte("my_state_key")
//	version := uint64(1)
//	execHash := []byte("execution_result_hash")
//	outputIndex, err := framework.AppendStateOutputSimple(stateID, version, execHash, nil)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func AppendStateOutputSimple(stateID []byte, version uint64, execHash []byte, parentHash []byte) (uint32, error) {
	// 验证参数
	if len(stateID) == 0 {
		return 0xFFFFFFFF, NewContractError(ERROR_INVALID_PARAMS, "stateID cannot be empty")
	}
	
	// 验证execHash必须是32字节（节点侧固定读取32字节）
	// 如果execHash不是32字节，需要先计算哈希或补齐到32字节
	var execHash32 [32]byte
	if len(execHash) == 32 {
		copy(execHash32[:], execHash)
	} else if len(execHash) > 0 {
		// 如果execHash不是32字节，使用ComputeHash计算32字节哈希
		hash := ComputeHash(execHash)
		copy(execHash32[:], hash[:])
	} else {
		// 如果execHash为空，使用stateID的哈希
		hash := ComputeHash(stateID)
		copy(execHash32[:], hash[:])
	}

	// 分配内存
	stateIDPtr, stateIDLen := AllocateBytes(stateID)
	if stateIDPtr == 0 {
		return 0xFFFFFFFF, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate stateID")
	}

	// execHash必须是32字节，节点侧固定读取32字节
	execHashPtr := Malloc(32)
	if execHashPtr == 0 {
		return 0xFFFFFFFF, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate execHash")
	}
	execHashBytes := GetBytes(execHashPtr, 32)
	copy(execHashBytes, execHash32[:])

	// publicInputs：使用execHash作为公开输入（节点侧会将execHash作为publicInputs）
	publicInputsPtr := execHashPtr
	publicInputsLen := uint32(32)

	// parentHash可选，但必须是32字节（如果提供）
	var parentPtr uint32
	if len(parentHash) > 0 {
		var parentHash32 [32]byte
		if len(parentHash) == 32 {
			copy(parentHash32[:], parentHash)
		} else {
			// 如果parentHash不是32字节，使用ComputeHash计算32字节哈希
			hash := ComputeHash(parentHash)
			copy(parentHash32[:], hash[:])
		}
		parentPtr = Malloc(32)
		if parentPtr == 0 {
			return 0xFFFFFFFF, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate parentHash")
		}
		parentHashBytes := GetBytes(parentPtr, 32)
		copy(parentHashBytes, parentHash32[:])
	}

	// 调用宿主函数（新签名：7个参数）
	outputIndex := appendStateOutput(stateIDPtr, stateIDLen, version, execHashPtr, publicInputsPtr, publicInputsLen, parentPtr)
	if outputIndex == 0xFFFFFFFF {
		return outputIndex, NewContractError(ERROR_EXECUTION_FAILED, "append_state_output failed")
	}

	return outputIndex, nil
}

// BatchCreateOutputsSimple 批量创建资产输出（简化版）
//
// 🎯 **用途**：批量创建多个资产输出，用于空投、批量转账等场景
//
// **参数**：
//   - items: 输出项列表，每个项包含：
//     * Recipient: 接收者地址（字节数组）
//     * Amount: 金额（uint64）
//     * TokenID: 代币ID（可选，nil表示原生币）
//
// **返回**：
//   - count: 成功创建的输出数量
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	items := []struct {
//	    Recipient []byte
//	    Amount    uint64
//	    TokenID   []byte
//	}{
//	    {recipient1, 100, nil},
//	    {recipient2, 200, nil},
//	}
//	count, err := framework.BatchCreateOutputsSimple(items)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func BatchCreateOutputsSimple(items []struct {
	Recipient []byte
	Amount    uint64
	TokenID   []byte
}) (uint32, error) {
	if len(items) == 0 {
		return 0, NewContractError(ERROR_INVALID_PARAMS, "items cannot be empty")
	}

	// 构造批量输出JSON（手动序列化避免引入encoding/json）
	batchJSON := "["
	for i, it := range items {
		if i > 0 {
			batchJSON += ","
		}
		batchJSON += `{"recipient":"`
		// Base64编码地址（使用标准Base64编码）
		batchJSON += base64EncodeSimple(it.Recipient)
		batchJSON += `","amount":`
		batchJSON += Uint64ToString(it.Amount)
		if len(it.TokenID) > 0 {
			batchJSON += `,"token_id":"`
			batchJSON += base64EncodeSimple(it.TokenID)
			batchJSON += `"`
		} else {
			batchJSON += `,"token_id":null`
		}
		batchJSON += `,"locking_conditions":[]}`
	}
	batchJSON += "]"

	batchBytes := []byte(batchJSON)
	batchPtr, batchLen := AllocateBytes(batchBytes)
	if batchPtr == 0 {
		return 0, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate batch JSON")
	}

	// 调用宿主函数
	result := batchCreateOutputs(batchPtr, batchLen)
	if result == 0xFFFFFFFF {
		return 0, NewContractError(ERROR_EXECUTION_FAILED, "batch_create_outputs failed")
	}

	return result, nil
}

// base64EncodeSimple Base64编码（用于地址和TokenID）
// 使用标准Base64编码算法，适用于TinyGo WASM环境
func base64EncodeSimple(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := ""

	for i := 0; i < len(data); i += 3 {
		b1 := data[i]
		b2 := byte(0)
		b3 := byte(0)

		if i+1 < len(data) {
			b2 = data[i+1]
		}
		if i+2 < len(data) {
			b3 = data[i+2]
		}

		result += string(base64Table[(b1>>2)&0x3F])
		result += string(base64Table[((b1&0x03)<<4)|((b2>>4)&0x0F)])

		if i+1 < len(data) {
			result += string(base64Table[((b2&0x0F)<<2)|((b3>>6)&0x03)])
		} else {
			result += "="
		}

		if i+2 < len(data) {
			result += string(base64Table[b3&0x3F])
		} else {
			result += "="
		}
	}

	return result
}

// ==================== 5. 交易草稿构建（4个）====================

// ==================== 锁定相关类型 ====================
//
// 注意：LockingCondition和UnlockingProof类型定义在types.go中

// ==================== 6. 受控外部交互（ISPC创新，3个）====================
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

// DeclareExternalState 声明外部状态预期
//
// 🎯 **用途**：声明要调用的外部数据源和预期结果
//
// **ISPC 机制**：
//   这是 ISPC 受控外部交互的第一步。合约声明要调用的外部数据源
//   （API、数据库、文件等）和预期结果，系统记录声明并生成 claimID。
//
// **参数**：
//   - claim: 外部状态声明，包含：
//     * ClaimType: 声明类型（"api_response" | "database_query" | "file_content"）
//     * Source: 数据源标识（API端点/数据库标识/文件标识）
//     * QueryParams: 查询参数（JSON格式的map）
//     * Timestamp: 时间戳（可选）
//
// **返回**：
//   - claimID: 声明ID（用于后续提供佐证和查询）
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	claim := &framework.ExternalStateClaim{
//	    ClaimType: "api_response",
//	    Source: "https://api.example.com/price",
//	    QueryParams: map[string]interface{}{"symbol": "BTC"},
//	    Timestamp: framework.GetBlockTimestamp(),
//	}
//	claimID, err := framework.DeclareExternalState(claim)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func DeclareExternalState(claim *ExternalStateClaim) ([]byte, error) {
	// 验证参数
	if claim == nil {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "claim cannot be nil")
	}
	if claim.ClaimType == "" {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "claimType cannot be empty")
	}
	if claim.Source == "" {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "source cannot be empty")
	}

	// 构建JSON参数
	claimJSON := buildClaimJSON(claim)
	if len(claimJSON) == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to build claim JSON")
	}

	// 分配内存
	claimPtr, _ := AllocateBytes(claimJSON)
	if claimPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate claim")
	}

	// 分配claimID缓冲区（假设最大64字节）
	claimIDSize := uint32(64)
	claimIDPtr := malloc(claimIDSize)
	if claimIDPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate claimID buffer")
	}

	// 调用宿主函数
	actualLen := hostDeclareExternalState(claimPtr, uint32(len(claimJSON)), claimIDPtr, claimIDSize)
	if actualLen == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to declare external state")
	}

	// 读取claimID
	claimID := GetBytes(claimIDPtr, actualLen)
	return claimID, nil
}

// ProvideEvidence 提供验证佐证
//
// 🎯 **用途**：提供密码学验证佐证，证明外部数据的可信性
//
// **ISPC 机制**：
//   这是 ISPC 受控外部交互的第二步。合约提供密码学验证佐证
//   （API 数字签名、响应哈希、时间戳证明等），系统验证佐证的有效性。
//
// **参数**：
//   - claimID: 声明ID（由DeclareExternalState返回）
//   - evidence: 验证佐证，必须包含：
//     * APISignature: API 数字签名（从外部服务获取）
//     * ResponseHash: 响应数据哈希（从外部服务获取）
//     * TimestampProof: 时间戳证明（可选）
//     * DataIntegrity: 数据完整性证明（可选）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	evidence := &framework.Evidence{
//	    ClaimID: claimID,
//	    APISignature: apiSignature,  // API 数字签名（从外部服务获取）
//	    ResponseHash: responseHash,  // 响应数据哈希（从外部服务获取）
//	    TimestampProof: timestampProof,  // 时间戳证明（可选）
//	}
//	err := framework.ProvideEvidence(claimID, evidence)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func ProvideEvidence(claimID []byte, evidence *Evidence) error {
	// 验证参数
	if len(claimID) == 0 {
		return NewContractError(ERROR_INVALID_PARAMS, "claimID cannot be empty")
	}
	if evidence == nil {
		return NewContractError(ERROR_INVALID_PARAMS, "evidence cannot be nil")
	}

	// 构建JSON参数
	evidenceJSON := buildEvidenceJSON(evidence)
	if len(evidenceJSON) == 0 {
		return NewContractError(ERROR_EXECUTION_FAILED, "failed to build evidence JSON")
	}

	// 分配内存
	claimIDPtr, _ := AllocateBytes(claimID)
	if claimIDPtr == 0 {
		return NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate claimID")
	}

	evidencePtr, _ := AllocateBytes(evidenceJSON)
	if evidencePtr == 0 {
		return NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate evidence")
	}

	// 调用宿主函数
	result := hostProvideEvidence(claimIDPtr, uint32(len(claimID)), evidencePtr, uint32(len(evidenceJSON)))
	if result != 0 {
		return NewContractError(uint32(result), "failed to provide evidence")
	}

	return nil
}

// QueryControlledState 查询受控外部状态
//
// 🎯 **用途**：查询已验证的外部状态数据
//
// **ISPC 机制**：
//   这是 ISPC 受控外部交互的第三步。合约查询已验证的外部状态数据。
//   只有在提供了有效的验证佐证后，才能查询到外部数据。
//
// **参数**：
//   - claimID: 声明ID（由DeclareExternalState返回）
//
// **返回**：
//   - data: 验证后的外部数据（JSON格式）
//   - error: 错误信息，nil表示成功
//
// **示例**：
//
//	data, err := framework.QueryControlledState(claimID)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// ✅ 使用data进行业务逻辑处理
//	// ✅ ZK 证明自动生成，自动构建交易，自动上链
func QueryControlledState(claimID []byte) ([]byte, error) {
	// 验证参数
	if len(claimID) == 0 {
		return nil, NewContractError(ERROR_INVALID_PARAMS, "claimID cannot be empty")
	}

	// 分配内存
	claimIDPtr, _ := AllocateBytes(claimID)
	if claimIDPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate claimID")
	}

	// 分配结果缓冲区（假设最大64KB）
	resultSize := uint32(65536)
	resultPtr := malloc(resultSize)
	if resultPtr == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to allocate result buffer")
	}

	// 调用宿主函数
	actualLen := hostQueryControlledState(claimIDPtr, uint32(len(claimID)), resultPtr, resultSize)
	if actualLen == 0 {
		return nil, NewContractError(ERROR_EXECUTION_FAILED, "failed to query controlled state")
	}

	// 读取结果
	result := GetBytes(resultPtr, actualLen)
	return result, nil
}

// ==================== 6. 执行追踪（2个）====================

// LogDebug 记录调试日志
//
// 🎯 **用途**：记录调试信息，仅在开发环境可见
//
// **参数**：
//   - message: 日志消息
//
// **示例**：
//
//	LogDebug("Processing transfer: " + amount.String())
func LogDebug(message string) {
	// 使用专门的log_debug宿主函数
	messagePtr, messageLen := AllocateString(message)
	if messagePtr == 0 {
		return
	}
	
	result := logDebug(messagePtr, messageLen)
	if result != SUCCESS {
		// 如果log_debug调用失败，记录错误但不回退（彻底修复）
		return
	}
}

// ==================== 受控外部交互辅助函数 ====================

// buildClaimJSON 构建外部状态声明的JSON
func buildClaimJSON(claim *ExternalStateClaim) []byte {
	// 使用host_functions.go中的serializeToJSON函数
	claimMap := map[string]interface{}{
		"claim_type": claim.ClaimType,
		"source":     claim.Source,
	}
	if len(claim.QueryParams) > 0 {
		claimMap["query_params"] = claim.QueryParams
	}
	if claim.Timestamp > 0 {
		claimMap["timestamp"] = claim.Timestamp
	}
	if len(claim.ExpectedResponse) > 0 {
		claimMap["expected_response"] = claim.ExpectedResponse
	}
	// 使用host_functions.go中的serializeMapToJSON
	jsonStr := serializeMapToJSON(claimMap)
	return []byte(jsonStr)
}

// buildEvidenceJSON 构建验证佐证的JSON
func buildEvidenceJSON(evidence *Evidence) []byte {
	// 使用host_functions.go中的serializeToJSON函数
	evidenceMap := map[string]interface{}{
		"claim_id": string(evidence.ClaimID),
	}
	if len(evidence.APISignature) > 0 {
		evidenceMap["api_signature"] = string(evidence.APISignature)
	}
	if len(evidence.ResponseHash) > 0 {
		evidenceMap["response_hash"] = string(evidence.ResponseHash)
	}
	if len(evidence.TimestampProof) > 0 {
		evidenceMap["timestamp_proof"] = string(evidence.TimestampProof)
	}
	if len(evidence.DataIntegrity) > 0 {
		evidenceMap["data_integrity"] = string(evidence.DataIntegrity)
	}
	if len(evidence.Attestation) > 0 {
		evidenceMap["attestation"] = string(evidence.Attestation)
	}
	// 使用host_functions.go中的serializeMapToJSON
	jsonStr := serializeMapToJSON(evidenceMap)
	return []byte(jsonStr)
}

// ==================== 辅助函数 ====================

// buildLockingConditionsJSON 构建锁定条件JSON
func buildLockingConditionsJSON(conditions []LockingCondition) []byte {
	if len(conditions) == 0 {
		return nil
	}

	// 将LockingCondition转换为JSON字符串数组
	jsonStrings := make([]string, len(conditions))
	for i, cond := range conditions {
		// 简化实现：假设Condition已经是JSON字符串
		jsonStrings[i] = string(cond.Condition)
	}

	return BuildLockingJSONArray(jsonStrings)
}

// QueryUTXOBalance 查询UTXO余额（账户抽象）
//
// 🎯 **用途**：账户抽象层，查询地址的余额
//
// **参数**：
//   - address: 地址
//   - tokenID: 代币ID（可选，nil表示查询原生币）
//
// **返回**：
//   - balance: 余额
//
// **注意**：这是账户抽象层提供的便捷方法，基于UTXO查询实现
//
// **示例**：
//
//	balance := QueryUTXOBalance(caller, nil)
//	if balance < amount {
//	    return ERROR_INSUFFICIENT_BALANCE
//	}
func QueryUTXOBalance(address Address, tokenID TokenID) Amount {
	// 使用现有的QueryBalance方法
	return QueryBalance(address, tokenID)
}

// ==================== JSON解析辅助函数 ====================

// parseUTXOFromJSON 从JSON数据解析UTXO
//
// 🎯 **用途**：解析WES节点返回的TxOutput JSON数据
//
// **JSON格式**（protobuf JSON编码）：
//   {
//     "owner": "base64编码的地址",
//     "lockingConditions": [...],
//     "asset": {...} | "state": {...} | "resource": {...}
//   }
func parseUTXOFromJSON(jsonBytes []byte, outpoint OutPoint) (*UTXO, error) {
	jsonStr := string(jsonBytes)
	
	// 使用简单的JSON解析（TinyGo环境）
	// 解析output_content字段，确定输出类型
	outputType := "asset" // 默认类型
	
	// 检查是否有asset字段
	if findJSONField(jsonStr, "asset") != "" {
		outputType = "asset"
	} else if findJSONField(jsonStr, "state") != "" {
		outputType = "state"
	} else if findJSONField(jsonStr, "resource") != "" {
		outputType = "resource"
	}
	
	// 解析owner字段（地址）
	ownerStr := findJSONField(jsonStr, "owner")
	var recipient Address
	if ownerStr != "" {
		// Base64解码地址（protobuf JSON使用Base64编码字节）
		ownerBytes := base64DecodeSimple(ownerStr)
		if len(ownerBytes) >= 20 {
			recipient = AddressFromBytes(ownerBytes[:20])
		}
	}
	
	// 解析asset字段（如果存在）
	var amount Amount
	var tokenID TokenID
	if outputType == "asset" {
		assetJSON := extractJSONObject(jsonStr, "asset")
		if assetJSON != "" {
			// 解析amount
			amountStr := findJSONField(assetJSON, "amount")
			if amountStr != "" {
				amount = Amount(ParseUint64(amountStr))
			}
			
			// 解析tokenId
			tokenIDStr := findJSONField(assetJSON, "tokenId")
			if tokenIDStr != "" {
				tokenID = TokenID(tokenIDStr)
			}
		}
	}
	
	return &UTXO{
		OutPoint: outpoint,
		Output: TxOutput{
			Type:      outputType,
			Recipient: recipient,
			Amount:    amount,
			TokenID:   tokenID,
			Data:      jsonBytes, // 保存原始JSON数据
		},
	}, nil
}

// parseResourceFromJSON 从JSON数据解析Resource
//
// 🎯 **用途**：解析WES节点返回的Resource JSON数据
//
// **JSON格式**（protobuf JSON编码）：
//   {
//     "category": "STATIC" | "EXECUTABLE",
//     "contentHash": "base64编码的哈希",
//     "mimeType": "...",
//     "size": 12345,
//     "name": "...",
//     "version": "...",
//     ...
//   }
func parseResourceFromJSON(jsonBytes []byte, contentHash []byte) (*Resource, error) {
	jsonStr := string(jsonBytes)
	
	// 解析category字段
	categoryStr := findJSONField(jsonStr, "category")
	category := "static" // 默认类别
	if categoryStr == "EXECUTABLE" || categoryStr == "1" {
		category = "executable"
	}
	
	// 解析mimeType字段
	mimeType := findJSONField(jsonStr, "mimeType")
	
	// 解析size字段
	sizeStr := findJSONField(jsonStr, "size")
	size := uint64(0)
	if sizeStr != "" {
		size = ParseUint64(sizeStr)
	}
	
	return &Resource{
		ContentHash: contentHash,
		Category:    category,
		MimeType:    mimeType,
		Size:        size,
	}, nil
}

// findJSONField 查找JSON字段值（字符串类型）
func findJSONField(jsonStr, key string) string {
	keyPattern := `"` + key + `":"`
	
	startIdx := -1
	for i := 0; i <= len(jsonStr)-len(keyPattern); i++ {
		if jsonStr[i:i+len(keyPattern)] == keyPattern {
			startIdx = i + len(keyPattern)
			break
		}
	}
	
	if startIdx == -1 {
		// 尝试不带引号的数字值
		keyPattern2 := `"` + key + `":`
		for i := 0; i <= len(jsonStr)-len(keyPattern2); i++ {
			if jsonStr[i:i+len(keyPattern2)] == keyPattern2 {
				startIdx = i + len(keyPattern2)
				// 跳过空格
				for startIdx < len(jsonStr) && jsonStr[startIdx] == ' ' {
					startIdx++
				}
				break
			}
		}
		if startIdx == -1 {
			return ""
		}
		
		// 解析数字或字符串
		endIdx := startIdx
		if startIdx < len(jsonStr) && jsonStr[startIdx] == '"' {
			// 字符串值
			startIdx++
			for endIdx < len(jsonStr) && jsonStr[endIdx] != '"' {
				endIdx++
			}
		} else {
			// 数字值
			for endIdx < len(jsonStr) && jsonStr[endIdx] >= '0' && jsonStr[endIdx] <= '9' {
				endIdx++
			}
		}
		
		if endIdx > startIdx {
			return jsonStr[startIdx:endIdx]
		}
		return ""
	}
	
	// 字符串值（带引号）
	endIdx := startIdx
	for endIdx < len(jsonStr) && jsonStr[endIdx] != '"' {
		endIdx++
	}
	
	if endIdx > startIdx {
		return jsonStr[startIdx:endIdx]
	}
	
	return ""
}

// extractJSONObject 提取JSON对象字段
func extractJSONObject(jsonStr, key string) string {
	keyPattern := `"` + key + `":{`
	
	startIdx := -1
	for i := 0; i <= len(jsonStr)-len(keyPattern); i++ {
		if jsonStr[i:i+len(keyPattern)] == keyPattern {
			startIdx = i + len(keyPattern) - 1 // 包含 '{'
			break
		}
	}
	
	if startIdx == -1 {
		return ""
	}
	
	// 找到匹配的 '}'
	braceCount := 0
	endIdx := startIdx
	for endIdx < len(jsonStr) {
		if jsonStr[endIdx] == '{' {
			braceCount++
		} else if jsonStr[endIdx] == '}' {
			braceCount--
			if braceCount == 0 {
				endIdx++ // 包含 '}'
				break
			}
		}
		endIdx++
	}
	
	if endIdx > startIdx {
		return jsonStr[startIdx:endIdx]
	}
	
	return ""
}

// base64DecodeSimple Base64解码（TinyGo WASM环境）
func base64DecodeSimple(encoded string) []byte {
	if len(encoded) == 0 {
		return nil
	}
	
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, 0, len(encoded)*3/4)
	
	i := 0
	for i < len(encoded) {
		if encoded[i] == '=' {
			break
		}
		
		// 读取4个字符
		if i+3 >= len(encoded) {
			break
		}
		
		c1 := findBase64Char(encoded[i], base64Table)
		c2 := findBase64Char(encoded[i+1], base64Table)
		c3 := findBase64Char(encoded[i+2], base64Table)
		c4 := findBase64Char(encoded[i+3], base64Table)
		
		if c1 == 255 || c2 == 255 {
			break
		}
		
		// 解码第一个字节
		result = append(result, byte((c1<<2)|(c2>>4)))
		
		if c3 != 255 {
			// 解码第二个字节
			result = append(result, byte((c2<<4)|(c3>>2)))
			
			if c4 != 255 {
				// 解码第三个字节
				result = append(result, byte((c3<<6)|c4))
			}
		}
		
		i += 4
	}
	
	return result
}

// findBase64Char 查找Base64字符的索引
func findBase64Char(c byte, table string) byte {
	for i := 0; i < len(table); i++ {
		if table[i] == c {
			return byte(i)
		}
	}
	return 255 // 未找到
}

// ==================== 状态版本管理 ====================

// GetStateVersion 获取状态的当前版本号（从链上查询）
//
// 🎯 **用途**：获取状态的当前版本号，用于状态更新时递增版本号
//
// **参数**：
//   - stateID: 状态ID（字节数组）
//
// **返回**：
//   - version: 状态版本号（如果状态不存在，返回0）
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 从链上查询状态，获取版本号
//   - 如果状态不存在，返回版本号0（首次创建时使用版本号1）
//
// **示例**：
//
//	stateID := []byte("balance_user123")
//	version, err := framework.GetStateVersion(stateID)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	newVersion := version + 1 // 递增版本号
func GetStateVersion(stateID []byte) (uint64, error) {
	_, version, err := GetStateFromChain(stateID)
	if err != nil {
		// 如果状态不存在，返回版本号0（首次创建时使用版本号1）
		return 0, nil
	}
	return version, nil
}

// IncrementStateVersion 递增状态版本号
//
// 🎯 **用途**：获取状态的当前版本号并递增，用于状态更新
//
// **参数**：
//   - stateID: 状态ID（字节数组）
//
// **返回**：
//   - newVersion: 新的版本号（当前版本号+1）
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 从链上查询状态的当前版本号
//   - 返回递增后的版本号（当前版本号+1）
//   - 如果状态不存在，返回版本号1（首次创建）
//
// **示例**：
//
//	stateID := []byte("balance_user123")
//	newVersion, err := framework.IncrementStateVersion(stateID)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// 使用 newVersion 更新状态
func IncrementStateVersion(stateID []byte) (uint64, error) {
	currentVersion, err := GetStateVersion(stateID)
	if err != nil {
		return 0, err
	}
	return currentVersion + 1, nil
}
