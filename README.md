# Missio

Missio(密書)は、アプリケーションの秘匿情報ファイルを一括エクスポート・インポートするためのコマンドラインツールです。

## 使い方

```sh
$ missio export original_dir destination_dir
```

## 動作

たとえばこのようなディレクトリ構造とファイル配置になっているorigin_dirがあるとして

```
original_dir
└── github.com
    ├── org1
    │   ├── repo11
    │   │   └── .env
    │   └── repo12
    │       └── config
    │           ├── credentials
    │           │   └── development.key
    │           └── master.key
    └── org2
        └── repo21
            └── .env.development
```

アプリケーションの秘匿情報を扱うファイルで、バージョン管理されていないものだけを抽出して、destination_dirにコピーします。

## 設定ファイル

`missio.yml`ファイルを使用して、抽出するファイルと除外するファイルを設定できます。
設定ファイルが存在しない場合は、デフォルトの設定が使用されます。

サンプル設定ファイル`missio.example.yml`を参考に、プロジェクトに合わせた設定を作成してください。

```yaml
# 抽出するファイルパターン
include:
  # ファイル名のパターン
  names:
    - .env
    - .env.
    - .envrc
    - master.key
    - credentials.yml.enc
    - id_rsa
    - id_ecdsa
    - id_ed25519

  # 拡張子のパターン
  extensions:
    - .key
    - .pem
    - .crt
    - .p12
    - .pfx
    - .jks
    - .keystore

  # パスのパターン
  paths:
    - config/credentials/*
    - config/master.key

# 除外するファイルパターン
exclude:
  # ファイル名のパターン
  names:
    - .env.example
    - .env.sample
    - .env.template
    - example
    - sample
    - template
    - test
    - spec

  # 拡張子のパターン
  extensions:
    - .md
    - .txt
    - .example
    - .sample
    - .template

  # パスのパターン
  paths:
    - test/*
    - spec/*
    - examples/*
```

### パターンの説明

- `names`: ファイル名に含まれる文字列を指定します（大文字小文字は区別しません）
- `extensions`: ファイルの拡張子を指定します（先頭のドットを含めて指定）
- `paths`: ファイルパスのパターンを指定します（`*`でワイルドカード指定可能）

設定ファイルでは、デフォルトの設定を上書きまたは追加できます。
設定ファイルで指定されていない項目は、デフォルトの設定が使用されます。
