package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestExamplesBuild 测试所有示例可以成功构建
func TestExamplesBuild(t *testing.T) {
	// 检查 TinyGo 是否安装
	_, err := exec.LookPath("tinygo")
	if err != nil {
		t.Skip("TinyGo not installed, skipping build tests")
	}

	// 获取示例目录
	examplesDir := filepath.Join("..", "examples")

	examples := []string{"hello-world", "simple-token"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			exampleDir := filepath.Join(examplesDir, example)
			buildScript := filepath.Join(exampleDir, "build.sh")

			// 检查构建脚本是否存在
			if _, err := os.Stat(buildScript); os.IsNotExist(err) {
				t.Fatalf("Build script not found: %s", buildScript)
			}

			// 执行构建
			cmd := exec.Command("bash", "build.sh")
			cmd.Dir = exampleDir
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Build failed for %s:\n%s", example, string(output))
				return
			}

			t.Logf("Build succeeded for %s:\n%s", example, string(output))

			// 检查 WASM 文件是否生成
			var wasmFile string
			if example == "hello-world" {
				wasmFile = "hello.wasm"
			} else {
				wasmFile = "token.wasm"
			}

			wasmPath := filepath.Join(exampleDir, wasmFile)
			info, err := os.Stat(wasmPath)
			if err != nil {
				t.Errorf("WASM file not found: %s", wasmPath)
				return
			}

			if info.Size() == 0 {
				t.Errorf("WASM file is empty: %s", wasmPath)
			}

			t.Logf("WASM file size: %d bytes", info.Size())
		})
	}
}

// TestExamplesStructure 测试示例目录结构
func TestExamplesStructure(t *testing.T) {
	examplesDir := filepath.Join("..", "examples")

	examples := []struct {
		name          string
		requiredFiles []string
	}{
		{
			name: "hello-world",
			requiredFiles: []string{
				"hello.go",
				"build.sh",
			},
		},
		{
			name: "simple-token",
			requiredFiles: []string{
				"token.go",
				"build.sh",
				"go.mod",
				"README.md",
			},
		},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			exampleDir := filepath.Join(examplesDir, example.name)

			for _, file := range example.requiredFiles {
				filePath := filepath.Join(exampleDir, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Required file missing: %s", file)
				}
			}
		})
	}
}

// TestSDKStructure 测试SDK目录结构完整性
func TestSDKStructure(t *testing.T) {
	requiredDirs := []string{
		"../framework",
		"../examples",
		"../scripts",
		"../docs",
		"../tests",
	}

	for _, dir := range requiredDirs {
		t.Run(dir, func(t *testing.T) {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Required directory missing: %s", dir)
			}
		})
	}

	requiredFiles := []string{
		"../README.md",
		"../go.mod",
		"../framework/contract_base.go",
		"../framework/host_functions.go",
		"../scripts/build-all.sh",
	}

	for _, file := range requiredFiles {
		t.Run(file, func(t *testing.T) {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("Required file missing: %s", file)
			}
		})
	}
}
