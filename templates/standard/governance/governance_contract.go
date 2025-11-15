//go:build tinygo || (js && wasm)

package main

import (
	"github.com/weisyn/contract-sdk-go/framework"
)

// GovernanceContract 治理合约模板
//
// 核心功能：
// - 提案创建
// - 投票
// - 提案执行
// - 权限管理
type GovernanceContract struct {
	framework.ContractBase
}

// ProposalStatus 提案状态
const (
	ProposalPending   = 0 // 待投票
	ProposalApproved  = 1 // 已通过
	ProposalRejected  = 2 // 已拒绝
	ProposalExecuted  = 3 // 已执行
	ProposalCancelled = 4 // 已取消
)

// ================================================================================================
// 核心导出方法（Host ABI v1.1规范）
// ================================================================================================

// Initialize 初始化治理合约
//
// 参数：
//   - votingPeriod: 投票周期（区块数）
//   - quorum: 法定人数（票数）
//   - threshold: 通过阈值（百分比，如60表示60%）
//
// 输出：
//   - StateOutput: 治理参数配置
//   - StateOutput: admin（管理员地址）
//   - Event: GovernanceInitialized
//
//export Initialize
func Initialize() uint32 {
	contract := &GovernanceContract{}

	admin := []byte(contract.GetCaller())
	votingPeriod := uint64(1000) // 示例：1000个区块
	quorum := uint64(100)        // 示例：至少100票
	threshold := uint64(60)      // 示例：60%通过

	// 1. 设置治理参数
	configStateID := []byte("config")
	configData := encodeGovConfig(votingPeriod, quorum, threshold)
	if _, err := framework.AppendStateOutputSimple(configStateID, 1, configData, nil); err != nil {
		contract.EmitLog("error", "Failed to set config")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 2. 设置管理员
	adminStateID := []byte("admin")
	if _, err := framework.AppendStateOutputSimple(adminStateID, 1, admin, nil); err != nil {
		contract.EmitLog("error", "Failed to set admin")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 初始化提案计数器
	proposalCountStateID := []byte("proposalCount")
	if _, err := framework.AppendStateOutputSimple(proposalCountStateID, 1, uint64ToBytes(0), nil); err != nil {
		contract.EmitLog("error", "Failed to set proposalCount")
		return framework.ERROR_EXECUTION_FAILED
	}

	contract.EmitEvent("GovernanceInitialized", admin)
	return framework.SUCCESS
}

// CreateProposal 创建提案
//
// 参数：
//   - title: 提案标题
//   - description: 提案描述
//   - actionData: 执行动作数据
//
// 输出：
//   - StateOutput: proposal_{id}（提案记录）
//   - StateOutput: proposalCount（提案计数器）
//   - Event: ProposalCreated
//
//export CreateProposal
func CreateProposal() uint32 {
	contract := &GovernanceContract{}

	proposer := []byte(contract.GetCaller())
	currentHeight := contract.GetBlockHeight()

	// TODO: 解析提案参数
	title := []byte("Example Proposal")
	description := []byte("This is an example proposal")

	// 1. 获取下一个提案ID
	proposalCountData := contract.GetState("proposalCount")
	proposalCount := bytesToUint64(proposalCountData)
	proposalID := proposalCount + 1

	// 2. 获取投票周期
	configData := contract.GetState("config")
	votingPeriod, _, _ := decodeGovConfig(configData)
	endHeight := currentHeight + votingPeriod

	// 3. 创建提案记录
	proposalData := encodeProposal(proposer, currentHeight, endHeight, ProposalPending, title, description)
	proposalStateID := append([]byte("proposal_"), uint64ToBytes(proposalID)...)
	if _, err := framework.AppendStateOutputSimple(proposalStateID, 1, proposalData, nil); err != nil {
		contract.EmitLog("error", "Failed to create proposal")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 更新提案计数器
	proposalCountStateID := []byte("proposalCount")
	if _, err := framework.AppendStateOutputSimple(proposalCountStateID, 1, uint64ToBytes(proposalID), nil); err != nil {
		contract.EmitLog("error", "Failed to update proposalCount")
		return framework.ERROR_EXECUTION_FAILED
	}

	contract.EmitEvent("ProposalCreated", append(proposer, uint64ToBytes(proposalID)...))
	return framework.SUCCESS
}

// Vote 投票
//
// 参数：
//   - proposalID: 提案ID
//   - support: 是否支持（1=支持，0=反对）
//
// 输出：
//   - StateOutput: vote_{proposalID}_{voter}（投票记录）
//   - StateOutput: voteCount_{proposalID}（票数统计）
//   - Event: Voted
//
//export Vote
func Vote() uint32 {
	contract := &GovernanceContract{}

	voter := []byte(contract.GetCaller())
	currentHeight := contract.GetBlockHeight()

	// TODO: 解析投票参数
	proposalID := uint64(1)
	support := uint64(1) // 1=支持，0=反对

	// 1. 检查提案是否存在
	proposalStateID := append([]byte("proposal_"), uint64ToBytes(proposalID)...)
	proposalData := contract.GetState(string(proposalStateID))
	if len(proposalData) == 0 {
		contract.EmitLog("error", "Proposal not found")
		return framework.ERROR_NOT_FOUND
	}

	// 2. 检查投票期是否有效
	_, _, endHeight, status, _, _ := decodeProposal(proposalData)
	if currentHeight > endHeight {
		contract.EmitLog("error", "Voting period ended")
		return framework.ERROR_EXECUTION_FAILED
	}
	if status != ProposalPending {
		contract.EmitLog("error", "Proposal not in pending status")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 检查是否已投票
	voteStateID := append(append([]byte("vote_"), uint64ToBytes(proposalID)...), voter...)
	existingVote := contract.GetState(string(voteStateID))
	if len(existingVote) > 0 {
		contract.EmitLog("error", "Already voted")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 记录投票
	voteData := encodeVote(voter, support, currentHeight)
	if _, err := framework.AppendStateOutputSimple(voteStateID, 1, voteData, nil); err != nil {
		contract.EmitLog("error", "Failed to record vote")
		return framework.ERROR_EXECUTION_FAILED
	}

	// 5. 更新票数统计
	voteCountStateID := append([]byte("voteCount_"), uint64ToBytes(proposalID)...)
	voteCountData := contract.GetState(string(voteCountStateID))
	yesVotes, noVotes := decodeVoteCount(voteCountData)

	if support == 1 {
		yesVotes++
	} else {
		noVotes++
	}

	newVoteCountData := encodeVoteCount(yesVotes, noVotes)
	if _, err := framework.AppendStateOutputSimple(voteCountStateID, 1, newVoteCountData, nil); err != nil {
		contract.EmitLog("error", "Failed to update vote count")
		return framework.ERROR_EXECUTION_FAILED
	}

	contract.EmitEvent("Voted", append(voter, uint64ToBytes(proposalID)...))
	return framework.SUCCESS
}

// ExecuteProposal 执行提案（TODO: v1.2实现）
//
//export ExecuteProposal
func ExecuteProposal() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// CancelProposal 取消提案（TODO: v1.2实现）
//
//export CancelProposal
func CancelProposal() uint32 {
	return framework.ERROR_NOT_IMPLEMENTED
}

// ================================================================================================
// 辅助函数
// ================================================================================================

func encodeGovConfig(votingPeriod, quorum, threshold uint64) []byte {
	result := make([]byte, 24)
	copy(result[0:8], uint64ToBytes(votingPeriod))
	copy(result[8:16], uint64ToBytes(quorum))
	copy(result[16:24], uint64ToBytes(threshold))
	return result
}

func decodeGovConfig(data []byte) (votingPeriod, quorum, threshold uint64) {
	if len(data) < 24 {
		return 0, 0, 0
	}
	votingPeriod = bytesToUint64(data[0:8])
	quorum = bytesToUint64(data[8:16])
	threshold = bytesToUint64(data[16:24])
	return
}

func encodeProposal(proposer []byte, startHeight, endHeight, status uint64, title, description []byte) []byte {
	// 简化编码：proposer(20) + startHeight(8) + endHeight(8) + status(8) + title_len(8) + title + desc
	titleLen := uint64(len(title))
	result := make([]byte, 52+len(title)+len(description))
	copy(result[0:20], proposer)
	copy(result[20:28], uint64ToBytes(startHeight))
	copy(result[28:36], uint64ToBytes(endHeight))
	copy(result[36:44], uint64ToBytes(status))
	copy(result[44:52], uint64ToBytes(titleLen))
	copy(result[52:52+len(title)], title)
	copy(result[52+len(title):], description)
	return result
}

func decodeProposal(data []byte) (proposer []byte, startHeight, endHeight, status uint64, title, description []byte) {
	if len(data) < 52 {
		return nil, 0, 0, 0, nil, nil
	}
	proposer = data[0:20]
	startHeight = bytesToUint64(data[20:28])
	endHeight = bytesToUint64(data[28:36])
	status = bytesToUint64(data[36:44])
	titleLen := bytesToUint64(data[44:52])
	if len(data) < int(52+titleLen) {
		return proposer, startHeight, endHeight, status, nil, nil
	}
	title = data[52 : 52+titleLen]
	description = data[52+titleLen:]
	return
}

func encodeVote(voter []byte, support, voteHeight uint64) []byte {
	result := make([]byte, 36)
	copy(result[0:20], voter)
	copy(result[20:28], uint64ToBytes(support))
	copy(result[28:36], uint64ToBytes(voteHeight))
	return result
}

func encodeVoteCount(yesVotes, noVotes uint64) []byte {
	result := make([]byte, 16)
	copy(result[0:8], uint64ToBytes(yesVotes))
	copy(result[8:16], uint64ToBytes(noVotes))
	return result
}

func decodeVoteCount(data []byte) (yesVotes, noVotes uint64) {
	if len(data) < 16 {
		return 0, 0
	}
	yesVotes = bytesToUint64(data[0:8])
	noVotes = bytesToUint64(data[8:16])
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
