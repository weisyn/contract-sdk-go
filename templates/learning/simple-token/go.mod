module github.com/weisyn/v1/contracts/templates/learning/simple-token

go 1.24

// WES 代币学习合约模板
//
// 这是一个教学用的代币合约示例，展示如何在WES上开发代币应用
// 使用 Go 1.24 以兼容 TinyGo 0.39.0

require github.com/weisyn/contract-sdk-go v0.0.0

// 使用本地SDK（开发模式）
replace github.com/weisyn/contract-sdk-go => ../../../_sdks/contract-sdk-go
