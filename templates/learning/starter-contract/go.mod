module github.com/weisyn/v1/contracts/templates/learning/starter-contract

go 1.24

// WES 入门合约模板
//
// 这是一个教学用的空白合约模板，展示如何从零开始开发自定义合约
// 使用 Go 1.24 以兼容 TinyGo 0.39.0

require github.com/weisyn/contract-sdk-go v0.0.0

// 使用本地SDK（开发模式）
replace github.com/weisyn/contract-sdk-go => ../../../_sdks/contract-sdk-go
