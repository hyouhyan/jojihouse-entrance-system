# ベースイメージとして公式のGolangイメージを使用
FROM golang:1.24-alpine

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

# ポート8080を公開
EXPOSE 8080

# 実行可能ファイルをデフォルトのコマンドとして設定
CMD ["go", "run", "cmd/main.go"]