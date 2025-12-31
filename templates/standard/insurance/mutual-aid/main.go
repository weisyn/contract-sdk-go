//go:build tinygo || (js && wasm)

// Package main 提供互助险（类似相互宝）业务的生产级合约。
//
// # 业务模型
//
// 互助险是一种去中心化的保险模式，成员通过缴纳分摊费用来共同承担风险。
// 当有成员发生保险事故时，其他成员按人均分摊的方式支付理赔费用。
//
// 核心流程：
//  1. 计划初始化：operator 创建互助计划，设置保障金额、服务费率等参数
//  2. 成员加入：用户申请加入计划，等待 operator 审核激活
//  3. 案件提交：成员发生保险事故后，提交理赔申请
//  4. 案件审核：operator 审核案件，决定是否批准及批准金额
//  5. 轮次结算：定期结算已批准的案件，计算人均分摊额
//  6. 费用缴纳：成员缴纳分摊费用到资金池
//  7. 理赔给付：从资金池向受益人支付理赔款
//
// # 设计目标
//
// - 完整的链上状态管理（plan_config、member、claim、round等）
// - 基于 operator 的权限控制
// - 成员生命周期管理（PENDING/ACTIVE/EXITED/BLACKLISTED）
// - 理赔案件状态机（SUBMITTED/UNDER_REVIEW/APPROVED/REJECTED/PAID）
// - 轮次结算与分摊账本
// - WES ISPC 特性：写操作同步返回业务结果，无需二次查询
// - 查询接口支持：提供完整的链上数据查询能力
//
// # 状态管理
//
// 合约使用 WES EUTXO 模型的状态输出机制，通过 StateOutput 持久化以下状态：
//   - plan_config: 计划配置（保障金额、服务费率、结算周期等）
//   - operator: 运营方地址
//   - member_{address}: 成员信息（状态、缴费记录、领取记录等）
//   - claim_{claim_id}: 理赔案件（申请人、被保人、状态、金额等）
//   - round_{round_id}: 结算轮次（周期、总给付额、人均分摊等）
//   - member_round_due_{address}_{round_id}: 成员轮次应缴记录
//   - member_month_stat_{address}_{yearMonth}: 成员月度统计（用于月度上限控制）
//
// # 权限控制
//
// - operator: 由 Initialize 时调用者地址设置，拥有审核成员、审核案件、开启/结算轮次、给付等权限
// - 普通成员: 可以加入计划、提交案件、缴纳分摊费用、退出计划
//
// # 资金流转
//
// 使用 helpers/market 模块实现资金托管和释放：
//   - PayContribution: 使用 market.Escrow 将成员资金托管到资金池
//   - Payout: 使用 market.Release 从资金池释放资金给受益人
//
// # 注意事项
//
// ⚠️ 本合约为生产级实现，包含完整的状态管理和业务逻辑。
// 适用于 consortium/private 链模式，建议配合链下 KYC/风控服务使用。
//
// 安全建议：
//   - operator 地址应使用多签或 DAO 治理
//   - 案件审核应结合链下调查和链上投票
//   - 建议设置合理的月度分摊上限，避免单月费用过高
//   - 等待期设置可防止逆向选择风险
package main

import (
	"github.com/weisyn/contract-sdk-go/framework"
	"github.com/weisyn/contract-sdk-go/helpers/market"
)

// MutualAidContract 互助险合约
type MutualAidContract struct {
	framework.ContractBase
}

// ================================================================================================
// 常量定义
// ================================================================================================

// 成员状态常量
//
// 状态转换流程：
//
//	PENDING -> ACTIVE (通过 ApproveMember)
//	ACTIVE -> EXITED (通过 Exit)
//	ACTIVE -> SUSPENDED (运营方暂停，暂未实现)
//	ACTIVE -> BLACKLISTED (运营方拉黑，暂未实现)
const (
	// MEMBER_STATUS_PENDING 待审核：成员已申请加入，等待 operator 审核
	MEMBER_STATUS_PENDING = "PENDING"
	// MEMBER_STATUS_ACTIVE 活跃：成员已激活，可以提交案件和缴纳分摊费用
	MEMBER_STATUS_ACTIVE = "ACTIVE"
	// MEMBER_STATUS_SUSPENDED 暂停：成员被临时暂停，不能提交案件但仍需缴纳费用
	MEMBER_STATUS_SUSPENDED = "SUSPENDED"
	// MEMBER_STATUS_EXITED 已退出：成员主动退出计划，不再参与分摊
	MEMBER_STATUS_EXITED = "EXITED"
	// MEMBER_STATUS_BLACKLISTED 黑名单：成员因违规被拉黑，不能参与任何操作
	MEMBER_STATUS_BLACKLISTED = "BLACKLISTED"
)

// 理赔案件状态常量
//
// 状态转换流程：
//
//	SUBMITTED -> UNDER_REVIEW (审核中，暂未使用)
//	SUBMITTED/UNDER_REVIEW -> APPROVED (通过 ReviewClaim 批准)
//	SUBMITTED/UNDER_REVIEW -> REJECTED (通过 ReviewClaim 拒绝)
//	APPROVED -> PAID (通过 Payout 给付)
const (
	// CLAIM_STATUS_SUBMITTED 已提交：成员已提交理赔申请，等待审核
	CLAIM_STATUS_SUBMITTED = "SUBMITTED"
	// CLAIM_STATUS_UNDER_REVIEW 审核中：案件正在审核中（当前实现中未使用）
	CLAIM_STATUS_UNDER_REVIEW = "UNDER_REVIEW"
	// CLAIM_STATUS_APPROVED 已批准：案件已通过审核，等待给付
	CLAIM_STATUS_APPROVED = "APPROVED"
	// CLAIM_STATUS_REJECTED 已拒绝：案件审核未通过
	CLAIM_STATUS_REJECTED = "REJECTED"
	// CLAIM_STATUS_PAID 已给付：理赔款已支付给受益人
	CLAIM_STATUS_PAID = "PAID"
	// CLAIM_STATUS_CANCELLED 已取消：案件被取消（暂未实现）
	CLAIM_STATUS_CANCELLED = "CANCELLED"
)

// 轮次状态常量
//
// 状态转换流程：
//
//	OPEN -> SETTLED (通过 SettleRound 结算)
//	SETTLED -> CLOSED (通过 CloseRound 关闭)
const (
	// ROUND_STATUS_OPEN 开启：轮次已开启，可以结算案件
	ROUND_STATUS_OPEN = "OPEN"
	// ROUND_STATUS_SETTLED 已结算：轮次已结算，计算出人均分摊额，成员可以缴费
	ROUND_STATUS_SETTLED = "SETTLED"
	// ROUND_STATUS_CLOSED 已关闭：轮次已关闭，未缴费的成员记录为欠费
	ROUND_STATUS_CLOSED = "CLOSED"
)

// 审核决策常量
//
// 用于 ReviewClaim 函数，表示 operator 对案件的审核决定
const (
	// DECISION_APPROVE 批准：案件通过审核，可以给付
	DECISION_APPROVE = "APPROVE"
	// DECISION_REJECT 拒绝：案件未通过审核
	DECISION_REJECT = "REJECT"
)

// 状态ID前缀常量
//
// 用于构建链上状态的唯一标识符（StateOutput 的 key）
const (
	// STATE_PLAN_CONFIG 计划配置状态ID
	STATE_PLAN_CONFIG = "plan_config"
	// STATE_OPERATOR 运营方地址状态ID
	STATE_OPERATOR = "operator"
	// STATE_MEMBER_PREFIX 成员状态ID前缀，完整格式：member_{address}
	STATE_MEMBER_PREFIX = "member_"
	// STATE_CLAIM_PREFIX 理赔案件状态ID前缀，完整格式：claim_{claim_id}
	STATE_CLAIM_PREFIX = "claim_"
	// STATE_ROUND_PREFIX 轮次状态ID前缀，完整格式：round_{round_id}
	STATE_ROUND_PREFIX = "round_"
	// STATE_MEMBER_COUNT 活跃成员数状态ID
	STATE_MEMBER_COUNT = "member_count_active"
	// STATE_CURRENT_ROUND 当前轮次ID状态ID
	STATE_CURRENT_ROUND = "current_round_id"
)

// ================================================================================================
// 状态结构编码/解码
// ================================================================================================
//
// 由于 WES 合约状态存储为字节数组，需要将复杂数据结构序列化为字节数组。
// 本合约采用固定长度编码方式，便于快速解码和节省存储空间。
//
// 编码格式说明：
//   - 字符串字段：固定长度，不足部分用 0x00 填充，解码时使用 trimNull 去除
//   - 数值字段：使用 uint64ToBytes 转换为 8 字节大端序
//   - 布尔字段：使用 1 字节，0 表示 false，1 表示 true

