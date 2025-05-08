package main

import (
	"fmt"
	"log"
	"os"
	"proto2gorm/parser"
	"proto2gorm/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: proto2gorm path/to/file.proto")
		return
	}
	protoPath := os.Args[1]
	fmt.Println("🚀 Parsing .proto file...")
	err := parser.ParseProtoToSQL(protoPath)
	if err != nil {
		log.Fatalf("❌ Failed to parse proto: %v", err)
	}
	fmt.Println("🧩 Executing SQL in SQLite and generating GORM models...")
	if err := runner.ExecuteAndGenerate(); err != nil {
		log.Fatalf("❌ Failed to generate models: %v", err)
	}
	fmt.Println("✅ All done! Models saved in ./gen")
}
