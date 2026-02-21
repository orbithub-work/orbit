package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "modernc.org/sqlite"
)

type Asset struct {
	bun.BaseModel `bun:"table:assets"`
	ID            string `bun:",pk"`
	Status        string
}

func main() {
	sqldb, err := sql.Open("sqlite", "file:data/db/media_assistant.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())

	var results []struct {
		Status string
		Count  int
	}

	err = db.NewSelect().
		Model((*Asset)(nil)).
		ColumnExpr("status, count(*) as count").
		Group("status").
		Scan(context.Background(), &results)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Asset Status Distribution:")
	for _, r := range results {
		fmt.Printf("%s: %d\n", r.Status, r.Count)
	}

	var taskResults []struct {
		Status string
		Count  int
	}
	err = db.NewSelect().
		Table("media_tasks").
		ColumnExpr("status, count(*) as count").
		Group("status").
		Scan(context.Background(), &taskResults)

	if err == nil {
		fmt.Println("\nTask Status Distribution:")
		for _, r := range taskResults {
			fmt.Printf("%s: %d\n", r.Status, r.Count)
		}
	}

	var indexedAssets []struct {
		Path      string
		MediaMeta string
	}
	err = db.NewSelect().
		Model((*Asset)(nil)).
		Column("path", "media_meta").
		Where("status = ?", "INDEXED").
		Limit(1).
		Scan(context.Background(), &indexedAssets)

	if err == nil && len(indexedAssets) > 0 {
		fmt.Println("\nSample INDEXED Asset Meta:")
		fmt.Println(indexedAssets[0].Path)
		fmt.Println(indexedAssets[0].MediaMeta)
	}
}
