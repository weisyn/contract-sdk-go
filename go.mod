module github.com/weisyn/contract-sdk-go

go 1.24.0

// WES Smart Contract SDK for Go
//
// 这是一个独立的SDK模块，用于TinyGo编译的WASM合约
// 使用 Go 1.24 以兼容 TinyGo 0.31.0+
//
// 特点：
// 1. 零外部依赖（纯Go标准库）
// 2. 针对WASM/WASI优化
// 3. 轻量级API设计
// 4. 完整的Host ABI v1.0.0支持
