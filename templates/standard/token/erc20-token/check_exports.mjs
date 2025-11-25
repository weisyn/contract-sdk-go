import { readFileSync } from 'fs';

const wasmPath = process.argv[2] || 'main.wasm';
console.log(`\n🔍 检查 WASM 文件: ${wasmPath}\n`);

try {
  // 读取 WASM 文件
  const wasmBytes = readFileSync(wasmPath);
  console.log(`📦 WASM 文件大小: ${wasmBytes.length} 字节\n`);

  // 使用 WebAssembly API 编译模块
  const module = await WebAssembly.compile(wasmBytes);
  const exports = WebAssembly.Module.exports(module);

  // 内部函数集合（与工作台保持一致）
  const internalFunctions = new Set([
    'malloc',
    'calloc',
    'realloc',
    'free',
    '_start',
    '_initialize',
  ]);

  console.log('=== 所有导出项 ===');
  const allExports = [];
  exports.forEach((e) => {
    console.log(`  - ${e.name} (${e.kind})`);
    allExports.push({ name: e.name, kind: e.kind });
  });

  // 过滤出业务导出函数
  const exportedFunctions = exports
    .filter((e) => {
      if (e.kind !== 'function') return false;
      const name = e.name;
      if (!name || typeof name !== 'string' || name.length === 0) return false;
      if (internalFunctions.has(name)) return false;
      if (name.startsWith('_')) return false;
      return true;
    })
    .map((e) => e.name);

  console.log('\n=== 业务导出函数（过滤后）===');
  if (exportedFunctions.length === 0) {
    console.log('  ❌ 未找到业务导出函数！');
    console.log('\n可能的原因：');
    console.log('  1. TinyGo 编译时未正确导出函数');
    console.log('  2. 函数名被优化或重命名');
    console.log('  3. 编译选项影响了导出');
  } else {
    exportedFunctions.forEach((name) => {
      console.log(`  ✅ ${name}`);
    });
    console.log(`\n总计: ${exportedFunctions.length} 个业务导出函数`);
  }

  // 检查预期的函数
  const expectedFunctions = [
    'Initialize',
    'Transfer',
    'Mint',
    'Burn',
    'Approve',
    'Airdrop',
    'Freeze',
  ];

  console.log('\n=== 预期函数检查 ===');
  const missingFunctions = expectedFunctions.filter(
    (name) => !exportedFunctions.includes(name)
  );
  const foundFunctions = expectedFunctions.filter((name) =>
    exportedFunctions.includes(name)
  );

  foundFunctions.forEach((name) => {
    console.log(`  ✅ ${name}`);
  });
  if (missingFunctions.length > 0) {
    console.log('\n缺失的函数:');
    missingFunctions.forEach((name) => {
      console.log(`  ❌ ${name}`);
    });
  }

  // 额外的导出函数（不在预期列表中）
  const extraFunctions = exportedFunctions.filter(
    (name) => !expectedFunctions.includes(name)
  );
  if (extraFunctions.length > 0) {
    console.log('\n额外的导出函数:');
    extraFunctions.forEach((name) => {
      console.log(`  ℹ️  ${name}`);
    });
  }

  console.log('\n=== 总结 ===');
  console.log(`预期函数: ${expectedFunctions.length}`);
  console.log(`找到函数: ${foundFunctions.length}`);
  console.log(`缺失函数: ${missingFunctions.length}`);
  console.log(`额外函数: ${extraFunctions.length}`);
  console.log(`总业务函数: ${exportedFunctions.length}`);

  if (missingFunctions.length === 0 && exportedFunctions.length > 0) {
    console.log('\n✅ 所有预期函数都已正确导出！');
    process.exit(0);
  } else {
    console.log('\n⚠️  存在问题，请检查上述输出');
    process.exit(1);
  }
} catch (error) {
  console.error('❌ 错误:', error.message);
  console.error(error.stack);
  process.exit(1);
}
