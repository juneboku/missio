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
func NewScanner(rootDir string, verbose bool, maxDepth int) *Scanner {
	// .gitignoreファイルを読み込む
	gitIgnorePath := filepath.Join(rootDir, ".gitignore")
	ignore, _ := ignore.CompileIgnoreFile(gitIgnorePath)

	return &Scanner{
		rootDir: rootDir,
		ignore:  ignore,
		logger:  NewLogger(rootDir, verbose, maxDepth),
	}
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
		if s.isSecretFile(path) {
			files = append(files, relPath) // 相対パスを保存
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
func (s *Scanner) isSecretFile(path string) bool {
	filename := filepath.Base(path)
	ext := filepath.Ext(path)

	// 秘匿ファイルのパターン
	secretPatterns := []string{
		// 環境変数ファイル
		".env",
		".env.",
		".envrc",
		// キーファイル
		".key",
		"id_rsa",
		"id_ecdsa",
		"id_ed25519",
		// 認証情報
		"credentials",
		"credential",
		"secret",
		"secrets",
		"master.key",
		"password",
		"token",
		// AWS関連
		"aws_access_key_id",
		"aws_secret_access_key",
		// その他
		"oauth",
		"private",
	}

	// 除外パターン
	excludePatterns := []string{
		".example",
		".sample",
		".template",
		"test",
		"spec",
		".md",
		".txt",
	}

	// ファイル名を小文字に変換
	lowerFilename := strings.ToLower(filename)

	// 除外パターンをチェック
	for _, pattern := range excludePatterns {
		if strings.Contains(lowerFilename, pattern) {
			return false
		}
	}

	// 秘匿ファイルパターンをチェック
	for _, pattern := range secretPatterns {
		if strings.Contains(lowerFilename, pattern) {
			return true
		}
	}

	// 特定の拡張子をチェック
	secretExts := []string{".pem", ".crt", ".key", ".p12", ".pfx", ".jks", ".keystore"}
	for _, secretExt := range secretExts {
		if ext == secretExt {
			return true
		}
	}

	return false
}
