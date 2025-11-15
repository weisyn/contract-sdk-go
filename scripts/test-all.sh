#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SDK_ROOT="$SCRIPT_DIR/.."

echo "🧪 运行 WES Contract SDK Go 测试套件..."
echo ""

# 1. 单元测试
echo "▶ 运行单元测试 (framework)..."
cd "$SDK_ROOT/framework"
go test -v -cover
echo ""

# 2. 集成测试
echo "▶ 运行集成测试 (build & structure)..."
cd "$SDK_ROOT/tests"
go test -v
echo ""

echo "✅ 所有测试通过！"
