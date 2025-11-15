//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== HostABI 原语封装（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// TxAddInput 添加交易输入
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func TxAddInput(outpoint framework.OutPoint, isReferenceOnly bool, unlockingProof framework.UnlockingProof) (uint32, error) {
	// 验证参数
	if len(outpoint.TxHash) != 32 {
		return 0, framework.NewContractError(framework.ERROR_INVALID_PARAMS, "txHash must be 32 bytes")
	}

	// 分配txID缓冲区
	txIDPtr, _ := framework.AllocateBytes(outpoint.TxHash)
	if txIDPtr == 0 {
		return 0, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate txID")
	}

	// 分配proof缓冲区
	var proofPtr, proofLen uint32
	if len(unlockingProof.ProofData) > 0 {
		proofPtr, proofLen = framework.AllocateBytes(unlockingProof.ProofData)
		if proofPtr == 0 {
			return 0, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate proof")
		}
	}

	// 转换isReferenceOnly为uint32（1=true, 0=false）
	var isRefOnly uint32
	if isReferenceOnly {
		isRefOnly = 1
	}

	// 调用宿主函数
	inputIndex := appendTxInput(txIDPtr, 32, outpoint.Index, isRefOnly, proofPtr, proofLen)
	if inputIndex == 0xFFFFFFFF {
		return 0, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to add input")
	}

	return inputIndex, nil
}

// TxAddAssetOutput 添加资产输出（HostABI原语）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
//
// 设计说明：
// - 统一使用 CreateAssetOutputWithLock，无论是否有锁定条件
// - 这样可以保证返回值一致性（始终返回输出索引）
// - 当没有锁定条件时，传入空的锁定条件数组
func TxAddAssetOutput(owner framework.Address, amount framework.Amount, tokenID framework.TokenID, lockingConditions []framework.LockingCondition) (uint32, error) {
	// 统一使用 CreateAssetOutputWithLock，保证返回值一致性
	lockingJSON := buildLockingConditionsJSON(lockingConditions)
	return CreateAssetOutputWithLock(owner, amount, tokenID, lockingJSON)
}

// TxAddResourceOutput 添加资源输出（HostABI原语）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
//
// 注意：当前实现中，category 和 metadata 参数暂未使用
// 这些信息应该编码到 resourceBytes 中，或者通过其他方式传递
func TxAddResourceOutput(contentHash []byte, category string, owner framework.Address, lockingConditions []framework.LockingCondition, metadata []byte) (uint32, error) {
	// 当前使用现有的AppendResourceOutput方法
	// TODO: 如果需要支持 category 和 metadata，需要更新 AppendResourceOutput 函数签名
	// 或者将这些信息编码到 resourceBytes 中
	lockingJSON := buildLockingConditionsJSON(lockingConditions)
	return AppendResourceOutput(contentHash, owner, lockingJSON)
}

// TxAddStateOutput 添加状态输出（HostABI原语）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func TxAddStateOutput(stateID []byte, stateVersion uint64, executionResultHash []byte, publicInputs []byte, parentStateHash []byte) (uint32, error) {
	// 使用现有的AppendStateOutputSimple方法
	return AppendStateOutputSimple(stateID, stateVersion, executionResultHash, parentStateHash)
}

// ==================== 宿主函数声明 ====================

//go:wasmimport env append_tx_input
func appendTxInput(txIDPtr uint32, txIDLen uint32, index uint32, isReferenceOnly uint32, proofPtr uint32, proofLen uint32) uint32

// ==================== 辅助函数 ====================

// buildLockingConditionsJSON 构建锁定条件JSON
func buildLockingConditionsJSON(conditions []framework.LockingCondition) []byte {
	if len(conditions) == 0 {
		return nil
	}

	// 将LockingCondition转换为JSON字符串数组
	jsonStrings := make([]string, len(conditions))
	for i, cond := range conditions {
		// 简化实现：假设Condition已经是JSON字符串
		jsonStrings[i] = string(cond.Condition)
	}

	return framework.BuildLockingJSONArray(jsonStrings)
}

