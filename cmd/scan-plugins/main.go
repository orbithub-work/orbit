package main

import (
	"fmt"
	"os"
	"path/filepath"

	"media-assistant-os/internal/services"
)

func main() {
	// Get plugins directory
	pluginsDir := "./plugins"
	if len(os.Args) > 1 {
		pluginsDir = os.Args[1]
	}

	// Make absolute path
	absPath, err := filepath.Abs(pluginsDir)
	if err != nil {
		fmt.Printf("❌ Failed to resolve path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔍 智归档OS 插件扫描器")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("扫描目录: %s\n\n", absPath)

	// Create scanner
	scanner := services.NewPluginScanner(absPath)

	// Scan all plugins
	manifests, err := scanner.ScanAll()
	if err != nil {
		fmt.Printf("❌ Scan failed: %v\n", err)
		os.Exit(1)
	}

	// Print report
	scanner.PrintDiscoveryReport(manifests)

	// Validate each manifest
	fmt.Println("🔍 验证插件配置...")
	hasErrors := false
	for _, manifest := range manifests {
		if err := scanner.ValidateManifest(&manifest); err != nil {
			fmt.Printf("❌ %s: %v\n", manifest.ID, err)
			hasErrors = true
		} else {
			fmt.Printf("✅ %s: 配置有效\n", manifest.ID)
		}
	}

	if hasErrors {
		fmt.Println("\n⚠️  部分插件配置有误，请检查")
		os.Exit(1)
	}

	fmt.Println("\n✅ 所有插件配置有效！")
}
