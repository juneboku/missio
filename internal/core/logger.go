package core

import (
	"fmt"
	"os"
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
	maxDepth     int    // 表示する最大階層数
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

// formatFileSize はファイルサイズを人間が読みやすい形式に変換します
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
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

	// 相対パスを取得
	relPath, err := filepath.Rel(l.rootDir, path)
	if err != nil {
		relPath = path
	}

	// ファイル情報を取得
	info, err := os.Stat(path)
	var sizeStr string
	if err == nil && !info.IsDir() {
		sizeStr = fmt.Sprintf(" (%s)", formatFileSize(info.Size()))
	}

	// インデントを追加してログを出力
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s%s\n", indent, relPath, sizeStr)
}

// LogSummary はスキャンのサマリーを出力します
func (l *Logger) LogSummary(files []string) {
	fmt.Printf("\nスキャン完了:\n")
	fmt.Printf("  スキャンしたファイル数: %d\n", l.scannedFiles)
	fmt.Printf("  検出した秘匿ファイル数: %d\n", len(files))
	fmt.Printf("  処理時間: %dms\n", time.Since(l.startTime).Milliseconds())

	if len(files) > 0 {
		fmt.Printf("\n検出された秘匿ファイル:\n")
		for _, relPath := range files {
			absPath := filepath.Join(l.rootDir, relPath)
			info, err := os.Stat(absPath)
			if err != nil {
				continue
			}
			fmt.Printf("  %s (%s)\n", relPath, formatFileSize(info.Size()))
		}
	}
}
