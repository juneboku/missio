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
    │   ├── repo11
    │   │   └── .env
    │   └── repo12
    │       └── config
    │           ├── credentials
    │           │   └── development.key
    │           └── master.key
    └── org2
        └── repo21
            └── .env.development
```

アプリケーションの秘匿情報を扱うファイルで、バージョン管理されていないものだけを抽出して、destination_dirにコピーします。
