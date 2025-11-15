//go:build tinygo || (js && wasm)

package governance

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// Propose 合约内创建提案操作
//
// 🎯 **用途**：在合约代码中创建治理提案
//
// **参数**：
//   - proposer: 提案者地址
//   - proposalID: 提案ID（由合约生成）
//   - proposalData: 提案数据
//
// **返回**：
//   - error: 错误信息，nil表示成功
//
// **注意**：
//   - 提案状态通过StateOutput记录
//   - 权限控制和提案格式验证是业务逻辑，需要在合约代码中实现
//
// **示例**：
//
//	func Propose() uint32 {
//	    caller := framework.GetCaller()
//
//	    // 权限检查（业务逻辑）
//	    if !isAuthorizedProposer(caller) {
//	        return framework.ERROR_UNAUTHORIZED
//	    }
//
//	    proposalID := generateProposalID(caller)
//	    proposalData := []byte("proposal content")
//
//	    err := governance.Propose(
//	        caller,
//	        proposalID,
//	        proposalData,
//	    )
//	    if err != nil {
//	        return framework.ERROR_EXECUTION_FAILED
//	    }
//	    return framework.SUCCESS
//	}
func Propose(proposer framework.Address, proposalID []byte, proposalData []byte) error {
	// 1. 参数验证
	if err := validateProposeParams(proposer, proposalID, proposalData); err != nil {
		return err
	}

	// 2. 构建提案状态ID
	stateID := buildProposalStateID(proposalID)

	// 3. 计算提案状态哈希
	execHash := computeProposalHash(stateID, proposalData)

	// 4. 构建交易（使用internal包链式API）
	// 使用StateOutput记录提案状态
	success, _, errCode := framework.BeginTransaction().
		AddStateOutput(stateID, 1, execHash).
		Finalize()

	if !success {
		return framework.NewContractError(errCode, "propose failed")
	}

	// 5. 发出提案事件
	caller := framework.GetCaller()
	event := framework.NewEvent("Propose")
	event.AddAddressField("proposer", proposer)
	event.AddField("proposal_id", string(proposalID))
	event.AddField("proposal_data", string(proposalData))
	event.AddAddressField("caller", caller)
	framework.EmitEvent(event)

	return nil
}

// Vote 合约内投票操作
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

// Vote 函数已移至 vote.go，请使用 governance.Vote()
// validateVoteParams, buildVoteStateID, computeVoteHash 等辅助函数也已移至 vote.go

// validateProposeParams 验证提案参数
func validateProposeParams(proposer framework.Address, proposalID []byte, proposalData []byte) error {
	// 验证地址
	zeroAddr := framework.Address{}
	if proposer == zeroAddr {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"proposer address cannot be zero",
		)
	}

	// 验证提案ID
	if len(proposalID) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"proposalID cannot be empty",
		)
	}

	// 验证提案数据
	if len(proposalData) == 0 {
		return framework.NewContractError(
			framework.ERROR_INVALID_PARAMS,
			"proposalData cannot be empty",
		)
	}

	return nil
}

// buildProposalStateID 构建提案状态ID
func buildProposalStateID(proposalID []byte) []byte {
	stateID := "proposal:" + string(proposalID)
	return []byte(stateID)
}

// computeProposalHash 计算提案状态哈希
// 使用framework.ComputeHash计算真实哈希值
func computeProposalHash(stateID []byte, proposalData []byte) []byte {
	// 组合所有数据用于哈希计算
	data := make([]byte, 0, len(stateID)+len(proposalData))
	data = append(data, stateID...)
	data = append(data, proposalData...)

	// 使用framework提供的真实哈希函数
	hash := framework.ComputeHash(data)
	return hash.ToBytes()
}
