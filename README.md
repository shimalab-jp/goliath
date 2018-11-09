# 'goliath' ゴリアテ

はじめに
------

このフレームワークは、ゲーム等のプレイヤー管理が必要なアプリ向けの RESTful API 開発用フレームワークです。
Go言語のパッケージが 'go' から始まるのが多いのと、先日ラピュタをみてたので「goliath(ゴリアテ)」と名付けました。
このフレームワークで、アクセスしてくる人を「見ろ！人がゴミのようだ！」と言えるように使い倒して頂けると作った甲斐があります。


特徴
----

* 基本的なプレイヤー管理の実装(生成、パスワード再発行、端末移動)まで済んでいます。
* プレイヤー生成時に、データベースの振り分け(データベース番号の生成)を自動で行います。
* アカウント作成ログ、HAUログが取られています。
* APIスイッチ機能で、任意のタイミングでAPIを実行停止にすることが出来ます。
* 多言語対応の為のメッセージ定義機能があります。(クライアントのAccept-Languageに従って返す)


今後予定している更新
---------------

* データベースアクセスの改善
* 複数データベースを跨いだトランザクションの実現
* DebugモードでのSQL実行プランの自動取得
* エラー処理の改善
* json 以外のフォーマットへの対応


開発環境
---------

* macOS 10.14, Ubuntu 18.04 LTS
* golang 1.11
* nginx 1.14
* memcached 1.5
* mysql 5.7, mariadb 10.2

※基本、golangが可動すれば動くはずですが、動作を保証するものではありません。


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

```Go
package main

import (
        "fmt"
        "log"

        "github.com/shimalab-jp/goliath"
)

// PreExecute, PostExecute を下記の定義て実装する事により、ExecutionHooks となります。
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
    // ユーザープログラムの初期化処理など、必要であればこの辺で実行しておくと良いです。
    if err == nil {
        // todo: user initialization
    }

    // [必須]
    // ユーザーのRESTリソースを追加してください。
    // ※システムリソースや、ユーザーリソースとパスが被る場合は、起動時にエラーとなりますのでご注意ください。
    if err == nil {
        err = goliath.AppendResource("", nil)
    }

    // [オプション]
    // APIの実行前、実行後に、独自の処理を行いたい場合は、
    // ExecutionHooker を作成する事で、処理をフックする事ができます。
    if err == nil {
        goliath.SetHooks(&ExampleHooks{})
    }

    // [必須]
    // Listenを開始してください。
    if err == nil {
        err = goliath.Listen()
    }

    if err != nil {
        log.Ee(err)
    }
}
```


ライセンス
--------

このパッケージは、 Apache License 2.0 の下に公開しています。
Apache License 2.0 については、 LICENSE ファイルをご参照ください。


