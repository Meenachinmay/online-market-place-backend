# Soda Interview Service (日本語訳)

Go言語で構築された包括的なマイクロサービスベースのバックエンドシステムです。紹介ブログや報酬機能を備えたECプラットフォームを実装しており、クリーンアーキテクチャ、gRPC通信、PostgreSQLを使用したトランザクション整合性を実証しています。

## 🚀 概要 (Overview)

Soda Interview Serviceは、「Soda Store」というオンラインストアのバックエンドをシミュレートしています。ユーザーは紹介ブログを通じて商品を発見し、購入することができます。このシステムは、購入者とコンテンツクリエイター（ブログ著者）の両方に、デジタル通貨（「Soda Balance」）に換金可能な報酬ポイント（「Soda Points」）を付与することでインセンティブを与えます。

## ✨ 主な機能 (Key Features)

### 1. 商品管理 (Product Management)
- 価格や報酬設定を含む商品詳細の取得。
- 利用可能な商品の一覧表示。

### 2. 紹介ブログシステム (Referral Blog System)
- 著者は特定の商品にリンクしたブログを作成可能。
- コンテンツと商品の関係を追跡し、売上を帰属させます。

### 3. 注文処理と報酬 (Order Processing & Rewards)
- **注文処理**: 購入者、商品、紹介ブログを紐付け、安全に注文を処理します。
- **購入者への報酬**: 購入者は、特定の商品を**初めて購入**した際に**Soda Points**を獲得します。
- **著者への報酬**: ブログ著者は、自身の紹介ブログを通じて売上が発生するたびに**Soda Points**を獲得します。

### 4. Sodaファイナンス / ウォレット (Soda Finance)
- **ウォレット管理**: ユーザーごとのウォレットを自動的に作成・維持します。
- **ポイントシステム**: 獲得したSoda Pointsを追跡します。
- **通貨換算**: ユーザーはSoda PointsをSoda Balance（円）に換金できます。
  - **換算レート**: 2ポイント = 1円。
  - **しきい値**: **1000 Soda Points**以上を保有している場合のみ換算可能です。

## 🏗 アーキテクチャ (Architecture)

このプロジェクトは、関心の分離と保守性を確保するために**クリーンアーキテクチャ (Clean Architecture) / ヘキサゴナルアーキテクチャ (Hexagonal Architecture)**を採用しています。

### 1. プロトコル層 (`foundation/proto`)
**Protocol Buffers (Protobuf)** を使用してAPI規約を定義します。
- **Order Service**: `PlaceOrder`
- **Product Service**: `GetProduct`, `ListProducts`
- **Blog Service**: `CreateBlog`, `GetBlog`
- **Finance Service**: `GetWallet`, `ConvertPoints`

### 2. トランスポート層 (`app/services/soda-interview-grpc`)
gRPCサーバーの実装を含みます (`internal/transport/grpc`)。
- **Handlers**: ビジネスエンティティへのマッピングやコアロジックの呼び出しを行う具体的な実装（例: `order/handlers.go`）。

### 3. ビジネスコア (`business/core`)
純粋なビジネスロジックを含むアプリケーションの中核です。
- **Order Core**: トランザクションの調整、初回購入の判定、報酬配布のトリガーなどを処理します。
- **Finance Core**: 換算ルール（しきい値、レート）の適用やウォレットの状態管理を行います。

### 4. データ層 (`business/data`)
データベースとのやり取りを処理します。
- **Schema**: `goose` で管理されるPostgreSQLマイグレーション。
- **Stores**: `sqlc` で生成された型安全なSQLクエリ。
- **Transactional Support**: `pgxpool` を使用し、複数のドメインストアにまたがる堅牢なトランザクション管理を実現します。

## 🛠 技術スタック (Tech Stack)

- **言語**: [Go](https://go.dev/) (Golang)
- **通信**: [gRPC](https://grpc.io/)
- **データベース**: [PostgreSQL](https://www.postgresql.org/)
- **ドライバ**: [jackc/pgx](https://github.com/jackc/pgx)
- **SQL生成**: [sqlc](https://sqlc.dev/)
- **マイグレーション**: [goose](https://github.com/pressly/goose)
- **設定**: [Viper](https://github.com/spf13/viper)
- **コンテナ化**: [Docker](https://www.docker.com/) & Docker Compose

---
*Soda Interview Project用に生成されました。*
