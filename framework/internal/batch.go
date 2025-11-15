//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== 批量操作相关（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// BatchCreateOutputs 批量创建资产输出（JSON 承载 DTO）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func BatchCreateOutputs(batchJSON string) (uint32, error) {
	if batchJSON == "" {
		return 0, nil
	}
	batchPtr, batchLen := framework.AllocateString(batchJSON)
	if batchPtr == 0 {
		return 0, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate batch json")
	}
	n := batchCreateOutputs(batchPtr, batchLen)
	return n, nil
}

// BatchCreateOutputsSimple 批量创建资产输出（简化版，无锁定条件）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func BatchCreateOutputsSimple(items []struct {
	Recipient []byte
	Amount    uint64
	TokenID   []byte
}) (uint32, error) {
	// 构造 DTO JSON（手动序列化避免引入encoding/json）
	batchJSON := "["
	for i, it := range items {
		if i > 0 {
			batchJSON += ","
		}
		batchJSON += `{"recipient":"`
		batchJSON += base64Encode(it.Recipient)
		batchJSON += `","amount":`
		batchJSON += uint64ToStr(it.Amount)
		if len(it.TokenID) > 0 {
			batchJSON += `,"token_id":"`
			batchJSON += base64Encode(it.TokenID)
			batchJSON += `"`
		} else {
			batchJSON += `,"token_id":null`
		}
		batchJSON += `,"locking_conditions":[]}`
	}
	batchJSON += "]"

	batchBytes := []byte(batchJSON)
	batchPtr, batchLen := framework.AllocateBytes(batchBytes)
	if batchPtr == 0 {
		return 0, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate batch bytes")
	}

	n := batchCreateOutputs(batchPtr, batchLen)
	return n, nil
}

// ==================== 宿主函数声明 ====================

//go:wasmimport env batch_create_outputs
func batchCreateOutputs(batchPtr uint32, batchLen uint32) uint32

