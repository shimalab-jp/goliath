;(function() {
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

    /** API Result Code **/
    const RESULT_OK                    = 200;  // 正常終了
    const RESULT_BAD_REQUEST           = 400;  // 不正なリクエスト
    const RESULT_UNAUTHORIZED          = 401;  // 認証されていない
    const RESULT_PAYMENT_REQUIRED      = 402;  // 課金されてない
    const RESULT_FORBIDDEN             = 403;  // APIの実行権限がない
    const RESULT_NOT_FOUND             = 404;  // APIが見つからない
    const RESULT_METHOD_NOT_ALLOWED    = 405;  // メソッドは許可されていない
    const RESULT_NOT_ACCEPTABLE        = 406;  // 受理出来無いリクエスト（一般的なコマンドエラー）
    const RESULT_REQUEST_TIME_OUT      = 408;  // タイムアウト
    const RESULT_CONFLICT              = 409;  // 競合エラー（報酬を既に受領済みとか）
    const RESULT_GONE                  = 410;  // 不正なトークン
    const RESULT_INTERNAL_SERVER_ERROR = 500;  // 内部エラー
    const RESULT_NOT_IMPLEMENTED       = 501;  // 未実装
    const RESULT_SERVICE_UNAVAILABLE   = 503;  // メンテナンス中
    const RESULT_REQUIRE_UPDATE        = 600;  // アップデート要求
    const RESULT_USER_ERROR            = 700;  // ユーザーエラー
    const RESULT_SYSTEM_ERROR          = 800;  // システムエラー
    const RESULT_FATAL_ERROR           = 900;  // 致命的なエラー

    let util = {
        /**
         * 文字列sから、指定文字cをトリムします。
         * @param s 文字列
         * @param c トリム文字
         * @returns string
         */
        trim_char: function (s, c) {
            return s.replace(new RegExp("^" + c + "+|" + c + "+$", "g"), '');
        },

        /**
         * htmlをエスケープします
         * @param str
         * @returns {void | string | *}
         */
        html_escape: function(str) {
            if (!str) return;
            return str.replace(/[<>&"'`]/g, (match) => {
                const escape = {
                    '<': '&lt;',
                    '>': '&gt;',
                    '&': '&amp;',
                    '"': '&quot;',
                    "'": '&#39;',
                    '`': '&#x60;'
                };
                return escape[match];
            });
        },

        replace_all: function(txt, replace, with_this) {
            if (typeof txt == "string") {
                return txt.replace(new RegExp(replace, 'g'), with_this);
            }
            return txt;
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
         * [テストフォームのtree view用]
         */
        _tree_view_count: 0,

        /**
         * エラーコンテンツを表示します。
         * @param s エラー内容
         */
        set_error_content: function(s) {
            let html = "";
            html += '<div class="contents"><h1 class="api_title">Error</h1>';
            html += '<div class="contents">Failed to load index for reference.<br />Please confirm the following response data. </div>';
            html += '<h2>Response Data</h2>';
            html += '<div class="contents">';
            html += util.html_escape(s);
            html += '</div></div>';
            $("#content_frame").html(html);
        },

        /**
         * リサイズします。
         */
        resize: function() {
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
        get_resource: function(version, path, method) {
            if (!reference._data) {
                reference.set_error_content("Resources data is null.");
                return null;
            }

            for (let ver in reference._data.Resources) {
                for (let group in reference._data.Resources[ver]) {
                    let group_define = reference._data.Resources[ver][group];
                    for (let p in group_define) {
                        if (path != p) continue;
                        for (let m in group_define[p].Methods) {
                            if (method != m) continue;
                            return group_define[p].Methods[method];
                        }
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
        get_type_name: function(type_code) {
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
        is_numeric_type: function(type_code) {
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

        /**
         * 表示用のパラメータの入力条件を取得します。
         * @param param
         * @returns {string}
         */
        get_condition: function(param) {
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

        /**
         * パラメータの値を文字列として取得します。
         * @param param
         * @returns {string}
         */
        get_parameter_value: function (param) {
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

        /**
         * バージョン番号に対するurlを取得します。
         * @param version
         * @returns {string}
         */
        get_version_url: function(version) {
            let version_url = "";
            for (let index in reference._data.Versions) {
                if (version == reference._data.Versions[index].Version) {
                    version_url = "/" + util.trim_char(reference._data.Versions[index].Url, "/");
                    break;
                }
            }
            return "/" + util.trim_char(version_url, "/");
        },

        /**
         * [テストフォームのtree view用] jsonの型名を取得
         * @param v
         * @returns {*}
         */
        get_json_type_name: function(v) {
            if (typeof(v) === "undefined") {
                return "undefined";
            }
            return typeof(v);
        },

        /**
         * [テストフォームのtree view用] jsonの値を安全に取得
         * @param v
         * @returns {string}
         */
        get_json_safe_value: function(v) {
            if (v === null) {
                return "null";
            }
            if (typeof(v) === "undefined") {
                return "undefined";
            }
            return String(v);
        },

        /**
         * [テストフォームのtree view用] テーブル作成サブルーチン
         * @param parent_id
         * @param data
         * @returns {string}
         */
        render_result_table_sub: function(parent_id, data) {
            let html = "";
            if (data instanceof Array || data instanceof Object) {
                for (let i in data) {
                    let id = ++reference._tree_view_count;
                    let tt_id = parent_id.length > 0 ? parent_id + "-" + String(id) : String(id);

                    if (parent_id.length > 0) {
                        html += "<tr data-tt-id=\"" + tt_id + "\" data-tt-parent-id=\"" + parent_id + "\">";
                    }
                    else {
                        html += "<tr data-tt-id=\"" + tt_id + "\">";
                    }

                    if (data[i] instanceof Array || data[i] instanceof Object) {
                        html += '<td class="ui-widget-content"><span class=\'folder\'>' + String(i) + "</span></td>";
                        html += '<td class="ui-widget-content">' + reference.get_json_type_name(data[i]) + "</td>";
                        html += '<td class="ui-widget-content">-</td>';
                        html += '</tr>';
                        html += reference.render_result_table_sub(tt_id, data[i]);
                    }
                    else {
                        html += '<td class="ui-widget-content"><span class="file">' + String(i) + '</span></td>';
                        html += '<td class="ui-widget-content">' + reference.get_json_type_name(data[i]) + '</td>';
                        html += '<td class="ui-widget-content">' + reference.get_json_safe_value(data[i]) + '</td>';
                        html += '</tr>';
                    }
                }
            }

            return html;
        },

        /**
         * [テストフォームのtree view用] 全展開
         */
        expand_result_tree_view: function() {
            $("#result_tree_view").treetable('expandAll');
        },

        /**
         * [テストフォームのtree view用] 全縮小
         */
        collapse_result_tree_view: function() {
            $("#result_tree_view").treetable('collapseAll');
        },

        /**
         * [テストフォームのtree view用] 結果テーブル作成
         * @param data
         * @returns {string}
         */
        create_result_table: function(data) {
            reference._tree_view_count = 0;

            let html = "";
            html += '<table id="result_tree_view" class="ui-widget list">';
            html += '<caption style="text-align:left">';
            html += '<input type="button" id="expand_button" value="Expand all" class="ui-button ui-corner-all ui-widget" />';
            html += '<input type="button" id="collapse_button" value="Collapse all" class="ui-button ui-corner-all ui-widget" />';
            html += '</caption>';
            html += '<thead><tr>';
            html += '<th class="ui-widget-header header">Name or Array Index</th>';
            html += '<th class="ui-widget-header header">Type</th>';
            html += '<th class="ui-widget-header header">Value</th>';
            html += '</tr></thead>';
            html += '<tbody>';
            html += reference.render_result_table_sub("", data);
            html += '</tbody></table>';

            // ツリービューを表示
            $("#table_view").html(html);
            $("#result_tree_view").treetable({ expandable: true });

            // クリックイベントを設定
            $("#expand_button").off('click');
            $("#expand_button").on('click', reference.expand_result_tree_view);
            $("#collapse_button").off('click');
            $("#collapse_button").on('click', reference.collapse_result_tree_view);
        },

        /**
         * テストフォームの結果コードをメッセージに変換します。
         * @param code
         * @returns {*|string}
         */
        get_api_result_text: function(code) {
            let ret = "";
            switch (code) {
                case RESULT_OK:                    ret = sprintf('%03d %s', code, 'OK');                    break;
                case RESULT_BAD_REQUEST:           ret = sprintf('%03d %s', code, 'BAD_REQUEST');           break;
                case RESULT_UNAUTHORIZED:          ret = sprintf('%03d %s', code, 'UNAUTHORIZED');          break;
                case RESULT_PAYMENT_REQUIRED:      ret = sprintf('%03d %s', code, 'PAYMENT_REQUIRED');      break;
                case RESULT_FORBIDDEN:             ret = sprintf('%03d %s', code, 'FORBIDDEN');             break;
                case RESULT_NOT_FOUND:             ret = sprintf('%03d %s', code, 'NOT_FOUND');             break;
                case RESULT_METHOD_NOT_ALLOWED:    ret = sprintf('%03d %s', code, 'METHOD_NOT_ALLOWED');    break;
                case RESULT_NOT_ACCEPTABLE:        ret = sprintf('%03d %s', code, 'RESULT_NOT_ACCEPTABLE'); break;
                case RESULT_REQUEST_TIME_OUT:      ret = sprintf('%03d %s', code, 'REQUEST_TIME_OUT');      break;
                case RESULT_CONFLICT:              ret = sprintf('%03d %s', code, 'RESULT_CONFLICT');       break;
                case RESULT_GONE:                  ret = sprintf('%03d %s', code, 'RESULT_GONE');           break;
                case RESULT_INTERNAL_SERVER_ERROR: ret = sprintf('%03d %s', code, 'INTERNAL_SERVER_ERROR'); break;
                case RESULT_NOT_IMPLEMENTED:       ret = sprintf('%03d %s', code, 'NOT_IMPLEMENTED');       break;
                case RESULT_SERVICE_UNAVAILABLE:   ret = sprintf('%03d %s', code, 'SERVICE_UNAVAILABLE');   break;
                case RESULT_REQUIRE_UPDATE:        ret = sprintf('%03d %s', code, 'REQUIRE_UPDATE');        break;
                case RESULT_USER_ERROR:            ret = sprintf('%03d %s', code, 'USER_ERROR');            break;
                case RESULT_SYSTEM_ERROR:          ret = sprintf('%03d %s', code, 'SYSTEM_ERROR');          break;
                case RESULT_FATAL_ERROR:           ret = sprintf('%03d %s', code, 'FATAL_ERROR');           break;
                default:
                    ret = sprintf('%03d %s', code, 'UNKOWN');
            }
            return ret;
        },

        /**
         * テストフォーム POSTデータを
         * @returns {*|string}
         */
        create_post_message: function() {
            let elements = [];
            $('input').each(function(index, element){ elements.push(element); });
            $('textarea').each(function(index, element){ elements.push(element); });
            $('select').each(function(index, element){ elements.push(element); });

            let params = {};
            for (let i = 0; i < elements.length; i++) {
                let element = elements[i];

                if (element.id.indexOf('params{"') >= 0) {
                    let param_name = util.replace_all(util.replace_all(util.replace_all(element.id, 'params{', ''), '"', ''), '}', '');

                    // 配列
                    if ($(element).attr('is_array') !== undefined) {
                        params[param_name] = eval("(" + $(element).val() + ")");
                    }

                    // 連想配列
                    else if ($(element).attr('is_map') !== undefined) {
                        params[param_name] = eval("(" + $(element).val() + ")");
                    }

                    // 数値型
                    else if ($(element).attr('is_number') !== undefined) {
                        let number_val = $(element).val().toLowerCase().trim();
                        if ($.isNumeric(number_val)) {
                            params[param_name] = Number(number_val);
                        } else {
                            params[param_name] = $(element).val();
                        }
                    }

                    // bool型
                    else if ($(element).attr('is_bool') !== undefined) {
                        let bool_val = $(element).val().toLowerCase().trim();
                        params[param_name] = bool_val === 'true'
                    }

                    // その他
                    else {
                        params[param_name] = "" + $(element).val();
                    }
                }
            }

            return JSON.stringify(params);
        },

        /**
         * テストフォーム 実行
         */
        execute: function() {
            // 入力無効化
            $("input").prop("disabled", true);
            $("select").prop("disabled", true);
            $("textarea").prop("disabled", true);

            // 結果を初期化
            $("#request_method").html("-");
            $("#request_url").html("-");
            $("#request_message").html("-");
            $("#response_http_status_code").html("-");
            $("#response_json_decode_result").html("-");
            $("#response_api_result_code").html("-");
            $("#response_error_message").html("-");
            $("#response_debug_message").html("-");
            $("#response_api_processing_time").html("-");
            $("#response_raw").html("-");
            $("#table_view").html("No execution.");

            // ユーザーエージェント取得
            let user_agent = $("#user_agent").val();
            if (!user_agent || user_agent.trim().length <= 0) {
                user_agent = "GOLIATH TEST CLIENT/0.0.0 PC DEVELOP";
            }

            // アカウントトークン取得
            let account_token = "";
            if ($("#account_token") && $("#account_token").length > 0) {
                account_token = $("#account_token").val();
            }

            // メソッド取得
            let method = $("#method").val();

            // URL取得
            let url = $("#url").val();

            // TODO: urlパラメータ対応

            // POSTメッセージを取得
            let post_message = null;
            if (method == "POST") {
                post_message = reference.create_post_message();
            }

            // リクエスト情報を出力
            $("#request_method").html(method);
            $("#request_url").html(url);
            $("#request_message").html(post_message);

            // 実行
            let start_time = new Date().getTime();
            $.ajax({
                type: method,
                url: url,
                headers: {
                    "Pragma": "no-cache",
                    "Cache-Control": "no-cache",
                    "Content-Type": "application/json",
                    "X-Goliath-User-Agent": user_agent,
                    "X-Goliath-Token": account_token,
                },
                dataType: 'json',
                data: post_message
            }).done(function(data, status, xhr) {

                if (!data || !data.ResultCode) return;

                // アカウントトークンを記憶
                if (data.Result && data.Result.AccountInfo && data.Result.AccountInfo.Token) {
                    $.cookie("account_token", data.Result.AccountInfo.Token, {expires:30,secure:false});
                }

                // UserAgentを記憶
                $.cookie("user_agent", document.getElementById("user_agent").value, {expires:30,secure:false});

                // 結果格納
                let end_time = new Date().getTime();
                let total_time = end_time - start_time;
                $("#response_http_status_code").html(xhr.status);
                $("#response_json_decode_result").html("Successful");
                $("#response_api_result_code").html(reference.get_api_result_text(data.ResultCode));
                $("#response_error_message").html(data.ErrorMessage);
                $("#response_debug_message").html(data.DebugMessage);
                $("#response_api_processing_time").html(sprintf("API:%dms / Total:%dms", (data.ProcessTime / 1000000), total_time));
                $("#response_raw").html(util.html_escape(xhr.responseText));
                $("#table_view").html("No execution.");

                // ツリービュー作成
                reference.create_result_table(data.Result);

            }).fail(function(xhr, status, error) {

                let end_time = new Date().getTime();
                let total_time = end_time - start_time;
                $("#response_http_status_code").html(xhr.status);
                $("#response_error_message").html(error);
                $("#response_json_decode_result").html("Failed");
                $("#response_api_processing_time").html(sprintf("API:-ms / Total:%dms", total_time));
                $("#response_raw").html(util.html_escape(xhr.responseText));

            }).always(function(xhr, textStatus) {

                $("input").prop("disabled", false);
                $("select").prop("disabled", false);
                $("textarea").prop("disabled", false);
                reference.resize();

            });
        },

        show_reference: function(e) {
            let id = "#" + e.currentTarget.id;
            let method = $(id).attr('method');
            let path = $(id).attr('path');
            let version = $(id).attr('version');
            let api_title = reference.get_version_url(version) + path
            let api_url = location.protocol + "//" + location.host + api_title;

            let info = reference.get_resource(version, path, method);
            if (!info) return;

            // APIタイトル
            let html = '<div class="contents"><h1 class="api_title">' + method + "&nbsp;" + api_title + "</h1>";
            html += '<input type="hidden" name="display_name" id="display_name" value="' + path + '" />';
            html += '<input type="hidden" name="resource_name" id="resource_name" value="' + path + '" />';
            html += '<div class="contents">';

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
                html += '<td class="ui-widget-content name">Authentication require</td>';
                html += '<td class="ui-widget-content value">' + info.RequireAuthentication + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Administration mode only</td>';
                html += '<td class="ui-widget-content value">' + info.IsAdminModeOnly + "</td>";
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Debug mode only</td>';
                html += '<td class="ui-widget-content value">' + info.IsDebugModeOnly + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Executable in maintenance</td>';
                html += '<td class="ui-widget-content value">' + info.RunInMaintenance + '</td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Test form URL</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="url" id="url" value="' + api_url + '" /></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Method</td>';
                html += '<td class="ui-widget-content value">' + method + '<input type="hidden" name="method" id="method" value="' + method + "\" /></td>";
                html += '</tr>';

                let user_agent = $.cookie("user_agent");
                if (!user_agent || user_agent == "null") user_agent = reference._default_user_agent;
                html += '<tr>';
                html += '<td class="ui-widget-content name">User Agent</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="user_agent" id="user_agent" value="' + user_agent + "\" /></td>";
                html += '</tr>';
                html += '</table>';
                html += '<div class="notice">Goliathでは、クライアントのバージョンチェックをUserAgentにて行っています。';
                html += '<div class="notice">正しいUserAgentを指定していない場合は、不正なクライアントとしてアクセス拒否される場合がありますので、クライアント側はconfigに定義したパターンのUserAgentを指定してください。</div>';
                html += '<div class="notice">※UserAgentの変更が出来ないクライアントの場合は、ヘッダに \'<strong class="orange">X-Goliath-User-Agent</strong>\' という名前でUserAgentを指定してください。</div>';
                html += '</div>';
            }

            // Authentication
            if (info.RequireAuthentication) {
                let account_token = $.cookie("account_token");
                if (!account_token) account_token = "";
                html += '<h2>Authentication</h2>';
                html += '<div class="contents">';
                html += '<div class="notice">このAPIは認証が要求されています。</div>';
                html += '<div class="notice">リクエスト時のヘッダ情報に \'<strong class="orange">X-Goliath-Token</strong>\' という名前でアカウントトークンを渡してください。</div>';
                html += '<div class="notice">※リファレンス内では、 \'<strong class="blue">account/regist</strong>\' API、または \'<strong class="blue">account/trans</strong>\' APIを実行するとcookieに保存され、自動的に入力されます。</div>';
                html += '<table class="ui-widget list">';
                html += '<tr>';
                html += '<th class="ui-widget-header header">Name</th>';
                html += '<th class="ui-widget-header header">Value</th>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content name">Account Token</td>';
                html += '<td class="ui-widget-content value"><input type="text" name="account_token" id="account_token" value="' + account_token + "\" /></td>";
                html += '</tr>';
                html += '</table>';
                html += '</div>';
            }

            // todo: URL Parameters

            // POST parameters
            if (method === "POST") {
                // パラメータ数をカウント
                let param_count = 0;
                for (let v in info.PostParameters) {param_count++;}

                // サブタイトル
                html += '<h2>' + method + ' Parameter' + (param_count > 1 ? 's' : '') + '</h2>';
                html += "<div class=\"contents\">";

                // パラメータリスト
                if (param_count <= 0) {
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

                        if (parameter.Type == TYPE_ARRAY) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_array="true" value="' + reference.get_parameter_value(parameter) + '" /></td>';
                        }
                        else if (parameter.Type === TYPE_MAP) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_map="true" value="' + reference.get_parameter_value(parameter) + '" /></td>';
                        }
                        else if (reference.is_numeric_type(parameter.Type)) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_number="true" value="' + reference.get_parameter_value(parameter) + '" /></td>';
                        }
                        else if (parameter.Type == TYPE_BOOL) {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' is_bool="true" value="' + reference.get_parameter_value(parameter) + '" /></td>';
                        }
                        else if (parameter.Type == TYPE_STRING && parameter.IsMultilineString) {
                            html += '<td class="ui-widget-content value"><textarea name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\'>' + reference.get_parameter_value(parameter) + '</textarea></td>';
                        }
                        else {
                            html += '<td class="ui-widget-content value"><input type="text" name=\'params{"' + parameter_name + '"}\' id=\'params{"' + parameter_name + '"}\' value="' + reference.get_parameter_value(parameter) + '" /></td>';
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

            // テストフォーム
            {
                html += "<h2>Test form</h2>";
                html += '<div class="contents">';
                html += '<input type="button" id="exec_button" value="Execute" class="ui-button ui-corner-all ui-widget" />';

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
                html += '<td class="ui-widget-content item"><strong class="green">[REQUEST]</strong> Method</td>';
                html += '<td class="ui-widget-content resultListValue" id="request_method"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="green">[REQUEST]</strong> URL</td>';
                html += '<td class="ui-widget-content resultListValue" id="request_url"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="green">[REQUEST]</strong> Request message</td>';
                html += '<td class="ui-widget-content value" id="request_message"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> HTTP status code</td>';
                html += '<td class="ui-widget-content value" id="response_http_status_code"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> JSON decode result</td>';
                html += '<td class="ui-widget-content value" id="response_json_decode_result"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> API result code</td>';
                html += '<td class="ui-widget-content value" id="response_api_result_code"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> Error message</td>';
                html += '<td class="ui-widget-content value" id="response_error_message"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> Debug message</td>';
                html += '<td class="ui-widget-content value" id="response_debug_message"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> API processing time</td>';
                html += '<td class="ui-widget-content value" id="response_api_processing_time"></td>';
                html += '</tr>';
                html += '<tr>';
                html += '<td class="ui-widget-content item"><strong class="orange">[RESPONSE]</strong> RAW</td>';
                html += '<td class="ui-widget-content value" id="response_raw"></td>';
                html += '</tr>';
                html += '</table>';
                html += '</div>';
                html += '<div id="result2"><div id="table_view">No execution.</div></div>';
                html += '</div>';

                html += '</div>';
                html += '</div>';
            }


            html += '</div>';

            // 出力
            $("#content_frame").html(html);

            // テストフォームをタブ化
            $("#tabs").tabs();

            // クリックイベントを設定
            $("#exec_button").off('click');
            $("#exec_button").on('click', reference.execute);
        },

        /**
         * リソース一覧からメニューを作成します。
         */
        create_menu: function() {
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
            for (let version in resources) {
                let version_resources = resources[version];

                // バージョンの基本URLを取得
                let version_url = reference.get_version_url(version);

                // バージョンの定義がされていない場合はリファレンス表示対象外
                if (version_url.length <= 0) continue;

                // メニュー作成
                for (let group in version_resources) {
                    let group_define = version_resources[group];
                    html += "<div class='api_group'><h4><a href='#'><span class='api_group_name'>" + version_url + group + "</span></a></h4><div><ol>";

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
        initialize: function() {
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

    $(function() {
        // リファレンスの構成情報を取得
        reference.initialize();

        // 画面のリサイズイベント
        window.onresize = reference.resize;
    });
}());
