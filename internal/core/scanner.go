package core

import (
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Scanner は指定されたディレクトリ内の秘匿ファイルを検出します
type Scanner struct {
	rootDir string
	ignore  *ignore.GitIgnore
	logger  *Logger
	config  *Config
}

// 除外するディレクトリ名のリスト
var excludeDirs = []string{
	".git",          // Gitリポジトリ
	"node_modules",  // Node.js依存関係
	"vendor",        // 依存関係（Go, PHP, Ruby等）
	"tmp",           // 一時ファイル
	"cache",         // キャッシュファイル
	"log",           // ログファイル
	"logs",          // ログファイル
	"coverage",      // テストカバレッジ
	"dist",          // ビルド成果物
	"build",         // ビルド成果物
	".bundle",       // Bundler
	".gradle",       // Gradle
	".idea",         // IntelliJ IDEA
	".vscode",       // Visual Studio Code
	"__pycache__",   // Python
	".pytest_cache", // Python
}

// NewScanner は新しいScannerインスタンスを作成します
func NewScanner(rootDir string, verbose bool, maxDepth int) (*Scanner, error) {
	// .gitignoreファイルを読み込む
	gitIgnorePath := filepath.Join(rootDir, ".gitignore")
	ignore, _ := ignore.CompileIgnoreFile(gitIgnorePath)

	// 設定ファイルを読み込む
	config, err := LoadConfig(rootDir)
	if err != nil {
		return nil, err
	}

	return &Scanner{
		rootDir: rootDir,
		ignore:  ignore,
		logger:  NewLogger(rootDir, verbose, maxDepth),
		config:  config,
	}, nil
}

// Scan はディレクトリを走査し、秘匿ファイルのリストを返します
func (s *Scanner) Scan() ([]string, error) {
	var files []string

	err := filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 除外ディレクトリをスキップ
			for _, dir := range excludeDirs {
				if info.Name() == dir {
					return filepath.SkipDir
				}
			}
			s.logger.LogProgress(path)
			return nil
		}

		s.logger.IncrementScanned()
		s.logger.LogProgress(path)

		// パスを相対パスに変換
		relPath, err := filepath.Rel(s.rootDir, path)
		if err != nil {
			return err
		}

		// gitignoreでマッチしないファイルはスキップ
		if s.ignore != nil && !s.ignore.MatchesPath(relPath) {
			return nil
		}

		// 秘匿ファイルかどうかをチェック
		if s.isSecretFile(relPath) {
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.LogSummary(files)
	return files, nil
}

// isSecretFile はファイルが秘匿情報を含むかどうかを判定します
func (s *Scanner) isSecretFile(relPath string) bool {
	filename := filepath.Base(relPath)
	ext := filepath.Ext(relPath)
	lowerFilename := strings.ToLower(filename)
	lowerPath := strings.ToLower(relPath)

	// 除外パターンをチェック
	for _, pattern := range s.config.Exclude.Names {
		if strings.Contains(lowerFilename, strings.ToLower(pattern)) {
			return false
		}
	}

	for _, pattern := range s.config.Exclude.Extensions {
		if strings.EqualFold(ext, pattern) {
			return false
		}
	}

	for _, pattern := range s.config.Exclude.Paths {
		if matched, _ := filepath.Match(pattern, lowerPath); matched {
			return false
		}
	}

	// 秘匿ファイルパターンをチェック
	for _, pattern := range s.config.Include.Names {
		if strings.Contains(lowerFilename, strings.ToLower(pattern)) {
			return true
		}
	}

	for _, pattern := range s.config.Include.Extensions {
		if strings.EqualFold(ext, pattern) {
			return true
		}
	}

	for _, pattern := range s.config.Include.Paths {
		if matched, _ := filepath.Match(pattern, lowerPath); matched {
			return true
		}
	}

	return false
}
