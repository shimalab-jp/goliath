// [CONST] リスト表示でのページングサイズ
var PAGING_SIZE = 10;

// [GLOBAL ARGS] モーダルダイアログのサイズ
var _dialogWidth = 0, _dialogHeight = 0;

/**********************************************************************
** [GLOBAL]
** 全置換します。
***********************************************************************
** Parameters
**   txt       : 置換対象文字列
**   replace   : 検索文字列
**   with_this : 置換文字列
** Return
**   置換後の文字列
**********************************************************************/
function replaceAll(txt, replace, with_this) {
    if (typeof txt == "string") {
        return txt.replace(new RegExp(replace, 'g'), with_this);
    }
    return txt;
}

/**********************************************************************
** [GLOBAL]
** 現在のページのQueryStringの値を取得します。
***********************************************************************
** Parameters
**   key : QueryStringのキー
** Return
**   keyに対応するQueryStringの値
**********************************************************************/
function getQueryString(key) {
    var returnValue = null;

    var k2 = key;
    k2 = k2.replace(/[\[]/,"\\\[");
    k2 = k2.replace(/[\]]/,"\\\]");

    var re = new RegExp("[\\?&]" + k2 + "=([^&#]*)");  
    var qs = re.exec(window.location.href);  
    if (qs) returnValue = qs[1];

    return returnValue;
}

/**********************************************************************
** [GLOBAL]
** 指定されたタグ名全ての disabled プロパティを true に設定します。
***********************************************************************
** Parameters
**   tagName : タグ名。
**********************************************************************/
function disable(tagName) {
    var elements = this.document.getElementsByTagName(tagName);
    for (var i = 0; i < elements.length; i++) {
        elements[i].disabled = true;
    }
}

/**********************************************************************
** [GLOBAL]
** 指定されたIDのタグの disabled プロパティを true に設定します。
***********************************************************************
** Parameters
**   id : ID
**********************************************************************/
function disableById(id) {
    var element = document.getElementById(id);
    if (element) {
        element.disabled = true;
    }
}

/**********************************************************************
** [GLOBAL]
** 指定されたタグ名全ての disabled プロパティを false に設定します。
***********************************************************************
** Parameters
**   tagName : タグ名。
**********************************************************************/
function enable(tagName) {
    var elements = this.document.getElementsByTagName(tagName);
    for (var i = 0; i < elements.length; i++) {
        elements[i].disabled = false;
    }
}

/**********************************************************************
** [GLOBAL]
** 指定されたIDのタグの disabled プロパティを false に設定します。
***********************************************************************
** Parameters
**   id : ID
**********************************************************************/
function enableById(id) {
    var element = document.getElementById(id);
    if (element) {
        element.disabled = false;
    }
}

/**********************************************************************
** [GLOBAL]
** 指定された値でラジオボタンを選択状態に設定します。
** 指定された値が存在しない場合は、既定値が選択状態に設定されます。
***********************************************************************
** Parameters
**   name         : タグ名
**   value        : 値
**   defaultValue : 既定値
**********************************************************************/
function selectRadio(name, value, defaultValue) {
    var elements = document.getElementsByName(name);
    for (var i = 0; i < elements.length; i++) {
        if (elements[i] && elements[i].value && elements[i].value == defaultValue) {
            elements[i].checked = true;
            break;
        }
    }
    for (var i = 0; i < elements.length; i++) {
        if (elements[i] && elements[i].value && elements[i].value == value) {
            elements[i].checked = true;
            break;
        }
    }
}

/**********************************************************************
** [GLOBAL]
** 指定された値でドロップダウンリストを選択状態に設定します。
** 指定された値が存在しない場合は、既定値が選択状態に設定されます。
***********************************************************************
** Parameters
**   name         : ドロップダウンリスト（selectタグ）のID
**   value        : 値
**   defaultValue : 既定値
**********************************************************************/
function selectDropdown(id, value, defaultValue) {
    var element = document.getElementById(id);
    for (var i = 0; i < element.options.length; i++) {
        if (element.options[i] && element.options[i].value && element.options[i].value == defaultValue) {
            element.selectedIndex = i;
            break;
        }
    }
    for (var i = 0; i < element.options.length; i++) {
        if (element.options[i] && element.options[i].value && element.options[i].value == value) {
            element.selectedIndex = i;
            break;
        }
    }
}

/**********************************************************************
** [GLOBAL]
** 指定された値でドロップダウンリストを選択状態に設定します。
** 指定された値が存在しない場合は、既定値が選択状態に設定されます。
***********************************************************************
** Parameters
**   name         : ドロップダウンリスト（selectタグ）のID
**   value        : 値
**   defaultValue : 既定値
**********************************************************************/
function selectDropdown(id, value, defaultValue) {
    var element = document.getElementById(id);
    for (var i = 0; i < element.options.length; i++) {
        if (element.options[i] && element.options[i].value && element.options[i].value == defaultValue) {
            element.selectedIndex = i;
            break;
        }
    }
    for (var i = 0; i < element.options.length; i++) {
        if (element.options[i] && element.options[i].value && element.options[i].value == value) {
            element.selectedIndex = i;
            break;
        }
    }
}

