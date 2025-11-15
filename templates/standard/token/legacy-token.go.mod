module github.com/weisyn/v1/contracts/templates/token

go 1.24.0

toolchain go1.24.7

// 合约 SDK 模块（本地开发通过 replace 指向源码，实现细节不对外暴露）
replace github.com/weisyn/contract-sdk-go => ../../../_sdks/contract-sdk-go

require github.com/weisyn/contract-sdk-go v0.0.0
