package main

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// StandardToken ERC20风格的标准代币合约模板
type StandardToken struct {
	framework.ContractBase
}

// ================================================================================================
// 核心导出方法（Host ABI v1.1规范）
// ================================================================================================

// Initialize 初始化代币合约
//
// 参数（通过GetContractInitParams获取）：
//   - name: 代币名称
//   - symbol: 代币符号
//   - decimals: 小数位数
//   - initialSupply: 初始供应量
//
// 输出：
//   - StateOutput: totalSupply（总供应量）
//   - StateOutput: balance_{owner}（所有者余额）
//   - Event: TokenInitialized
//
//export Initialize
func Initialize() uint32 {
	contract := &StandardToken{}

	// TODO: 解析初始化参数
	// params := contract.GetContractInitParams()

	owner := []byte(contract.GetCaller())
	initialSupply := uint64(1000000) // 示例：100万代币

	// 1. 设置总供应量（使用 AppendStateOutputSimple）
	totalSupplyStateID := []byte("totalSupply")
	totalSupplyValue := uint64ToBytes(initialSupply)
	if _, err := framework.AppendStateOutputSimple(totalSupplyStateID, 1, totalSupplyValue, nil); err != nil {
		contract.EmitLog("error", "Failed to set totalSupply")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 2. 设置所有者余额
	ownerBalanceStateID := append([]byte("balance_"), owner...)
	if _, err := framework.AppendStateOutputSimple(ownerBalanceStateID, 1, totalSupplyValue, nil); err != nil {
		contract.EmitLog("error", "Failed to set owner balance")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 发出初始化事件
	contract.EmitEvent("TokenInitialized", owner)

	return framework.SUCCESS
}

// Transfer 转账代币
//
// 参数（通过GetContractInitParams获取）：
//   - to: 接收者地址
//   - amount: 转账金额
//
// 输出：
//   - StateOutput: balance_{from}（发送者余额）
//   - StateOutput: balance_{to}（接收者余额）
//   - Event: Transfer
//
//export Transfer
func Transfer() uint32 {
	contract := &StandardToken{}

	// TODO: 解析转账参数
	from := []byte(contract.GetCaller())
	to := []byte("recipient_address") // TODO: 从参数获取
	amount := uint64(100)             // TODO: 从参数获取

	// 1. 检查发送者余额
	fromBalanceStateID := append([]byte("balance_"), from...)
	fromBalanceData := contract.GetState(string(fromBalanceStateID))
	fromBalance := bytesToUint64(fromBalanceData)

	if fromBalance < amount {
		contract.EmitLog("error", "Insufficient balance")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 2. 更新发送者余额
	newFromBalance := fromBalance - amount
	if _, err := framework.AppendStateOutputSimple(fromBalanceStateID, 1, uint64ToBytes(newFromBalance), nil); err != nil {
		contract.EmitLog("error", "Failed to update sender balance")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 更新接收者余额
	toBalanceStateID := append([]byte("balance_"), to...)
	toBalanceData := contract.GetState(string(toBalanceStateID))
	toBalance := bytesToUint64(toBalanceData)
	newToBalance := toBalance + amount

	if _, err := framework.AppendStateOutputSimple(toBalanceStateID, 1, uint64ToBytes(newToBalance), nil); err != nil {
		contract.EmitLog("error", "Failed to update receiver balance")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 发出转账事件
	contract.EmitEvent("Transfer", append(append(from, to...), uint64ToBytes(amount)...))

	return framework.SUCCESS
}

// BalanceOf 查询余额（只读）
//
// 参数：
//   - address: 查询地址
//
// 返回：
//   - balance: 余额（通过SetReturnData返回）
//
//export BalanceOf
func BalanceOf() uint32 {
	contract := &StandardToken{}

	// TODO: 解析地址参数
	address := []byte(contract.GetCaller())

	balanceStateID := append([]byte("balance_"), address...)
	balanceData := contract.GetState(string(balanceStateID))
	balance := bytesToUint64(balanceData)

	contract.SetReturnData(uint64ToBytes(balance))
	return framework.SUCCESS
}

// Approve 授权额度（TODO: v1.2实现）
//
//export Approve
func Approve() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// TransferFrom 代理转账（TODO: v1.2实现）
//
//export TransferFrom
func TransferFrom() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// ================================================================================================
// 辅助函数
// ================================================================================================

func uint64ToBytes(n uint64) []byte {
	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[7-i] = byte(n >> (i * 8))
	}
	return result
}

func bytesToUint64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	var result uint64
	for i := 0; i < 8; i++ {
		result |= uint64(b[7-i]) << (i * 8)
	}
	return result
}

func main() {}
