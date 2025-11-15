//go:build tinygo || (js && wasm)

package rwa

import (
	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/external"
	"github.com/weisyn/contract-sdk-go/helpers/token"
)

// ValidateAndTokenizeResult 验证并代币化结果
type ValidateAndTokenizeResult struct {
	TokenID         framework.TokenID
	Validated       bool
	ValidationProof []byte // 验证过程的ZK证明
	Valuation       uint64
	ValuationProof  []byte // 估值过程的ZK证明
}

// ValidateAndTokenize 验证并代币化资产
//
// 🎯 **用途**：通过ISPC受控机制验证资产并代币化，替代传统预言机
//
// **ISPC 创新点**：
//   传统区块链需要中心化的预言机服务来获取外部数据（资产验证、估值等）。
//   WES ISPC 通过"受控声明+佐证+验证"机制，让合约可以直接调用外部 API，
//   无需传统预言机。这是 ISPC 的核心创新之一。
//
// **ISPC 工作原理**：
//   1. 声明外部状态预期：告诉系统"我要调用验证API和估值API，预期得到这样的数据"
//   2. 提供验证佐证：提供 API 数字签名、响应哈希等密码学佐证
//   3. 运行时验证：ISPC 运行时验证佐证的有效性
//   4. 记录到执行轨迹：外部调用被记录到执行轨迹
//   5. 生成 ZK 证明：执行轨迹自动生成 ZK 证明（包含外部交互验证）
//   6. 验证节点验证证明：其他节点验证证明，无需重复调用外部 API
//   7. 自动上链：验证通过后，结果自动上链
//
// **参数**：
//   - assetID: 资产ID
//   - documents: 资产文档（JSON格式）
//   - validatorAPI: 验证服务API端点（如 "https://validator.example.com/api/validate"）
//   - validatorEvidence: 验证机构签名等佐证（包含 APISignature、ResponseHash 等）
//   - valuationAPI: 估值服务API端点（如 "https://valuation.example.com/api/value"）
//   - valuationEvidence: 估值服务签名等佐证（包含 APISignature、ResponseHash 等）
//
// **返回**：
//   - result: 验证并代币化结果，包含：
//     * TokenID: 生成的代币ID
//     * Validated: 是否验证通过
//     * ValidationProof: 验证过程的ZK证明
//     * Valuation: 资产估值
//     * ValuationProof: 估值过程的ZK证明
//   - error: 错误信息，nil表示成功
//
// **与传统区块链的对比**：
//   传统区块链：
//     - 需要预言机服务调用外部验证和估值 API
//     - 预言机将结果喂入链上
//     - 合约使用预言机提供的数据
//     - 问题：预言机是中心化瓶颈，需要支付费用，存在延迟
//
//   WES ISPC：
//     - 直接调用外部验证和估值 API
//     - 单次调用，多点验证，自动生成 ZK 证明
//     - 无需传统预言机，直接获取外部数据
//     - 执行后自动上链，用户直接获得结果
//
// **示例**：
//
//	result, err := rwa.ValidateAndTokenize(
//	    "real_estate_001",
//	    assetDocuments,
//	    "https://validator.example.com/api/validate",
//	    &framework.Evidence{
//	        APISignature: validatorSignature,  // API 数字签名（从外部服务获取）
//	        ResponseHash: validationResponseHash,  // 响应数据哈希（从外部服务获取）
//	    },
//	    "https://valuation.example.com/api/value",
//	    &framework.Evidence{
//	        APISignature: valuationSignature,
//	        ResponseHash: valuationResponseHash,
//	    },
//	)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// ✅ 用户直接获得业务结果，无需知道交易的存在
//	// ✅ ZK 证明自动生成，自动构建交易，自动上链
func ValidateAndTokenize(
	assetID string,
	documents []byte,
	validatorAPI string,
	validatorEvidence *framework.Evidence,
	valuationAPI string,
	valuationEvidence *framework.Evidence,
) (*ValidateAndTokenizeResult, error) {
	// 1. 参数验证
	if assetID == "" {
		return nil, framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"assetID cannot be empty",
		)
	}
	if validatorEvidence == nil || valuationEvidence == nil {
		return nil, framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"evidence cannot be nil",
		)
	}

	// 2. 通过ISPC受控机制验证资产
	validationParams := map[string]interface{}{
		"asset_id": assetID,
		"documents": string(documents),
	}
	validationData, err := external.ValidateAndQuery(
		"api_response",
		validatorAPI,
		validationParams,
		validatorEvidence,
	)
	if err != nil {
		return nil, err
	}

	// 解析验证结果（真实JSON解析实现）
	// 验证数据格式：{"validated": true/false, "proof": "..."}
	validated := false
	if len(validationData) > 0 {
		// 使用framework的JSON解析功能
		params := framework.NewContractParams(validationData)
		validatedStr := params.ParseJSON("validated")
		validated = validatedStr == "true" || validatedStr == "1"
	}

	// 3. 通过ISPC受控机制获取资产估值
	valuationParams := map[string]interface{}{
		"asset_id": assetID,
	}
	valuationData, err := external.ValidateAndQuery(
		"api_response",
		valuationAPI,
		valuationParams,
		valuationEvidence,
	)
	if err != nil {
		return nil, err
	}

	// 解析估值结果（真实JSON解析实现）
	// 估值数据格式：{"value": 1000000, "currency": "USD", "proof": "..."}
	valuation := uint64(0)
	if len(valuationData) > 0 {
		params := framework.NewContractParams(valuationData)
		valueStr := params.ParseJSON("value")
		if valueStr != "" {
			valuation = framework.ParseUint64(valueStr)
		}
		// 如果解析失败，使用默认值0（表示估值失败）
		if valuation == 0 {
			return nil, framework.NewContractError(
				framework.ERROR_INVALID_PARAMS,
				"failed to parse valuation value",
			)
		}
	}

	// 4. 执行代币化（使用helpers/token）
	caller := framework.GetCaller()
	tokenID := framework.TokenID("RWA_" + assetID)
	err = token.Mint(caller, tokenID, framework.Amount(valuation))
	if err != nil {
		return nil, err
	}

	// 5. 返回结果（包含验证和估值的证明）
	// 注意：validationData和valuationData已经包含ISPC生成的ZK证明
	return &ValidateAndTokenizeResult{
		TokenID:         tokenID,
		Validated:       validated,
		ValidationProof: validationData, // ISPC自动生成的ZK证明
		Valuation:       valuation,
		ValuationProof:  valuationData, // ISPC自动生成的ZK证明
	}, nil
}

