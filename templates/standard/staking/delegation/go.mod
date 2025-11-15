module github.com/weisyn/contract-sdk-go/examples/staking/delegation

go 1.24.0

toolchain go1.24.7

// 本地开发时，使用 replace 指向上级 SDK 目录
// 提取到独立仓库后，这个 replace 将被移除
replace github.com/weisyn/contract-sdk-go => ../../..

require github.com/weisyn/contract-sdk-go v0.0.0