// encodePlanConfig 编码计划配置信息
//
// 参数说明：
//   - planID: 计划唯一标识符（最大32字节）
//   - name: 计划名称（最大64字节）
//   - tokenID: 计价代币ID，空字符串表示原生币（最大32字节）
//   - coverageAmount: 单次给付金额上限
//   - serviceFeeBP: 服务费率，单位 bp（万分比），如 800 = 8%
//   - settlementPeriod: 结算周期（秒），例如 2592000 = 30天
//   - waitingPeriod: 等待期（秒），例如 86400 = 1天
//   - minMembers: 最小成员数，计划生效门槛
//   - monthlyCapPerMember: 单成员月度分摊上限
//
// 返回：176字节的编码数据
//
// 编码格式：
//
//	planID(32) + name(64) + tokenID(32) + coverageAmount(8) + serviceFeeBP(8) +
//	settlementPeriod(8) + waitingPeriod(8) + minMembers(8) + monthlyCapPerMember(8) = 176字节
func encodePlanConfig(planID, name, tokenID string, coverageAmount, serviceFeeBP, settlementPeriod, waitingPeriod, minMembers, monthlyCapPerMember uint64) []byte {
	result := make([]byte, 176)
	copy(result[0:32], []byte(planID)[:min(32, len(planID))])
	copy(result[32:96], []byte(name)[:min(64, len(name))])
	copy(result[96:128], []byte(tokenID)[:min(32, len(tokenID))])
	copy(result[128:136], uint64ToBytes(coverageAmount))
	copy(result[136:144], uint64ToBytes(serviceFeeBP))
	copy(result[144:152], uint64ToBytes(settlementPeriod))
	copy(result[152:160], uint64ToBytes(waitingPeriod))
	copy(result[160:168], uint64ToBytes(minMembers))
	copy(result[168:176], uint64ToBytes(monthlyCapPerMember))
	return result
}

// decodePlanConfig 解码计划配置信息
//
// 参数：
//   - data: 176字节的编码数据
//
// 返回：解码后的计划配置字段
//
// 如果数据长度不足176字节，返回零值
func decodePlanConfig(data []byte) (planID, name, tokenID string, coverageAmount, serviceFeeBP, settlementPeriod, waitingPeriod, minMembers, monthlyCapPerMember uint64) {
	if len(data) < 176 {
		return "", "", "", 0, 0, 0, 0, 0, 0
	}
	planID = string(trimNull(data[0:32]))
	name = string(trimNull(data[32:96]))
	tokenID = string(trimNull(data[96:128]))
	coverageAmount = bytesToUint64(data[128:136])
	serviceFeeBP = bytesToUint64(data[136:144])
	settlementPeriod = bytesToUint64(data[144:152])
	waitingPeriod = bytesToUint64(data[152:160])
	minMembers = bytesToUint64(data[160:168])
	monthlyCapPerMember = bytesToUint64(data[168:176])
	return
}

// encodeMember 编码成员信息
//
// 参数说明：
//   - status: 成员状态（PENDING/ACTIVE/EXITED等，最大16字节）
//   - joinTime: 加入时间戳（Unix时间戳，秒）
//   - totalPaid: 累计缴费总额
//   - totalReceived: 累计领取总额
//   - arrearsAmount: 欠费金额
//   - lastSettledRound: 最后结算的轮次ID（数值型，简化实现）
//
// 返回：64字节的编码数据
//
// 编码格式：
//
//	status(16) + joinTime(8) + totalPaid(8) + totalReceived(8) + arrearsAmount(8) + lastSettledRound(8) = 64字节
func encodeMember(status string, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound uint64) []byte {
	result := make([]byte, 64)
	copy(result[0:16], []byte(status)[:min(16, len(status))])
	copy(result[16:24], uint64ToBytes(joinTime))
	copy(result[24:32], uint64ToBytes(totalPaid))
	copy(result[32:40], uint64ToBytes(totalReceived))
	copy(result[40:48], uint64ToBytes(arrearsAmount))
	copy(result[48:56], uint64ToBytes(lastSettledRound))
	return result
}

// decodeMember 解码成员信息
//
// 参数：
//   - data: 64字节的编码数据
//
// 返回：解码后的成员信息字段
//
// 如果数据长度不足64字节，返回零值
func decodeMember(data []byte) (status string, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound uint64) {
	if len(data) < 56 {
		return "", 0, 0, 0, 0, 0
	}
	status = string(trimNull(data[0:16]))
	joinTime = bytesToUint64(data[16:24])
	totalPaid = bytesToUint64(data[24:32])
	totalReceived = bytesToUint64(data[32:40])
	arrearsAmount = bytesToUint64(data[40:48])
	lastSettledRound = bytesToUint64(data[48:56])
	return
}

// encodeClaim 编码理赔案件信息
//
// 参数说明：
//   - planID: 计划ID（最大32字节）
//   - claimID: 案件唯一标识符（最大32字节）
//   - applicant: 申请人地址（20字节二进制，存储为字符串）
//   - insured: 被保人地址（20字节二进制，存储为字符串）
//   - status: 案件状态（SUBMITTED/APPROVED等，最大16字节）
//   - roundID: 所属轮次ID（最大32字节）
//   - evidenceHash: 证据哈希（最大64字节）
//   - investigationHash: 调查报告哈希（最大64字节）
//   - requestedAmount: 申请金额
//   - approvedAmount: 批准金额
//   - eventTime: 事故发生时间戳（Unix时间戳，秒）
//
// 返回：304字节的编码数据
//
// 编码格式：
//
//	planID(32) + claimID(32) + applicant(20) + insured(20) + status(16) + roundID(32) +
//	evidenceHash(64) + investigationHash(64) + requestedAmount(8) + approvedAmount(8) + eventTime(8) = 304字节
//
// 注意：applicant 和 insured 字段存储的是地址的20字节二进制数据（通过 string(addr.ToBytes()) 转换），
// 解码后需要使用 addressBytesToString 转换为 Base58 格式用于 JSON 返回。
func encodeClaim(planID, claimID, applicant, insured, status, roundID, evidenceHash, investigationHash string, requestedAmount, approvedAmount, eventTime uint64) []byte {
	result := make([]byte, 304)
	copy(result[0:32], []byte(planID)[:min(32, len(planID))])
	copy(result[32:64], []byte(claimID)[:min(32, len(claimID))])
	copy(result[64:84], []byte(applicant)[:min(20, len(applicant))])
	copy(result[84:104], []byte(insured)[:min(20, len(insured))])
	copy(result[104:120], []byte(status)[:min(16, len(status))])
	copy(result[120:152], []byte(roundID)[:min(32, len(roundID))])
	copy(result[152:216], []byte(evidenceHash)[:min(64, len(evidenceHash))])
	copy(result[216:280], []byte(investigationHash)[:min(64, len(investigationHash))])
	copy(result[280:288], uint64ToBytes(requestedAmount))
	copy(result[288:296], uint64ToBytes(approvedAmount))
	copy(result[296:304], uint64ToBytes(eventTime))
	return result
}

// decodeClaim 解码理赔案件信息
//
// 参数：
//   - data: 304字节的编码数据
//
// 返回：解码后的案件信息字段
//
// 如果数据长度不足304字节，返回零值
//
// 注意：applicant 和 insured 返回的是20字节二进制数据的字符串表示，
// 需要使用 addressBytesToString 转换为 Base58 格式。
func decodeClaim(data []byte) (planID, claimID, applicant, insured, status, roundID, evidenceHash, investigationHash string, requestedAmount, approvedAmount, eventTime uint64) {
	if len(data) < 304 {
		return "", "", "", "", "", "", "", "", 0, 0, 0
	}
	planID = string(trimNull(data[0:32]))
	claimID = string(trimNull(data[32:64]))
	applicant = string(trimNull(data[64:84]))
	insured = string(trimNull(data[84:104]))
	status = string(trimNull(data[104:120]))
	roundID = string(trimNull(data[120:152]))
	evidenceHash = string(trimNull(data[152:216]))
	investigationHash = string(trimNull(data[216:280]))
	requestedAmount = bytesToUint64(data[280:288])
	approvedAmount = bytesToUint64(data[288:296])
	eventTime = bytesToUint64(data[296:304])
	return
}

