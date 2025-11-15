# Hello World 合约示例

最简单的WES智能合约示例，展示SDK的基本用法。

## 功能

- ✅ **Initialize** - 初始化合约，发出欢迎消息
- ✅ **SayHello** - 问候调用者，返回个性化消息
- ✅ **GetInfo** - 查询合约信息（只读）

## 快速开始

### 构建

```bash
# 确保已安装 TinyGo >= 0.31
tinygo version

# 编译合约
bash build.sh
```

编译成功后会生成 `main.wasm` 文件。

### 测试编译

```bash
# 使用 Go 标准编译器测试逻辑（非WASM）
go build .
```

## 合约接口

### Initialize()

初始化合约并发出欢迎事件。

**调用示例**:
```bash
wes contract deploy main.wasm --function Initialize
```

**返回**: 成功返回 0

**事件**:
```json
{
  "name": "ContractInitialized",
  "data": {
    "message": "Hello World Contract Initialized",
    "timestamp": "1234567890"
  }
}
```

### SayHello()

向调用者问候并返回个性化消息。

**调用示例**:
```bash
wes contract call <contract-address> --function SayHello
```

**返回**: 成功返回 0

**返回数据**: `"Hello, <caller-address>!"`

**事件**:
```json
{
  "name": "Greeting",
  "data": {
    "caller": "<caller-address>",
    "message": "Hello, <caller-address>!"
  }
}
```

### GetInfo()

查询合约基本信息（只读操作）。

**调用示例**:
```bash
wes contract query <contract-address> --function GetInfo
```

**返回数据**:
```json
{
  "name": "HelloWorld",
  "version": "1.0.0",
  "author": "WES Team"
}
```

## 代码结构

```go
type HelloContract struct {
    framework.ContractBase
}

//export Initialize
func Initialize() uint32 {
    // 初始化逻辑
}

//export SayHello
func SayHello() uint32 {
    // 业务逻辑
}

//export GetInfo
func GetInfo() uint32 {
    // 查询逻辑
}
```

## 核心概念

### 1. 合约结构

所有合约都应嵌入 `framework.ContractBase` 以获得基础功能：

```go
type HelloContract struct {
    framework.ContractBase
}
```

### 2. 导出函数

使用 `//export` 注释标记对外可调用的函数：

```go
//export SayHello
func SayHello() uint32 {
    // ...
    return 0  // 返回错误码
}
```

### 3. 宿主函数

通过 SDK 调用链上环境：

```go
caller := contract.GetCaller()        // 获取调用者
timestamp := contract.GetTimestamp()  // 获取时间戳
contract.EmitEvent("name", data)      // 发出事件
```

### 4. 错误处理

统一使用错误码：

```go
return 0  // SUCCESS
return 1  // ERROR_UNKNOWN
return 2  // ERROR_INVALID_PARAMS
return 3  // ERROR_UNAUTHORIZED
```

## 扩展示例

基于此示例，可以扩展为：

1. **计数器合约** - 添加状态管理
2. **投票合约** - 添加权限控制
3. **留言板** - 添加事件日志

## 部署与调用（使用 CLI）

### 前置条件

1. **启动节点**:
   ```bash
   # 从项目根目录运行
   make run-testing  # 或 ./bin/testing
   ```

2. **编译合约**:
   ```bash
   cd examples/hello-world
   bash build.sh
   # 生成 main.wasm
   ```

### 使用 WES CLI 部署

```bash
# 交互式部署（推荐）
wes

# 选择: 合约管理 -> 部署合约
# 按提示输入：
#   - 钱包名称
#   - 钱包密码
#   - WASM文件路径: examples/hello-world/main.wasm
#   - 合约名称: HelloWorld
#   - ABI版本: 1.0.0（按需）
```

部署成功后会返回 **合约ID**（content hash），保存它用于调用。

### 使用 CLI 调用合约

```bash
# 交互式调用
wes

# 选择: 合约管理 -> 调用合约
# 按提示输入：
#   - 钱包名称
#   - 钱包密码  
#   - 合约ID: <部署时返回的content_hash>
#   - 方法名: SayHello
#   - 参数: (留空或按格式输入)
```

### 使用 CLI 查询合约

```bash
# 查询合约信息（只读）
wes

# 选择: 合约管理 -> 查询合约
# 按提示输入：
#   - 合约ID: <content_hash>
#   - 方法名: GetInfo
```

## 相关文档

- [WES Contract SDK 文档](../../README.md)
- [Simple Token 示例](../simple-token/README.md)
- [Host ABI 规范](../../../_docs/specs/contracts/HOST_ABI_V1_MINIMAL.md)

## 许可证

Apache-2.0 License

