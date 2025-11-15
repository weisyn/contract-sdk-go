//go:build tinygo || (js && wasm)

package main

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// StakingContract 质押合约模板
//
// 核心功能：
// - 代币质押（锁定）
// - 质押解锁（提取）
// - 奖励分发
// - 委托质押
type StakingContract struct {
	framework.ContractBase
}

// ================================================================================================
// 核心导出方法（Host ABI v1.1规范）
// ================================================================================================

// Initialize 初始化质押合约
//
// 参数：
//   - stakingToken: 质押代币地址
//   - rewardToken: 奖励代币地址
//   - minStakingAmount: 最小质押金额
//   - lockDuration: 锁定区块数
//
// 输出：
//   - StateOutput: 质押参数配置
//   - Event: StakingInitialized
//
//export Initialize
func Initialize() uint32 {
	contract := &StakingContract{}

	// TODO: 解析初始化参数
	owner := []byte(contract.GetCaller())
	minStakingAmount := uint64(1000) // 示例：最小质押1000代币

	// 1. 设置质押参数
	configStateID := []byte("config")
	configData := encodeConfig(minStakingAmount, 100) // 锁定100个区块
	if _, err := framework.AppendStateOutputSimple(configStateID, 1, configData, nil); err != nil {
		contract.EmitLog("error", "Failed to set config")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 2. 设置所有者
	ownerStateID := []byte("owner")
	if _, err := framework.AppendStateOutputSimple(ownerStateID, 1, owner, nil); err != nil {
		contract.EmitLog("error", "Failed to set owner")
		return framework.ERROR_EXECUTION_FAILED
	}

	contract.EmitEvent("StakingInitialized", owner)
	return framework.SUCCESS
}

// Stake 质押代币
//
// 参数：
//   - amount: 质押金额
//
// 输出：
//   - StateOutput: stake_{staker}（质押记录）
//   - StateOutput: totalStaked（总质押量）
//   - AssetOutput: 锁定的代币（带HeightLock）
//   - Event: Staked
//
//export Stake
func Stake() uint32 {
	contract := &StakingContract{}

	staker := []byte(contract.GetCaller())
	amount := uint64(1000) // TODO: 从参数获取

	// 1. 检查最小质押金额
	configData := contract.GetState("config")
	minAmount, _ := decodeConfig(configData)
	if amount < minAmount {
		contract.EmitLog("error", "Amount below minimum")
		return framework.ERROR_INVALID_PARAMS
	}

	// 2. 记录质押信息
	stakeStateID := append([]byte("stake_"), staker...)
	currentHeight := contract.GetBlockHeight()
	unlockHeight := currentHeight + 100 // 锁定100个区块

	stakeData := encodeStake(amount, currentHeight, unlockHeight)
	if _, err := framework.AppendStateOutputSimple(stakeStateID, 1, stakeData, nil); err != nil {
		contract.EmitLog("error", "Failed to record stake")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 更新总质押量
	totalStakedData := contract.GetState("totalStaked")
	totalStaked := bytesToUint64(totalStakedData)
	newTotalStaked := totalStaked + amount

	totalStakedStateID := []byte("totalStaked")
	if _, err := framework.AppendStateOutputSimple(totalStakedStateID, 1, uint64ToBytes(newTotalStaked), nil); err != nil {
		contract.EmitLog("error", "Failed to update totalStaked")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 创建锁定输出（TODO: v1.2 使用 HeightLock）
	// lockingConditions := buildHeightLock(unlockHeight, staker)
	// framework.CreateAssetOutputWithLock(contract.GetContractAddress(), amount, nil, lockingConditions)

	contract.EmitEvent("Staked", append(staker, uint64ToBytes(amount)...))
	return framework.SUCCESS
}

// Unstake 解除质押
//
// 参数：无（自动读取调用者的质押记录）
//
// 输出：
//   - StateOutput: stake_{staker}（更新质押记录）
//   - StateOutput: totalStaked（更新总质押量）
//   - AssetOutput: 解锁的代币
//   - Event: Unstaked
//
//export Unstake
func Unstake() uint32 {
	contract := &StakingContract{}

	staker := []byte(contract.GetCaller())
	currentHeight := contract.GetBlockHeight()

	// 1. 读取质押记录
	stakeStateID := append([]byte("stake_"), staker...)
	stakeData := contract.GetState(string(stakeStateID))
	if len(stakeData) == 0 {
		contract.EmitLog("error", "No stake found")
		return framework.ERROR_NOT_FOUND
	}

	amount, _, unlockHeight := decodeStake(stakeData)

	// 2. 检查是否达到解锁高度
	if currentHeight < unlockHeight {
		contract.EmitLog("error", "Stake still locked")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 删除质押记录（设置为0）
	if _, err := framework.AppendStateOutputSimple(stakeStateID, 2, []byte{}, nil); err != nil {
		contract.EmitLog("error", "Failed to clear stake")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 更新总质押量
	totalStakedData := contract.GetState("totalStaked")
	totalStaked := bytesToUint64(totalStakedData)
	newTotalStaked := totalStaked - amount

	totalStakedStateID := []byte("totalStaked")
	if _, err := framework.AppendStateOutputSimple(totalStakedStateID, 1, uint64ToBytes(newTotalStaked), nil); err != nil {
		contract.EmitLog("error", "Failed to update totalStaked")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 5. 创建解锁输出
	// 注意：实际应用中应该使用 helpers/staking 模块的 Unstake 函数
	// 这里展示底层实现方式，使用 TransactionBuilder 创建资产输出
	// 
	// 推荐做法（使用 helpers/staking）：
	//   validatorAddr := framework.Address{} // TODO: 从状态获取验证者地址
	//   stakerAddr := framework.AddressFromBytes(staker)
	//   if err := staking.Unstake(stakerAddr, validatorAddr, framework.TokenID(""), framework.Amount(amount)); err != nil {
	//       contract.EmitLog("error", "Failed to unstake")
	//       return framework.ERROR_EXECUTION_FAILED
	//   }
	//
	// 底层实现方式（仅用于演示）：
	//   使用 TransactionBuilder 创建资产输出
	stakerAddr := framework.AddressFromBytes(staker)
	success, _, errCode := framework.BeginTransaction().
		AddAssetOutput(stakerAddr, framework.TokenID(""), framework.Amount(amount)).
		Finalize()
	if !success {
		contract.EmitLog("error", "Failed to create output")
		return errCode
	}

	contract.EmitEvent("Unstaked", append(staker, uint64ToBytes(amount)...))
	return framework.SUCCESS
}

// ClaimRewards 领取奖励（TODO: v1.2实现）
//
//export ClaimRewards
func ClaimRewards() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// Delegate 委托质押（TODO: v1.2实现）
//
//export Delegate
func Delegate() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// ================================================================================================
// 辅助函数
// ================================================================================================

func encodeConfig(minAmount, lockDuration uint64) []byte {
	result := make([]byte, 16)
	copy(result[0:8], uint64ToBytes(minAmount))
	copy(result[8:16], uint64ToBytes(lockDuration))
	return result
}

func decodeConfig(data []byte) (minAmount, lockDuration uint64) {
	if len(data) < 16 {
		return 0, 0
	}
	minAmount = bytesToUint64(data[0:8])
	lockDuration = bytesToUint64(data[8:16])
	return
}

func encodeStake(amount, stakeHeight, unlockHeight uint64) []byte {
	result := make([]byte, 24)
	copy(result[0:8], uint64ToBytes(amount))
	copy(result[8:16], uint64ToBytes(stakeHeight))
	copy(result[16:24], uint64ToBytes(unlockHeight))
	return result
}

func decodeStake(data []byte) (amount, stakeHeight, unlockHeight uint64) {
	if len(data) < 24 {
		return 0, 0, 0
	}
	amount = bytesToUint64(data[0:8])
	stakeHeight = bytesToUint64(data[8:16])
	unlockHeight = bytesToUint64(data[16:24])
	return
}

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
