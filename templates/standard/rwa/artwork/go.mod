module github.com/weisyn/contract-sdk-go/examples/rwa/artwork

go 1.24.0

toolchain go1.24.7

// 本地开发时，使用 replace 指向本地SDK
replace github.com/weisyn/contract-sdk-go => ../../..

require github.com/weisyn/contract-sdk-go v0.0.0
