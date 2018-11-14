;(function () {

    'use strict';

    /** golang types **/
    const TYPE_BOOL    =  1;
    const TYPE_INT     =  2;
    const TYPE_INT8    =  3;
    const TYPE_INT16   =  4;
    const TYPE_INT32   =  5;
    const TYPE_INT64   =  6;
    const TYPE_UINT    =  7;
    const TYPE_UINT8   =  8;
    const TYPE_UINT16  =  9;
    const TYPE_UINT32  = 10;
    const TYPE_UINT64  = 11;
    const TYPE_FLOAT32 = 13;
    const TYPE_FLOAT64 = 14;
    const TYPE_ARRAY   = 17;
    const TYPE_MAP     = 21;
    const TYPE_STRING  = 24;


    const HTTP_OK = 200;  // 正常終了
    const HTTP_BAD_REQUEST = 400;  // 不正なリクエスト
    const HTTP_UNAUTHORIZED = 401;  // 認証されていない
    const HTTP_PAYMENT_REQUIRED = 402;  // 課金されてない
    const HTTP_FORBIDDEN = 403;  // APIの実行権限がない
    const HTTP_NOT_FOUND = 404;  // APIが見つからない
    const HTTP_NOT_ACCEPTABLE = 406;  // 受理出来無いリクエスト（一般的なコマンドエラー）
    const HTTP_CONFLICT = 409;  // 競合エラー（報酬を既に受領済みとか）
    const HTTP_INTERNAL_SERVER_ERROR = 500;  // 内部エラー
    const HTTP_NOT_IMPLEMENTED = 501;  // 未実装
    const HTTP_SERVICE_UNAVAILABLE = 503;  // メンテナンス中
    const RF_RESULT_USER_ERROR = 700;  // ユーザーエラー
    const RF_RESULT_SYSTEM_ERROR = 800;  // システムエラー
    const RF_RESULT_FATAL_ERROR = 900;  // 致命的なエラー

    const MAX_WEB_SERVER_PROCESSES = 500;     // 1サーバ当たりのWebサーバーのプロセス数
    const MAX_FRONT_MEMORY = 4294967296;    // 搭載メモリ
    const FRONT_SYSTEM_USAGE = 2147483648;    // システム使用メモリ
    const LIMIT_MEMORY = (MAX_FRONT_MEMORY - FRONT_SYSTEM_USAGE) / MAX_WEB_SERVER_PROCESSES;  // php1プロセス当たりが利用可能なメモリ
    const WARNING_MEMORY = LIMIT_MEMORY * 0.7;    // 警告を出力するメモリ量

    let _content_data = [];
    let _ssid = "";

    let util = {
        /**
         * 文字列sから、指定文字cをトリムします。
         * @param s 文字列
         * @param c トリム文字
         * @returns string
         */
        trim_char: function (s, c) {
            return s.replace(new RegExp("^" + c + "+|" + c + "+$", "g"), '');
        }
    };

    let reference = {
        /**
         * リファレンスデータ
         */
        _data: null,

        /**
         * 既定のユーザーエージェント
         */
        _default_user_agent: "",

        /**
         * エラーコンテンツを表示します。
         * @param s エラー内容
         */
        set_error_content: function (s) {
            let html = "";
            html += "<div class=\"contents\"><h1>Error</h1>";
            html += "Failed to load index for reference.<br />Please confirm the following response data. ";
            html += "<h2>Response Data</h2>";
            html += "<div class=\"contents\">";
            html += s;
            html += "</div></div>";
            $("#content_frame").html(html);
        },

        /**
         * リサイズします。
         */
        resize: function () {
            let client_height = window.innerHeight ? window.innerHeight : document.documentElement.clientHeight;
            let client_width = window.innerWidth ? window.innerWidth : document.documentElement.clientWidth;

            let header_height = $("#header").outerHeight();
            let left_menu_width = $("#left_menu").outerWidth();

            let new_height = client_height - header_height;
            if (new_height < 5) new_height = 5;

            let new_width = client_width - left_menu_width;
            if (new_width < 5) new_width = 5;

            $("#left_menu").css({
                "height": sprintf("%dpx", new_height)
            });
            $("#content").css({
                "height": sprintf("%dpx", new_height),
                "width": sprintf("%dpx", new_width)
            });
        },

        /**
         * 指定パス、メソッドのリソース情報を取得します
         * @param path string
         * @param method string
         * @returns object
         */
        get_resource: function (path, method) {
            if (!reference._data) {
                reference.set_error_content("Resources data is null.");
                return null;
            }

            for (let group in reference._data.Resources) {
                let group_define = reference._data.Resources[group];
                for (let p in group_define) {
                    if (path != p) continue;
                    for (let m in group_define[p].Methods) {
                        if (method != m) continue;
                        return group_define[p].Methods[method];
                    }
                }
            }

            reference.set_error_content(method + " " + path + " is not found.");
            return null;
        },

        /**
         * golangの型コードから型名称を取得します。
         * @param type_code　golangの型コード
         * @returns string 型名称
         */
        get_type_name: function (type_code) {
            switch (type_code) {
                case TYPE_BOOL:    return "bool";
                case TYPE_INT:     return "int";
                case TYPE_INT8:    return "int8";
                case TYPE_INT16:   return "int16";
                case TYPE_INT32:   return "int32";
                case TYPE_INT64:   return "int64";
                case TYPE_UINT:    return "uint";
                case TYPE_UINT8:   return "uint8";
                case TYPE_UINT16:  return "uint16";
                case TYPE_UINT32:  return "uint32";
                case TYPE_UINT64:  return "uint64";
                case TYPE_FLOAT32: return "float32";
                case TYPE_FLOAT64: return "float64";
                case TYPE_ARRAY:   return "array";
                case TYPE_MAP:     return "map";
                case TYPE_STRING:  return "string";
                default:           return "unkown";
            }
        },

        /**
         * golangの型コードが数値を表しているかどうかを取得します。
         * @param type_code
         * @returns {boolean} 型コードが数値を表している場合はtrue、それ以外の場合はfalse。
         */
        is_numeric_type: function (type_code) {
            switch (type_code) {
                case TYPE_BOOL:    return false;
                case TYPE_INT:     return true;
                case TYPE_INT8:    return true;
                case TYPE_INT16:   return true;
                case TYPE_INT32:   return true;
                case TYPE_INT64:   return true;
                case TYPE_UINT:    return true;
                case TYPE_UINT8:   return true;
                case TYPE_UINT16:  return true;
                case TYPE_UINT32:  return true;
                case TYPE_UINT64:  return true;
                case TYPE_FLOAT32: return true;
                case TYPE_FLOAT64: return true;
                case TYPE_ARRAY:   return false;
                case TYPE_MAP:     return false;
                case TYPE_STRING:  return false;
                default:           return false;
            }
        },

        get_condition: function (param) {
            let ret = "";

            if (param !== null && param.Range && param.Range.length > 0) {
                if (ret.length > 0) ret += " & ";
                ret += "Range( " + param.Range[0] + "~" + param.Range[param.Range.length - 1] + " )";
            }

            if (param !== null && param.Select && param.Select.length > 0) {
                if (ret.length > 0) ret += " & ";
                ret += "Select( ";
                for (let i in param.Select) {
                    if (i > 0) ret += ", ";
                    ret += param.Select[i];
                }
                ret += " )";
            }

            if (param !== null && param.Regex && param.Regex.length > 0) {
                if (ret.length > 0) ret += " & ";
                ret += "RegEx( " + param.Regex + " )";
            }

            if (ret.length <= 0) {
                ret = "-";
            }

            return ret;
        },

        get_paramater_value: function (param) {
            if (param && param.Default) {
                return "" + param.Default;
            }
            else if (param && reference.is_numeric_type(param.Type)) {
                return "0";
            }
            else {
                return "";
            }
        },

        show_reference: function (e) {
            let id = "#" + e.currentTarget.id;
            let method = $(id).attr('method');
            let path = $(id).attr('path');
            let version = $(id).attr('version');

            let info = reference.get_resource(path, method);
            if (!info) return;

            let html = '<div class="contents"><h1 class="api_title">' + method + "&nbsp;" + path + "</h1>";
            html += '<input type="hidden" name="display_name" id="display_name" value="' + path + '" />';
            html += '<input type="hidden" name="resource_name" id="resource_name" value="' + path + '" />';

            // summary
            {
                html += '<h2>Summary</h2>';
                html += '<div class="contents">';
                html += info.Summary;
                if (info.IsDebugModeOnly) {
                    html += '<br />';
                    html += '<span class="warning">このAPIはデバッグ中でのみ実行可能です。</span>';
                }
                if (info.IsAdminModeOnly) {
                    html += '<br />';
                    html += '<span class="warning">このAPIは管理者ログイン中のみ実行可能です。</span>';
                }
                html += "</div>";
            }

            // description
            if (info.Description && info.Description.length > 0) {
                html += '<h2>Description</h2>';
                html += '<div class="contents">' + info.Description + '</div>';
            }

            // API info
            {
                html += '<h2>API Information</h2>';
                html += '<div class="contents">';
                html += '<table class="ui-widget list">';
                html += '<tr>';
                html += '<th class="ui-widget-header header">Name</th>';
                html += '<th class="ui-widget-header header">Value</th>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Require Authentication</td>';
                html += '<td class="ui-widget-content value">' + info.RequireAuthentication + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Is Admin Mode Only</td>';
                html += '<td class="ui-widget-content value">' + info.IsAdminModeOnly + "</td>";
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Is Debug Mode Only</td>';
                html += '<td class="ui-widget-content value">' + info.IsDebugModeOnly + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Run in Maintenance</td>';
                html += '<td class="ui-widget-content value">' + info.RunInMaintenance + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Test form URL</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="url" id="url" value="' + "" + '" /></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Method</td>';
                html += '<td class="ui-widget-content value">' + method + "<input type=\"hidden\" name=\"method\" id=\"method\" value=\"" + method + "\" /></td>";
                html += '</tr>';

                let user_agent = $.cookie("user_agent");
                if (!user_agent || user_agent == "null") user_agent = reference._default_user_agent;
                html += '<tr>';
                html += '<td class="ui-widget-content name">User Agent</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="user_agent" id="user_agent" value="' + user_agent + "\" /></td>";
                html += '</tr>';
                html += '</table>';
                html += '</div>';
            }

            // Authentication
            if (info.RequireAuthentication) {
                let ssid = $.cookie("ssid");
                if (!ssid) ssid = "";
                html += '<h2>Authentication</h2>';
                html += '<div class="contents">';
                html += '<div class="notice">このAPIは認証が要求されています。</div>';
                html += '<div class="notice">※セッションIDは、\'account/regist\'API、または\'account/auth\'APIを実行すると自動的に入力されます。</div>';
                html += '<table class="ui-widget list">';
                html += '<tr>';
                html += '<th class="ui-widget-header header">Name</th>';
                html += '<th class="ui-widget-header header">Value</th>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Session ID</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="ssid" id="ssid" value="' + ssid + "\" /></td>";
                html += '</tr>';
                html += '</table>';
                html += '</div>';
            }

            // todo: URL Parameters

            // POST parameters
            if (method === "POST") {
                // サブタイトル
                html += '<h2>' + method + ' Parameter' + (info.PostParameters.length > 1 ? 's' : '') + '</h2>';
                html += "<div class=\"contents\">";

                // パラメータリスト
                if (info.PostParameters.length <= 0) {
                    html += "No parameter.";
                }
                else {
                    html += '<table class="ui-widget list">';
                    html += '<tr>';
                    html += '<th class="ui-widget-header header">Name</th>';
                    html += '<th class="ui-widget-header header">Type</th>';
                    html += '<th class="ui-widget-header header">Required</th>';
                    html += '<th class="ui-widget-header header">Conditions</th>';
                    html += '<th class="ui-widget-header header">Value</th>';
                    html += '<th class="ui-widget-header header">Description</th>';
                    html += '</tr>';
                    for (let parameter_name in info.PostParameters) {
                        let parameter = info.PostParameters[parameter_name];
                        html += '<tr>';
                        html += '<td class="ui-widget-content name">' + parameter_name + '</td>';
                        html += '<td class="ui-widget-content type">' + reference.get_type_name(parameter.Type) + '</td>';
                        html += '<td class="ui-widget-content required">' + parameter.Require + '</td>';
                        html += '<td class="ui-widget-content condition">' + reference.get_condition(parameter) + '</td>';

                        if (parameter.Type == TYPE_ARRAY || parameter.Type === TYPE_MAP) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_map="true" value="' + reference.get_paramater_value(parameter) + '" /></td>';
                        }
                        else if (reference.is_numeric_type(parameter.Type)) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_number="true" value="' + reference.get_paramater_value(parameter) + '" /></td>';
                        }
                        else if (parameter.Type == TYPE_BOOL) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_bool="true" value="' + reference.get_paramater_value(parameter) + '" /></td>';
                        }
                        else if (parameter.Type == TYPE_STRING && parameter.IsMultilineString) {
                            html += '<td class="ui-widget-content value"><textarea name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\'>' + reference.get_paramater_value(parameter) + '</textarea></td>';
                        }
                        else {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' value="' + reference.get_paramater_value(parameter) + '" /></td>';
                        }

                        html += '<td class="ui-widget-content desc">' + parameter.Description + '</td>';
                        html += '</tr>';
                    }
                    html += '</table>';
                }
                html += '</div>';
            }

            // returns
            {
                html += '<h2>Returns</h2>';
                html += '<div class="contents">';
                html += '<h4>common</h4>';
                html += '<table class="ui-widget list">';
                html += '<tr>';
                html += '<th class="ui-widget-header header">Name</th>';
                html += '<th class="ui-widget-header header">Type</th>';
                html += '<th class="ui-widget-header header">Description</th>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">resultCode</td>';
                html += '<td class="ui-widget-content type">int</td>';
                html += '<td class="ui-widget-content desc">API実行結果コード。HTTP互換。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">commandName</td>';
                html += '<td class="ui-widget-content type">string</td>';
                html += '<td class="ui-widget-content desc">実行コマンド名。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">errorMessage</td>';
                html += '<td class="ui-widget-content type">string</td>';
                html += '<td class="ui-widget-content desc">ユーザー向けエラーメッセージ。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">debugMessage</td>';
                html += '<td class="ui-widget-content type">string</td>';
                html += '<td class="ui-widget-content desc">開発者向けのデバッグ用システムメッセージ。本番では出力されません。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">result</td>';
                html += '<td class="ui-widget-content type">array</td>';
                html += '<td class="ui-widget-content desc">コマンドの実行結果を格納する連想配列。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">memory</td>';
                html += '<td class="ui-widget-content type">array</td>';
                html += '<td class="ui-widget-content desc">コマンド実行時のメモリの使用状況。本番では出力されません。</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">processTime</td>';
                html += '<td class="ui-widget-content type">float</td>';
                html += '<td class="ui-widget-content desc">コマンドの処理時間。本番では出力されません。</td>';
                html += '</tr>';
                html += '</table>';
                if (info.Returns) {
                    html += '<h4>result details</h4>';
                    html += '<table class="ui-widget list">';
                    html += '<tr>';
                    html += '<th class="ui-widget-header header">Name</th>';
                    html += '<th class="ui-widget-header header">Type</th>';
                    html += '<th class="ui-widget-header header">Description</th>';
                    html += "</tr>";
                    for (let name in info.Returns) {
                        html += '<tr>';
                        html += '<td class="ui-widget-content name">' + name + "</td>";
                        html += '<td class="ui-widget-content type">' + reference.get_type_name(info.Returns[name].Type) + "</td>";
                        html += '<td class="ui-widget-content desc">' + info.Returns[name].Description + "</td>";
                        html += '</tr>';
                    }
                    html += '</table>';
                }
                html += '</div>';
            }

            // test form
            {
                html += "<h2>Test form</h2>";
                html += '<div class="contents">';
                html += '<input type="button" value="Execute" class="ui-button ui-corner-all ui-widget" onclick="do_test()" />';

                html += '<div id="tabs">';
                html += '<ul>';
                html += '<li><a href="#result1">Basic View</a></li>';
                html += '<li><a href="#result2">Result Tree View</a></li>';
                html += '</ul>';

                html += '<div id="result1">';
                html += '<table class="ui-widget result">';
                html += '<tr>';
                html += '<th class="ui-widget-header header">Result Items</th>';
                html += '<th class="ui-widget-header header">Values</th>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">URL</td>';
                html += '<td class="ui-widget-content resultListValue" id="requestUrl"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">Request message</td>';
                html += '<td class="ui-widget-content value" id="requestMsg"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">HTTP status code</td>';
                html += '<td class="ui-widget-content value" id="httpStatusCode"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">JSON to object result</td>';
                html += '<td class="ui-widget-content value" id="jsonConvert"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><b>[result]</b> API result code</td>';
                html += '<td class="ui-widget-content value" id="resultCode"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><b>[result]</b> Error message</td>';
                html += '<td class="ui-widget-content value" id="errorMessage"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><b>[result]</b> Debug message</td>';
                html += '<td class="ui-widget-content value" id="debugMessage"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><b>[result]</b> Peak usage memory</td>';
                html += '<td class="ui-widget-content value" id="peakUsageMemory"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><b>[result]</b> API processing time</td>';
                html += '<td class="ui-widget-content value" id="cmdProcessTime"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">Total processing time</td>';
                html += '<td class="ui-widget-content value" id="totalProcessTime"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item">Response data</td>';
                html += '<td class="ui-widget-content value" id="json"></td>';
                html += '</tr>';
                html += '</table>';
                html += '</div>';
                html += '<div id="result2"><div id="tableview">No execution.</div></div>';
                html += '</div>';

                html += '</div>';
                html += '</div>';
            }

            // 出力
            $("#content_frame").html(html);
            $("#tabs").tabs();
        },

        /**
         * リソース一覧からメニューを作成します。
         */
        create_menu: function () {
            if (!reference._data) {
                reference.set_error_content("Resources data is null.");
                return;
            }

            // タイトルなどをセット
            $("title").text(reference._data.Name);
            $("#api_name").html(reference._data.Name);
            $("#env_name").html(sprintf("[<span id='%s'>%s</span>]", reference._data.EnvClass, reference._data.EnvName));
            $("#header_right").html(reference._data.Logo);

            // メニューの作成
            let index = 0;
            let html = "";
            let resources = reference._data.Resources;
            for (let group in resources) {
                let group_define = resources[group];
                html += "<div class='api_group'><h4><a href='#'><span class='api_group_name'>" + group + "</span></a></h4><div><ol>";

                for (let path in group_define) {
                    let resource_define = group_define[path];

                    for (let method in resource_define.Methods) {
                        let method_define = resource_define.Methods[method];

                        html += "<li><a href='#'><span class='api_path' id='api_res_" + ++index + "' method='" + method + "' path='" + path + "' version='" + "1" + "'>"
                            + method + "&nbsp;" + path
                            + "</span></a><br><span class='api_summary'>" + method_define.Summary + "</span></li>";
                    }
                }

                html += "</ol></div></div>";
            }

            // #accordion内にhtmlを入れて、アコーディオン化
            $("#accordion").html(html);
            $("#accordion").accordion({header: "h4", autoHeight: false, navigation: true});

            // クリックイベントを設定
            $('.api_path').off('click');
            $(".api_path").on('click', reference.show_reference);
        },

        /**
         * 処理化処理を開始します。
         */
        initialize: function () {
            let path = util.trim_char(location.pathname.replace("index.html", ""), "/");
            let config_url = location.protocol + "//" + location.host + "/" + path + "/config.json";
            $.ajax({
                type: "GET",
                url: config_url,
                headers: {
                    "Pragma": "no-cache",
                    "Cache-Control": "no-cache",
                    "Content-Type": "application/json",
                },
                dataType: 'json',
            }).done(function (data, status, xhr) {
                if (!data || !data.EnvClass) {
                    reference.set_error_content(xhr.responseText);
                    return;
                }
                reference._data = data;
                reference._default_user_agent = data.UserAgent;

                reference.create_menu();
            }).fail(function (xhr, status, error) {
                reference.set_error_content(xhr.responseText);
            }).always(function (jqXHR, textStatus) {
                reference.resize();
            });
        },
    };


    $(function () {
        // リファレンスの構成情報を取得
        reference.initialize();

        // 画面のリサイズイベント
        window.onresize = reference.resize;
    });
}());