/**********************************************************************
** [GLOBAL]
** year/month/dayで個別指定された生年月日を表示形式に変換します。
***********************************************************************
** Parameters
**   year  : 年
**   month : 月
**   day   : 日
** Return
**   フォーマットされた生年月日。 yyyy-MM-dd (xx歳)
**********************************************************************/
function formatBirthday(year, month, day) {
    var returnValue = affixZero(year, 4) + "-" + affixZero(month, 2) + "-" + affixZero(day, 2);
    var b = parseInt(affixZero(year, 4) + affixZero(month, 2) + affixZero(day, 2));
    var today = new Date();
    var n = parseInt(affixZero(today.getFullYear(), 4) + affixZero(today.getMonth() + 1, 2) + affixZero(today.getDate(), 2));
    var age = parseInt((n - b) / 10000);
    returnValue = returnValue + " (" + age.toString() + "歳)";
    return returnValue;
}

/**********************************************************************
** [GLOBAL]
** 指定された数値が指定桁未満の場合、前にゼロを付与して指定桁にします。
***********************************************************************
** Parameters
**   number : 数値
**   length : 桁数
** Return
**   前にゼロが付与された文字列。
**   numberに 5, lengthに 3 が指定された場合、"005" が返ります。
**   number が length 以上の場合は number がそのまま返ります。
**********************************************************************/
function affixZero(number, length) {
    var returnValue = number.toString();
    if (returnValue.length < length) {
        for (var i = number.toString().length; i < length; i++) {
            returnValue = "0" + returnValue;
        }
    }
    return returnValue;
}

/**********************************************************************
** [GLOBAL]
** シリアル値の日付を yyyy-MM-dd HH:mm:ss 形式に変換します。
***********************************************************************
** Parameters
**   serialValue : シリアル値
** Return
**   yyyy-MM-dd HH:mm:ss
**********************************************************************/
function toFullDateTimeString(serialValue) {
    var dt = new Date(serialValue);
    return dt.getFullYear() + "-" + affixZero(dt.getMonth() + 1, 2) + "-" + affixZero(dt.getDate(), 2) + " " + affixZero(dt.getHours(), 2) + ":" + affixZero(dt.getMinutes(), 2) + ":" + affixZero(dt.getSeconds(), 2);
}

