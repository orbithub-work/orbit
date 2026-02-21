package main

import (
	"context"
	"fmt"
	"log"

	"media-assistant-os/internal/db"
	"media-assistant-os/internal/services"
)

func main() {
	targetFile := "/Users/a/Projects/media-assistant-rust/test_assets_100k/摄影师/精选作品/精选照片_21092.jpg"
	assetID := "ce86e461-51e7-45d5-afee-a6f7bfc07f9a"

	// 1. 计算指纹
	fp, _, _, _, err := services.ComputeAdaptiveFingerprint(targetFile)
	if err != nil {
		log.Fatalf("Failed to compute fingerprint: %v", err)
	}
	fmt.Printf("Computed fingerprint: %s\n", fp)

	// 2. 更新数据库
	d, err := db.Open("/Users/a/Projects/media-assistant-rust/data")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer d.Close()

	_, err = d.ORM().NewUpdate().
		Table("assets").
		Set("fingerprint = ?", fp).
		Set("status = ?", "READY").
		Where("id = ?", assetID).
		Exec(context.Background())

	if err != nil {
		log.Fatalf("Failed to update DB: %v", err)
	}

	fmt.Println("Successfully updated asset with fingerprint and set status to READY")
}
