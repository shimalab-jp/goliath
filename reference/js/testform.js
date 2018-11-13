;(function () {

    'use strict';

    let HTTP_OK                    = 200;  // 正常終了
    let HTTP_BAD_REQUEST           = 400;  // 不正なリクエスト
    let HTTP_UNAUTHORIZED          = 401;  // 認証されていない
    let HTTP_PAYMENT_REQUIRED      = 402;  // 課金されてない
    let HTTP_FORBIDDEN             = 403;  // APIの実行権限がない
    let HTTP_NOT_FOUND             = 404;  // APIが見つからない
    let HTTP_NOT_ACCEPTABLE        = 406;  // 受理出来無いリクエスト（一般的なコマンドエラー）
    let HTTP_CONFLICT              = 409;  // 競合エラー（報酬を既に受領済みとか）
    let HTTP_INTERNAL_SERVER_ERROR = 500;  // 内部エラー
    let HTTP_NOT_IMPLEMENTED       = 501;  // 未実装
    let HTTP_SERVICE_UNAVAILABLE   = 503;  // メンテナンス中
    let RF_RESULT_USER_ERROR       = 700;  // ユーザーエラー
    let RF_RESULT_SYSTEM_ERROR     = 800;  // システムエラー
    let RF_RESULT_FATAL_ERROR      = 900;  // 致命的なエラー

    let MAX_WEB_SERVER_PROCESSES = 500;     // 1サーバ当たりのWebサーバーのプロセス数
    let MAX_FRONT_MEMORY   = 4294967296;    // 搭載メモリ
    let FRONT_SYSTEM_USAGE = 2147483648;    // システム使用メモリ
    let LIMIT_MEMORY       = (MAX_FRONT_MEMORY - FRONT_SYSTEM_USAGE) / MAX_WEB_SERVER_PROCESSES;  // php1プロセス当たりが利用可能なメモリ
    let WARNING_MEMORY     = LIMIT_MEMORY * 0.7;    // 警告を出力するメモリ量

    let _content_data = [];
    let _ssid = "";

    let load_reference_data = function() {
        let url = location.protocol + "//" + location.host + "/v1/reference.json";
        $.ajax({
            type: "GET",
            url: url,
            headers: {
                "Pragma": "no-cache",
                "Cache-Control": "no-cache",
                "Content-Type": "application/json",
            },
            dataType: 'json',
            success: function(json, status, xhr) {
                if (!json || !json.Result || !json.Result.EnvClass) {
                    let html = "";
                    html += "<div class=\"contents\"><h1>Error</h1>";
                    html += "Failed to load index for reference.<br />Please confirm the following response data. ";
                    html += "<h2>Response Data</h2>";
                    html += "<div class=\"contents\">";
                    html += xhr.responseText;
                    html += "</div></div>";
                    $("#contentFrame").html(html);
                    return
                }

                $("#api_name").html(json.Result.Name)
                $("#env_name").html(sprintf("[<span id=\"%s\">%s</span>]", json.Result.EnvClass, json.Result.EnvName));
                $("#header_right").html(json.Result.Logo)

/*


                    let menuHtml = "";

                    for (var namespace in index["result"]["index"]) {
                        if (namespace === "reference") continue;

                        if (namespace.length > 0) {
                            menuHtml += "<div class='leftPeinMenu' id='" + namespace + "'><h4><a href=\"#\">" + namespace + "</a></h4><div><ol>";
                        }
                        else {
                            menuHtml += "<div class='leftPeinMenu' id='" + namespace + "'><h4><a href=\"#\">root</a></h4><div><ol>";
                        }
                        if (index["result"]["index"][namespace]) {
                            for (var resource_name in index["result"]["index"][namespace]) {
                                //try {
                                for (var method_index in index["result"]["index"][namespace][resource_name]["support_methods"]) {
                                    var method = index["result"]["index"][namespace][resource_name]["support_methods"][method_index];
                                    if (index["result"]["index"][namespace][resource_name][method]) {
                                        var display_name = index["result"]["index"][namespace][resource_name]["display_name"];
                                        var content_data = index["result"]["index"][namespace][resource_name][method];
                                        if (content_data["is_ready"]) {
                                            menuHtml += "<li class=\"menu_item\"><a href=\"#\" onclick=\"show_reference('" + method + "', '" + resource_name + "')\" class=\"menu_item\">"
                                                + method + " " + display_name
                                                + "</a> <span class=\"cmdready\">[R]</span><br><span class=\"menu_description\">"
                                                + content_data["summary"]
                                                + "</li>";
                                        }
                                        else {
                                            menuHtml += "<li class=\"menu_item\"><a href=\"#\" onclick=\"show_reference('" + method + "', '" + resource_name + "')\" class=\"menu_item\">"
                                                + method + " " + display_name
                                                + "</a><br><span class=\"menu_description\">"
                                                + content_data["summary"]
                                                + "</span></li>";
                                        }
                                        if (!Array.isArray(_content_data[resource_name])) {
                                            _content_data[resource_name] = [];
                                        }
                                        _content_data[resource_name]["display_name"] = display_name;
                                        _content_data[resource_name][method] = content_data;
                                    }

                                }
                                //} catch (e2) {
                                //}
                            }
                            menuHtml += "</ol></div></div>";
                        }
                    }

                    menuHtml += "</div>";
                    menuHtml += "</div>";

                    $("#accordion").accordion({header: "h4", autoHeight: false, navigation: true});
                    $("#leftMenu").html(menuHtml);

*/

                /*
                }
                catch (e) {
                        var html = "";
                        html += "<div class=\"contents\"><h1>Error</h1>";
                        html += "Failed to load index for reference.<br />Please confirm the following response data. ";
                        html += "<h2>Response Data</h2>";
                        html += "<div class=\"contents\">";
                        html += xhr.responseText;
                        html += "</div>";
                        html += "<h2>Exception</h2>";
                        html += "<div class=\"contents\">";
                        html += e;
                        html += "</div></div>";
                        document.getElementById("contentFrame").innerHTML = html;
                }
                */
            }
        });
    };

    function resize() {
        var header = document.getElementById('header');
        var leftMenu = document.getElementById('leftMenu');
        var content = document.getElementById('content');

        var headerHeight = header.offsetHeight;
        var navHeight = document.getElementById('nav') !== null ? document.getElementById('nav').offsetHeight : 0;
        var leftMenuWidth = leftMenu.offsetWidth;
        var clientHeight = window.innerHeight ? window.innerHeight : document.documentElement.clientHeight;
        var clientWidth = (window.innerWidth ? window.innerWidth : document.documentElement.clientWidth);
        var setHeight = clientHeight - headerHeight - navHeight;
        if (setHeight < 5) setHeight = 5;
        var setWidth = clientWidth - leftMenuWidth;
        if (setWidth < 5) setWidth = 5;

        leftMenu.style.marginTop = navHeight.toString() + "px";
        leftMenu.style.height = setHeight.toString() + "px";

        content.style.marginTop = navHeight.toString() + "px";
        content.style.height = setHeight.toString() + "px";
        content.style.marginLeft = leftMenuWidth.toString() + "px";
        content.style.width = setWidth.toString() + "px";
    }

function show_reference(method, resource_name) {
    var display_name = _content_data[resource_name]["display_name"];
    var data = _content_data[resource_name][method];
    if (data !== null) {
        var html = "<div class=\"contents\"><h1>" + method + " " + display_name + "</h1>";
        html += "<input type=\"hidden\" name=\"display_name\" id=\"display_name\" value=\"" + display_name + "\" />";
        html += "<input type=\"hidden\" name=\"resource_name\" id=\"resource_name\" value=\"" + resource_name + "\" />";

        // summary
        html += "<h2>Summary</h2>";
        html += "<div class=\"contents\">";
        html += data["summary"];
        if (data["isDevOnly"] === true) {
            html += "<br />";
            html += "<span class=\"warning\">このコマンドは開発環境でのみ実行可能です。</span>";
        }
        html += "</div>";

        // description
        if (data["description"] && data["description"].length > 0) {
            html += "<h2>Description</h2>";
            html += "<div class=\"contents\">" + data["description"] + "</div>";
        }

        // API info
        html += "<h2>API Information</h2>";
        html += "<div class=\"contents\">";
        html += "<table class=\"ui-widget list\">";
        html += "<tr>";
        html += "<th class=\"ui-widget-header header\">Name</th>";
        html += "<th class=\"ui-widget-header header\">Value</th>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">認証要求</td>";
        html += "<td class=\"ui-widget-content value\">" + data["require_authentication"] + "</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">管理者専用</td>";
        html += "<td class=\"ui-widget-content value\">" + data["is_admin_only"] + "</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">開発時専用</td>";
        html += "<td class=\"ui-widget-content value\">" + data["is_dev_only"] + "</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">実行レベル</td>";
        html += "<td class=\"ui-widget-content value\">" + get_execution_level_text(data["execution_level"]) + "</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">Test form URL</td>";
        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name=\"url\" id=\"url\" value=\"" + _api_url + "\" /></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">メソッド</td>";
        html += "<td class=\"ui-widget-content value\">" + method + "<input type=\"hidden\" name=\"method\" id=\"method\" value=\"" + method + "\" /></td>";
        html += "</tr>";

        var user_agent = $.cookie("user_agent");
        if (!user_agent || user_agent == "null") user_agent = _user_agent_name + "/1.0.0(LOCAL,PC)";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">User Agent</td>";
        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name=\"user_agent\" id=\"user_agent\" value=\"" + user_agent + "\" /></td>";
        html += "</tr>";
        html += "</table>";
        html += "</div>";

        // Authentication
        if (data["require_authentication"]) {
            var ssid = $.cookie("ssid");
            html += "<h2>Authentication</h2>";
            html += "<div class=\"contents\">";
            html += "<div class='notice'>このコマンドは認証が要求されています。</div>";
            html += "<div class='notice'>※セッションIDは、account.SignInコマンドを実行すると自動的に入力されます。</div>";
            html += "<table class=\"ui-widget list\">";
            html += "<tr>";
            html += "<th class=\"ui-widget-header header\">Name</th>";
            html += "<th class=\"ui-widget-header header\">Value</th>";
            html += "</tr>";
            html += "<tr>";
            html += "<td class=\"ui-widget-content name\">Session ID</td>";
            html += "<td class=\"ui-widget-content value\"><input type=\"text\" name=\"ssid\" id=\"ssid\" value=\"" + ssid + "\" /></td>";
            html += "</tr>";
            html += "</table>";
            html += "</div>";
        }

        // URL Parameters
        var url_param_cnt = 0;
        for (var parameter_index in data["url_parameters"]) { url_param_cnt++; }
        html += "<h2>URL Parameter" + (url_param_cnt > 1 ? "s" : "") + "</h2>";
        html += "<div class=\"contents\">";
        if (data["url_parameters"].length <= 0) {
            html += "No parameter.";
        }
        else {
            html += "<table class=\"ui-widget list\">";
            html += "<tr>";
            html += "<th class=\"ui-widget-header header\">Name</th>";
            html += "<th class=\"ui-widget-header header\">Type</th>";
            html += "<th class=\"ui-widget-header header\">Required</th>";
            html += "<th class=\"ui-widget-header header\">Conditions</th>";
            html += "<th class=\"ui-widget-header header\">Value</th>";
            html += "<th class=\"ui-widget-header header\">Description</th>";
            html += "</tr>";
            for (var parameter_index in data["url_parameters"]) {
                var parameter = data["url_parameters"][parameter_index];
                html += "<tr>";
                html += "<td class=\"ui-widget-content name\">" + parameter_index + "</td>";
                html += "<td class=\"ui-widget-content type\">" + parameter["type"] + "</td>";
                html += "<td class=\"ui-widget-content required\">" + ("nullable" in parameter ? !parameter["nullable"] : "false") + "</td>";
                html += "<td class=\"ui-widget-content condition\">" + get_condition(parameter) + "</td>";

                if (parameter["type"] === "{array}" || parameter["type"] === "map") {
                    html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' isMap=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                }
                else if (parameter["type"] === "array" || parameter["type"] === "[array]") {
                    html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' isArray=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                }
                else if (parameter["type"] === "int" || parameter["type"] === "integer" || parameter["type"] === "float" || parameter["type"] === "number") {
                    html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' isNumber=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                }
                else if (parameter["type"].toLowerCase() === "bool" || parameter["type"].toLowerCase() === "bool") {
                    html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' isBool=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                }
                else if (parameter["type"] === "row") {
                    html += "<td class=\"ui-widget-content value\"><textarea name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' isRow=\"true\">" + get_paramater_value(parameter) + "</textarea></td>";
                }
                else if ("multiline" in parameter && parameter["multiline"] === true) {
                    html += "<td class=\"ui-widget-content value\"><textarea name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "'>" + get_paramater_value(parameter) + "</textarea></td>";
                }
                else {
                    html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='url_param_" + parameter_index + "' id='url_param_" + parameter_index + "' value=\"" + get_paramater_value(parameter) + "\" /></td>";
                }
                /*
                else if (param["inputType"] == "select") {
                    var opt = "";
                    if ("isMultiple" in param["type"] && param["type"]["isMultiple"]) {
                        opt += "multiple";
                    }
                    if ("size" in param["type"]) {
                        opt += "size=" + param["type"]["size"];
                    }
                    html += "<td class=\"ui-widget-content value\"><select " + opt + " name=\"" + parameterName + "\" id=\"" + parameterName + "\">";
                    for (var idx2 in param["value"]) {
                        var val = param["value"][idx2]["value"];
                        var txt = param["value"][idx2]["text"];
                        html += "<option value=\"" + val + "\">" + txt + "</option>";
                    }
                    html += "</select></td>";
                }
                else if (param["inputType"] == "option") {
                    html += "<td class=\"ui-widget-content value\">";
                    for (var idx2 in param["value"]) {
                        var val = param["value"][idx2]["value"];
                        var txt = param["value"][idx2]["text"];
                        html += "<label><input type=\"radio\" name=\"" + parameterName + "\" id=\"" + parameterName + "\" value=\"" + val + "\" />" + txt + "</label>";
                    }
                    html += "</select></td>";
                }
                */
                html += "<td class=\"ui-widget-content desc\">" + parameter["description"] + "</td>";
                html += "</tr>";
            }
            html += "</table>";
        }
        html += "</div>";

        // POST/PUT/DELETE parameters
        if (method != "GET") {
            var paramCnt = 0;
            for (var parameter_name in data["parameters"]) { paramCnt++; }
            html += "<h2>" + method + " Parameter" + (paramCnt > 1 ? "s" : "") + "</h2>";
            html += "<div class=\"contents\">";
            if (data["parameters"].length <= 0) {
                html += "No parameter.";
            }
            else {
                html += "<table class=\"ui-widget list\">";
                html += "<tr>";
                html += "<th class=\"ui-widget-header header\">Name</th>";
                html += "<th class=\"ui-widget-header header\">Type</th>";
                html += "<th class=\"ui-widget-header header\">Required</th>";
                html += "<th class=\"ui-widget-header header\">Conditions</th>";
                html += "<th class=\"ui-widget-header header\">Value</th>";
                html += "<th class=\"ui-widget-header header\">Description</th>";
                html += "</tr>";
                for (var parameter_name in data["parameters"]) {
                    var parameter = data["parameters"][parameter_name];
                    html += "<tr>";
                    html += "<td class=\"ui-widget-content name\">" + parameter_name + "</td>";
                    html += "<td class=\"ui-widget-content type\">" + parameter["type"] + "</td>";
                    html += "<td class=\"ui-widget-content required\">" + ("nullable" in parameter ? !parameter["nullable"] : "false") + "</td>";
                    html += "<td class=\"ui-widget-content condition\">" + get_condition(parameter) + "</td>";

                    if (parameter["type"] === "{array}" || parameter["type"] === "map") {
                        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' isMap=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                    }
                    else if (parameter["type"] === "array" || parameter["type"] === "[array]") {
                        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' isArray=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                    }
                    else if (parameter["type"] === "int" || parameter["type"] === "integer" || parameter["type"] === "float" || parameter["type"] === "number") {
                        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' isNumber=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                    }
                    else if (parameter["type"].toLowerCase() === "bool" || parameter["type"].toLowerCase() === "bool") {
                        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' isBool=\"true\" value=\"" + get_paramater_value(parameter) + "\" /></td>";
                    }
                    else if (parameter["type"] === "row") {
                        html += "<td class=\"ui-widget-content value\"><textarea name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' isRow=\"true\">" + get_paramater_value(parameter) + "</textarea></td>";
                    }
                    else if ("multiline" in parameter && parameter["multiline"] === true) {
                        html += "<td class=\"ui-widget-content value\"><textarea name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}'>" + get_paramater_value(parameter) + "</textarea></td>";
                    }
                    else {
                        html += "<td class=\"ui-widget-content value\"><input type=\"text\" name='params{\"" + parameter_name + "\"}' id='params{\"" + parameter_name + "\"}' value=\"" + get_paramater_value(parameter) + "\" /></td>";
                    }
                    /*
                    else if (param["inputType"] == "select") {
                        var opt = "";
                        if ("isMultiple" in param["type"] && param["type"]["isMultiple"]) {
                            opt += "multiple";
                        }
                        if ("size" in param["type"]) {
                            opt += "size=" + param["type"]["size"];
                        }
                        html += "<td class=\"ui-widget-content value\"><select " + opt + " name=\"" + parameterName + "\" id=\"" + parameterName + "\">";
                        for (var idx2 in param["value"]) {
                            var val = param["value"][idx2]["value"];
                            var txt = param["value"][idx2]["text"];
                            html += "<option value=\"" + val + "\">" + txt + "</option>";
                        }
                        html += "</select></td>";
                    }
                    else if (param["inputType"] == "option") {
                        html += "<td class=\"itemValue\">";
                        for (var idx2 in param["value"]) {
                            var val = param["value"][idx2]["value"];
                            var txt = param["value"][idx2]["text"];
                            html += "<label><input type=\"radio\" name=\"" + parameterName + "\" id=\"" + parameterName + "\" value=\"" + val + "\" />" + txt + "</label>";
                        }
                        html += "</select></td>";
                    }
                    */
                    html += "<td class=\"ui-widget-content desc\">" + parameter["description"] + "</td>";
                    html += "</tr>";
                }
                html += "</table>";
            }
            html += "</div>";
        }

        // returns
        html += "<h2>Returns</h2>";
        html += "<div class=\"contents\">";
        html += "<h4>common</h4>";
        html += "<table class=\"ui-widget list\">";
        html += "<tr>";
        html += "<th class=\"ui-widget-header header\">Name</th>";
        html += "<th class=\"ui-widget-header header\">Type</th>";
        html += "<th class=\"ui-widget-header header\">Description</th>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">resultCode</td>";
        html += "<td class=\"ui-widget-content type\">int</td>";
        html += "<td class=\"ui-widget-content desc\">API実行結果コード。HTTP互換。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">commandName</td>";
        html += "<td class=\"ui-widget-content type\">string</td>";
        html += "<td class=\"ui-widget-content desc\">実行コマンド名。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">errorMessage</td>";
        html += "<td class=\"ui-widget-content type\">string</td>";
        html += "<td class=\"ui-widget-content desc\">ユーザー向けエラーメッセージ。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">debugMessage</td>";
        html += "<td class=\"ui-widget-content type\">string</td>";
        html += "<td class=\"ui-widget-content desc\">開発者向けのデバッグ用システムメッセージ。本番では出力されません。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">result</td>";
        html += "<td class=\"ui-widget-content type\">array</td>";
        html += "<td class=\"ui-widget-content desc\">コマンドの実行結果を格納する連想配列。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">memory</td>";
        html += "<td class=\"ui-widget-content type\">array</td>";
        html += "<td class=\"ui-widget-content desc\">コマンド実行時のメモリの使用状況。本番では出力されません。</td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content name\">processTime</td>";
        html += "<td class=\"ui-widget-content type\">float</td>";
        html += "<td class=\"ui-widget-content desc\">コマンドの処理時間。本番では出力されません。</td>";
        html += "</tr>";
        html += "</table>";
        if ("returns" in data) {
            html += "<h4>result details</h4>";
            html += "<table class=\"ui-widget list\">";
            html += "<tr>";
            html += "<th class=\"ui-widget-header header\">Name</th>";
            html += "<th class=\"ui-widget-header header\">Type</th>";
            html += "<th class=\"ui-widget-header header\">Description</th>";
            html += "</tr>";
            for (var idx in data["returns"]) {
                html += "<tr>";
                html += "<td class=\"ui-widget-content name\">" + idx + "</td>";
                html += "<td class=\"ui-widget-content type\">" + data["returns"][idx]["type"] + "</td>";
                html += "<td class=\"ui-widget-content desc\">" + data["returns"][idx]["description"] + "</td>";
                html += "</tr>";
            }
            html += "</table>";
        }
        html += "</div>";

        // test form
        html += "<h2>Test form</h2>";
        html += "<div class=\"contents\">";
        html += "<input type=\"button\" value=\"Execute\" class=\"ui-button ui-corner-all ui-widget\" onclick=\"do_test()\" />";

        html += "<div id=\"tabs\">";
        html += "<ul>";
        html += "<li><a href=\"#result1\">Basic View</a></li>";
        html += "<li><a href=\"#result2\">Result Tree View</a></li>";
        html += "</ul>";

        html += "<div id=\"result1\">";
        html += "<table class=\"ui-widget result\">";
        html += "<tr>";
        html += "<th class=\"ui-widget-header header\">Result Items</th>";
        html += "<th class=\"ui-widget-header header\">Values</th>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">URL</td>";
        html += "<td class=\"ui-widget-content resultListValue\" id=\"requestUrl\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">Request message</td>";
        html += "<td class=\"ui-widget-content value\" id=\"requestMsg\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">HTTP status code</td>";
        html += "<td class=\"ui-widget-content value\" id=\"httpStatusCode\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">JSON to object result</td>";
        html += "<td class=\"ui-widget-content value\" id=\"jsonConvert\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\"><b>[result]</b> API result code</td>";
        html += "<td class=\"ui-widget-content value\" id=\"resultCode\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\"><b>[result]</b> Error message</td>";
        html += "<td class=\"ui-widget-content value\" id=\"errorMessage\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\"><b>[result]</b> Debug message</td>";
        html += "<td class=\"ui-widget-content value\" id=\"debugMessage\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\"><b>[result]</b> Peak usage memory</td>";
        html += "<td class=\"ui-widget-content value\" id=\"peakUsageMemory\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\"><b>[result]</b> API processing time</td>";
        html += "<td class=\"ui-widget-content value\" id=\"cmdProcessTime\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">Total processing time</td>";
        html += "<td class=\"ui-widget-content value\" id=\"totalProcessTime\"></td>";
        html += "</tr>";
        html += "<tr>";
        html += "<td class=\"ui-widget-content item\">Response data</td>";
        html += "<td class=\"ui-widget-content value\" id=\"json\"></td>";
        html += "</tr>";
        html += "</table>";
        html += "</div>";
        html += "<div id=\"result2\"><div id=\"tableview\">No execution.</div></div>";
        html += "</div>";

        html += "</div>";
        html += "</div>";

        document.getElementById("contentFrame").innerHTML = html;

        $("#tabs").tabs();

        if (data["requireAuthentication"]) {
            FB.init({ appId:_platform_app_id, cookie:true, status:true, xfbml:true });
        }

        // account.SignIn コマンドの場合は、トークンを自動入力
        if (resource_name === "api/account/auth") {
            var token = $.cookie("token");
            document.getElementById("params{\"token\"}").value = token;
        }
    }
}

function get_condition(param) {
    var returnValue = "";
    if (param !== null && "range" in param) {
        if (returnValue.length > 0) returnValue += " & ";
        returnValue += "Range( " + param["range"][0] + "~" + param["range"][param["range"].length - 1] + " )";
    }
    if (param !== null && "select" in param) {
        if (returnValue.length > 0) returnValue += " & ";
        returnValue += "Select( ";
        for (var i in param["select"]) {
            if (i > 0) returnValue += ", ";
            returnValue += param["select"][i];
        }
        returnValue += " )";
    }
    if (param !== null && "regex" in param) {
        if (returnValue.length > 0) returnValue += " & ";
        returnValue += "RegEx( " + param["regex"] + " )";
    }
    if (returnValue.length <= 0) {
        returnValue = "-";
    }
    return returnValue;
}

function get_paramater_value(param) {
    if (param !== null && "default" in param) {
        return param["default"];
    }
    else if (param !== null && "type" in param && (param["type"] === "int" || param["type"] === "long" || param["type"] === "float" || param["type"] === "double")) {
        return "0";
    }
    else {
        return "";
    }
}

function get_execution_level_text(code) {
    var returnValue = "";
    switch (code) {
        case API_EXEC_LV_0: returnValue = 'API_EXEC_LV_0 (実行不可)';                   break;
        case API_EXEC_LV_1: returnValue = 'API_EXEC_LV_1 (通常時のみ実行可)';            break;
        case API_EXEC_LV_2: returnValue = 'API_EXEC_LV_2 (簡易メンテナンス時まで実行可)'; break;
        case API_EXEC_LV_9: returnValue = 'API_EXEC_LV_9 (常時実行可)';                  break;
    }
    return returnValue;
}

function get_api_result_text(code) {
    var returnValue = "";
    switch (code) {
        case HTTP_OK:                    returnValue = 'OK';                    break;
        case HTTP_BAD_REQUEST:           returnValue = 'BAD_REQUEST';           break;
        case HTTP_UNAUTHORIZED:          returnValue = 'UNAUTHORIZED';          break;
        case HTTP_PAYMENT_REQUIRED:      returnValue = 'PAYMENT_REQUIRED';      break;
        case HTTP_FORBIDDEN:             returnValue = 'FORBIDDEN';             break;
        case HTTP_NOT_FOUND:             returnValue = 'NOT_FOUND';             break;
        case HTTP_NOT_ACCEPTABLE:        returnValue = 'HTTP_NOT_ACCEPTABLE';   break;
        case HTTP_CONFLICT:              returnValue = 'HTTP_CONFLICT';         break;
        case HTTP_INTERNAL_SERVER_ERROR: returnValue = 'INTERNAL_SERVER_ERROR'; break;
        case HTTP_NOT_IMPLEMENTED:       returnValue = 'NOT_IMPLEMENTED';       break;
        case HTTP_SERVICE_UNAVAILABLE:   returnValue = 'SERVICE_UNAVAILABLE';   break;
        case RF_RESULT_USER_ERROR:       returnValue = 'USER_ERROR';            break;
        case RF_RESULT_SYSTEM_ERROR:     returnValue = 'SYSTEM_ERROR';          break;
        case RF_RESULT_FATAL_ERROR:      returnValue = 'FATAL_ERROR';           break;
    }
    return returnValue;
}

function create_request_message() {
    var inputTags = document.getElementsByTagName("input");
    var textareaTags = document.getElementsByTagName("textarea");
    var selectTags = document.getElementsByTagName("select");

    var tags = new Array();
    for (var i = 0; i < inputTags.length;    i++) { tags.push(inputTags[i]);    }
    for (var i = 0; i < textareaTags.length; i++) { tags.push(textareaTags[i]); }
    for (var i = 0; i < selectTags.length;   i++) { tags.push(selectTags[i]);   }

    var data = new Array();

    for (var i = 0; i < tags.length; i++) {
        var tag = tags[i];

        if (tag.name.length > 0) {
            if (tag.name === "user_agent")    continue;
            if (tag.name === "display_name")  continue;
            if (tag.name === "resource_name") continue;
            if (tag.name === "method")        continue;
            if (tag.name === "url")           continue;
            if (tag.name === "url_params")    continue;
            if (tag.name === "ssid")          continue;
            if (tag.type === "button")        continue;

            if (tag.id.indexOf("[") >= 0) {
                var name = tag.id.substring(0, tag.id.indexOf("["));
                var index = Number(replaceAll(replaceAll(tag.id.substring(tag.id.indexOf("[")), "\\[", ""), "\\]", ""));
                if (!(name in data))
                    data[name] = new Array();

                var isArray = tag.getAttribute("isArray");
                if (isArray !== null && Boolean(isArray)) {
                    data[name][index] = tag.value.split(",");
                    for (var j in data[name][index]) {
                        if (data[name][index][j] && !isNaN(data[name][index][j])) {
                            data[name][index][j] = Number(data[name][index][j]);
                        }
                        else {
                            if (data[name][index][j].length <= 0)
                                data[name][index][j] = null;
                            else
                                data[name][index][j] = data[name][index][j].replace(/^\s+|\s+$/g, "");
                        }
                    }
                }
                else {
                    data[name][index] = tag.value;
                }
            }
            else if (tag.id.indexOf("{") >= 0) {
                var name = tag.id.substring(0, tag.id.indexOf("{"));
                var key = replaceAll(replaceAll(replaceAll(tag.id.substring(tag.id.indexOf("{")), "\\{", ""), "\\}", ""), "\"", "");
                if (!(name in data))
                    data[name] = new Array();

                var isMap = tag.getAttribute("isMap");
                var isArray = tag.getAttribute("isArray");
                var isNumber = tag.getAttribute("isNumber");
                var isBool = tag.getAttribute("isBool");
                var isRow = tag.getAttribute("isRow");
                if (isArray !== null && Boolean(isArray)) {
                    data[name][key] = tag.value.split(",");
                    for (var j in data[name][key]) {
                        if (data[name][key][j] && !isNaN(data[name][key][j])) {
                            data[name][key][j] = Number(data[name][key][j]);
                        }
                        else {
                            if (data[name][key][j].length <= 0)
                                data[name][key][j] = null;
                            else
                                data[name][key][j] = data[name][key][j].replace(/^\s+|\s+$/g, "");
                        }
                    }
                }
                else if (isMap !== null && Boolean(isMap)) {
                    data[name][key] = eval("(" + tag.value + ")");
                    /*
                    var items = tag.value.split(",");
                    data[name][key] = new Array();
                    for (var j in items) {
                        var kv = items[j].split("=>");
                        if (kv.length >= 2) {
                            for (var k = 0; k < 2; k++) {
                                if (kv[k] && !isNaN(kv[k])) {
                                    kv[k] = Number(kv[k]);
                                }
                                else {
                                    if (kv[k].length <= 0)
                                        kv[k] = null;
                                    else
                                        kv[k] = kv[k].replace(/^\s+|\s+$/g, "");
                                }
                            }
                            data[name][key][kv[0]] = kv[1];
                        }
                    }
                    */
                }
                else if (isNumber !== null && Boolean(isNumber)) {
                    if (tag.value.length > 0) {
                        data[name][key] = Number(tag.value);
                    }
                }
                else if (isBool !== null && Boolean(isBool)) {
                    if (tag.value.length > 0 && (tag.value.toLowerCase() === "true" || tag.value.toLowerCase() === "false")) {
                        data[name][key] = tag.value.toLowerCase() === "true";
                    }
                }
                else if (isRow !== null && Boolean(isRow)) {
                    data[name][key] = eval("(" + tag.value + ")");
                }
                else if (tag.value === true || tag.value === false || tag.value === "true" || tag.value === "false" || tag.value === "True" || tag.value === "False" || tag.value === "TRUE" || tag.value === "false") {
                    data[name][key] = tag.value.toLowerCase() === "true";
                }
                else {
                    data[name][key] = "" + tag.value;
                }
            }
            else {
                if (tag.type === "hidden" || tag.type === "text" || tag.type === "password" || (tag.type === "radio" && tag.checked)) {
                    data[tag.name] = tag.value;
                }
                else if (tag.type === "checkbox" && tag.checked) {
                    if (!data[tag.name])
                        data[tag.name] = new Array();
                    data[tag.name].push(tag.value);
                }
            }
        }
    }

    if (!("params" in data))
        data["params"] = new Array();

    var jsonStr = to_json(data);

    return jsonStr;
}

function to_json(obj) {
    var returnValue = '';

    if (obj instanceof Array || obj instanceof Object) {
        var isFirst = true;
        var isHash = false;
        for (var key in obj) {
            isHash = isNaN(key);
            break;
        }

        if (isHash) {
            returnValue += "{";
            for (var key in obj) {
                if (!isFirst)
                    returnValue += ", ";
                returnValue += "\"" + key + "\":" + to_json(obj[key]);
                isFirst = false;
            }
            returnValue += "}";
        }
        else {
            returnValue += "[";
            for (var key in obj) {
                if (!isFirst) returnValue += ", ";
                returnValue += to_json(obj[key]);
                isFirst = false;
            }
            returnValue += "]";
        }
    }
    else {
        val = obj;
        if (val === undefined || val === null) {
            returnValue += "null";
        }
        else if (val === true || val === false || val === "true" || val === "false" || val === "True" || val === "False" || val === "TRUE" || val === "false") {
            returnValue += "" + val;
        }
        else if (typeof val != "number") {
            val = replaceAll(val, "\\\\", "\\\\");
            val = replaceAll(val, "\"", "\\\"");
            returnValue += "\"" + val + "\"";
        }
        else {
            if (val.toString().length <= 0) {
                returnValue += '""';
            }
            else {
                returnValue += val.toString();
            }
        }
    }

    return returnValue;
}

function do_test() {
    disable("input");
    disable("select");
    disable("textarea");

    document.getElementById("requestUrl").innerHTML = "-";
    document.getElementById("requestMsg").innerHTML = "-";
    document.getElementById("httpStatusCode").innerHTML = "-";
    document.getElementById("jsonConvert").innerHTML = "-";
    document.getElementById("resultCode").innerHTML = "-";
    document.getElementById("errorMessage").innerHTML = "-";
    document.getElementById("debugMessage").innerHTML = "-";
    document.getElementById("peakUsageMemory").innerHTML = "-";
    document.getElementById("cmdProcessTime").innerHTML = "-";
    document.getElementById("totalProcessTime").innerHTML = "-";
    document.getElementById("json").innerHTML = "-";
    document.getElementById("tableview").innerHTML = "";

    var user_agent = $("#user_agent").val();
    if (!user_agent || user_agent.trim().length <= 0) {
        user_agent = "KITTEN TEST CLIENT/0.0.0(DEVELOP,PC)";
    }

    var session_id = "";
    if ($("#ssid") && $("#ssid").length > 0) {
        session_id = $("#ssid").val();
    }

    var method = document.getElementById("method").value;
    var resource_name = document.getElementById("resource_name").value;
    var aurl = document.getElementById("url").value + document.getElementById("display_name").value + ".json";

    var data = _content_data[resource_name][method];
    if (data !== null && data["url_parameters"] !== null) {
        var is_first = true;
        for (var parameter_name in data["url_parameters"]) {
            var url_param = document.getElementById("url_param_" + parameter_name);
            if (url_param && url_param.value.trim().length > 0) {
                if (is_first) {
                    aurl += "?" + parameter_name + "=" + url_param.value.trim();
                    is_first = false;
                }
                else {
                    aurl += "&" + parameter_name + "=" + url_param.value.trim();
                }
            }
        }
    }

    var request_json = null;
    if (method != "GET") {
        request_json = create_request_message();
    }

    document.getElementById("requestUrl").innerHTML = method + " " + aurl;
    document.getElementById("requestMsg").innerHTML = request_json;

    var startTime = new Date().getTime();
    $.ajax({
        type: method,
        url: aurl,
        headers: {
            "Pragma": "no-cache",
            "Cache-Control": "no-cache",
            "Content-Type": "application/json",
            "X-Kitten-User-Agent": user_agent,
            "X-Kitten-SSID": session_id
        },
        data: request_json,
        error: function(request, status, thrown){
            enable("input");
            enable("select");
            enable("textarea");
        },
        success: function(data, status, xhr){
            enable("input");
            enable("select");
            enable("textarea");

            var endTime = new Date().getTime();
            var procTime = (endTime - startTime).toString() + " ms";
            var evalSuccess = "Failed.";

            var obj = null;
            if (data) {
                obj = data;
                evalSuccess = "Successful.";
            }

            var peakUsageMemory = obj && obj["memory"] ? obj["memory"]["peak_usage"] : 0;
            var peakUsageMemoryHtml = "";
            if (peakUsageMemory === 0) peakUsageMemoryHtml = "-";
            else if (peakUsageMemory >= WARNING_MEMORY) peakUsageMemoryHtml = "<span class='result_warning'>" + peakUsageMemory + " Byte</span>";
            else peakUsageMemoryHtml = peakUsageMemory + " Byte";

            document.getElementById("httpStatusCode").innerHTML = xhr.status;
            document.getElementById("jsonConvert").innerHTML = evalSuccess;
            document.getElementById("resultCode").innerHTML = obj && obj["result_code"] ? obj["result_code"] + " - " + get_api_result_text(obj["resultCode"]) : "-";
            document.getElementById("errorMessage").innerHTML = obj && obj["error_message"] ? obj["error_message"] : "-";
            document.getElementById("debugMessage").innerHTML = obj && obj["debug_message"] ? obj["debug_message"] : "-";
            document.getElementById("peakUsageMemory").innerHTML = peakUsageMemoryHtml;
            document.getElementById("cmdProcessTime").innerHTML = obj && obj["process_time"] ? obj["process_time"] + " ms" : "-";
            document.getElementById("totalProcessTime").innerHTML = procTime;
            document.getElementById("json").innerHTML =  escapeInputItem(xhr.responseText);
            document.getElementById("tableview").innerHTML = obj && obj["result"] ? create_result_table(obj["result"]) : "";

            // アカウントトークンを記憶(BOM_CUSTOMISE)
            if (obj && obj["result"] && obj["result"]["token"]) {
                var token = obj["result"]["token"];
                $.cookie("token", token, {expires:30,secure:false});
            }
            else if (obj && obj["result"] && obj["result"]["user_auth_data"] && obj["result"]["user_auth_data"]["token"]) {
                var token = obj["result"]["user_auth_data"]["token"];
                $.cookie("token", token, {expires:30,secure:false});
            }

            // セッションIDを記憶(BOM_CUSTOMISE)
            if (obj && obj["result"] && obj["result"]["ssid"]) {
                var ssid = obj["result"]["ssid"];
                $.cookie("ssid", ssid, {expires:30,secure:false});
            }
            else if (obj && obj["result"] && obj["result"]["user_auth_data"] && obj["result"]["user_auth_data"]["ssid"]) {
                var ssid = obj["result"]["user_auth_data"]["ssid"];
                $.cookie("ssid", ssid, {expires:30,secure:false});
            }

            // UserAgentを記憶
            $.cookie("user_agent", document.getElementById("user_agent").value, {expires:30,secure:false});

            $("#result_tree_view").treetable({ expandable: true });
        }
    });
}

function get_type_name(v) {
    if (typeof(v) === "undefined") {
        return "Undefined";
    }
    return typeof(v);
}

function get_safe_value(v) {
    if (v === null) {
        return "Null";
    }
    if (typeof(v) === "undefined") {
        return "Undefined";
    }
    return String(v);
}

/**
 * ツリービュー関連
 */
var _tvcount = 0;
function render_table_sub(parentid, data) {
    var html = "";
    if (data instanceof Array || data instanceof Object) {
        for (var i in data) {
            var id = ++_tvcount;
            var tt_id = parentid.length > 0 ? parentid + "-" + String(id) : String(id);

            if (parentid.length > 0) {
                html += "<tr data-tt-id=\"" + tt_id + "\" data-tt-parent-id=\"" + parentid + "\">";
            }
            else {
                html += "<tr data-tt-id=\"" + tt_id + "\">";
            }

            if (data[i] instanceof Array || data[i] instanceof Object) {
                html += "<td><span class='folder'>" + String(i) + "</span></td>";
                html += "<td>" + get_type_name(data[i]) + "</td>";
                html += "<td>-</td>";
                html += "</tr>";
                html += render_table_sub(tt_id, data[i]);
            }
            else {
                html += "<td><span class='file'>" + String(i) + "</span></td>";
                html += "<td>" + get_type_name(data[i]) + "</td>";
                html += "<td>" + get_safe_value(data[i]) + "</td>";
                html += "</tr>";
            }
        }
    }

    return html;
}

function expand_tree() {
    $("#result_tree_view").treetable('expandAll');
}

function collapse_tree() {
    $("#result_tree_view").treetable('collapseAll');
}

function create_result_table(data) {
    var html = "";
    _tvcount = 0;

    html += "<table id=\"result_tree_view\">";
    html += "<caption style=\"text-align:left\">";
    html += "<input type=\"button\" value=\"Expand all\" class=\"ui-button ui-corner-all ui-widget\" onclick=\"expand_tree(); return false;\" />";
    html += "<input type=\"button\" value=\"Collapse all\" class=\"ui-button ui-corner-all ui-widget\" onclick=\"collapse_tree(); return false;\" />";
    html += "</caption>";
    html += "<thead><tr>";
    html += "<th>Name or Array Index</th>";
    html += "<th id=\"tvtable_type_header\">Type</th>";
    html += "<th>Value</th>";
    html += "</tr></thead>";
    html += "<tbody>";
    html += render_table_sub("", data);
    html += "</tbody></table>";
    return html;
}




    $(function() {
        // リファレンスデータの読込
        load_reference_data();

        // 画面のリサイズ
        resize();
        window.onresize = resize;
    });

}());