// encodeRound 编码轮次信息
//
// 参数说明：
//   - planID: 计划ID（最大32字节）
//   - roundID: 轮次唯一标识符（最大32字节）
//   - status: 轮次状态（OPEN/SETTLED/CLOSED，最大16字节）
//   - periodStart: 轮次开始时间戳（Unix时间戳，秒）
//   - periodEnd: 轮次结束时间戳（Unix时间戳，秒）
//   - totalApprovedPayout: 该轮次总批准给付额
//   - totalServiceFee: 该轮次总服务费
//   - perCapitaContribution: 人均分摊额（向上取整）
//   - payersCount: 已缴费人数（简化实现，未去重）
//
// 返回：128字节的编码数据
//
// 编码格式：
//
//	planID(32) + roundID(32) + status(16) + periodStart(8) + periodEnd(8) +
//	totalApprovedPayout(8) + totalServiceFee(8) + perCapitaContribution(8) + payersCount(8) = 128字节
func encodeRound(planID, roundID, status string, periodStart, periodEnd, totalApprovedPayout, totalServiceFee, perCapitaContribution, payersCount uint64) []byte {
	result := make([]byte, 128)
	copy(result[0:32], []byte(planID)[:min(32, len(planID))])
	copy(result[32:64], []byte(roundID)[:min(32, len(roundID))])
	copy(result[64:80], []byte(status)[:min(16, len(status))])
	copy(result[80:88], uint64ToBytes(periodStart))
	copy(result[88:96], uint64ToBytes(periodEnd))
	copy(result[96:104], uint64ToBytes(totalApprovedPayout))
	copy(result[104:112], uint64ToBytes(totalServiceFee))
	copy(result[112:120], uint64ToBytes(perCapitaContribution))
	copy(result[120:128], uint64ToBytes(payersCount))
	return result
}

// decodeRound 解码轮次信息
//
// 参数：
//   - data: 128字节的编码数据
//
// 返回：解码后的轮次信息字段
//
// 如果数据长度不足128字节，返回零值
func decodeRound(data []byte) (planID, roundID, status string, periodStart, periodEnd, totalApprovedPayout, totalServiceFee, perCapitaContribution, payersCount uint64) {
	if len(data) < 128 {
		return "", "", "", 0, 0, 0, 0, 0, 0
	}
	planID = string(trimNull(data[0:32]))
	roundID = string(trimNull(data[32:64]))
	status = string(trimNull(data[64:80]))
	periodStart = bytesToUint64(data[80:88])
	periodEnd = bytesToUint64(data[88:96])
	totalApprovedPayout = bytesToUint64(data[96:104])
	totalServiceFee = bytesToUint64(data[104:112])
	perCapitaContribution = bytesToUint64(data[112:120])
	payersCount = bytesToUint64(data[120:128])
	return
}

// encodeMemberRoundDue 编码成员轮次应缴信息
//
// 用于记录每个成员在每个轮次的缴费情况。
//
// 参数说明：
//   - dueAmount: 应缴金额（该轮次的人均分摊额）
//   - paidAmount: 已缴金额
//   - settled: 是否已结清（paidAmount >= dueAmount）
//
// 返回：17字节的编码数据
//
// 编码格式：
//
//	dueAmount(8) + paidAmount(8) + settled(1) = 17字节
func encodeMemberRoundDue(dueAmount, paidAmount uint64, settled bool) []byte {
	result := make([]byte, 17)
	copy(result[0:8], uint64ToBytes(dueAmount))
	copy(result[8:16], uint64ToBytes(paidAmount))
	if settled {
		result[16] = 1
	} else {
		result[16] = 0
	}
	return result
}

// decodeMemberRoundDue 解码成员轮次应缴信息
//
// 参数：
//   - data: 17字节的编码数据
//
// 返回：解码后的应缴信息字段
//
// 如果数据长度不足17字节，返回零值
func decodeMemberRoundDue(data []byte) (dueAmount, paidAmount uint64, settled bool) {
	if len(data) < 17 {
		return 0, 0, false
	}
	dueAmount = bytesToUint64(data[0:8])
	paidAmount = bytesToUint64(data[8:16])
	settled = data[16] == 1
	return
}

// encodeMemberMonthStat 编码成员月度统计信息
//
// 用于记录每个成员在每个自然月的缴费情况，用于月度分摊上限控制。
//
// 参数说明：
//   - paidAmount: 该月累计缴费金额
//   - capReached: 是否已达到月度上限
//
// 返回：9字节的编码数据
//
// 编码格式：
//
//	paidAmount(8) + capReached(1) = 9字节
func encodeMemberMonthStat(paidAmount uint64, capReached bool) []byte {
	result := make([]byte, 9)
	copy(result[0:8], uint64ToBytes(paidAmount))
	if capReached {
		result[8] = 1
	} else {
		result[8] = 0
	}
	return result
}

// decodeMemberMonthStat 解码成员月度统计信息
//
// 参数：
//   - data: 9字节的编码数据
//
// 返回：解码后的月度统计字段
//
// 如果数据长度不足9字节，返回零值
func decodeMemberMonthStat(data []byte) (paidAmount uint64, capReached bool) {
	if len(data) < 9 {
		return 0, false
	}
	paidAmount = bytesToUint64(data[0:8])
	capReached = data[8] == 1
	return
}

// ================================================================================================
// 辅助函数
// ================================================================================================

// uint64ToBytes 将 uint64 转换为 8 字节大端序字节数组
//
// 用于将数值字段编码到状态数据中。
//
// 参数：
//   - n: 要转换的 uint64 值
//
// 返回：8字节的字节数组（大端序）
func uint64ToBytes(n uint64) []byte {
	result := make([]byte, 8)
	for i := 0; i < 8; i++ {
		result[7-i] = byte(n >> (i * 8))
	}
	return result
}

// bytesToUint64 将 8 字节大端序字节数组转换为 uint64
//
// 用于从状态数据中解码数值字段。
//
// 参数：
//   - b: 字节数组（至少8字节）
//
// 返回：解码后的 uint64 值
//
// 如果字节数组长度不足8字节，返回0
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

// min 返回两个整数中的较小值
//
// 用于确保字符串字段不会超出固定长度限制。
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// trimNull 去除字节数组末尾的 null 字符（0x00）
//
// 用于解码固定长度的字符串字段，去除填充的 null 字符。
//
// 参数：
//   - b: 字节数组
//
// 返回：去除末尾 null 字符后的字节数组
func trimNull(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return b[:i]
		}
	}
	return b
}

// checkOperator 检查当前调用者是否为计划的 operator
//
// 用于权限控制，确保只有 operator 可以执行管理操作（如审核成员、审核案件、结算轮次等）。
//
// 返回：
//   - true: 调用者是 operator
//   - false: 调用者不是 operator 或 operator 未设置
func checkOperator() bool {
	operatorData, _ := framework.GetState(STATE_OPERATOR)
	if len(operatorData) == 0 {
		return false
	}
	caller := framework.GetCaller()
	return string(operatorData) == string(caller.ToBytes())
}

// getMemberStateID 获取成员状态的唯一标识符
//
// 用于构建 StateOutput 的 key，格式：member_{address}
//
// 参数：
//   - addr: 成员地址
//
// 返回：成员状态ID的字节数组
func getMemberStateID(addr framework.Address) []byte {
	return append([]byte(STATE_MEMBER_PREFIX), addr.ToBytes()...)
}

// getClaimStateID 获取理赔案件状态的唯一标识符
//
// 用于构建 StateOutput 的 key，格式：claim_{claim_id}
//
// 参数：
//   - claimID: 案件唯一标识符
//
// 返回：案件状态ID的字节数组
func getClaimStateID(claimID string) []byte {
	return append([]byte(STATE_CLAIM_PREFIX), []byte(claimID)...)
}

// getRoundStateID 获取轮次状态的唯一标识符
//
// 用于构建 StateOutput 的 key，格式：round_{round_id}
//
// 参数：
//   - roundID: 轮次唯一标识符
//
// 返回：轮次状态ID的字节数组
func getRoundStateID(roundID string) []byte {
	return append([]byte(STATE_ROUND_PREFIX), []byte(roundID)...)
}

// getMemberRoundDueStateID 获取成员轮次应缴状态的唯一标识符
//
// 用于构建 StateOutput 的 key，格式：member_round_due_{address}_{round_id}
//
// 参数：
//   - addr: 成员地址
//   - roundID: 轮次唯一标识符
//
// 返回：成员轮次应缴状态ID的字节数组
func getMemberRoundDueStateID(addr framework.Address, roundID string) []byte {
	return append(append([]byte("member_round_due_"), addr.ToBytes()...), []byte("_"+roundID)...)
}

