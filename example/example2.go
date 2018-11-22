package example

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

// API用の構造体を作成します
type Example2 struct {
    // ResourceBaseを引き継ぐと実装が楽です
    rest.ResourceBase
}

func (res Example2) Define() *rest.ResourceDefine {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "GET": {
                Summary: "減算(UrlParameters版)",
                Description: "UrlParametersでパラメータを受け取り、減算します。\n" +
                    "UrlParametersは、以下の特徴があります。\n" +
                    "* 名前ではなく、「API名の後の何番目(Index)の引数」という形式で取得します。\n" +
                    "* Indexは0から順に割り振られます。\n" +
                    "* このAPI例では、GETにて /example/example2/5/3 というアクセスを行うと、5 - 3 が処理され、2が返ります。\n" +
                    "* Requireは最後のIndexの値だけを指定する事が可能です。\n" +
                    "* また、最後のIndexの値には IsMultiple を指定する事ができます。\nIsMultiple をtrueに設定すると/example/example2/5/3/1/1 様に最後のパラメータを複数指定出来るようになります。\n" +
                    "* UrlParameters を定義しない場合は、値の妥当性チェックはされませんが、実行時に取得する事ができます。\n\n" +
                    "なお、現在リファレンス上では IsMultiple パラメータの入力には対応していません。",
                UrlParameters: []rest.UrlParameter{
                    {
                        Type:        reflect.Int32,
                        Description: "左の値",
                        Require:     true},
                    {
                        Type:        reflect.Int32,
                        Description: "右の値",
                        Require:     true,
                        IsMultiple:  true}},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters: map[string]rest.PostParameter{},
                Returns: map[string]rest.Return{
                    "Result": {
                        Type:        reflect.Int32,
                        Description: "減算結果"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Example2) Get(request *rest.Request, response *rest.Response) error {
    var result int32
    var val int32
    var next = true
    var index = 0

    // パラメータを取得しつつ計算
    result, next = request.GetUrlParamInt32(index, 0)
    for next {
        index++
        val, next = request.GetUrlParamInt32(index, 0)
        if next {
            result -= val
        }
    }

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"Result": result}

    return nil
}
