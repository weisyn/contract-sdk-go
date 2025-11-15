//go:build tinygo || (js && wasm)

package token

import (
	"strconv"

	"github.com/weisyn/contract-sdk-go/framework"
)

// MintWithState 基于状态存储的代币铸造
//
// 🎯 **用途**：在合约代码中铸造新代币，使用状态存储模式
//
// **参数**：
//   - to: 接收者地址
//   - tokenID: 代币ID（可选，用于多代币合约）
//   - amount: 铸造数量
//   - balanceStateKey: 余额状态键（例如："balance_user123"）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 从链上读取当前余额
//   - 增加余额
//   - 更新状态（自动递增版本号）
//   - 发出铸造事件
//
// **示例**：
//
//	func Mint() uint32 {
//	    caller := framework.GetCaller()
//	    balanceKey := "balance_" + caller.String()
//	    
//	    err := token.MintWithState(
//	        caller,
//	        framework.TokenID(""),
//	        framework.Amount(1000),
//	        balanceKey,
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func MintWithState(to framework.Address, tokenID framework.TokenID, amount framework.Amount, balanceStateKey string) error {
	// 1. 参数验证
	if err := validateMintParams(to, tokenID, amount); err != nil {
		return err
	}

	// 2. 从链上读取当前余额
	stateID := []byte(balanceStateKey)
	currentBalanceData, version, err := framework.GetStateFromChain(stateID)
	if err != nil {
		// 如果状态不存在，版本号为0，余额为0
		version = 0
		currentBalanceData = []byte("0")
	}

	// 3. 解析余额
	currentBalance := parseBalanceFromBytes(currentBalanceData)

	// 4. 计算新余额
	newBalance := currentBalance + uint64(amount)

	// 5. 递增版本号
	newVersion := version + 1

	// 6. 更新状态
	execHash := framework.GetTxHash()
	_, err = framework.AppendStateOutputSimple(stateID, newVersion, []byte(strconv.FormatUint(newBalance, 10)), execHash.ToBytes())
	if err != nil {
		return framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to update balance state")
	}

	// 7. 发出铸造事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Mint")
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	event.AddAddressField("minter", caller)
	framework.EmitEvent(event)

	return nil
}

// TransferWithState 基于状态存储的代币转账
//
// 🎯 **用途**：在合约代码中执行转账，使用状态存储模式
//
// **参数**：
//   - from: 发送者地址
//   - to: 接收者地址
//   - tokenID: 代币ID（可选，用于多代币合约）
//   - amount: 转账金额
//   - balanceStateKeyPrefix: 余额状态键前缀（例如："balance_"）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 从链上读取发送者和接收者余额
//   - 检查发送者余额是否充足
//   - 扣除发送者余额
//   - 增加接收者余额
//   - 更新状态（自动递增版本号）
//   - 发出转账事件
//
// **示例**：
//
//	func Transfer() uint32 {
//	    caller := framework.GetCaller()
//	    params := framework.GetContractParams()
//	    toStr := params.ParseJSON("to")
//	    
//	    to, err := framework.ParseAddressBase58(toStr)
//	    if err != nil {
//	        return framework.ERROR_INVALID_PARAMS
//	    }
//	    
//	    amountStr := params.ParseJSON("amount")
//	    amount, _ := strconv.ParseUint(amountStr, 10, 64)
//	    
//	    err = token.TransferWithState(
//	        caller,
//	        to,
//	        framework.TokenID(""),
//	        framework.Amount(amount),
//	        "balance_",
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func TransferWithState(from, to framework.Address, tokenID framework.TokenID, amount framework.Amount, balanceStateKeyPrefix string) error {
	// 1. 参数验证
	if err := validateTransferParams(from, to, amount); err != nil {
		return err
	}

	// 2. 构建状态键
	fromKey := balanceStateKeyPrefix + from.String()
	toKey := balanceStateKeyPrefix + to.String()

	// 3. 从链上读取发送者余额
	fromStateID := []byte(fromKey)
	fromBalanceData, fromVersion, err := framework.GetStateFromChain(fromStateID)
	if err != nil {
		// 如果状态不存在，余额为0
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance",
		)
	}
	fromBalance := parseBalanceFromBytes(fromBalanceData)

	// 4. 检查余额
	if fromBalance < uint64(amount) {
		return framework.NewContractError(
			framework.ERROR_INSUFFICIENT_BALANCE,
			"insufficient balance",
		)
	}

	// 5. 从链上读取接收者余额
	toStateID := []byte(toKey)
	toBalanceData, toVersion, err := framework.GetStateFromChain(toStateID)
	if err != nil {
		// 如果状态不存在，版本号为0，余额为0
		toVersion = 0
		toBalanceData = []byte("0")
	}
	toBalance := parseBalanceFromBytes(toBalanceData)

	// 6. 计算新余额
	newFromBalance := fromBalance - uint64(amount)
	newToBalance := toBalance + uint64(amount)

	// 7. 更新发送者余额
	execHash := framework.GetTxHash()
	_, err = framework.AppendStateOutputSimple(fromStateID, fromVersion+1, []byte(strconv.FormatUint(newFromBalance, 10)), execHash.ToBytes())
	if err != nil {
		return framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to update sender balance state")
	}

	// 8. 更新接收者余额
	_, err = framework.AppendStateOutputSimple(toStateID, toVersion+1, []byte(strconv.FormatUint(newToBalance, 10)), execHash.ToBytes())
	if err != nil {
		return framework.NewContractError(framework.ERROR_EXECUTION_FAILED, "failed to update receiver balance state")
	}

	// 9. 发出转账事件
	event := framework.NewEvent("Transfer")
	event.AddAddressField("from", from)
	event.AddAddressField("to", to)
	event.AddStringField("token_id", string(tokenID))
	event.AddUint64Field("amount", uint64(amount))
	framework.EmitEvent(event)

	return nil
}

// GetBalanceFromState 从状态中读取余额
//
// 🎯 **用途**：查询指定地址的代币余额（状态存储模式）
//
// **参数**：
//   - address: 地址
//   - stateKey: 状态键（例如："balance_user123"）
//
// **返回**：
//   - balance: 余额
//   - error: 错误信息，nil表示成功
//
// **说明**：
//   - 从链上查询状态，解析余额
//   - 如果状态不存在，返回余额0
//
// **示例**：
//
//	balance, err := token.GetBalanceFromState(address, "balance_"+address.String())
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
func GetBalanceFromState(address framework.Address, stateKey string) (framework.Amount, error) {
	stateID := []byte(stateKey)
	balanceData, _, err := framework.GetStateFromChain(stateID)
	if err != nil {
		// 如果状态不存在，返回余额0
		return 0, nil
	}

	balance := parseBalanceFromBytes(balanceData)
	return framework.Amount(balance), nil
}

// parseBalanceFromBytes 从字节数组解析余额
func parseBalanceFromBytes(data []byte) uint64 {
	if len(data) == 0 {
		return 0
	}

	// 移除尾部的零字节
	data = trimTrailingZeros(data)
	if len(data) == 0 {
		return 0
	}

	// 解析字符串为uint64
	balanceStr := string(data)
	balance, err := strconv.ParseUint(balanceStr, 10, 64)
	if err != nil {
		return 0
	}

	return balance
}

// trimTrailingZeros 移除尾部的零字节
func trimTrailingZeros(data []byte) []byte {
	// 从后往前查找第一个非零字节
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return []byte{} // 全部是零
}

