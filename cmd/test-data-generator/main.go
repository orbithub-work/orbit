package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// TestDataGenerator 测试数据生成器
type TestDataGenerator struct {
	targetDir string
	lean      bool
}

// NewTestDataGenerator 创建测试数据生成器
func NewTestDataGenerator(targetDir string, lean bool) *TestDataGenerator {
	return &TestDataGenerator{
		targetDir: targetDir,
		lean:      lean,
	}
}

// GenerateCreatorEnvironment 生成创作者环境
func (g *TestDataGenerator) GenerateCreatorEnvironment(targetCount int) error {
	fmt.Printf("🚀 生成创作者测试环境到: %s\n", g.targetDir)
	fmt.Println("================================================")

	// 1. 创建目录结构
	if err := g.createDirectoryStructure(); err != nil {
		return err
	}

	// 2. 生成素材文件
	if err := g.generateAssetFiles(targetCount); err != nil {
		return err
	}

	fmt.Println("================================================")
	fmt.Println("✅ 创作者测试环境生成完成！")
	return nil
}

// createDirectoryStructure 创建创作者目录结构
func (g *TestDataGenerator) createDirectoryStructure() error {
	fmt.Println("📁 创建目录结构...")

	dirs := []string{
		// 摄影师
		"摄影师/作品集/2024/1月-人像",
		"摄影师/作品集/2024/2月-风景",
		"摄影师/作品集/2024/3月-街拍",
		"摄影师/RAW文件",
		"摄影师/修图素材",
		"摄影师/精选作品",

		// 小博主
		"博主/视频项目/小红书",
		"博主/视频项目/抖音",
		"博主/封面设计",
		"博主/文案草稿",
		"博主/发布内容",

		// 内容创作者
		"创作者/图片素材/背景",
		"创作者/图片素材/图标",
		"创作者/视频素材/片段",
		"创作者/音频素材/音乐",
		"创作者/音频素材/音效",
		"创作者/字体库",
		"创作者/模板库/PSD",
		"创作者/模板库/AI",

		// 项目
		"项目/小红书运营/2024-01",
		"项目/小红书运营/2024-02",
		"项目/抖音视频/2024-01",
		"项目/公众号/2024-01",

		// 归档
		"归档/2023",
		"归档/2022",
	}

	for _, dir := range dirs {
		path := filepath.Join(g.targetDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %w", dir, err)
		}
	}

	fmt.Printf("✅ 已创建 %d 个目录\n", len(dirs))
	return nil
}

