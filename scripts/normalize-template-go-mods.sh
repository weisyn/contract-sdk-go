#!/bin/bash
# normalize-template-go-mods.sh
# 批量规范化所有模板的 go.mod 文件
# 移除 replace 指令，使用发布版 SDK 版本号

set -e

SDK_VERSION="${1:-v0.1.0-alpha}"  # 默认版本，可通过参数传入
TEMPLATES_DIR="templates"

if [ ! -d "$TEMPLATES_DIR" ]; then
    echo "❌ Error: templates directory not found: $TEMPLATES_DIR"
    exit 1
fi

echo "🔧 Normalizing go.mod files with SDK version: $SDK_VERSION"
echo ""

# 计数器
total=0
updated=0
skipped=0

find "$TEMPLATES_DIR" -name "go.mod" -type f | sort | while read -r go_mod; do
    total=$((total + 1))
    echo "Processing: $go_mod"
    
    # 检查是否包含 replace 或 v0.0.0
    has_replace=$(grep -c "replace github.com/weisyn/contract-sdk-go =>" "$go_mod" 2>/dev/null || echo "0")
    has_v000=$(grep -c "require github.com/weisyn/contract-sdk-go v0.0.0" "$go_mod" 2>/dev/null || echo "0")
    
    if [ "$has_replace" -eq 0 ] && [ "$has_v000" -eq 0 ]; then
        echo "  ⏭️  Already normalized, skipping"
        skipped=$((skipped + 1))
        continue
    fi
    
    # 创建临时文件
    tmp_file=$(mktemp)
    
    # 处理文件
    # 1. 移除 replace 行
    # 2. 移除 replace 相关的注释
    # 3. 替换 require 版本号
    sed -E \
        -e '/^replace github\.com\/weisyn\/contract-sdk-go =>/d' \
        -e '/^\/\/ 本地开发时，使用 replace/d' \
        -e '/^\/\/ 提取到独立仓库后，这个 replace 将被移除/d' \
        -e "s|require github\.com/weisyn/contract-sdk-go v0\.0\.0|require github.com/weisyn/contract-sdk-go $SDK_VERSION|g" \
        "$go_mod" > "$tmp_file"
    
    # 检查是否有变更
    if ! diff -q "$go_mod" "$tmp_file" > /dev/null 2>&1; then
        # 替换原文件
        mv "$tmp_file" "$go_mod"
        echo "  ✅ Updated"
        updated=$((updated + 1))
    else
        rm "$tmp_file"
        echo "  ⏭️  No changes needed"
        skipped=$((skipped + 1))
    fi
done

echo ""
echo "📊 Summary:"
echo "  Total files: $total"
echo "  Updated: $updated"
echo "  Skipped: $skipped"
echo ""
echo "✅ Done!"

