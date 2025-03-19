package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/junebako/missio/internal/core"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		handleList()
	case "export":
		handleExport()
	default:
		fmt.Printf("未知のコマンド: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("使用方法: missio <command> [options] <directory>")
	fmt.Println("\nコマンド:")
	fmt.Println("  list    秘匿ファイルを一覧表示")
	fmt.Println("  export  秘匿ファイルをエクスポート")
	fmt.Println("\nオプション:")
	fmt.Println("  -verbose  詳細な出力を表示")
	fmt.Println("  -depth n  表示する最大階層数（デフォルト: 2, -1で制限なし）")
}

func handleList() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	verbose := listCmd.Bool("verbose", false, "詳細な出力を表示")
	maxDepth := listCmd.Int("depth", 2, "表示する最大階層数（-1で制限なし）")

	listCmd.Parse(os.Args[2:])
	args := listCmd.Args()
	if len(args) < 1 {
		fmt.Println("ディレクトリを指定してください")
		os.Exit(1)
	}

	scanner, err := core.NewScanner(args[0], *verbose, *maxDepth)
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}

	_, err = scanner.Scan()
	if err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}

func handleExport() {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	verbose := exportCmd.Bool("verbose", false, "詳細な出力を表示")
	maxDepth := exportCmd.Int("depth", 2, "表示する最大階層数（-1で制限なし）")

	exportCmd.Parse(os.Args[2:])
	args := exportCmd.Args()
	if len(args) < 2 {
		fmt.Println("使用方法: missio export [options] <source_dir> <dest_dir>")
		os.Exit(1)
	}

	sourceDir := args[0]
	destDir := args[1]

	exporter := core.NewExporter(sourceDir, destDir, *verbose, *maxDepth)
	if err := exporter.Export(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}
