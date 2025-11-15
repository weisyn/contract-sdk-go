#!/bin/bash

# 编译治理代币合约
# 使用 TinyGo 编译为 WASM

set -e

echo "🔨 编译治理代币合约..."

tinygo build -o main.wasm \
    -target=wasi \
    -scheduler=none \
    -no-debug \
    -opt=2 \
    main.go

if [ $? -eq 0 ]; then
    echo "✅ 编译成功: main.wasm"
    ls -lh main.wasm
else
    echo "❌ 编译失败"
    exit 1
fi

