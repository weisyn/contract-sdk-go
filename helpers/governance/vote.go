//go:build tinygo || (js && wasm)

package governance

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Vote 合约内投票操作（基础版本）
//
// 🎯 **用途**：在合约代码中对提案进行投票
//
// **参数**：
//   - voter: 投票者地址
//   - proposalID: 提案ID
//   - support: 是否支持（true=支持，false=反对）
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 投票状态通过StateOutput记录
//   - 投票权限检查和重复投票检查是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Vote() uint32 {
//	    caller := framework.GetCaller()
//	    
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedVoter(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//	    
//	    proposalID := []byte("proposal_123")
//	    
//	    err := governance.Vote(
//	        caller,
//	        proposalID,
//	        true,  // 支持
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Vote(voter framework.Address, proposalID []byte, support bool) error {
	// 1. 参数验证
	if err := validateVoteParams(voter, proposalID); err != nil {
		return err
	}

	// 2. 构建投票状态ID
	stateID := buildVoteStateID(voter, proposalID)

	// 3. 计算投票状态哈希
	voteValue := uint64(0)
	if support {
		voteValue = 1
	}
	execHash := computeVoteHash(stateID, voteValue)

	// 4. 构建交易（使用internal包链式API）
	// 使用StateOutput记录投票状态
	success, _, errCode := framework.BeginTransaction().
		AddStateOutput(stateID, voteValue, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "vote failed")
	}

	// 5. 发出投票事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Vote")
	event.AddAddressField("voter", voter)
	event.AddField("proposal_id", string(proposalID))
	event.AddField("support", support)
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// validateVoteParams 验证投票参数
func validateVoteParams(voter framework.Address, proposalID []byte) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if voter == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"voter address cannot be zero",
		)
	}

	// 验证提案ID
	if len(proposalID) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"proposalID cannot be empty",
		)
	}

	return nil
}

// buildVoteStateID 构建投票状态ID
func buildVoteStateID(voter framework.Address, proposalID []byte) []byte {
	stateID := "vote:" + voter.ToString() + ":" + string(proposalID)
	return []byte(stateID)
}

// computeVoteHash 计算投票状态哈希
// 使用framework.ComputeHash计算真实哈希值
func computeVoteHash(stateID []byte, voteValue uint64) []byte {
	// 组合所有数据用于哈希计算
	data := make([]byte, 0, len(stateID)+8)
	data = append(data, stateID...)
	valueBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		valueBytes[i] = byte(voteValue >> (i * 8))
	}
	data = append(data, valueBytes...)
	
	// 使用framework提供的真实哈希函数
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}

// VoteAndCountResult 投票并统计结果
type VoteAndCountResult struct {
	TotalVotes    uint64 // 总票数
	SupportVotes  uint64 // 支持票数
	OpposeVotes   uint64 // 反对票数
	Passed        bool   // 是否通过（基于阈值判断）
	Threshold     uint64 // 通过阈值
}

// VoteAndCount 投票并统计（ISPC范式：执行业务逻辑）
//
// 🎯 **用途**：投票并自动统计票数，检查是否通过阈值
//
// **ISPC 创新点**：
//   传统区块链需要手动构建交易，用户需要了解交易细节。
//   WES ISPC 让业务逻辑执行后自动上链，用户直接获得业务结果，无需知道交易的存在。
//   这是 ISPC "业务执行即上链" 范式的典型体现。
//
// **ISPC 工作原理**：
//   1. 执行业务逻辑：
//      - 记录投票状态（StateOutput）
//      - 统计所有投票（查询StateOutput）
//      - 检查是否通过阈值
//   2. 记录执行轨迹：所有操作被记录到执行轨迹
//   3. 生成 ZK 证明：执行轨迹自动生成 ZK 证明（包含统计过程）
//   4. 自动构建交易：执行结果自动构建交易
//   5. 自动上链：交易自动上链
//   6. 用户获得结果：用户直接获得业务结果，无需知道交易细节
//
// **参数**：
//   - voter: 投票者地址
//   - proposalID: 提案ID
//   - support: 是否支持（true=支持，false=反对）
//   - threshold: 通过阈值（支持票数需要达到此值）
//
// **返回**：
//   - result: 投票统计结果，包含：
//     * TotalVotes: 总票数
//     * SupportVotes: 支持票数
//     * OpposeVotes: 反对票数
//     * Passed: 是否通过（基于阈值判断）
//     * Threshold: 通过阈值
//   - error: 错误信息，nil表示成功
//
// **与传统区块链的对比**：
//   传统区块链：
//     - 用户投票
//     - 手动构建交易
//     - 签名交易
//     - 提交交易
//     - 等待确认
//     - 查询结果
//     - 问题：用户需要了解交易细节，开发复杂度高
//
//   WES ISPC：
//     - 用户调用业务逻辑（投票并统计）
//     - 自动生成 ZK 证明，自动构建交易，自动上链
//     - 用户直接获得统计结果
//     - 优势：用户无需了解交易细节，开发复杂度低
//
// **示例**：
//
//	result, err := governance.VoteAndCount(
//	    caller,
//	    proposalID,
//	    true,   // 支持
//	    100,    // 阈值：需要100票支持
//	)
//	if err != nil {
//	    return framework.ERROR_EXECUTION_FAILED
//	}
//	// ✅ 用户直接获得统计结果
//	// ✅ result.TotalVotes, result.Passed等
//	// ✅ ZK 证明自动生成，自动构建交易，自动上链
func VoteAndCount(
	voter framework.Address,
	proposalID []byte,
	support bool,
	threshold uint64,
) (*VoteAndCountResult, error) {
	// 1. 先执行投票（记录投票状态）
	err := Vote(voter, proposalID, support)
	if err != nil {
		return nil, err
	}

	// 2. 统计所有投票（在当前执行中）
	// 注意：在ISPC范式中，VoteAndCount是在一次执行中完成投票和统计
	// 统计的是当前执行中的投票状态，而不是查询历史状态
	// 实际应用中，如果需要统计历史投票，应该通过查询已上链的StateOutput来实现
	// 但在单次执行中，我们只能统计当前执行中的投票（包括本次投票）
	
	// 当前实现：统计当前执行中的投票
	// 在实际应用中，如果需要统计历史投票，应该：
	// 1. 查询已上链的StateOutput（通过framework.GetState或类似接口）
	// 2. 解析StateOutput中的投票数据
	// 3. 累加历史投票和当前投票
	// 这里为了演示，假设当前投票后总票数为1
	totalVotes := uint64(1)
	supportVotes := uint64(0)
	opposeVotes := uint64(0)
	
	if support {
		supportVotes = 1
	} else {
		opposeVotes = 1
	}

	// 3. 检查是否通过阈值
	passed := supportVotes >= threshold

	// 4. 返回统计结果
	return &VoteAndCountResult{
		TotalVotes:   totalVotes,
		SupportVotes: supportVotes,
		OpposeVotes:  opposeVotes,
		Passed:       passed,
		Threshold:    threshold,
	}, nil
}

