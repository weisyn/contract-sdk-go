# WES Error Specification 实施文档

## 概述

本文档记录了 `contract-sdk-go.git`（Go 合约 SDK）中 WES Error Specification 的实施情况。

## 实施范围

### 1. 错误码定义

**文件**: `framework/errors.go`, `framework/contract_base.go`

- 定义了标准的合约错误码（`uint32`），与 HostABI 错误码对齐
- 提供了 `ContractError` 类型用于错误处理
- 错误码包括：
  - `SUCCESS` (0): 成功
  - `ERROR_INVALID_PARAMS` (1): 参数无效
  - `ERROR_INSUFFICIENT_BALANCE` (2): 余额不足
  - `ERROR_UNAUTHORIZED` (3): 未授权
  - `ERROR_NOT_FOUND` (4): 资源不存在
  - `ERROR_ALREADY_EXISTS` (5): 资源已存在
  - `ERROR_EXECUTION_FAILED` (6): 执行失败
  - `ERROR_INVALID_STATE` (7): 状态无效
  - `ERROR_TIMEOUT` (8): 超时
  - `ERROR_NOT_IMPLEMENTED` (9): 未实现
  - `ERROR_PERMISSION_DENIED` (10): 权限不足
  - `ERROR_UNKNOWN` (999): 未知错误

### 2. 错误码映射

**文件**: `framework/error_mapping.go`

- 提供了合约错误码到 WES 错误码的映射函数
- `ContractErrorCodeToWESCode()`: 映射到 WES 错误码字符串
- `ContractErrorCodeToUserMessage()`: 映射到用户友好的消息
- `ContractErrorCodeToHTTPStatus()`: 映射到 HTTP 状态码

**注意**: 此文件仅在非合约环境中编译（`!tinygo && !(js && wasm)`），用于文档和测试。

### 3. 错误处理流程

合约 SDK 的错误处理流程：

1. **合约执行时**：合约返回错误码（`uint32`）
2. **区块链服务层**：捕获错误码并转换为 Problem Details
3. **客户端**：接收 Problem Details 格式的错误响应

## 错误码映射表

| 合约错误码 | WES 错误码 | HTTP 状态码 | 用户消息 |
|-----------|-----------|-----------|---------|
| `SUCCESS` (0) | - | 200 | - |
| `ERROR_INVALID_PARAMS` (1) | `COMMON_VALIDATION_ERROR` | 400 | 参数验证失败，请检查输入参数。 |
| `ERROR_INSUFFICIENT_BALANCE` (2) | `BC_INSUFFICIENT_BALANCE` | 422 | 余额不足，无法完成交易。 |
| `ERROR_UNAUTHORIZED` (3) | `COMMON_VALIDATION_ERROR` | 401 | 未授权操作，请检查权限。 |
| `ERROR_NOT_FOUND` (4) | `BC_CONTRACT_NOT_FOUND` | 404 | 资源不存在。 |
| `ERROR_ALREADY_EXISTS` (5) | `COMMON_VALIDATION_ERROR` | 409 | 资源已存在。 |
| `ERROR_EXECUTION_FAILED` (6) | `BC_CONTRACT_INVOCATION_FAILED` | 422 | 合约执行失败，请检查合约逻辑。 |
| `ERROR_INVALID_STATE` (7) | `BC_CONTRACT_INVOCATION_FAILED` | 422 | 合约状态无效，请检查合约状态。 |
| `ERROR_TIMEOUT` (8) | `COMMON_TIMEOUT` | 408 | 执行超时，请稍后重试。 |
| `ERROR_NOT_IMPLEMENTED` (9) | `BC_CONTRACT_INVOCATION_FAILED` | 501 | 功能未实现。 |
| `ERROR_PERMISSION_DENIED` (10) | `COMMON_VALIDATION_ERROR` | 403 | 权限不足，无法执行此操作。 |
| `ERROR_UNKNOWN` (999) | `COMMON_INTERNAL_ERROR` | 500 | 未知错误，请稍后重试或联系管理员。 |

## 使用示例

### 在合约中返回错误

```go
package main

import "github.com/weisyn/contract-sdk-go/framework"

func Transfer(ctx *framework.Context, to framework.Address, amount framework.Amount) uint32 {
    // 检查余额
    balance := ctx.GetBalance(ctx.GetCaller())
    if balance < amount {
        return framework.ERROR_INSUFFICIENT_BALANCE
    }
    
    // 检查参数
    if to.IsZero() {
        return framework.ERROR_INVALID_PARAMS
    }
    
    // 执行转账
    // ...
    
    return framework.SUCCESS
}
```

### 错误码映射示例

```go
// 在非合约环境中（用于文档和测试）
import "github.com/weisyn/contract-sdk-go/framework"

// 将合约错误码映射到 WES 错误码
wesCode := framework.ContractErrorCodeToWESCode(framework.ERROR_INSUFFICIENT_BALANCE)
// wesCode = "BC_INSUFFICIENT_BALANCE"

// 获取用户友好的消息
userMsg := framework.ContractErrorCodeToUserMessage(framework.ERROR_INSUFFICIENT_BALANCE)
// userMsg = "余额不足，无法完成交易。"

// 获取 HTTP 状态码
status := framework.ContractErrorCodeToHTTPStatus(framework.ERROR_INSUFFICIENT_BALANCE)
// status = 422
```

## 与区块链服务层的集成

合约执行错误会被区块链服务层（`weisyn.git`）捕获并转换为 Problem Details：

1. **合约返回错误码**：例如 `ERROR_INSUFFICIENT_BALANCE` (2)
2. **服务层转换**：映射到 `BC_INSUFFICIENT_BALANCE`，创建 Problem Details
3. **客户端接收**：客户端 SDK 解析 Problem Details，显示用户友好的错误消息

## 注意事项

1. **合约环境限制**：合约代码在 WASM 环境中运行，不能直接使用 Problem Details
2. **错误码对齐**：合约错误码必须与 HostABI 错误码对齐
3. **映射一致性**：错误码映射必须与区块链服务层的映射一致

## 参考

- [WES Error Specification](../../weisyn.git/docs/error-spec/README.md)
- [HostABI 错误码定义](../../weisyn.git/internal/core/ispc/hostabi/errors.go)

