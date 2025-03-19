package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Exporter は秘匿ファイルをエクスポートする機能を提供します
type Exporter struct {
	sourceDir   string
	destDir     string
	logger      *Logger
	copiedFiles int
	totalSize   int64
	startTime   time.Time
}

// NewExporter は新しいExporterインスタンスを作成します
func NewExporter(sourceDir, destDir string, verbose bool, maxDepth int) *Exporter {
	return &Exporter{
		sourceDir: sourceDir,
		destDir:   destDir,
		logger:    NewLogger(sourceDir, verbose, maxDepth),
		startTime: time.Now(),
	}
}

// Export は秘匿ファイルを抽出してコピーします
func (e *Exporter) Export() error {
	// スキャナーを作成
	scanner, err := NewScanner(e.sourceDir, e.logger.verbose, e.logger.maxDepth)
	if err != nil {
		return fmt.Errorf("スキャナーの作成に失敗しました: %w", err)
	}

	// 秘匿ファイルを検出
	files, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("ファイルのスキャンに失敗しました: %w", err)
	}

	// 宛先ディレクトリを作成
	if err := os.MkdirAll(e.destDir, 0755); err != nil {
		return fmt.Errorf("宛先ディレクトリの作成に失敗しました: %w", err)
	}

	fmt.Printf("\nエクスポートを開始します...\n")

	// 各ファイルをコピー
	for _, relPath := range files {
		srcPath := filepath.Join(e.sourceDir, relPath)
		dstPath := filepath.Join(e.destDir, relPath)

		// コピー先のディレクトリを作成
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %s: %w", dstDir, err)
		}

		// ファイルをコピー
		size, err := e.copyFile(srcPath, dstPath)
		if err != nil {
			return fmt.Errorf("ファイルのコピーに失敗しました: %s: %w", relPath, err)
		}

		e.copiedFiles++
		e.totalSize += size

		if e.logger.verbose {
			fmt.Printf("  コピー完了: %s (%s)\n", relPath, formatFileSize(size))
		}
	}

	e.logSummary()
	return nil
}

// copyFile はファイルをコピーし、コピーしたサイズを返します
func (e *Exporter) copyFile(src, dst string) (int64, error) {
	// ソースファイルを開く
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	// コピー先ファイルを作成
	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	// ファイルの権限を保持
	srcInfo, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// ファイルをコピー
	size, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return 0, err
	}

	// 権限を設定
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return 0, err
	}

	return size, nil
}

// logSummary はエクスポートのサマリーを出力します
func (e *Exporter) logSummary() {
	duration := time.Since(e.startTime)
	fmt.Printf("\nエクスポート完了:\n")
	fmt.Printf("  コピーしたファイル数: %d\n", e.copiedFiles)
	fmt.Printf("  合計サイズ: %s\n", formatFileSize(e.totalSize))
	fmt.Printf("  処理時間: %v\n", duration.Round(time.Millisecond))
	fmt.Printf("  出力先: %s\n", e.destDir)
}