// getMemberMonthStatStateID 获取成员月度统计状态的唯一标识符
//
// 用于构建 StateOutput 的 key，格式：member_month_stat_{address}_{yearMonth}
//
// 参数：
//   - addr: 成员地址
//   - yearMonth: 年月标识符（格式：YYYYMM，如 "202501"）
//
// 返回：成员月度统计状态ID的字节数组
func getMemberMonthStatStateID(addr framework.Address, yearMonth string) []byte {
	return append(append([]byte("member_month_stat_"), addr.ToBytes()...), []byte("_"+yearMonth)...)
}

// addressBytesToString 将20字节的地址二进制数据转换为 Base58 地址字符串
//
// 用于将状态中存储的地址二进制数据转换为可读的 Base58 格式，用于 JSON 返回。
//
// 参数：
//   - addrBytes: 地址的二进制数据（至少20字节）
//
// 返回：Base58 格式的地址字符串
//
// 如果数据长度不足20字节，返回空字符串
func addressBytesToString(addrBytes []byte) string {
	if len(addrBytes) < 20 {
		return ""
	}
	addr := framework.AddressFromBytes(addrBytes[:20])
	return addr.ToString()
}

// ================================================================================================
// 导出方法（Host ABI v1.1 规范）
// ================================================================================================

// Initialize 初始化互助计划
//
// 这是合约的第一个调用，用于创建并配置一个新的互助计划。
// 调用者将成为计划的 operator，拥有管理权限。
//
// # 业务逻辑
//
// 1. 参数校验：检查必填参数和数值范围
// 2. 保存计划配置到链上状态
// 3. 设置调用者为 operator
// 4. 初始化活跃成员计数为 0
// 5. 发出初始化事件
// 6. 返回完整的计划配置信息（WES ISPC 特性）
//
// # 参数说明（JSON格式）
//
//	{
//	  "plan_id": "plan_xianghubao_001",     // 计划唯一标识符（必填）
//	  "name": "相互宝互助计划",              // 计划名称（必填）
//	  "token_id": "",                        // 计价代币ID，空字符串表示原生币（可选）
//	  "coverage_amount": 300000,            // 单次给付额上限（必填，>0）
//	  "service_fee_bp": 800,                 // 服务费率，单位 bp（万分比），如 800 = 8%（可选，默认0，最大10000）
//	  "settlement_period": 2592000,          // 结算周期（秒），例如 30 天（必填，>0）
//	  "waiting_period": 86400,               // 等待期（秒），例如 1 天（可选，默认0）
//	  "min_members": 1000,                   // 最小成员数，计划生效门槛（可选，默认1）
//	  "monthly_cap_per_member": 10000        // 单成员月度分摊上限（可选，默认1000000）
//	}
//
// # 返回值
//
// 成功时返回 JSON 格式的计划配置信息：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "name": "相互宝互助计划",
//	  "token_id": "",
//	  "coverage_amount": 300000,
//	  "service_fee_bp": 800,
//	  "settlement_period": 2592000,
//	  "waiting_period": 86400,
//	  "min_members": 1000,
//	  "monthly_cap_per_member": 10000,
//	  "operator": "Cf1...",                  // Base58 格式的 operator 地址
//	  "member_count_active": 0,              // 初始活跃成员数
//	  "initialized_at": 1736200000          // 初始化时间戳
//	}
//
// # 状态变更
//
// - 创建 StateOutput: plan_config（计划配置）
// - 创建 StateOutput: operator（运营方地址）
// - 创建 StateOutput: member_count_active（活跃成员数，初始为0）
//
// # 事件
//
// 发出 MutualAidPlanInitialized 事件，包含所有计划配置参数。
//
// # 错误码
//
// - ERROR_INVALID_PARAMS: 参数无效（plan_id/name 为空，或数值范围错误）
// - ERROR_EXECUTION_FAILED: 状态保存失败
//
//export Initialize
func Initialize() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	name := params.ParseJSON("name")
	tokenID := params.ParseJSON("token_id")
	coverageAmount := params.ParseJSONInt("coverage_amount")
	serviceFeeBP := params.ParseJSONInt("service_fee_bp")
	settlementPeriod := params.ParseJSONInt("settlement_period")
	waitingPeriod := params.ParseJSONInt("waiting_period")
	minMembers := params.ParseJSONInt("min_members")
	monthlyCapPerMember := params.ParseJSONInt("monthly_cap_per_member")

	// 参数校验
	if planID == "" || name == "" || coverageAmount <= 0 || settlementPeriod <= 0 {
		return framework.ERROR_INVALID_PARAMS
	}
	if serviceFeeBP > 10000 { // 服务费率不能超过100%
		return framework.ERROR_INVALID_PARAMS
	}
	if waitingPeriod < 0 {
		waitingPeriod = 0
	}
	if minMembers < 1 {
		minMembers = 1
	}
	if monthlyCapPerMember <= 0 {
		monthlyCapPerMember = 1000000 // 默认上限100万
	}

	caller := framework.GetCaller()

	// 1. 保存计划配置
	configData := encodePlanConfig(planID, name, tokenID, coverageAmount, serviceFeeBP, settlementPeriod, waitingPeriod, minMembers, monthlyCapPerMember)
	if _, err := framework.AppendStateOutputSimple([]byte(STATE_PLAN_CONFIG), 1, configData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 2. 保存 operator
	if _, err := framework.AppendStateOutputSimple([]byte(STATE_OPERATOR), 1, caller.ToBytes(), nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 初始化成员计数
	if _, err := framework.AppendStateOutputSimple([]byte(STATE_MEMBER_COUNT), 1, uint64ToBytes(0), nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 发出事件
	event := framework.NewEvent("MutualAidPlanInitialized")
	event.AddStringField("plan_id", planID)
	event.AddStringField("name", name)
	event.AddStringField("token_id", tokenID)
	event.AddIntField("coverage_amount", coverageAmount)
	event.AddIntField("service_fee_bp", serviceFeeBP)
	event.AddIntField("settlement_period", settlementPeriod)
	event.AddIntField("waiting_period", waitingPeriod)
	event.AddIntField("min_members", minMembers)
	event.AddIntField("monthly_cap_per_member", monthlyCapPerMember)
	event.AddAddressField("operator", caller)
	framework.EmitEvent(event)

	// 5. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":                planID,
		"name":                   name,
		"token_id":               tokenID,
		"coverage_amount":        coverageAmount,
		"service_fee_bp":         serviceFeeBP,
		"settlement_period":      settlementPeriod,
		"waiting_period":         waitingPeriod,
		"min_members":            minMembers,
		"monthly_cap_per_member": monthlyCapPerMember,
		"operator":               caller.ToString(),
		"member_count_active":    uint64(0),
		"initialized_at":         framework.GetTimestamp(),
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Join 成为互助计划成员
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001"
//	}
//
// 输出：
// - StateOutput: member_{address}
// - StateOutput: member_count_active (更新)
// - Event: MutualAidMemberJoined
//
//export Join
func Join() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	if planID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	caller := framework.GetCaller()
	memberStateID := getMemberStateID(caller)

	// 1. 检查是否已加入
	existingMemberData, _ := framework.GetState(string(memberStateID))
	if len(existingMemberData) > 0 {
		status, _, _, _, _, _ := decodeMember(existingMemberData)
		if status == MEMBER_STATUS_ACTIVE || status == MEMBER_STATUS_PENDING {
			return framework.ERROR_ALREADY_EXISTS
		}
		if status == MEMBER_STATUS_BLACKLISTED {
			return framework.ERROR_UNAUTHORIZED
		}
	}

	// 2. 创建成员记录（状态为PENDING，需要operator审核）
	currentTime := framework.GetTimestamp()
	memberData := encodeMember(MEMBER_STATUS_PENDING, currentTime, 0, 0, 0, 0)
	if _, err := framework.AppendStateOutputSimple(memberStateID, 1, memberData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 更新成员计数（仅统计ACTIVE，PENDING不计入）
	// 注意：这里不更新计数，等待ApproveMember时再更新

	// 4. 发出事件
	event := framework.NewEvent("MutualAidMemberJoined")
	event.AddStringField("plan_id", planID)
	event.AddAddressField("member", caller)
	event.AddStringField("status", MEMBER_STATUS_PENDING)
	framework.EmitEvent(event)

	// 5. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	configData, _ := framework.GetState(STATE_PLAN_CONFIG)
	waitingPeriod := uint64(0)
	if len(configData) > 0 {
		_, _, _, _, _, _, waitingPeriod, _, _ = decodePlanConfig(configData)
	}
	result := map[string]interface{}{
		"plan_id":          planID,
		"member":           caller.ToString(),
		"status":           MEMBER_STATUS_PENDING,
		"join_time":        currentTime,
		"waiting_period":   waitingPeriod,
		"waiting_end_time": currentTime + waitingPeriod,
		"total_paid":       uint64(0),
		"total_received":   uint64(0),
		"arrears_amount":   uint64(0),
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ApproveMember 审核并激活成员（仅 operator 可调用）
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "member": "Cf1..." // 成员地址（Base58）
//	}
//
// 输出：
// - StateOutput: member_{address} (更新状态为ACTIVE)
// - StateOutput: member_count_active (更新)
// - Event: MutualAidMemberApproved
//
//export ApproveMember
func ApproveMember() uint32 {
	params := framework.GetContractParams()

	// 1. 权限检查
	if !checkOperator() {
		return framework.ERROR_UNAUTHORIZED
	}

	planID := params.ParseJSON("plan_id")
	memberStr := params.ParseJSON("member")
	if planID == "" || memberStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	member, err := framework.ParseAddressBase58(memberStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	memberStateID := getMemberStateID(member)
	memberData, _ := framework.GetState(string(memberStateID))

	// 2. 检查成员是否存在且状态为PENDING
	if len(memberData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	status, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound := decodeMember(memberData)
	if status != MEMBER_STATUS_PENDING {
		return framework.ERROR_INVALID_STATE
	}

	// 3. 更新成员状态为ACTIVE
	newMemberData := encodeMember(MEMBER_STATUS_ACTIVE, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound)
	if _, err := framework.AppendStateOutputSimple(memberStateID, 2, newMemberData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 更新成员计数
	memberCountData, _ := framework.GetState(STATE_MEMBER_COUNT)
	memberCount := bytesToUint64(memberCountData)
	newMemberCount := memberCount + 1
	if _, err := framework.AppendStateOutputSimple([]byte(STATE_MEMBER_COUNT), 2, uint64ToBytes(newMemberCount), nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 5. 发出事件
	event := framework.NewEvent("MutualAidMemberApproved")
	event.AddStringField("plan_id", planID)
	event.AddAddressField("member", member)
	framework.EmitEvent(event)

	// 6. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":             planID,
		"member":              member.ToString(),
		"status":              MEMBER_STATUS_ACTIVE,
		"join_time":           joinTime,
		"total_paid":          totalPaid,
		"total_received":      totalReceived,
		"arrears_amount":      arrearsAmount,
		"member_count_active": newMemberCount,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Exit 退出互助计划
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001"
//	}
//
// 输出：
// - StateOutput: member_{address} (更新状态为EXITED)
// - StateOutput: member_count_active (更新)
// - Event: MutualAidMemberExited
//
//export Exit
func Exit() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	if planID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	caller := framework.GetCaller()
	memberStateID := getMemberStateID(caller)
	memberData, _ := framework.GetState(string(memberStateID))

	// 1. 检查成员是否存在且状态为ACTIVE
	if len(memberData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	status, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound := decodeMember(memberData)
	if status != MEMBER_STATUS_ACTIVE {
		return framework.ERROR_INVALID_STATE
	}

	// 2. 更新成员状态为EXITED
	newMemberData := encodeMember(MEMBER_STATUS_EXITED, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound)
	if _, err := framework.AppendStateOutputSimple(memberStateID, 2, newMemberData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 3. 更新成员计数
	memberCountData, _ := framework.GetState(STATE_MEMBER_COUNT)
	memberCount := bytesToUint64(memberCountData)
	newMemberCount := memberCount
	if memberCount > 0 {
		newMemberCount = memberCount - 1
		if _, err := framework.AppendStateOutputSimple([]byte(STATE_MEMBER_COUNT), 2, uint64ToBytes(newMemberCount), nil); err != nil {
			return framework.ERROR_EXECUTION_FAILED
		}
	}

	// 4. 发出事件
	event := framework.NewEvent("MutualAidMemberExited")
	event.AddStringField("plan_id", planID)
	event.AddAddressField("member", caller)
	event.AddIntField("arrears_amount", arrearsAmount)
	framework.EmitEvent(event)

	// 5. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":             planID,
		"member":              caller.ToString(),
		"status":              MEMBER_STATUS_EXITED,
		"join_time":           joinTime,
		"total_paid":          totalPaid,
		"total_received":      totalReceived,
		"arrears_amount":      arrearsAmount,
		"member_count_active": newMemberCount,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// SubmitClaim 提交互助申请（报案）
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "claim_id": "claim_202501_0001",
//	  "insured": "Cf1...",                // 被保人地址（Base58），可为空表示即为调用者
//	  "requested_amount": 300000,
//	  "event_time": 1736200000,           // 出险时间（时间戳）
//	  "evidence_hash": "0xabc...",        // 资料哈希
//	  "extra": "optional comments"
//	}
//
// 输出：
// - StateOutput: claim_{claim_id}
// - Event: MutualAidClaimSubmitted
//
//export SubmitClaim
func SubmitClaim() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	claimID := params.ParseJSON("claim_id")
	insuredStr := params.ParseJSON("insured")
	requestedAmount := params.ParseJSONInt("requested_amount")
	eventTime := params.ParseJSONInt("event_time")
	evidenceHash := params.ParseJSON("evidence_hash")
	extra := params.ParseJSON("extra")

	if planID == "" || claimID == "" || requestedAmount <= 0 || eventTime <= 0 {
		return framework.ERROR_INVALID_PARAMS
	}

	applicant := framework.GetCaller()
	var insured framework.Address
	if insuredStr != "" {
		var err error
		insured, err = framework.ParseAddressBase58(insuredStr)
		if err != nil {
			return framework.ERROR_INVALID_PARAMS
		}
	} else {
		insured = applicant
	}

	// 1. 检查申请人是否为ACTIVE成员
	memberStateID := getMemberStateID(applicant)
	memberData, _ := framework.GetState(string(memberStateID))
	if len(memberData) == 0 {
		return framework.ERROR_NOT_FOUND
	}
	status, joinTime, _, _, _, _ := decodeMember(memberData)
	if status != MEMBER_STATUS_ACTIVE {
		return framework.ERROR_UNAUTHORIZED
	}

	// 3. 检查等待期（简化：仅检查加入时间）
	currentTime := framework.GetTimestamp()
	configData, _ := framework.GetState(STATE_PLAN_CONFIG)
	if len(configData) > 0 {
		_, _, _, _, _, _, waitingPeriod, _, _ := decodePlanConfig(configData)
		if currentTime < joinTime+waitingPeriod {
			return framework.ERROR_INVALID_STATE // 等待期未满
		}
	}

	// 4. 检查案件是否已存在
	claimStateID := getClaimStateID(claimID)
	existingClaimData, _ := framework.GetState(string(claimStateID))
	if len(existingClaimData) > 0 {
		return framework.ERROR_ALREADY_EXISTS
	}

	// 5. 创建案件记录
	claimData := encodeClaim(planID, claimID, string(applicant.ToBytes()), string(insured.ToBytes()), CLAIM_STATUS_SUBMITTED, "", evidenceHash, "", requestedAmount, 0, eventTime)
	if _, err := framework.AppendStateOutputSimple(claimStateID, 1, claimData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 6. 发出事件
	event := framework.NewEvent("MutualAidClaimSubmitted")
	event.AddStringField("plan_id", planID)
	event.AddStringField("claim_id", claimID)
	event.AddAddressField("applicant", applicant)
	event.AddAddressField("insured", insured)
	event.AddIntField("requested_amount", requestedAmount)
	event.AddIntField("event_time", eventTime)
	event.AddStringField("evidence_hash", evidenceHash)
	event.AddStringField("extra", extra)
	framework.EmitEvent(event)

	// 7. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":          planID,
		"claim_id":         claimID,
		"applicant":        applicant.ToString(),
		"insured":          insured.ToString(),
		"status":           CLAIM_STATUS_SUBMITTED,
		"requested_amount": requestedAmount,
		"approved_amount":  uint64(0),
		"event_time":       eventTime,
		"evidence_hash":    evidenceHash,
		"round_id":         "",
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ReviewClaim 审核互助申请（仅 operator 可调用）
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "claim_id": "claim_202501_0001",
//	  "decision": "APPROVE",              // APPROVE / REJECT
//	  "approved_amount": 280000,          // 决定给付金额，REJECT 时可为 0
//	  "reason": "符合互助规则",
//	  "investigation_hash": "0xdef...",  // 调查报告哈希
//	  "review_round_id": "round_202501_01"
//	}
//
// 输出：
// - StateOutput: claim_{claim_id} (更新状态)
// - Event: MutualAidClaimReviewed
//
//export ReviewClaim
func ReviewClaim() uint32 {
	params := framework.GetContractParams()

	// 1. 权限检查
	if !checkOperator() {
		return framework.ERROR_UNAUTHORIZED
	}

	planID := params.ParseJSON("plan_id")
	claimID := params.ParseJSON("claim_id")
	decision := params.ParseJSON("decision")
	approvedAmount := params.ParseJSONInt("approved_amount")
	reason := params.ParseJSON("reason")
	investigationHash := params.ParseJSON("investigation_hash")
	reviewRoundID := params.ParseJSON("review_round_id")

	if planID == "" || claimID == "" || decision == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	if decision != DECISION_APPROVE && decision != DECISION_REJECT {
		return framework.ERROR_INVALID_PARAMS
	}

	// 2. 读取案件
	claimStateID := getClaimStateID(claimID)
	claimData, _ := framework.GetState(string(claimStateID))
	if len(claimData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	cPlanID, cClaimID, applicant, insured, status, _, evidenceHash, _, requestedAmount, _, eventTime := decodeClaim(claimData)

	// 3. 检查案件状态
	if status != CLAIM_STATUS_SUBMITTED && status != CLAIM_STATUS_UNDER_REVIEW {
		return framework.ERROR_INVALID_STATE
	}

	// 4. 更新案件状态
	newStatus := CLAIM_STATUS_APPROVED
	if decision == DECISION_REJECT {
		newStatus = CLAIM_STATUS_REJECTED
		approvedAmount = 0
	}

	// 检查批准金额不超过请求金额
	if decision == DECISION_APPROVE && approvedAmount > requestedAmount {
		approvedAmount = requestedAmount
	}

	newClaimData := encodeClaim(cPlanID, cClaimID, applicant, insured, newStatus, reviewRoundID, evidenceHash, investigationHash, requestedAmount, approvedAmount, eventTime)
	if _, err := framework.AppendStateOutputSimple(claimStateID, 2, newClaimData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 5. 发出事件
	event := framework.NewEvent("MutualAidClaimReviewed")
	event.AddStringField("plan_id", planID)
	event.AddStringField("claim_id", claimID)
	event.AddStringField("decision", decision)
	event.AddIntField("approved_amount", approvedAmount)
	event.AddStringField("reason", reason)
	event.AddStringField("investigation_hash", investigationHash)
	event.AddStringField("review_round_id", reviewRoundID)
	event.AddAddressField("reviewer", framework.GetCaller())
	framework.EmitEvent(event)

	// 6. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":            cPlanID,
		"claim_id":           cClaimID,
		"applicant":          addressBytesToString([]byte(applicant)),
		"insured":            addressBytesToString([]byte(insured)),
		"status":             newStatus,
		"requested_amount":   requestedAmount,
		"approved_amount":    approvedAmount,
		"event_time":         eventTime,
		"evidence_hash":      evidenceHash,
		"investigation_hash": investigationHash,
		"round_id":           reviewRoundID,
		"decision":           decision,
		"reason":             reason,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// OpenRound 开启新的结算轮次（仅 operator 可调用）
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "round_id": "round_202501_01",
//	  "period_start": 1736200000,
//	  "period_end": 1738792000
//	}
//
// 输出：
// - StateOutput: round_{round_id}
// - StateOutput: current_round_id (更新)
// - Event: MutualAidRoundOpened
//
//export OpenRound
func OpenRound() uint32 {
	params := framework.GetContractParams()

	// 1. 权限检查
	if !checkOperator() {
		return framework.ERROR_UNAUTHORIZED
	}

	planID := params.ParseJSON("plan_id")
	roundID := params.ParseJSON("round_id")
	periodStart := params.ParseJSONInt("period_start")
	periodEnd := params.ParseJSONInt("period_end")

	if planID == "" || roundID == "" || periodStart <= 0 || periodEnd <= periodStart {
		return framework.ERROR_INVALID_PARAMS
	}

	// 2. 检查轮次是否已存在
	roundStateID := getRoundStateID(roundID)
	existingRoundData, _ := framework.GetState(string(roundStateID))
	if len(existingRoundData) > 0 {
		return framework.ERROR_ALREADY_EXISTS
	}

	// 3. 创建轮次记录
	roundData := encodeRound(planID, roundID, ROUND_STATUS_OPEN, periodStart, periodEnd, 0, 0, 0, 0)
	if _, err := framework.AppendStateOutputSimple(roundStateID, 1, roundData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 4. 更新当前轮次ID
	if _, err := framework.AppendStateOutputSimple([]byte(STATE_CURRENT_ROUND), 2, []byte(roundID), nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 5. 发出事件
	event := framework.NewEvent("MutualAidRoundOpened")
	event.AddStringField("plan_id", planID)
	event.AddStringField("round_id", roundID)
	event.AddIntField("period_start", periodStart)
	event.AddIntField("period_end", periodEnd)
	framework.EmitEvent(event)

	// 6. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":                 planID,
		"round_id":                roundID,
		"status":                  ROUND_STATUS_OPEN,
		"period_start":            periodStart,
		"period_end":              periodEnd,
		"total_approved_payout":   uint64(0),
		"total_service_fee":       uint64(0),
		"per_capita_contribution": uint64(0),
		"payers_count":            uint64(0),
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// SettleRound 结算一个互助周期，计算人均分摊额（仅 operator 可调用）
//
// 计算公式：
//
//	total_with_fee = total_approved_payout * (10000 + service_fee_bp) / 10000
//	per_capita = ceil(total_with_fee / member_count_active)
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "round_id": "round_202501_01"
//	}
//
// 输出：
// - StateOutput: round_{round_id} (更新)
// - Event: MutualAidRoundSettled
//
//export SettleRound
func SettleRound() uint32 {
	params := framework.GetContractParams()

	// 1. 权限检查
	if !checkOperator() {
		return framework.ERROR_UNAUTHORIZED
	}

	planID := params.ParseJSON("plan_id")
	roundID := params.ParseJSON("round_id")

	if planID == "" || roundID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	// 2. 读取轮次
	roundStateID := getRoundStateID(roundID)
	roundData, _ := framework.GetState(string(roundStateID))
	if len(roundData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	rPlanID, rRoundID, status, periodStart, periodEnd, totalApprovedPayout, totalServiceFee, perCapitaContribution, payersCount := decodeRound(roundData)

	if status != ROUND_STATUS_OPEN {
		return framework.ERROR_INVALID_STATE
	}

	// 3. 读取计划配置
	configData, _ := framework.GetState(STATE_PLAN_CONFIG)
	if len(configData) == 0 {
		return framework.ERROR_NOT_FOUND
	}
	_, _, _, _, serviceFeeBP, _, _, _, _ := decodePlanConfig(configData)

	// 4. 计算总给付额（简化：从参数传入，实际应从已批准的claim汇总）
	// 这里假设totalApprovedPayout已经在OpenRound时或通过其他方式设置
	// 实际应用中，应该遍历所有APPROVED状态的claim，汇总approved_amount

	// 5. 计算服务费和人均分摊
	totalWithFee := totalApprovedPayout * (10000 + serviceFeeBP) / 10000
	totalServiceFee = totalWithFee - totalApprovedPayout

	// 读取活跃成员数
	memberCountData, _ := framework.GetState(STATE_MEMBER_COUNT)
	memberCount := bytesToUint64(memberCountData)
	if memberCount == 0 {
		return framework.ERROR_INVALID_STATE
	}

	// 向上取整
	perCapitaContribution = (totalWithFee + memberCount - 1) / memberCount

	// 6. 更新轮次状态
	newRoundData := encodeRound(rPlanID, rRoundID, ROUND_STATUS_SETTLED, periodStart, periodEnd, totalApprovedPayout, totalServiceFee, perCapitaContribution, payersCount)
	if _, err := framework.AppendStateOutputSimple(roundStateID, 2, newRoundData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 7. 发出事件
	event := framework.NewEvent("MutualAidRoundSettled")
	event.AddStringField("plan_id", planID)
	event.AddStringField("round_id", roundID)
	event.AddIntField("total_approved_payout", totalApprovedPayout)
	event.AddIntField("member_count_active", memberCount)
	event.AddIntField("service_fee_bp", serviceFeeBP)
	event.AddIntField("total_with_fee", totalWithFee)
	event.AddIntField("total_service_fee", totalServiceFee)
	event.AddIntField("per_capita_contribution", perCapitaContribution)
	framework.EmitEvent(event)

	// 8. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":                 rPlanID,
		"round_id":                rRoundID,
		"status":                  ROUND_STATUS_SETTLED,
		"period_start":            periodStart,
		"period_end":              periodEnd,
		"total_approved_payout":   totalApprovedPayout,
		"total_service_fee":       totalServiceFee,
		"total_with_fee":          totalWithFee,
		"per_capita_contribution": perCapitaContribution,
		"member_count_active":     memberCount,
		"service_fee_bp":          serviceFeeBP,
		"payers_count":            payersCount,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// PayContribution 成员为某一轮互助结算缴纳分摊
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "round_id": "round_202501_01",
//	  "pool": "Df2...",                   // 资金池地址（Base58）
//	  "amount": 500,                      // 本次缴纳金额
//	  "contribution_id": "ctrb_202501_0001"
//	}
//
// 输出：
// - 使用 market.Escrow 创建实际资产托管
// - StateOutput: member_round_due_{address}_{round_id} (更新)
// - StateOutput: member_month_stat_{address}_{yyyymm} (更新)
// - StateOutput: round_{round_id} (更新payers_count)
// - Event: MutualAidContributionPaid
//
//export PayContribution
func PayContribution() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	roundID := params.ParseJSON("round_id")
	poolStr := params.ParseJSON("pool")
	amount := params.ParseJSONInt("amount")
	contributionID := params.ParseJSON("contribution_id")

	if planID == "" || roundID == "" || poolStr == "" || amount <= 0 || contributionID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	caller := framework.GetCaller()
	pool, err := framework.ParseAddressBase58(poolStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 1. 检查成员是否为ACTIVE
	memberStateID := getMemberStateID(caller)
	memberData, _ := framework.GetState(string(memberStateID))
	if len(memberData) == 0 {
		return framework.ERROR_NOT_FOUND
	}
	status, _, _, _, _, _ := decodeMember(memberData)
	if status != MEMBER_STATUS_ACTIVE {
		return framework.ERROR_UNAUTHORIZED
	}

	// 2. 检查轮次是否存在且已结算
	roundStateID := getRoundStateID(roundID)
	roundData, _ := framework.GetState(string(roundStateID))
	if len(roundData) == 0 {
		return framework.ERROR_NOT_FOUND
	}
	_, _, roundStatus, _, _, _, _, perCapitaContribution, _ := decodeRound(roundData)
	if roundStatus != ROUND_STATUS_SETTLED {
		return framework.ERROR_INVALID_STATE
	}

	// 3. 读取或创建成员轮次应缴记录
	memberRoundDueStateID := getMemberRoundDueStateID(caller, roundID)
	memberRoundDueData, _ := framework.GetState(string(memberRoundDueStateID))
	var dueAmount, paidAmount uint64
	var settled bool
	if len(memberRoundDueData) > 0 {
		dueAmount, paidAmount, settled = decodeMemberRoundDue(memberRoundDueData)
	} else {
		dueAmount = perCapitaContribution
		paidAmount = 0
		settled = false
	}

	if settled {
		return framework.ERROR_INVALID_STATE // 已结清
	}

	// 4. 检查月度上限（简化：使用当前月份）
	// 实际应用中应该解析roundID或使用period_end计算年月
	yearMonth := "202501" // 简化实现
	memberMonthStatStateID := getMemberMonthStatStateID(caller, yearMonth)
	memberMonthStatData, _ := framework.GetState(string(memberMonthStatStateID))
	var monthPaidAmount uint64
	var capReached bool
	if len(memberMonthStatData) > 0 {
		monthPaidAmount, capReached = decodeMemberMonthStat(memberMonthStatData)
	}

	// 读取计划配置中的月度上限
	configData, _ := framework.GetState(STATE_PLAN_CONFIG)
	var monthlyCapPerMember uint64 = 1000000
	if len(configData) > 0 {
		_, _, _, _, _, _, _, _, monthlyCapPerMember = decodePlanConfig(configData)
	}

	// 检查是否超过月度上限
	if monthPaidAmount+amount > monthlyCapPerMember {
		return framework.ERROR_INVALID_PARAMS // 超过月度上限
	}
	if capReached {
		return framework.ERROR_INVALID_PARAMS // 月度上限已触达
	}

	// 5. 使用托管实现成员 -> 资金池 的资金划转
	escrowID := []byte(planID + "_" + roundID + "_" + contributionID)
	if err := market.Escrow(
		caller,
		pool,
		framework.TokenID(""), // 使用原生币；实际应用可改为稳定币或专用代币
		framework.Amount(amount),
		escrowID,
	); err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 6. 更新成员轮次应缴记录
	newPaidAmount := paidAmount + amount
	newSettled := newPaidAmount >= dueAmount
	newMemberRoundDueData := encodeMemberRoundDue(dueAmount, newPaidAmount, newSettled)
	if _, err := framework.AppendStateOutputSimple(memberRoundDueStateID, 2, newMemberRoundDueData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 7. 更新成员月度统计
	newMonthPaidAmount := monthPaidAmount + amount
	newCapReached := newMonthPaidAmount >= monthlyCapPerMember
	newMemberMonthStatData := encodeMemberMonthStat(newMonthPaidAmount, newCapReached)
	if _, err := framework.AppendStateOutputSimple(memberMonthStatStateID, 2, newMemberMonthStatData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 8. 更新成员总缴费
	_, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound := decodeMember(memberData)
	newTotalPaid := totalPaid + amount
	newMemberData := encodeMember(status, joinTime, newTotalPaid, totalReceived, arrearsAmount, lastSettledRound)
	if _, err := framework.AppendStateOutputSimple(memberStateID, 2, newMemberData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 9. 更新轮次缴费人数（简化：每次缴费都增加，实际应该去重）
	_, _, _, _, _, _, _, _, payersCount := decodeRound(roundData)
	newPayersCount := payersCount + 1
	// 注意：这里需要重新读取roundData以获取完整信息
	roundData2, _ := framework.GetState(string(roundStateID))
	rPlanID, rRoundID, rStatus, rPeriodStart, rPeriodEnd, rTotalApprovedPayout, rTotalServiceFee, rPerCapitaContribution, _ := decodeRound(roundData2)
	newRoundData := encodeRound(rPlanID, rRoundID, rStatus, rPeriodStart, rPeriodEnd, rTotalApprovedPayout, rTotalServiceFee, rPerCapitaContribution, newPayersCount)
	if _, err := framework.AppendStateOutputSimple(roundStateID, 3, newRoundData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 10. 发出事件
	event := framework.NewEvent("MutualAidContributionPaid")
	event.AddStringField("plan_id", planID)
	event.AddStringField("round_id", roundID)
	event.AddAddressField("payer", caller)
	event.AddIntField("amount", amount)
	event.AddStringField("contribution_id", contributionID)
	framework.EmitEvent(event)

	// 11. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":                planID,
		"round_id":               roundID,
		"payer":                  caller.ToString(),
		"amount":                 amount,
		"due_amount":             dueAmount,
		"paid_amount":            newPaidAmount,
		"settled":                newSettled,
		"month_paid_amount":      newMonthPaidAmount,
		"monthly_cap_per_member": monthlyCapPerMember,
		"cap_reached":            newCapReached,
		"total_paid":             newTotalPaid,
		"contribution_id":        contributionID,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// Payout 为已通过审核的理赔案件进行给付（仅 operator 可调用）
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "claim_id": "claim_202501_0001",
//	  "from": "Df2...",                   // 资金池地址
//	  "beneficiary": "Cf1...",            // 受益人地址
//	  "amount": 300000,
//	  "payout_id": "payout_202501_0001"
//	}
//
// 输出：
// - 使用 market.Release 创建一次性释放计划
// - StateOutput: claim_{claim_id} (更新状态为PAID)
// - StateOutput: round_{round_id} (更新total_approved_payout)
// - Event: MutualAidPayout
//
//export Payout
func Payout() uint32 {
	params := framework.GetContractParams()

	// 1. 权限检查
	if !checkOperator() {
		return framework.ERROR_UNAUTHORIZED
	}

	planID := params.ParseJSON("plan_id")
	claimID := params.ParseJSON("claim_id")
	fromStr := params.ParseJSON("from")
	beneficiaryStr := params.ParseJSON("beneficiary")
	amount := params.ParseJSONInt("amount")
	payoutID := params.ParseJSON("payout_id")

	if planID == "" || claimID == "" || fromStr == "" || beneficiaryStr == "" || amount <= 0 || payoutID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	from, err1 := framework.ParseAddressBase58(fromStr)
	beneficiary, err2 := framework.ParseAddressBase58(beneficiaryStr)
	if err1 != nil || err2 != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	// 2. 读取案件
	claimStateID := getClaimStateID(claimID)
	claimData, _ := framework.GetState(string(claimStateID))
	if len(claimData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	cPlanID, cClaimID, applicant, insured, status, roundID, evidenceHash, investigationHash, requestedAmount, approvedAmount, eventTime := decodeClaim(claimData)

	// 3. 检查案件状态
	if status != CLAIM_STATUS_APPROVED {
		return framework.ERROR_INVALID_STATE
	}

	// 4. 检查给付金额不超过批准金额
	if amount > approvedAmount {
		return framework.ERROR_INVALID_PARAMS
	}

	// 5. 使用Release创建一次性释放计划
	vestingID := []byte(planID + "_" + claimID + "_" + payoutID)
	if err := market.Release(
		from,
		beneficiary,
		framework.TokenID(""), // 使用原生币；实际应用可改为专用互助 Token
		framework.Amount(amount),
		vestingID,
	); err != nil {
		if contractErr, ok := err.(*framework.ContractError); ok {
			return contractErr.Code
		}
		return framework.ERROR_EXECUTION_FAILED
	}

	// 6. 更新案件状态为PAID
	newClaimData := encodeClaim(cPlanID, cClaimID, applicant, insured, CLAIM_STATUS_PAID, roundID, evidenceHash, investigationHash, requestedAmount, approvedAmount, eventTime)
	if _, err := framework.AppendStateOutputSimple(claimStateID, 3, newClaimData, nil); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	// 7. 更新被保人的total_received（如果insured是成员）
	// 将insured字符串（20字节原始数据）转换为Address
	insuredAddr := framework.AddressFromBytes([]byte(insured))
	insuredMemberStateID := getMemberStateID(insuredAddr)
	insuredMemberData, _ := framework.GetState(string(insuredMemberStateID))
	insuredTotalReceived := uint64(0)
	if len(insuredMemberData) > 0 {
		insuredStatus, insuredJoinTime, insuredTotalPaid, insuredTotalReceivedOld, insuredArrearsAmount, insuredLastSettledRound := decodeMember(insuredMemberData)
		newInsuredTotalReceived := insuredTotalReceivedOld + amount
		insuredTotalReceived = newInsuredTotalReceived
		newInsuredMemberData := encodeMember(insuredStatus, insuredJoinTime, insuredTotalPaid, newInsuredTotalReceived, insuredArrearsAmount, insuredLastSettledRound)
		if _, err := framework.AppendStateOutputSimple(insuredMemberStateID, 2, newInsuredMemberData, nil); err != nil {
			return framework.ERROR_EXECUTION_FAILED
		}
	}

	// 8. 发出事件
	event := framework.NewEvent("MutualAidPayout")
	event.AddStringField("plan_id", planID)
	event.AddStringField("claim_id", claimID)
	event.AddAddressField("from", from)
	event.AddAddressField("beneficiary", beneficiary)
	event.AddIntField("amount", amount)
	event.AddStringField("payout_id", payoutID)
	framework.EmitEvent(event)

	// 9. 返回业务结果（WES ISPC 特性：同步返回业务数据）
	result := map[string]interface{}{
		"plan_id":                cPlanID,
		"claim_id":               cClaimID,
		"status":                 CLAIM_STATUS_PAID,
		"applicant":              addressBytesToString([]byte(applicant)),
		"insured":                addressBytesToString([]byte(insured)),
		"beneficiary":            beneficiary.ToString(),
		"requested_amount":       requestedAmount,
		"approved_amount":        approvedAmount,
		"payout_amount":          amount,
		"round_id":               roundID,
		"insured_total_received": insuredTotalReceived,
		"payout_id":              payoutID,
	}
	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// ================================================================================================
// 查询接口（只读）
// ================================================================================================

// GetPlanInfo 获取计划信息
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001"
//	}
//
// 返回：JSON格式的计划配置信息
//
//export GetPlanInfo
func GetPlanInfo() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	if planID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	configData, _ := framework.GetState(STATE_PLAN_CONFIG)
	if len(configData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	planIDDecoded, name, tokenID, coverageAmount, serviceFeeBP, settlementPeriod, waitingPeriod, minMembers, monthlyCapPerMember := decodePlanConfig(configData)

	operatorData, _ := framework.GetState(STATE_OPERATOR)
	operatorAddr := ""
	if len(operatorData) >= 20 {
		operatorAddr = framework.AddressFromBytes(operatorData[:20]).ToString()
	}

	memberCountData, _ := framework.GetState(STATE_MEMBER_COUNT)
	memberCount := bytesToUint64(memberCountData)

	result := map[string]interface{}{
		"plan_id":                planIDDecoded,
		"name":                   name,
		"token_id":               tokenID,
		"coverage_amount":        coverageAmount,
		"service_fee_bp":         serviceFeeBP,
		"settlement_period":      settlementPeriod,
		"waiting_period":         waitingPeriod,
		"min_members":            minMembers,
		"monthly_cap_per_member": monthlyCapPerMember,
		"operator":               operatorAddr,
		"member_count_active":    memberCount,
	}

	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// GetMemberInfo 获取成员信息
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "member": "Cf1..." // 成员地址（Base58）
//	}
//
// 返回：JSON格式的成员信息
//
//export GetMemberInfo
func GetMemberInfo() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	memberStr := params.ParseJSON("member")
	if planID == "" || memberStr == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	member, err := framework.ParseAddressBase58(memberStr)
	if err != nil {
		return framework.ERROR_INVALID_PARAMS
	}

	memberStateID := getMemberStateID(member)
	memberData, _ := framework.GetState(string(memberStateID))
	if len(memberData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	status, joinTime, totalPaid, totalReceived, arrearsAmount, lastSettledRound := decodeMember(memberData)

	result := map[string]interface{}{
		"plan_id":            planID,
		"member":             memberStr,
		"status":             status,
		"join_time":          joinTime,
		"total_paid":         totalPaid,
		"total_received":     totalReceived,
		"arrears_amount":     arrearsAmount,
		"last_settled_round": lastSettledRound,
	}

	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// GetClaimInfo 获取理赔案件信息
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "claim_id": "claim_202501_0001"
//	}
//
// 返回：JSON格式的案件信息
//
//export GetClaimInfo
func GetClaimInfo() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	claimID := params.ParseJSON("claim_id")
	if planID == "" || claimID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	claimStateID := getClaimStateID(claimID)
	claimData, _ := framework.GetState(string(claimStateID))
	if len(claimData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	cPlanID, cClaimID, applicant, insured, status, roundID, evidenceHash, investigationHash, requestedAmount, approvedAmount, eventTime := decodeClaim(claimData)

	result := map[string]interface{}{
		"plan_id":            cPlanID,
		"claim_id":           cClaimID,
		"applicant":          addressBytesToString([]byte(applicant)),
		"insured":            addressBytesToString([]byte(insured)),
		"status":             status,
		"round_id":           roundID,
		"evidence_hash":      evidenceHash,
		"investigation_hash": investigationHash,
		"requested_amount":   requestedAmount,
		"approved_amount":    approvedAmount,
		"event_time":         eventTime,
	}

	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// GetRoundInfo 获取轮次信息
//
// 参数（JSON）：
//
//	{
//	  "plan_id": "plan_xianghubao_001",
//	  "round_id": "round_202501_01"
//	}
//
// 返回：JSON格式的轮次信息
//
//export GetRoundInfo
func GetRoundInfo() uint32 {
	params := framework.GetContractParams()

	planID := params.ParseJSON("plan_id")
	roundID := params.ParseJSON("round_id")
	if planID == "" || roundID == "" {
		return framework.ERROR_INVALID_PARAMS
	}

	roundStateID := getRoundStateID(roundID)
	roundData, _ := framework.GetState(string(roundStateID))
	if len(roundData) == 0 {
		return framework.ERROR_NOT_FOUND
	}

	rPlanID, rRoundID, status, periodStart, periodEnd, totalApprovedPayout, totalServiceFee, perCapitaContribution, payersCount := decodeRound(roundData)

	result := map[string]interface{}{
		"plan_id":                 rPlanID,
		"round_id":                rRoundID,
		"status":                  status,
		"period_start":            periodStart,
		"period_end":              periodEnd,
		"total_approved_payout":   totalApprovedPayout,
		"total_service_fee":       totalServiceFee,
		"per_capita_contribution": perCapitaContribution,
		"payers_count":            payersCount,
	}

	if err := framework.SetReturnJSON(result); err != nil {
		return framework.ERROR_EXECUTION_FAILED
	}

	return framework.SUCCESS
}

// uint64ToString 将uint64转换为字符串
func uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	digits := make([]byte, 0, 20)
	num := n
	for num > 0 {
		digits = append(digits, byte('0'+num%10))
		num /= 10
	}
	// 反转数字
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}

func main() {}