// generateAssetFiles 生成素材文件
func (g *TestDataGenerator) generateAssetFiles(targetCount int) error {
	fmt.Println("📄 生成素材文件...")

	// 定义文件生成模板
	type FileTemplate struct {
		dir  string
		ext  string
		size int
		desc string
	}

	var templates []FileTemplate

	// 1. 摄影师照片
	photoFiles := []FileTemplate{
		{"摄影师/作品集/2024/1月-人像", ".jpg", 3 * 1024 * 1024, "人像照片"},
		{"摄影师/作品集/2024/1月-人像", ".jpg", 4 * 1024 * 1024, "人像照片"},
		{"摄影师/作品集/2024/2月-风景", ".jpg", 5 * 1024 * 1024, "风景照片"},
		{"摄影师/作品集/2024/2月-风景", ".jpg", 6 * 1024 * 1024, "风景照片"},
		{"摄影师/作品集/2024/3月-街拍", ".jpg", 2 * 1024 * 1024, "街拍照片"},
		{"摄影师/RAW文件", ".cr2", 25 * 1024 * 1024, "RAW文件"},
		{"摄影师/RAW文件", ".nef", 30 * 1024 * 1024, "RAW文件"},
		{"摄影师/精选作品", ".jpg", 4 * 1024 * 1024, "精选照片"},
	}
	templates = append(templates, photoFiles...)

	// 2. 博主视频
	videoFiles := []FileTemplate{
		{"博主/视频项目/小红书", ".mp4", 50 * 1024 * 1024, "小红书视频"},
		{"博主/视频项目/小红书", ".mov", 80 * 1024 * 1024, "小红书视频"},
		{"博主/视频项目/抖音", ".mp4", 30 * 1024 * 1024, "抖音视频"},
		{"博主/视频项目/抖音", ".mov", 60 * 1024 * 1024, "抖音视频"},
	}
	templates = append(templates, videoFiles...)

	// 3. 音频素材
	audioFiles := []FileTemplate{
		{"博主/视频项目/小红书", ".mp3", 5 * 1024 * 1024, "背景音乐"},
		{"博主/视频项目/抖音", ".wav", 10 * 1024 * 1024, "音效"},
		{"创作者/音频素材/音乐", ".mp3", 3 * 1024 * 1024, "音乐"},
		{"创作者/音频素材/音效", ".wav", 8 * 1024 * 1024, "音效"},
	}
	templates = append(templates, audioFiles...)

	// 4. 封面和设计
	designFiles := []FileTemplate{
		{"博主/封面设计", ".png", 2 * 1024 * 1024, "封面"},
		{"博主/封面设计", ".jpg", 1536 * 1024, "封面"},
		{"创作者/图片素材/背景", ".png", 3 * 1024 * 1024, "背景"},
		{"创作者/图片素材/图标", ".png", 512 * 1024, "图标"},
	}
	templates = append(templates, designFiles...)

	// 5. 文档和文案
	docFiles := []FileTemplate{
		{"博主/文案草稿", ".txt", 16 * 1024, "文案"},
		{"博主/文案草稿", ".md", 32 * 1024, "文案"},
		{"项目/小红书运营/2024-01", ".xlsx", 64 * 1024, "数据"},
		{"项目/抖音视频/2024-01", ".pptx", 256 * 1024, "脚本"},
	}
	templates = append(templates, docFiles...)

	// 6. 字体和模板
	fontFiles := []FileTemplate{
		{"创作者/字体库", ".ttf", 512 * 1024, "字体"},
		{"创作者/字体库", ".otf", 640 * 1024, "字体"},
		{"创作者/模板库/PSD", ".psd", 10 * 1024 * 1024, "模板"},
		{"创作者/模板库/AI", ".ai", 8 * 1024 * 1024, "模板"},
	}
	templates = append(templates, fontFiles...)

	// 确定每种模板生成的数量
	// 如果 targetCount 很大，我们就循环使用模板

	count := 0
	lastProgress := -1

	for count < targetCount {
		// 随机选择一个模板，或者按顺序循环
		tpl := templates[count%len(templates)]

		// 稍微随机化大小
		var size int
		if g.lean {
			size = 1024 // 精简模式固定 1KB
		} else if tpl.size > 1024 {
			// +/- 20%
			diff := int(float64(tpl.size) * 0.2)
			b := make([]byte, 1)
			rand.Read(b)
			offset := int(b[0])%(diff*2) - diff
			size = tpl.size + offset
			if size < 100 {
				size = 100
			}
		} else {
			size = tpl.size
		}

		path := filepath.Join(g.targetDir, tpl.dir, fmt.Sprintf("%s_%05d%s", tpl.desc, count+1, tpl.ext))
		if err := g.createFile(path, size, tpl.ext); err != nil {
			return err
		}

		count++

		// 打印进度
		progress := count * 100 / targetCount
		if progress != lastProgress && progress%5 == 0 {
			fmt.Printf("\r⏳ 进度: %d%% (%d/%d)", progress, count, targetCount)
			lastProgress = progress
		}
	}

	fmt.Println()
	return nil
}

