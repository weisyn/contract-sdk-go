//go:build tinygo || (js && wasm)

package framework

import (
	"testing"
)

// TestAddress 测试地址类型
func TestAddress(t *testing.T) {
	addr := Address{0x01, 0x02, 0x03}

	// 注意：Address.String() 在非WASM环境下会panic
	// 因为它依赖内存分配
	t.Logf("Address created with first 3 bytes: %02x %02x %02x", addr[0], addr[1], addr[2])
}

// TestHash 测试哈希类型
func TestHash(t *testing.T) {
	hash := Hash{0x01, 0x02, 0x03, 0x04}
	t.Logf("Hash created with first 4 bytes: %02x %02x %02x %02x", hash[0], hash[1], hash[2], hash[3])
}

// TestEvent 测试事件创建
func TestEvent(t *testing.T) {
	event := NewEvent("TestEvent")

	if event.Name != "TestEvent" {
		t.Errorf("Event.Name = %s, want TestEvent", event.Name)
	}

	if event.Data == nil {
		t.Error("Event.Data should not be nil")
	}

	// 测试添加数据
	event.Data["key1"] = "value1"
	if event.Data["key1"] != "value1" {
		t.Errorf("Event.Data[key1] = %s, want value1", event.Data["key1"])
	}
}

// TestContractError 测试错误类型
func TestContractError(t *testing.T) {
	err := NewContractError(ERROR_INVALID_PARAMS, "test error")

	if err.Code != ERROR_INVALID_PARAMS {
		t.Errorf("ContractError.Code = %d, want %d", err.Code, ERROR_INVALID_PARAMS)
	}

	if err.Message != "test error" {
		t.Errorf("ContractError.Message = %s, want 'test error'", err.Message)
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("ContractError.Error() returned empty string")
	}
}

// TestErrorCodes 测试错误码常量
func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name string
		code uint32
	}{
		{"SUCCESS", SUCCESS},
		{"ERROR_UNKNOWN", ERROR_UNKNOWN},
		{"ERROR_INVALID_PARAMS", ERROR_INVALID_PARAMS},
		{"ERROR_UNAUTHORIZED", ERROR_UNAUTHORIZED},
		{"ERROR_EXECUTION_FAILED", ERROR_EXECUTION_FAILED},
	}

	// 验证错误码唯一性
	seen := make(map[uint32]string)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if prev, exists := seen[tt.code]; exists {
				t.Errorf("Duplicate error code %d: %s and %s", tt.code, prev, tt.name)
			}
			seen[tt.code] = tt.name
		})
	}
}

// TestContractParams 测试合约参数
func TestContractParams(t *testing.T) {
	testData := []byte(`{"key1":"value1"}`)
	params := NewContractParams(testData)

	rawData := params.GetRawData()
	if len(rawData) != len(testData) {
		t.Errorf("GetRawData length = %d, want %d", len(rawData), len(testData))
	}
}

// TestHostFunctions 测试全局宿主函数（stub模式）
func TestHostFunctions(t *testing.T) {
	t.Run("GetABIVersion", func(t *testing.T) {
		version := GetABIVersion()
		t.Logf("ABI Version: 0x%08X", version)
		// stub返回 0x00010000 (v1.0.0)
		if version != 0x00010000 {
			t.Errorf("GetABIVersion = 0x%08X, want 0x00010000", version)
		}
	})

	t.Run("CheckABICompatibility", func(t *testing.T) {
		// stub总是返回成功
		err := CheckABICompatibility(0x00010000)
		if err != nil {
			t.Errorf("CheckABICompatibility failed: %v", err)
		}
	})

	t.Run("GetCaller", func(t *testing.T) {
		// stub返回空地址
		addr := GetCaller()
		t.Logf("GetCaller (stub): %v", addr)
	})

	t.Run("GetTimestamp", func(t *testing.T) {
		// stub返回0
		ts := GetTimestamp()
		t.Logf("GetTimestamp (stub): %d", ts)
	})

	t.Run("GetBlockHeight", func(t *testing.T) {
		// stub返回0
		height := GetBlockHeight()
		t.Logf("GetBlockHeight (stub): %d", height)
	})
}

// TestStateOperationsStub 测试状态操作（stub模式）
func TestStateOperationsStub(t *testing.T) {
	// ⚠️ **已删除**：PutState 和 StateExists 测试
	// 原因：违背WES架构原则，EUTXO模型无全局状态存储

	t.Run("GetState", func(t *testing.T) {
		// stub返回空数据
		data, err := GetState("key1")
		if err != nil {
			t.Errorf("GetState failed: %v", err)
		}
		t.Logf("GetState (stub) returned %d bytes", len(data))
	})
}

// TestContractBase 测试 ContractBase 便捷方法
func TestContractBase(t *testing.T) {
	cb := &ContractBase{}

	t.Run("GetCaller", func(t *testing.T) {
		// 注意：会调用 Address.String()，在stub环境下会panic
		// 这里跳过测试
		t.Skip("GetCaller requires WASM environment")
	})

	t.Run("GetTimestamp", func(t *testing.T) {
		ts := cb.GetTimestamp()
		t.Logf("GetTimestamp: %d (stub)", ts)
	})

	t.Run("GetBlockHeight", func(t *testing.T) {
		height := cb.GetBlockHeight()
		t.Logf("GetBlockHeight: %d (stub)", height)
	})

	t.Run("SetReturnData", func(t *testing.T) {
		err := cb.SetReturnData([]byte("test"))
		if err != nil {
			t.Errorf("SetReturnData failed: %v", err)
		}
	})

	// ⚠️ **已删除**：PutState 测试
	// 原因：违背WES架构原则，EUTXO模型无全局状态存储

	t.Run("GetState", func(t *testing.T) {
		data := cb.GetState("key1")
		t.Logf("GetState returned %d bytes (stub)", len(data))
	})
}
