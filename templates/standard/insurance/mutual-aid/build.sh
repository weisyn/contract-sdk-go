#!/bin/bash

# WES 智能合约构建脚本 - Insurance Mutual Aid 示例
#
# 使用 TinyGo 编译 Go 合约为 WASM

set -e

echo "🔨 编译 Mutual Aid 互助险合约..."

# 检查 TinyGo 是否安装
if ! command -v tinygo &> /dev/null; then
    echo "❌ 错误: TinyGo 未安装"
    echo "请访问 https://tinygo.org/getting-started/install/ 安装 TinyGo"
    exit 1
fi

# 编译参数说明:
# -target=wasi        : 目标平台为 WASI (WebAssembly System Interface)
# -scheduler=none     : 禁用调度器(合约不需要并发)
# -no-debug           : 移除调试信息,减小体积
# -opt=2              : 优化级别 2 (平衡大小和性能)
# -gc=leaking         : 使用泄漏 GC (最简单,适合短生命周期合约)

tinygo build -o main.wasm \
  -target=wasi \
  -scheduler=none \
  -no-debug \
  -opt=2 \
  -gc=leaking \
  main.go

# 检查输出
if [ -f main.wasm ]; then
    SIZE=$(wc -c < main.wasm | tr -d ' ')
    echo "✅ 编译成功!"
    echo "📦 WASM 文件大小: $SIZE 字节"
    echo "📄 输出文件: main.wasm"
    
    # 显示 WASM 导出函数
    if command -v wasm-objdump &> /dev/null; then
        echo ""
        echo "📋 导出的函数:"
        wasm-objdump -x main.wasm | grep "export" | grep "func"
    fi
else
    echo "❌ 编译失败"
    exit 1
fi


