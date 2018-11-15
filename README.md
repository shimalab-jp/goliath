# goliath - RESTful API framework for golang

はじめに
------

このフレームワークは、ゲーム等のプレイヤー管理が必要なアプリ向けの RESTful API 開発用フレームワークです。
Go言語のパッケージが 'go' から始まるのが多いのと、先日ラピュタをみてたので「goliath(ゴリアテ)」と名付けました。
アプリ用サーバーに必要最低限の機能は実装していますが、その他に欲しい機能がありましたら、Issueに登録して頂けれれば検討はします。(実装するとは言ってない)


特徴
----

* 基本的なプレイヤー管理の実装(生成、パスワード再発行、端末移動)まで済んでいます。
* プレイヤー生成時に、データベースの振り分け(データベース番号の生成)を自動で行います。
* アカウント作成ログ、HAUログは取得しています。(集計はしていない)
* APIスイッチ機能で、任意のタイミングでAPIを実行停止にすることが出来ます。(テーブルを直接操作する必要があります)
* 多言語対応の為のメッセージ定義機能があります。(クライアントのAccept-Languageに従って返します)
* 作成したAPIを実行、テストするためのダイナミックリファレンス機能があります


実装検討中
--------

* データベースアクセスの改善
* 複数データベースを跨いだトランザクションの実現
* DebugモードでのSQL実行プランの自動取得
* AccountManagerのキャッシュ対応
* エラー処理の改善
* json 以外のフォーマットへの対応
* URLパラメータの改修(現在のQueryString形式からパス方式へ)
* 基本の管理ページ


開発環境
------

* macOS 10.14, Ubuntu 18.04 LTS
* golang 1.11
* nginx 1.14
* memcached 1.5
* mysql 5.7, mariadb 10.2
* chrome/safari

※基本、golangが可動すれば動くはずですが、動作を保証するものではありません。(ちゃんと動作テストしてね)


依存関係
------

goliath パッケージでは、errors, mysql, memcache, yamlパッケージを参照していますので、必要に応じてインストールしてください。
```
$ go get github.com/pkg/errors
$ go get github.com/go-sql-driver/mysql
$ go get github.com/bradfitz/gomemcache/memcache
$ go get gopkg.in/yaml.v2
```


インストール
---------

$GOHOME 直下で下記コマンドを実行し、goliathをインストールしてください。
```
$ go get github.com/shimalab-jp/goliath
```


使い方
------

### main.go
```Go
package main

import (
    "github.com/shimalab-jp/goliath"
    "github.com/shimalab-jp/goliath/example"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/resources/v1"
    "github.com/shimalab-jp/goliath/rest"
)

type ExampleHooks struct {}

// APIの実行前にコールされます
func (hooker *ExampleHooks) PreExecute(engine *rest.Engine, request *rest.Request, response *rest.Response) (error) {
    return nil
}

// APIの実行後にコールされます
func (hooker *ExampleHooks) PostExecute(engine *rest.Engine, request *rest.Request, response *rest.Response) (error) {
    return nil
}

func main() {
    var err error = nil

    // [必須]
    // configファイルのパスを指定してgoliathを初期化してください。
    // ※内部で os.ExpandEnv にてパスを展開していますので、環境変数も指定可能です。
    // ※環境毎にconfigファイルを分ける場合は、ここで実行環境に合ったconfigファイルを指定してください。
    if err == nil {
        err = goliath.Initialize("${GOPATH}/config_local.yaml")
    }

    // [オプション]
    // APIの実行前、実行後に、独自の処理を行いたい場合は、
    // ExecuteHookerを作成する事で、処理をフックする事ができます。
    if err == nil {
        goliath.SetHooks(&ExampleHooks{})
    }

    // [必須]
    // 公開するRESTリソースを追加してください
    // 認証を必要としないAPIのみを公開する場合は、basic resources の account 関連は削除で。
    if err == nil {
        // goliath basic resources(v1)
        if err == nil { err = goliath.AppendResource(1, "/account/regist",   &v1.AccountRegist{}) }
        if err == nil { err = goliath.AppendResource(1, "/account/password", &v1.AccountPassword{}) }
        if err == nil { err = goliath.AppendResource(1, "/account/trans",    &v1.AccountTrans{}) }
        if err == nil { err = goliath.AppendResource(1, "/debug/cache",      &v1.DebugCache{}) }

        // user defined resources(v1)
        if err == nil { err = goliath.AppendResource(1, "/example/example1", &example.Example1{}) }
    }

    // [必須] httpサーバーのListenを開始
    if err == nil {
        err = goliath.Listen()
    }

    if err != nil {
        log.Ee(err)
    }
}
```

### exmaple1.go
POST処理 実装例
```Go
package example

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

// API用の構造体を作成します
type Example1 struct {
    // ResourceBaseを引き継ぐと実装が楽です
    rest.ResourceBase
}

// Defineは必ず実装します。
// ここで、どんなAPIなのか、どういったパラメータを受け取り、チェックするのかを定義します。
// 今回は Type が Int32 に設定していますが、数値として評価出来ない値が渡されると、
// APIの実行前に REST エンジン側でチェックし、自動的に不正なパラメータを防ぎます。
func (res Example1) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/example/diff",
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:       "減算",
                Description:   "減算します",
                UrlParameters: []rest.UrlParameter{},
                PostParameters: map[string]rest.PostParameter{},
                Returns: map[string]rest.Return{
                    "Result": {
                        Type:        reflect.Int32,
                        Description: "減算結果"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

// Defineで定義したメソッドの実装を行います。
// 今回はPOSTなので、Postを実装します。
// パラメータは、 request の Get〜メソッドで取得する事ができます。
// 結果は response の Result に格納してください。
func (res Example1) Post(request *rest.Request, response *rest.Response) (error) {
    // パラメータを取得
    v1 := request.GetPostInt32("Value1", 0)
    v2 := request.GetPostInt32("Value2", 0)

    // 処理
    result := v1 - v2

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"Result": result}

    return nil
}
```

### exmaple2.go
```Go
```


ライセンス
--------

このパッケージのライセンスは、 Apache License 2.0 を適用するものとします。
Apache License 2.0 については、 LICENSE ファイルをご参照ください。

