//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== 状态输出相关（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// AppendStateOutput 追加 StateOutput（显式状态锚定）
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func AppendStateOutput(stateID []byte, version uint64, execHash []byte, zkProof []byte, parentHash []byte) (uint32, error) {
	// stateID
	stateIDPtr, stateIDLen := framework.AllocateBytes(stateID)
	if stateIDLen == 0 {
		return 0xFFFFFFFF, framework.NewContractError(framework.ERROR_INVALID_PARAMS, "invalid empty stateID")
	}
	
	// execHash（必需，必须是32字节）
	var execHash32 [32]byte
	if len(execHash) == 32 {
		copy(execHash32[:], execHash)
	} else if len(execHash) > 0 {
		// 如果execHash不是32字节，使用ComputeHash计算32字节哈希
		hash := framework.ComputeHash(execHash)
		copy(execHash32[:], hash[:])
	} else {
		// 如果execHash为空，使用stateID的哈希
		hash := framework.ComputeHash(stateID)
		copy(execHash32[:], hash[:])
	}
	
	execHashPtr := framework.Malloc(32)
	if execHashPtr == 0 {
		return 0xFFFFFFFF, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate execHash")
	}
	execHashBytes := framework.GetBytes(execHashPtr, 32)
	copy(execHashBytes, execHash32[:])
	
	// publicInputs：使用execHash作为公开输入
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
			hash := framework.ComputeHash(parentHash)
			copy(parentHash32[:], hash[:])
		}
		parentPtr = framework.Malloc(32)
		if parentPtr == 0 {
			return 0xFFFFFFFF, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to allocate parentHash")
		}
		parentHashBytes := framework.GetBytes(parentPtr, 32)
		copy(parentHashBytes, parentHash32[:])
	}

	// 调用宿主函数（新签名：7个参数）
	idx := appendStateOutput(stateIDPtr, stateIDLen, version, execHashPtr, publicInputsPtr, publicInputsLen, parentPtr)
	if idx == 0xFFFFFFFF {
		return idx, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "append_state_output failed")
	}
	return idx, nil
}

// AppendStateOutputSimple 追加 StateOutput
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func AppendStateOutputSimple(stateID []byte, version uint64, execHash []byte, parentHash []byte) (uint32, error) {
	// 直接调用原函数，zkProof传nil（会被忽略）
	return AppendStateOutput(stateID, version, execHash, nil, parentHash)
}

// ==================== 宿主函数声明 ====================

//go:wasmimport env append_state_output
func appendStateOutput(stateIDPtr uint32, stateIDLen uint32, version uint64, execHashPtr uint32, publicInputsPtr uint32, publicInputsLen uint32, parentHashPtr uint32) uint32

// ⚠️ **已迁移**：GetStateVersion 和 IncrementStateVersion 已迁移到 framework/hostabi.go
// 这些函数现在作为 framework 包的公开API提供

