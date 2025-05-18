# ベースイメージとして公式のGolangイメージを使用
FROM golang:1.23-alpine

# 作業ディレクトリを設定
WORKDIR /go/src

# Go Modulesの使用を許可
ENV GO111MODULE=on

# ローカルのモジュールキャッシュを最適化
COPY go.mod .
COPY go.sum .
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# Swaggerの生成
# Swaggerのインストール
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7
# ポート8080を公開
EXPOSE 8080

# ホットリロード系
RUN go install github.com/air-verse/air@latest

# 実行可能ファイルをデフォルトのコマンドとして設定
CMD ["air"]