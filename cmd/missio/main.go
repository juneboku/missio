package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/junebako/missio/internal/core"
)

func main() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	verbose := listCmd.Bool("verbose", false, "詳細な出力を表示")
	maxDepth := listCmd.Int("depth", 2, "表示する最大階層数（-1で制限なし）")

	if len(os.Args) < 2 {
		fmt.Println("使用方法: missio <command> [options] <directory>")
		fmt.Println("\nコマンド:")
		fmt.Println("  list    秘匿ファイルを一覧表示")
		fmt.Println("\nオプション:")
		fmt.Println("  -verbose  詳細な出力を表示")
		fmt.Println("  -depth n  表示する最大階層数（デフォルト: 2, -1で制限なし）")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		args := listCmd.Args()
		if len(args) < 1 {
			fmt.Println("ディレクトリを指定してください")
			os.Exit(1)
		}

		scanner := core.NewScanner(args[0], *verbose, *maxDepth)
		_, err := scanner.Scan()
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("未知のコマンド: %s\n", os.Args[1])
		os.Exit(1)
	}
}
