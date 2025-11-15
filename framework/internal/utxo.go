//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== UTXO操作相关（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// CreateUTXO 创建UTXO输出
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func CreateUTXO(recipient framework.Address, amount framework.Amount, tokenID framework.TokenID) error {
	recipientPtr, _ := framework.AllocateBytes(recipient.ToBytes())
	if recipientPtr == 0 {
		return framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate recipient address")
	}

	tokenIDPtr, tokenIDLen := framework.AllocateString(string(tokenID))
	if tokenIDPtr == 0 {
		return framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate token ID")
	}

	result := createUTXOOutput(recipientPtr, uint64(amount), tokenIDPtr, tokenIDLen)
	if result != framework.SUCCESS {
		return framework.NewContractError(result, "failed to create UTXO output")
	}

	return nil
}

// CreateAssetOutputWithLock 创建带锁定条件的资产输出
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func CreateAssetOutputWithLock(recipient framework.Address, amount framework.Amount, tokenID framework.TokenID, lockingBytes []byte) (uint32, error) {
	recipientPtr, _ := framework.AllocateBytes(recipient.ToBytes())
	tokenIDPtr, tokenIDLen := framework.AllocateString(string(tokenID))
	var lockPtr, lockLen uint32
	if len(lockingBytes) > 0 {
		lockPtr, lockLen = framework.AllocateBytes(lockingBytes)
	}
	idx := createAssetOutputWithLock(recipientPtr, 20, uint64(amount), tokenIDPtr, tokenIDLen, lockPtr, lockLen)
	if idx == 0xFFFFFFFF {
		return idx, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "create_asset_output_with_lock failed")
	}
	return idx, nil
}

// ⚠️ **已删除**：ExecuteUTXOTransferEx
//
// **原因**：违背WES"无业务语义"架构原则，功能不完整（仅创建输出，不处理UTXO选择和找零）
//
// **替代方案**：
// 1. 使用 helpers/token/Transfer - 包含完整的转账逻辑（推荐）
// 2. 使用 framework.BeginTransaction().Transfer().Finalize() - 链式API，包含完整业务逻辑
// 3. 直接使用原语函数 CreateAssetOutputWithLock（仅创建输出）

// ==================== 宿主函数声明 ====================

//go:wasmimport env create_utxo_output
func createUTXOOutput(recipientPtr uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32) uint32

//go:wasmimport env create_asset_output_with_lock
func createAssetOutputWithLock(recipientPtr uint32, recipientLen uint32, amount uint64, tokenIDPtr uint32, tokenIDLen uint32, lockPtr uint32, lockLen uint32) uint32

// ⚠️ **已移除**：execute_utxo_transfer_ex
// 原因：违背WES"无业务语义"架构原则
// 请使用原语函数：create_asset_output_with_lock