// ValidateAsset 验证资产
//
// 🎯 **用途**：通过ISPC受控机制验证资产，替代传统预言机
//
// **参数**：
//   - assetID: 资产ID
//   - documents: 资产文档
//   - validatorAPI: 验证服务API端点
//   - evidence: 验证机构签名等佐证
//
// **返回**：
//   - validated: 是否验证通过
//   - proof: 验证过程的ZK证明
//   - error: 错误信息
func ValidateAsset(
	assetID string,
	documents []byte,
	validatorAPI string,
	evidence *framework.Evidence,
) (bool, []byte, error) {
	// 构建验证参数
	params := map[string]interface{}{
		"asset_id": assetID,
		"documents": string(documents),
	}

	// 通过ISPC受控机制验证
	data, err := external.ValidateAndQuery("api_response", validatorAPI, params, evidence)
	if err != nil {
		return false, nil, err
	}

	// 解析验证结果（真实JSON解析实现）
	validated := false
	if len(data) > 0 {
		params := framework.NewContractParams(data)
		validatedStr := params.ParseJSON("validated")
		validated = validatedStr == "true" || validatedStr == "1"
	}

	return validated, data, nil
}

// ValueAsset 评估资产价值
//
// 🎯 **用途**：通过ISPC受控机制评估资产价值，替代传统预言机
//
// **参数**：
//   - assetID: 资产ID
//   - valuationAPI: 估值服务API端点
//   - evidence: 估值服务签名等佐证
//
// **返回**：
//   - value: 资产价值
//   - proof: 估值过程的ZK证明
//   - error: 错误信息
func ValueAsset(
	assetID string,
	valuationAPI string,
	evidence *framework.Evidence,
) (uint64, []byte, error) {
	// 构建估值参数
	params := map[string]interface{}{
		"asset_id": assetID,
	}

	// 通过ISPC受控机制估值
	data, err := external.ValidateAndQuery("api_response", valuationAPI, params, evidence)
	if err != nil {
		return 0, nil, err
	}

	// 解析估值结果（真实JSON解析实现）
	// 估值数据格式：{"value": 1000000, "currency": "USD", "proof": "..."}
	value := uint64(0)
	if len(data) > 0 {
		params := framework.NewContractParams(data)
		valueStr := params.ParseJSON("value")
		if valueStr != "" {
			value = framework.ParseUint64(valueStr)
		}
		// 如果解析失败，返回0（表示估值失败）
		if value == 0 {
			return 0, data, framework.NewContractError(
				framework.ERROR_INVALID_PARAMS,
				"failed to parse valuation value",
			)
		}
	}

	return value, data, nil
}

