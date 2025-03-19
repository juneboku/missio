package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// Logger はスキャンの進捗状況を表示するためのロガーです
type Logger struct {
	scannedFiles uint64
	startTime    time.Time
	verbose      bool
	maxDepth     int  // 表示する最大階層数
	rootDir      string // ルートディレクトリ
}

// NewLogger は新しいLoggerインスタンスを作成します
func NewLogger(rootDir string, verbose bool, maxDepth int) *Logger {
	return &Logger{
		startTime: time.Now(),
		verbose:   verbose,
		maxDepth:  maxDepth,
		rootDir:   rootDir,
	}
}

// IncrementScanned はスキャンしたファイル数をインクリメントします
func (l *Logger) IncrementScanned() {
	atomic.AddUint64(&l.scannedFiles, 1)
}

// GetScannedCount はスキャンしたファイル数を返します
func (l *Logger) GetScannedCount() uint64 {
	return atomic.LoadUint64(&l.scannedFiles)
}

// getPathDepth はパスの階層数を計算します
func (l *Logger) getPathDepth(path string) int {
	// ルートディレクトリからの相対パスを取得
	relPath, err := filepath.Rel(l.rootDir, path)
	if err != nil {
		return 0
	}

	// パスセパレータで分割して階層数を計算
	if relPath == "." {
		return 0
	}
	return len(strings.Split(relPath, string(filepath.Separator)))
}

// LogProgress は進捗状況を出力します
func (l *Logger) LogProgress(path string) {
	if !l.verbose {
		return
	}

	depth := l.getPathDepth(path)
	if l.maxDepth >= 0 && depth > l.maxDepth {
		return
	}

	// インデントを追加してログを出力
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s\n", indent, filepath.Base(path))
}

// LogSummary はスキャンのサマリーを出力します
func (l *Logger) LogSummary(secretFiles []string) {
	duration := time.Since(l.startTime)
	fmt.Printf("\nスキャン完了:\n")
	fmt.Printf("  スキャンしたファイル数: %d\n", l.GetScannedCount())
	fmt.Printf("  検出した秘匿ファイル数: %d\n", len(secretFiles))
	fmt.Printf("  処理時間: %v\n", duration.Round(time.Millisecond))
	if len(secretFiles) > 0 {
		fmt.Printf("\n検出された秘匿ファイル:\n")
		for _, file := range secretFiles {
			fmt.Printf("  %s\n", file)
		}
	}
}