/**********************************************************************
** [GLOBAL]
** 入力項目用にエスケープ処理をします。
***********************************************************************
** Parameters
**   str : 文字列
** Return
**   エスケープ処理をした結果
**********************************************************************/
function escapeInputItem(str) {
    var returnValue = str;
    returnValue = replaceAll(returnValue, "&", "&amp;");
    returnValue = replaceAll(returnValue, "<", "&lt;");
    returnValue = replaceAll(returnValue, ">", "&gt;");
    returnValue = replaceAll(returnValue, "\"", "&quot;");
    returnValue = replaceAll(returnValue, "'", "&#39;");
    returnValue = replaceAll(returnValue, "\n", "<br>");
    returnValue = replaceAll(returnValue, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;");
    return returnValue;
}

/**********************************************************************
** [GLOBAL]
** 入力項目用にエスケープ処理を解除します。
***********************************************************************
** Parameters
**   str : 文字列
** Return
**   エスケープ処理を解除した結果
**********************************************************************/
function unescapeInputItem(str) {
    var returnValue = str;
    returnValue = replaceAll(returnValue, "<br>", "\n");
    returnValue = replaceAll(returnValue, "<br />", "\n");
    returnValue = replaceAll(returnValue, "&nbsp;&nbsp;&nbsp;&nbsp;", "\t");
    returnValue = replaceAll(returnValue, "&quot;", "\"");
    returnValue = replaceAll(returnValue, "&#39;", "'");
    returnValue = replaceAll(returnValue, "&lt;", "<");
    returnValue = replaceAll(returnValue, "&gt;", ">");
    returnValue = replaceAll(returnValue, "&amp;", "&");
    return returnValue;
}

/**********************************************************************
** [GLOBAL]
** リストを元にページングボタンを表示します。
***********************************************************************
** Parameters
**   maxPages         : 最大ページ数
**   page             : 現在のページ番号
**   count            : 総データ件数
**   listCount        : 表示データ件数
**   pagingMethodName : ページを移動する際に呼び出すメソッド名
**   outputId         : 出力先タグID
**********************************************************************/
function drawPaging(maxPages, page, count, listCount, pagingMethodName, outputId) {
    var colCnt = 0;
    var resultHtml = "<table cellpadding=\"2\" cellspacing=\"0\" class=\"paging\">";
    resultHtml += "<tr>";

    if (page <= 1) {
        resultHtml += "<td class=\"pagingPrev\">&lt;&lt;前へ</td>";
    }
    else {
        resultHtml += "<td class=\"pagingPrevActive\"><div onclick='" + pagingMethodName + "(" + (page - 1).toString() + ")'>&lt;&lt;前へ</div></td>";
    }
    colCnt++;

    for (var i = page - PAGING_SIZE; i < page; i++) {
        if (i > 0 && i != page) {
            resultHtml += "<td class=\"pagingNumber\"><div onclick='" + pagingMethodName + "(" + i.toString() + ")'>" + i.toString() + "</div></td>";
            colCnt++;
        }
    }

    resultHtml += "<td class=\"pagingCurrent\">" + page.toString() + "</td>";
    colCnt++;

    for (var j = page + 1; j <= maxPages && j <= page + PAGING_SIZE; j++) {
        if (j > 0 && j != page) {
            resultHtml += "<td class=\"pagingNumber\"><div onclick='" + pagingMethodName + "(" + j.toString() + ")'>" + j.toString() + "</div></td>";
            colCnt++;
        }
    }

    if (maxPages <= page) {
        resultHtml += "<td class=\"pagingNext\">次へ&gt;&gt;</td>";
    }
    else {
        resultHtml += "<td class=\"pagingNextActive\"><div onclick='" + pagingMethodName + "(" + (page + 1).toString() + ")'>次へ&gt;&gt;</div></td>";
    }
    colCnt++;
    resultHtml += "</tr>";

    resultHtml += "<tr>";
    resultHtml += "<td colspan='"+ colCnt.toString() + "' class=\"pageInfo\">";
    resultHtml += page.toString() + "/" + maxPages.toString() + " ページ (全" + count.toString() + "件中" + listCount.toString() + "件表示)"; 
    resultHtml += "</small></td>";
    resultHtml += "</tr>";
    resultHtml += "</table>";

    var resultArea = document.getElementById(outputId);
    if (resultArea) resultArea.innerHTML = resultHtml;
}

/**********************************************************************
** [GLOBAL]
** 指定したURLをモーダルダイアログで表示します。
***********************************************************************
** Parameters
**   url    : URL
**   width  : 幅
**   height : 高さ
**********************************************************************/
function showInlineDialog(url, width, height) {
    _dialogWidth = width;
    _dialogheight = height;
    resizeInlineDialog();

    var overlay = document.getElementById('overlay');
    var modalContent = document.getElementById('modalContent');
    var modalIFrame = document.getElementById('modalIFrame');
    document.body.style.overflow = "hidden";
    modalIFrame.src = url;
    overlay.style.display = "block";
    modalContent.style.display = "block";
    window.onresize = resizeInlineDialog;
}

/**********************************************************************
** [GLOBAL]
** モーダルダイアログのリサイズ処理。
**********************************************************************/
function resizeInlineDialog() {
    var overlayWidth = 0, overlayHeight = 0;
    var displayWidth = 0, displayHeight = 0;
    if (window.innerWidth) {
        overlayWidth = (document.body.clientWidth > window.innerWidth ? document.body.clientWidth : window.innerWidth);
        overlayHeight = (document.body.clientHeight > window.innerHeight ? document.body.clientHeight : window.innerHeight);
        displayWidth = window.innerWidth;
        displayHeight = window.innerHeight;
    }
    else {
        overlayWidth = (document.body.clientWidth > document.documentElement.clientWidth ? document.body.clientWidth : document.documentElement.clientWidth);
        overlayHeight = (document.body.clientHeight > document.documentElement.clientHeight ? document.body.clientHeight : document.documentElement.clientHeight);
        displayWidth = document.documentElement.clientWidth;
        displayHeight = document.documentElement.clientHeight;
    }

    var l = 0; t = 0, w = _dialogWidth, h = _dialogheight;
    l = (displayWidth / 2) - (_dialogWidth / 2);
    if (l < 0) {
        l = 0;
        w = displayWidth;
    }
    if (window.scrollX) {
        l += window.scrollX;
    }
    else {
        l += document.documentElement.scrollLeft;
    }

    t = (displayHeight / 2) - (_dialogheight / 2);
    if (t < 0) {
        t = 0;
        h = displayHeight;
    }
    if (window.scrollY) {
        t += window.scrollY;
    }
    else {
        t += document.documentElement.scrollTop;
    }

    var overlay = document.getElementById('overlay');
    var modalContent = document.getElementById('modalContent');
    //var modalIFrame = document.getElementById('modalIFrame');
    overlay.style.left = 0;
    overlay.style.top = 0;
    overlay.style.width = overlayWidth.toString() + "px";
    overlay.style.height = overlayHeight.toString() + "px";
    modalContent.style.left = l.toString() + "px";
    modalContent.style.top = t.toString() + "px";
    modalContent.style.width = w.toString() + "px";
    modalContent.style.height = h.toString() + "px";
}

/**********************************************************************
** [GLOBAL]
** モーダルダイアログを閉じます。
***********************************************************************
** Parameters
**   isReload : 親画面をリロードする場合は true 、それ以外の場合は false
**********************************************************************/
function closeInlineDialog(isReload) {
    if (isReload && reload) reload();
    var overlay = document.getElementById('overlay');
    var modalContent = document.getElementById('modalContent');
    var modalIFrame = document.getElementById('modalIFrame');
    overlay.style.display = "none";
    modalContent.style.display = "none";
    modalIFrame.src = "";
    document.body.style.overflow = "visible";
    window.onresize = null;
}
