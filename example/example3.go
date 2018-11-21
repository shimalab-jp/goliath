package example

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

// API用の構造体を作成します
type Example3 struct {
    // ResourceBaseを引き継ぐと実装が楽です
    rest.ResourceBase
}

func (res Example3) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "GET": {
                Summary: "減算(QueryParameters版)",
                Description: "QueryParametersでパラメータを受け取り、減算します。\n" +
                    "QueryParametersは、以下の特徴があります。\n" +
                    "* POSTと同じ様に名前で値を取得します。\n" +
                    "* パラメータがURLに含まれるので、多くのWebサーバーのログに残ります。\n" +
                    "* URLでパラメータを渡すので、複雑だったり大きなデータ送信には向きません。\n",
                UrlParameters: []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{
                    "Value1": {
                        Type:        reflect.Int32,
                        Description: "左の値",
                        Require:     true},
                    "Value2": {
                        Type:        reflect.Int32,
                        Description: "右の値",
                        Require:     true}},
                PostParameters: map[string]rest.PostParameter{},
                Returns: map[string]rest.Return{
                    "Result": {
                        Type:        reflect.Int32,
                        Description: "減算結果"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Example3) Get(request *rest.Request, response *rest.Response) (error) {
    // パラメータを取得
    v1, _ := request.GetParamInt32(rest.QueryParam, "Value1", 0)
    v2, _ := request.GetParamInt32(rest.QueryParam, "Value2", 0)

    // 処理
    result := v1 - v2

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"Result": result}

    return nil
}
