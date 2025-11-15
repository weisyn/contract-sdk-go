//go:build tinygo || (js && wasm)

package external

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ValidateAndQuery 验证并查询外部状态
//
// 🎯 **用途**：通过ISPC受控机制验证并查询外部状态，替代传统预言机
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
// **参数**：
//   - claimType: 声明类型，可选值：
//     * "api_response": API 响应
//     * "database_query": 数据库查询
//     * "file_content": 文件内容
//   - source: 数据源标识：
//     * API 响应：API 端点 URL（如 "https://api.example.com/price"）
//     * 数据库查询：数据库标识（如 "db:main"）
//     * 文件内容：文件标识（如 "file:contract.pdf"）
//   - params: 查询参数（JSON 格式的 map）
//   - evidence: 验证佐证，必须包含：
//     * APISignature: API 数字签名（从外部服务获取）
//     * ResponseHash: 响应数据哈希（从外部服务获取）
//     * Timestamp: 时间戳（可选）
//     * Nonce: 随机数（可选）
//
// **返回**：
//   - data: 验证后的外部数据（JSON格式）
//   - error: 错误信息，nil表示成功
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
// **示例**：
//
//	// 查询API响应
//	data, err := external.ValidateAndQuery(
//	    "api_response",
//	    "https://api.example.com/price",
//	    map[string]interface{}{"symbol": "BTC"},
//	    &framework.Evidence{
//	        APISignature: apiSignature,  // API 数字签名（从外部服务获取）
//	        ResponseHash: responseHash,  // 响应数据哈希（从外部服务获取）
//	    },
//	)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// ✅ 使用data进行业务逻辑处理
//	// ✅ ZK 证明自动生成，自动构建交易，自动上链
func ValidateAndQuery(
	claimType string,
	source string,
	params map[string]interface{},
	evidence *framework.Evidence,
) ([]byte, error) {
	// 1. 验证参数
	if claimType == "" {
		return nil, framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"claimType cannot be empty",
		)
	}
	if source == "" {
		return nil, framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"source cannot be empty",
		)
	}
	if evidence == nil {
		return nil, framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"evidence cannot be nil",
		)
	}

	// 2. 声明外部状态预期
	claim := &framework.ExternalStateClaim{
		ClaimType:   claimType,
		Source:     source,
		QueryParams: params,
		Timestamp:  framework.GetTimestamp(),
	}

	claimID, err := framework.DeclareExternalState(claim)
	if err != nil {
		return nil, err
	}

	// 3. 设置claimID到evidence
	evidence.ClaimID = claimID

	// 4. 提供验证佐证
	err = framework.ProvideEvidence(claimID, evidence)
	if err != nil {
		return nil, err
	}

	// 5. 查询已验证的外部状态
	data, err := framework.QueryControlledState(claimID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CallAPI 受控外部API调用
//
// 🎯 **用途**：通过ISPC受控机制调用外部API，替代传统预言机
//
// **ISPC 创新点**：
//   这是 ISPC 受控外部交互的便捷封装函数，专门用于调用外部 API。
//   内部调用 ValidateAndQuery，使用 "api_response" 作为 claimType。
//
// **ISPC 工作原理**：
//   1. 声明API响应预期：告诉系统"我要调用这个 API，预期得到这样的数据"
//   2. 提供API签名和响应哈希作为佐证：提供密码学验证佐证
//   3. ISPC运行时验证佐证：验证佐证的有效性
//   4. 返回验证后的API响应数据：返回验证后的外部数据
//   5. 生成 ZK 证明：执行轨迹自动生成 ZK 证明（包含外部交互验证）
//   6. 验证节点验证证明：其他节点验证证明，无需重复调用外部 API
//
// **参数**：
//   - apiURL: API端点URL（如 "https://api.example.com/price"）
//   - method: HTTP方法，可选值：
//     * "GET": GET 请求
//     * "POST": POST 请求
//     * "PUT": PUT 请求
//     * "DELETE": DELETE 请求
//   - params: 请求参数（JSON 格式的 map）
//   - apiSignature: API数字签名（从外部服务获取）
//   - responseHash: 响应数据哈希（从外部服务获取）
//
// **返回**：
//   - data: API响应数据（JSON格式）
//   - error: 错误信息，nil表示成功
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
// **示例**：
//
//	data, err := external.CallAPI(
//	    "https://api.example.com/price",
//	    "GET",
//	    map[string]interface{}{"symbol": "BTC"},
//	    apiSignature,  // API 数字签名（从外部服务获取）
//	    responseHash,  // 响应数据哈希（从外部服务获取）
//	)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// ✅ 使用data进行业务逻辑处理
//	// ✅ ZK 证明自动生成，自动构建交易，自动上链
func CallAPI(
	apiURL string,
	method string,
	params map[string]interface{},
	apiSignature []byte,
	responseHash []byte,
) ([]byte, error) {
	// 构建查询参数（包含HTTP方法）
	queryParams := map[string]interface{}{
		"method": method,
	}
	if params != nil {
		for k, v := range params {
			queryParams[k] = v
		}
	}

	// 构建验证佐证
	evidence := &framework.Evidence{
		APISignature: apiSignature,
		ResponseHash: responseHash,
	}

	// 调用ValidateAndQuery
	return ValidateAndQuery("api_response", apiURL, queryParams, evidence)
}

// QueryDatabase 受控数据库查询
//
// 🎯 **用途**：通过ISPC受控机制查询数据库，替代传统预言机
//
// **参数**：
//   - dbIdentifier: 数据库标识
//   - query: 查询语句
//   - params: 查询参数
//   - stateHash: 数据库状态哈希
//   - merkleProof: 默克尔证明
//
// **返回**：
//   - data: 查询结果（JSON格式）
//   - error: 错误信息
func QueryDatabase(
	dbIdentifier string,
	query string,
	params []interface{},
	stateHash []byte,
	merkleProof []byte,
) ([]byte, error) {
	// 构建查询参数
	queryParams := map[string]interface{}{
		"query": query,
	}
	if params != nil {
		queryParams["params"] = params
	}

	// 构建验证佐证
	evidence := &framework.Evidence{
		DataIntegrity: merkleProof, // 使用默克尔证明作为数据完整性证明
		ResponseHash:  stateHash,    // 使用状态哈希作为响应哈希
	}

	// 调用ValidateAndQuery
	return ValidateAndQuery("database_query", dbIdentifier, queryParams, evidence)
}

