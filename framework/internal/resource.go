//go:build tinygo || (js && wasm)

package internal

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// ==================== 资源输出相关（内部包）====================
//
// ⚠️ **内部包**：此包仅供 helpers 层使用，外部开发者不应导入

// AppendResourceOutput 追加 ResourceOutput
//
// ⚠️ **内部接口**：仅供 helpers 层使用
func AppendResourceOutput(resourceBytes []byte, owner framework.Address, lockingBytes []byte) (uint32, error) {
	resourcePtr, resourceLen := framework.AllocateBytes(resourceBytes)
	if resourceLen == 0 {
		return 0xFFFFFFFF, framework.NewContractError(framework.ERROR_INVALID_PARAMS, "invalid empty resource bytes")
	}
	ownerPtr, _ := framework.AllocateBytes(owner.ToBytes())
	var lockPtr, lockLen uint32
	if len(lockingBytes) > 0 {
		lockPtr, lockLen = framework.AllocateBytes(lockingBytes)
	}
	idx := appendResourceOutput(resourcePtr, resourceLen, ownerPtr, 20, lockPtr, lockLen)
	if idx == 0xFFFFFFFF {
		return idx, framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "append_resource_output not implemented")
	}
	return idx, nil
}

// ==================== 宿主函数声明 ====================

//go:wasmimport env append_resource_output
func appendResourceOutput(resourcePtr uint32, resourceLen uint32, ownerPtr uint32, ownerLen uint32, lockPtr uint32, lockLen uint32) uint32