// createFile 创建模拟文件
func (g *TestDataGenerator) createFile(path string, size int, ext string) error {
	var content []byte

	switch ext {
	case ".jpg", ".jpeg":
		content = createValidImage("jpeg")
	case ".png":
		content = createValidImage("png")
	default:
		content = make([]byte, 0)
		// 添加文件头
		switch ext {
		case ".cr2", ".nef":
			content = append([]byte{0x49, 0x49, 0x2A, 0x00}, content...)
		case ".mp4":
			content = append([]byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}, content...)
		case ".mov":
			content = append([]byte{0x00, 0x00, 0x00, 0x14, 0x66, 0x74, 0x79, 0x70}, content...)
		case ".mp3":
			content = append([]byte{0xFF, 0xFB}, content...)
		case ".wav":
			content = append([]byte{0x52, 0x49, 0x46, 0x46}, content...)
		case ".txt":
			content = append([]byte("这是一个测试文档\n"), content...)
		case ".md":
			content = append([]byte("# 测试文档\n\n"), content...)
		}
	}

	// Pad to size
	currentSize := len(content)
	if currentSize < size {
		paddingSize := size - currentSize
		// Create padding in chunks to avoid large memory allocation for huge files
		// But for simplicity in this test tool, one alloc is okay up to ~100MB
		padding := make([]byte, paddingSize)
		_, _ = rand.Read(padding)
		content = append(content, padding...)
	}

	return os.WriteFile(path, content, 0644)
}

func createValidImage(format string) []byte {
	// Create a small 100x100 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with a color
	c := color.RGBA{100, 150, 200, 255}
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer
	if format == "jpeg" {
		_ = jpeg.Encode(&buf, img, nil)
	} else if format == "png" {
		_ = png.Encode(&buf, img)
	}
	return buf.Bytes()
}

// PrintStats 打印统计信息
func (g *TestDataGenerator) PrintStats() error {
	fmt.Println("\n📊 数据统计:")
	fmt.Println("------------------------------------------------")

	// 统计文件数量和总大小
	var totalSize int64
	var fileCount int64

	err := filepath.WalkDir(g.targetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err == nil {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("统计文件失败: %w", err)
	}

	fmt.Printf("  文件总数: %d\n", fileCount)
	fmt.Printf("  总大小: %.2f MB\n", float64(totalSize)/(1024*1024))
	if fileCount > 0 {
		fmt.Printf("  平均大小: %.2f KB\n", float64(totalSize)/(1024*float64(fileCount)))
	}

	return nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	outDir := flag.String("out", "./test_assets", "输出目录")
	fileCount := flag.Int("count", 100, "生成文件数量 (approx)")
	lean := flag.Bool("lean", false, "精简模式 (生成极小文件以节省空间)")
	flag.Parse()

	// 创建目标目录
	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Printf("❌ 创建目录失败: %v\n", err)
		os.Exit(1)
	}

	absDir, _ := filepath.Abs(*outDir)
	fmt.Printf("📁 目标目录: %s\n", absDir)
	fmt.Printf("🔢 目标数量: ~%d 个文件\n\n", *fileCount)

	// 创建生成器
	generator := NewTestDataGenerator(absDir, *lean)

	// 如果是精简模式，修改模板大小
	if *lean {
		fmt.Println("🍃 精简模式已开启，将生成极小文件...")
	}

	// 生成环境
	if err := generator.GenerateCreatorEnvironment(*fileCount); err != nil {
		fmt.Printf("❌ 生成环境失败: %v\n", err)
		os.Exit(1)
	}

	// 打印统计
	if err := generator.PrintStats(); err != nil {
		fmt.Printf("❌ 统计失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n🎉 测试环境生成成功！")
	fmt.Printf("📁 目录位置: %s\n", absDir)
	fmt.Println("\n这个环境模拟了以下创作者的素材库:")
	fmt.Println("  📸 摄影师 - 作品集、RAW文件、修图素材")
	fmt.Println("  🎬 小博主 - 视频项目、封面设计、文案")
	fmt.Println("  🎨 内容创作者 - 图片、视频、音频、字体、模板")
	fmt.Println("  📊 项目管理 - 小红书、抖音、公众号项目")
	fmt.Println("\n你可以:")
	fmt.Println("  1. 使用这个目录作为测试数据")
	fmt.Println("  2. 运行应用扫描这个目录")
	fmt.Println("  3. 查看应用如何管理这些素材")
}
