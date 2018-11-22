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
func (res Example1) Define() *rest.ResourceDefine {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:         "減算(PostParameters版)",
                Description: "PostParametersでパラメータを受け取り、減算します。\n" +
                    "PostParametersは、以下の特徴があります。\n" +
                    "* jsonで送信するので、jsonで表現できる配列や構造体を送信する事ができます。\n" +
                    "* 大きなサイズのメッセージを送る事ができます。\n" +
                    "* POST送信なので、一般的なWebサーバーでは送信内容はログに残りません。\n",
                UrlParameters:   []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters: map[string]rest.PostParameter{
                    "Value1": {
                        Type:        reflect.Int32,
                        Description: "左の値",
                        Require:     true},
                    "Value2": {
                        Type:        reflect.Int32,
                        Description: "右の値",
                        Require:     true}},
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
func (res Example1) Post(request *rest.Request, response *rest.Response) error {
    // パラメータを取得
    v1, _ := request.GetParamInt32(rest.PostParam, "Value1", 0)
    v2, _ := request.GetParamInt32(rest.PostParam, "Value2", 0)

    // 処理
    result := v1 - v2

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"Result": result}

    return nil
}
