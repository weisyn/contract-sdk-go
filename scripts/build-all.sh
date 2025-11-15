#!/bin/bash
set -e
# build all examples
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXAMPLES_DIR="$SCRIPT_DIR/../examples"

echo "🔨 批量编译所有示例..."
echo ""

for dir in "$EXAMPLES_DIR"/*; do
  if [ -d "$dir" ] && [ -f "$dir/build.sh" ]; then
    example_name=$(basename "$dir")
    echo "▶ 编译 $example_name ..."
    (cd "$dir" && bash build.sh)
    echo ""
  fi
done

echo "✅ 所有示例编译完成！"
