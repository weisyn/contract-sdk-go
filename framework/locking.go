//go:build tinygo || (js && wasm)

package framework

// 锁定条件（LockingConditions）TinyGo 侧辅助
//
// 说明：宿主当前接受的锁定条件编码为 JSON 数组，数组元素为
// protojson 编码后的单个 LockingCondition 对象（字符串拼接后形成数组）。
// 本辅助提供最小能力的拼装函数，避免在 TinyGo 环境依赖 Protobuf。

// BuildLockingJSONArray 将若干 protojson 字符串拼装为 JSON 数组字节
// 入参示例：
//   conds := []string{
//     `{"singleKeyLock": {"requiredAddressHash": "..."}}`,
//     `{"timeLock": {"unlockTimestamp": 1735689600, "baseLock": {...}}}`,
//   }
// 返回：[]byte("[<cond1>,<cond2>,...]")
func BuildLockingJSONArray(conditions []string) []byte {
    if len(conditions) == 0 {
        return nil
    }
    // 计算总长度：逗号 + 方括号
    // 简单串接，避免引入复杂 JSON 库
    // 形如: [cond0,cond1,...]
    // 注意：传入的每个条件字符串必须是合法的 JSON 对象文本
    totalLen := 2 // for '[' and ']'
    if len(conditions) > 1 {
        totalLen += len(conditions) - 1 // commas
    }
    for i := 0; i < len(conditions); i++ {
        totalLen += len(conditions[i])
    }

    buf := make([]byte, 0, totalLen)
    buf = append(buf, '[')
    for i := 0; i < len(conditions); i++ {
        if i > 0 {
            buf = append(buf, ',')
        }
        // 直接拼接原始 JSON 片段
        c := conditions[i]
        if len(c) > 0 {
            buf = append(buf, c...)
        }
    }
    buf = append(buf, ']')
    return buf
}
