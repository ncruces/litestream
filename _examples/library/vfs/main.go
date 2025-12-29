// Example: Basic Litestream Library Usage
//
// This example demonstrates the simplest way to use Litestream as a Go library.
// It replicates a SQLite database to the local filesystem.
//
// Run: go run main.go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/litestream"
	"github.com/ncruces/litestream/s3"
)

func main() {
	// Load configuration from environment
	bucket := os.Getenv("LITESTREAM_BUCKET")
	if bucket == "" {
		log.Fatal("LITESTREAM_BUCKET environment variable required")
	}
	path := os.Getenv("LITESTREAM_PATH")
	if path == "" {
		path = "litestream"
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	// 1. Create S3 replica client
	client := s3.NewReplicaClient()
	client.Bucket = bucket
	client.Path = path
	client.Region = region
	client.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	client.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	// 2. Create VFS replica client
	litestream.NewVFS(path, client, litestream.VFSOptions{
		PollInterval: 5 * time.Second,
	})

	// 3. Open your app's SQLite connection for reading
	db, err := driver.Open(fmt.Sprintf("file:%s?vfs=litestream", path))
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// 4. Open your app's SQLite connection for reading
	for {
		time.Sleep(time.Second)
		rows, err := db.Query("SELECT * FROM events")
		if err != nil {
			log.Fatalln(err)
		}

		for rows.Next() {
			var message string
			var created time.Time
			err := rows.Scan(&message, &created)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(message, created)
		}

		log.Println("===")
		rows.Close()
	}
}
