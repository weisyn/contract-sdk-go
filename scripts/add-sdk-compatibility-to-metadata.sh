#!/bin/bash
# add-sdk-compatibility-to-metadata.sh
# 为所有 Go 模板的 metadata.json 添加 SDK 版本兼容性字段

set -e

SDK_VERSION="${1:-0.1.0-alpha}"  # 默认版本，可通过参数传入
TEMPLATES_DIR="templates"

if [ ! -d "$TEMPLATES_DIR" ]; then
    echo "❌ Error: templates directory not found: $TEMPLATES_DIR"
    exit 1
fi

# 检查 jq 是否安装
if ! command -v jq &> /dev/null; then
    echo "❌ Error: jq is required but not installed"
    echo "   Install with: brew install jq (macOS) or apt-get install jq (Linux)"
    exit 1
fi

echo "🔧 Adding SDK compatibility fields to Go template metadata.json files"
echo "   SDK Version: $SDK_VERSION"
echo ""

# 计数器
total=0
updated=0
skipped=0

find "$TEMPLATES_DIR" -name "metadata.json" -type f | sort | while read -r metadata; do
    total=$((total + 1))
    echo "Processing: $metadata"
    
    # 检查是否已包含 sdkCompatibility 字段
    if jq -e '.sdkCompatibility' "$metadata" > /dev/null 2>&1; then
        echo "  ⏭️  Already has sdkCompatibility, skipping"
        skipped=$((skipped + 1))
        continue
    fi
    
    # 创建临时文件
    tmp_file=$(mktemp)
    
    # 使用 jq 添加字段
    # 如果已有 version 字段，保留；否则设置为 1.0.0
    jq \
        --arg go_version "$SDK_VERSION" \
        '.sdkCompatibility = {"go": ">='$SDK_VERSION' <0.2.0"} | 
         .sinceSdk = {"go": $go_version} |
         (.version //= "1.0.0")' \
        "$metadata" > "$tmp_file"
    
    # 检查是否有变更
    if ! diff -q "$metadata" "$tmp_file" > /dev/null 2>&1; then
        # 替换原文件
        mv "$tmp_file" "$metadata"
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

