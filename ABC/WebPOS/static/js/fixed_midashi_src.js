/* 2019-09-27 Chrome 61 bug fix */
function get_browser() {
	var ua=navigator.userAgent,tem,M=ua.match(/(opera|chrome|safari|firefox|msie|trident(?=\/))\/?\s*(\d+)/i) || [];
	if(/trident/i.test(M[1])){
		tem=/\brv[ :]+(\d+)/g.exec(ua) || [];
		return {name:'IE',version:(tem[1]||'')};
	}
	if(M[1]==='Chrome'){
		tem=ua.match(/\bOPR|Edge\/(\d+)/)
		if(tem!=null)   {return {name:'Opera', version:tem[1]};}
	}
	M=M[2]? [M[1], M[2]]: [navigator.appName, navigator.appVersion, '-?'];
	if((tem=ua.match(/version\/(\d+)/i))!=null) {M.splice(1,1,tem[1]);}
	return {
		name: M[0],
		version: M[1]
	};
}
var browserInfo = get_browser();
browserInfo.versionNo = Number(browserInfo.version);

/* Minified 化する前のソース */

/*
 * FixedMidashi JavaScript Library, version 1.9 (2015/11/15)
 * http://hp.vector.co.jp/authors/VA056612/fixed_midashi/
 * Copyright (C) 2012-2015 K.Koiso
 */

/* Exsample */
/*******************************************************************
<script type="text/javascript" src="hoge/fixed_midashi.js"></script>

<!-- body mode -->
<body onLoad="FixedMidashi.create();">
  <table _fixedhead="rows:1; cols:1">
    ...

<!-- div mode -->
<style type="text/css" media="screen">
  div.scroll_div { overflow: auto; }
</style>
<body onLoad="FixedMidashi.create();">
  <div class="scroll_div">
    <table _fixedhead="rows:1; cols:1">
      ...
*******************************************************************/

var FixedMidashi = new function()
{

////////////////////////////////////////////////////////////////////////////////////////////////////

var DISABLED = false; // 本機能の無効化フラグ

// 定数
// 各パラメータのデフォルト値
var ROWS = 1;
var COLS = 0;
var DIV_FULL_MODE = false;
var DIV_AUTO_SIZE = "both"; // both | width | height | none
var COPY_ID = true;
var BORDER_COLOR = null;
var BORDER_STYLE = null;
var BORDER_WIDTH = null;
var BOX_SHADOW = null;

var DIV_MIN_WIDTH = 150;
var DIV_MIN_HEIGHT = 150;

/*
 * div モードでも body のスクロールに反応するか否か(ver1.8)
 * 0: 反応しない
 * 1: div にスクロールバーがない場合のみ反応する
 * 2: 常に反応する(android の Safari だと正しい位置に表示されない)
 */
var DIV_BODY_SCROLL = 1;

var RADIO_PREFIX = "_FIXED_HEADER_";

var POS_FIXED    = 1;
var POS_ABSOLUTE = 2;
var POS_MIX      = 3;
var _positionMode = -1;

/*
 * Firefox で、[表示] → [ズーム] → [文字サイズだけ変更] にすると、
 * 拡大・縮小しても Resize イベントが通知されないため、リサイズ処理の契機がない。
 * ズーム変更の有無を次のタイマー間隔でチェックし、
 * 変更されていたらリサイズ処理を行うことにする。
 * ⇒ Safariも(ver1.3)
 */
var TIMER_WATCH_TABLESIZE = 3000;

////////////////////////////////////////////////////////////////////////////////////////////////////

var TID_HEADER = "H";
var TID_NUMBER = "N";
var TID_CORNER = "C";
var PX = "px";
var HEIGHT_MARGIN = 10;
var MIN_SIZE = 1;

/* ブラウザ情報 */
var _isIE = false;
var _IEver = 0;
var _isIE11 = false; // IE11以上(IE10までとは別のブラウザとして処理)
var _isFirefox = false;
var _isOpera = false;
var _isSafari = false;
var _isChrome = false;
var _isMobile = false;
// 「このブラウザだったら」という処理は、今となっては不要なものもある。

var _isBackCompat = false;
var _fixedHeaders = null;
var _fixedList = null;
var _body = null;
var _resizeTimerId = null;
var _execFlag = false;
var _IE_retryCount = 0;

////////////////////////////////////////////////////////////////////////////////////////////////////

/**
 * 見出し用のテーブルを作成
 */
this.create = function()
{
	if (DISABLED) return -1;
	if (!document.body.getBoundingClientRect) return -2; // 対象外ブラウザ
	if (!window.addEventListener && !window.attachEvent) return -3; // 対象外ブラウザ

	var stime = new Date().getTime();

	var isFirstCall = (_fixedHeaders == null);
	if (!isFirstCall) {
		// 2回目以降
		for (var i = 0; i < _fixedHeaders.length; i++) {
			_fixedHeaders[i].removeAllTables(false);
		}
	}

	//
	// _fixedhead が定義されている table を処理対象とする
	//
	var tables = document.body.getElementsByTagName("TABLE");
	var targetTables = new Array();
	for (var i = 0; i < tables.length; i++) {
		var table = tables[i];
		var opt = table.getAttribute("_fixedhead");
		if (opt == null) opt = table._fixedhead;
		if (opt == undefined) continue;
		if (table.rows.length == 0) continue;

		targetTables.push(table);
	}

	if (targetTables.length == 0) { // 対象テーブルなし
		return -4;
	}

	var dummy = null;
	if (isFirstCall) {
		// 初回
		var ua = navigator.userAgent.toLowerCase();
		_isIE = (ua.indexOf("msie") >= 0);
		if (_isIE) {
			var n = ua.indexOf("msie");
			var m = ua.indexOf(";", n);
			_IEver = Number(ua.substring(n + 5, m));
		}
		if (!_isIE) _isIE11 = (ua.indexOf("trident") >= 0); // IE11 以上の場合 _isIE は false で処理
		_isFirefox = (ua.indexOf("firefox") >= 0);
		_isOpera = (ua.indexOf("opera") >= 0); // TODO "OPR" になっているようだ
	//	_isOpera = (ua.indexOf(" opr/") >= 0);
		_isSafari = (ua.indexOf("safari") >= 0) && (ua.indexOf("chrome") < 0);
		_isChrome = (ua.indexOf("chrome") >= 0);
		_isMobile = (ua.indexOf("mobile") >= 0);

		//
		// (_isIE) RECTの値が桁違いの場合がある(2桁多い)。
		// frame を使用したページで発生する場合がある。
		// タイミングをずらして再実行すると正常に取得できる。
		//
		if (_isIE && (_rect(tables[0]).right >= 10000) && (_IE_retryCount < 10)) {
			setTimeout(FixedMidashi.create, 10);
			_IE_retryCount++;
			return -5;
		}

		_isBackCompat = (document.compatMode == "BackCompat");

		// body の clientWidth, clientHeight, scrollWidth, scrollHeight の取得元
		_body = _isBackCompat ? document.body : document.documentElement;

		if (_isIE && ((_IEver <= 7) || ((_IEver <= 9) && _isBackCompat))) {
			//
			// Operaは fixed をサポートしているが、スクロール中、画面がちらつく場合がある。
			// ⇒ ちらつかなくなっているので _isOpera は削除(ver1.3)
			// ただし、拡大表示するとちらつく。
			// IE10 は互換モードでも fixed が動くようだ。(ver1.4)
			// IE7(互換表示)の場合は、標準モードでも absolute にしておく。(ver1.6)
			//
			_positionMode = POS_ABSOLUTE; // "position: fixed" 未サポート
		} else
		if (_positionMode == -1) {
			_positionMode = POS_MIX; // default
			if (_isMobile) {
				//
				// タブレットの場合、POS_MIX だと見た目が良くない(POS_FIXED の方がまし)。
				// また、fixed にしたり absolute にしたりすると、正常に動作しないブラウザもあった(android Safari)。
				//
				_positionMode = POS_FIXED;
			}
		}

		if (_isIE && (_IEver == 8) && !_isBackCompat) { // (_isIE8 && 標準モード)
			_fixedList = new _FixedElementList();
		}
		if (_isFirefox) {
			dummy = _createObjectForFirefox();
		}
		if (_isChrome) {
			_createObjectForChrome();
		}

		// body.onResize のイベントリスナー設定
		_addEventListener(window, "resize", _onBodyResize);
		if (_isMobile) {
			// 画面回転イベント
			_addEventListener(window, "orientationchange", _onBodyResize);
		}

		if ((_isFirefox || _isSafari) && !_isMobile && (TIMER_WATCH_TABLESIZE >= 0)) {
			setInterval(_checkZoom, TIMER_WATCH_TABLESIZE);
		}

		// body.onScroll のイベントリスナー設定
		_addEventListener(window, "scroll", _onBodyScroll);
	}

	//
	// FixedHeader を生成して実行
	//
	_fixedHeaders = new Array();
	for (var i = 0; i < targetTables.length; i++) {
		var table = targetTables[i];
		var opt = table.getAttribute("_fixedhead");
		if (opt == null) opt = table._fixedhead;
		var fixedHeader = _createFixedHeader(table, opt, targetTables.length);
		_fixedHeaders.push(fixedHeader);
	}

	_execute("init");

	if (isFirstCall) {
		// 印刷時、固定テーブルを非表示にするスタイル
		// (_isIE) 先に行うと遅くなるようなので最後に行う。
		_createCSS("print", ".fixed_header_display_none_at_print { display: none; visibility: hidden; }");
	}

	if (dummy != null) {
		dummy.parentNode.removeChild(dummy);
	}

	return new Date().getTime() - stime; // 外部仕様は「戻り値なし」
};


/* Resize イベントが通知されない場合の処置 */
function _checkZoom()
{
	if (_fixedHeaders == null) return;
	for (var i = 0; i < _fixedHeaders.length; i++) {
		if (_fixedHeaders[i].checkZoom()) {
			_execute("resize");
			break;
		}
	}
}


/**
 * 作成した見出し用テーブルを削除
 */
this.remove = function()
{
	if (_fixedHeaders == null) return;
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].removeAllTables(true);
	}
	_fixedHeaders = new Array(); // null にはしない
};


/**
 * 「固定テーブルの入力オブジェクト」の値をソースと同じにする。
 * checkbox, radio, select, text
 */
this.syncValue = function(srcElement)
{
	if (_fixedHeaders == null) return;
	if (!srcElement) return;
	_copyValues(srcElement); // src の値を全 dst にコピー
};


/**
 * 「固定テーブルの要素」の style をソースと同じにする。
 */
this.syncStyle = function(srcElement, styleName)
{
	if (_fixedHeaders == null) return;
	if (!srcElement) return;

	if (_fixedList != null) { // (_isIE8 && 標準モード)
		var dst = _fixedList.getAll(srcElement);
		if (dst == null) return;
		for (var i = 0; i < dst.length; i++) {
			_copyStyle(srcElement, dst[i], styleName);
		}
		return;
	}

	var obj = srcElement.$FXH_FIXED_ELEMENT;
	if (obj == undefined) return;
	if (!obj.$IS_ARRAY) {
		_copyStyle(srcElement, obj, styleName);
	} else {
		for (var i = 0; i < obj.length; i++) {
			_copyStyle(srcElement, obj[i], styleName);
		}
	}
};


/**
 * 指定された要素(td, tr 等)に対する「固定テーブルの要素」を返す。
 * 固定行と固定列の重なり合うセルの場合の返却値は不定。
 */
this.getFixedElement = function(srcElement)
{
	if (_fixedHeaders == null) return null;
	if (!srcElement) return null;

	if (_fixedList != null) { // (_isIE8 && 標準モード)
		return _fixedList.get(srcElement);
	}

	var obj = srcElement.$FXH_FIXED_ELEMENT;
	if (!obj) return null;
	if (!obj.$IS_ARRAY) return obj;
	if (obj.length == 0) return null;
	return obj[obj.length - 1]; // 最後の要素
};


/**
 * 指定された要素(td, tr 等)に対する「固定テーブルの要素」を全て返す。
 */
this.getFixedElements = function(srcElement)
{
	if (_fixedHeaders == null) return null;
	if (!srcElement) return null;

	if (_fixedList != null) { // (_isIE8 && 標準モード)
		return _fixedList.getAll(srcElement);
	}

	var obj = srcElement.$FXH_FIXED_ELEMENT;
	if (!obj) return null;
	if (obj.$IS_ARRAY && (obj.length == 0)) return null;

	var array = new Array();
	if (!obj.$IS_ARRAY) {
		array.push(obj);
	} else {
		for (var i = 0; i < obj.length; i++) {
			array.push(obj[i]);
		}
	}
	return array;
};


/**
 * 指定された「固定テーブルの要素(td, tr 等)」に対する元要素を返す。
 */
this.getSourceElement = function(element)
{
	if (_fixedHeaders == null) return null;
	if (!element) return null;
	var src = element.$SOURCE_ELEMENT;
	return !src ? null : src;
};


/**
 * 指定された要素(td, tr 等)が「固定テーブルの要素」か否かを返す。
 */
this.isFixedElement = function(element)
{
	if (_fixedHeaders == null) return false;
	if (!element) return false;
	return (element.$SOURCE_ELEMENT != undefined);
};

////////////////////////////////////////////////////////////////////////////////////////////////////
// 以降、private function

/*
 * _FixedHeader インスタンスの生成
 */
function _createFixedHeader(tblSource, opt, tableCount)
{
	// テーブルのセル数
	var trList = tblSource.rows;
	var tdList = _cells(trList[0]);
	var columnCount = 0;
	for (var i = 0; i < tdList.length; i++) {
		columnCount += tdList[i].colSpan;
	}

	// パラメータ(_fixedhead 属性)初期値
	var fixedRows = ROWS;
	var fixedCols = COLS;
	var divMaxWidth = -1;
	var divMaxHeight = -1;
	var copyId = COPY_ID;
	var borderColor = BORDER_COLOR;
	var borderStyle = BORDER_STYLE;
	var borderWidth = BORDER_WIDTH;
	var boxShadow = BOX_SHADOW;
	var bgcolor = null;
	var opacity = 1.0;
	var divFullMode = DIV_FULL_MODE;
	var divAutoSize = DIV_AUTO_SIZE;
	var bodyHeaderId = null;
	var leftHeaderId = null;

	var thead = _getElementByTagName(tblSource, "THEAD");
	if (thead != null) { // THEAD がある場合は、その行数が初期値
		fixedRows = thead.rows.length;
	}

	// パラメータ(_fixedhead 属性)解析
	var opts = opt.split(";");
	for (var i = 0; i < opts.length; i++) {
		var buf = opts[i].split(":");
		if (buf.length != 2) continue;
		var name = _trim(buf[0]).toLowerCase();
		var value = _trim(buf[1]);

		switch (name) {
		case "rows": fixedRows = Number(value); break;
		case "cols": fixedCols = Number(value); break;
		case "div-max-width": divMaxWidth = _percent(value); break;
		case "div-max-height": divMaxHeight = _percent(value); break;
		case "div-full-mode": divFullMode = (value.toLowerCase() == "yes"); break;
		case "div-auto-size": divAutoSize = value.toLowerCase(); break; // (ver1.8)
		case "copyid": copyId = (value.toLowerCase() == "yes"); break;
		case "border-color": borderColor = value; break;
		case "border-style": borderStyle = value; break;
		case "border-width": borderWidth = value; break;
		case "box-shadow": boxShadow = value.replace(/ +/g, " "); break; // (ver1.9)
		case "bgcolor": bgcolor = value; break;
		case "opacity": opacity = Number(value); break;
		case "body-header-id": bodyHeaderId = value; break; // (ver1.7)
		case "body-left-header-id": leftHeaderId = value; break; // (ver1.8)
		}
	}

	if ((fixedRows < 0) || (fixedRows >= Math.min(trList.length, 11))) fixedRows = 0;
	if ((fixedCols < 0) || (fixedCols >= Math.min(columnCount, 11))) fixedCols = 0;
	if (!divMaxWidth || (divMaxWidth < 1) || (divMaxWidth > 100)) divMaxWidth = -1;
	if (!divMaxHeight || (divMaxHeight < 1) || (divMaxHeight > 100)) divMaxHeight = -1;
	if (tableCount > 1) divFullMode = false;
	if ((opacity < 0) || (opacity > 1)) opacity = 1.0;

	var testDiv = document.createElement("DIV");
	if ((borderColor != null) && !_setStyle(testDiv, "borderColor", borderColor)) {
		borderColor = null;
	}
	if ((borderStyle != null) && !_setStyle(testDiv, "borderStyle", borderStyle)) {
		borderStyle = null;
	}
	if ((borderWidth != null) && !_setStyle(testDiv, "borderWidth", borderWidth)) {
		borderWidth = null;
	}
	if ((boxShadow != null) && !_setStyle(testDiv, "boxShadow", boxShadow)) {
		boxShadow = null;
	}
	if ((bgcolor != null) && !_setStyle(testDiv, "backgroundColor", bgcolor)) {
		bgcolor = null;
	}

	var bodyHeader = null;
	var leftHeader = null;
	if ((bodyHeaderId != null) && !(_isIE && ((_IEver <= 7) || ((_IEver <= 9) && _isBackCompat)))) {
		bodyHeader = document.getElementById(bodyHeaderId);
	}
	if ((leftHeaderId != null) && !(_isIE && ((_IEver <= 7) || ((_IEver <= 9) && _isBackCompat)))) {
		leftHeader = document.getElementById(leftHeaderId);
	}

	// 親が div で、その overflow が auto または scroll なら div モード
	var divSource = null;
	if (tblSource.parentNode.tagName == "DIV") {
		var div = tblSource.parentNode;
		var style = div.currentStyle || document.defaultView.getComputedStyle(div, '');
		if (((style.overflowX == "auto") || (style.overflowX == "scroll")) ||
		    ((style.overflowY == "auto") || (style.overflowY == "scroll"))) {
			divSource = div;
		}
	}

//	if ((divSource != null) && _isIE) {
//		divSource.style.paddingRight = "0px";
//		divSource.style.marginRight = "0px";
//	}

	return new _FixedHeader(divSource, tblSource, columnCount,
		fixedRows, fixedCols, divMaxWidth, divMaxHeight, divFullMode, divAutoSize,
		copyId, borderColor, borderStyle, borderWidth, boxShadow, bgcolor, opacity, bodyHeader, leftHeader);
}


/*
 * body の onScroll のイベントハンドラ
 */
function _onBodyScroll()
{
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].onBodyScroll();
	}
}


/*
 * body の onResize のイベントハンドラ
 */
function _onBodyResize()
{
	//
	// (_isIE && 標準モード)プログラム内での divSource のサイズ変更の度に
	// body の onResize が発生してしまうため、
	// 本プログラム実行中の body の onResize は無視する。
	// この制御をしないと無限ループする。
	//
	if (_execFlag) return; // 処理中

	// resize イベントは連続的に発生する可能性があるのですぐ処理しない。
	if (_resizeTimerId != null) clearTimeout(_resizeTimerId);
	var func = function() { _execute("resize"); };
	_resizeTimerId = setTimeout(func, 500);
}


/*
 * 全対象テーブルに対して、見出し用のテーブルを作成。
 * (create() および body のりサイズ時に呼ばれる)
 */
function _execute(caller)
{
	_execFlag = true;
	_resizeTimerId = null;

	// 初期設定
	for (var i = 0; i < _fixedHeaders.length; i++) {
		if (caller == "resize") { // onResize 時
			_fixedHeaders[i].initOnResize();
		} else {
			_fixedHeaders[i].init();
		}
	}

	// div の自動リサイズ
	_resizeSourceDiv();

	// 見出し用テーブル作成
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].main();
	}

	_execFlag = false;
}


/*
 * div の自動りサイズ処理。
 * div が body をはみでないようにする。
 */
function _resizeSourceDiv()
{
	var existsDivMode = false;

	// (1) 全 div を非表示にする
	for (var i = 0; i < _fixedHeaders.length; i++) {
		if (_fixedHeaders[i].hideAllDivs(true)) {
			existsDivMode = true;
		}
	}

	if (!existsDivMode) return; // div モードのテーブルなし

	//
	// (2) div-max-width, div-max-height パラメータの処理
	//
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].resizeSourceDiv1();
	}

	//
	// この時点で body をはみ出ていない場合は、
	// resizeDivWidth() のために overflowX を hidden にしておく。
	//
	var overflowX = document.body.style.overflowX; // 保存
	var bodyStyle = document.body.currentStyle || document.defaultView.getComputedStyle(document.body, '');
	if ((_body.scrollWidth <= _body.clientWidth) && (bodyStyle.overflowX != "scroll")) {
		document.body.style.overflowX = "hidden";
	}

	//
	// (3) div の自動りサイズ
	//
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].resizeDivHeight();
	}
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].resizeDivWidth();
	}

	document.body.style.overflowX = overflowX; // 戻す

	//
	// (4) div サイズの微調整
	//
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].resizeSourceDiv2();
	}

	// (5) 全 div を再表示する
	for (var i = 0; i < _fixedHeaders.length; i++) {
		_fixedHeaders[i].hideAllDivs(false);
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * _FixedHeader 内部クラス
 */
function _FixedHeader(_divSource, _tblSource, _columnCount,
	_fixedRows, _fixedCols, _divMaxWidth, _divMaxHeight, _divFullMode, _divAutoSize,
	_copyId, _borderColor, _borderStyle, _borderWidth, _boxShadow, _bgcolor, _opacity, _bodyHeader, _leftHeader)
{
	var _isDivMode = (_divSource != null);
	var _tblHeader = null;
	var _tblNumber = null;
	var _tblCorner = null;

	var _tableWidth = 0;
	var _tableHeight = 0;
	var _colsWidthInfo = "";
	var _rowsHeightInfo = "";
	var _timerId_adjustLeft = null;
	var _timerId_adjustTop = null;
	var _cellPadding = new Array();
	var _cellBorder = new Array();
	var _cellPaddingV = 0;
	var _maxCellPadding = 0;
	var _zIndex = (_tblSource.style.zIndex ? _tblSource.style.zIndex : 0);

	// div mode 専用変数
	var _autoWidth = false;
	var _autoHeight = false;
	var _maxWidth = 0;
	var _maxHeight = 0;
	var _divLeft = 0;
	var _divTop = 0;
	var _divMarginLeft = 0;
	var _divMarginTop = 0;
	var _tblMarginLeft = 0;
	var _tblMarginTop = 0;
	var _divPaddingWidth = 0;
	var _divPaddingHeight = 0;
	var _initTableWidth = 0;
	var _tableDisplay = _tblSource.style.display;
	var _divScrollLeft = -1;
	var _divScrollTop = -1;
	var _tblDummy = null;

	// body mode 専用変数
	var _bodyScrollLeft = -1;
	var _bodyScrollTop = -1;
	// fjn 201706
	var _bodyTop = 101;
	var _bodyLeft = 0;
	var _timerId_bodyScroll = null;
	var _timerId_bodyScroll_abs = null;

	/*
	 * 初期処理
	 */
	this.init = function()
	{
		// 背景色が透明の場合は不透明にする
		if (_bgcolor == null) {
			var source = (_isDivMode ? _divSource : _tblSource);
			var bgcolor = _getBackgroundColor(source);
			if (bgcolor == null) { // 透明の場合は親と同じ色にする
				var node = source.parentNode;
				while (node) {
					bgcolor = _getBackgroundColor(node);
					if (bgcolor != null) break;
					if (node.tagName == "HTML") break;
					node = node.parentNode;
				}
				if (bgcolor == null) bgcolor = "white";
			}
			_bgcolor = bgcolor;
		}

		if (_isDivMode) {
			// div.onScroll のイベントリスナー設定
			_addEventListener(_divSource, "scroll", _onDivScroll);

			// divSource の padding サイズ等を取得
			_getDivStyle();

			_initTableWidth = _offsetWidth(_tblSource);
		}

		if (_bodyHeader) _bodyTop = Math.max(_rect(_bodyHeader).bottom, 0);
		if (_leftHeader) _bodyLeft = Math.max(_rect(_leftHeader).right, 0);
		_getCellPadding();
	};


	/*
	 * onResize 時の初期処理
	 */
	this.initOnResize = function()
	{
		_bodyScrollTop = -1;
		_bodyScrollLeft = -1;
		_divScrollTop = -1;
		_divScrollLeft = -1;
		if (_bodyHeader) _bodyTop = Math.max(_rect(_bodyHeader).bottom, 0);
		if (_leftHeader) _bodyLeft = Math.max(_rect(_leftHeader).right, 0);

		if (_isDivMode) {
			// div のサイズをクリア。
			if (_autoWidth || (_divMaxWidth > 0)) _divSource.style.width = "";
			if (_autoHeight || (_divMaxHeight > 0)) _divSource.style.height = "";

			_initTableWidth = _offsetWidth(_tblSource);
		}

		// ブラウザによっては、拡大するとセルパディングが微妙に変わるので取得しなおす。
		_getCellPadding();
	};


	/*
	 * 全ての固定テーブルを削除
	 */
	this.removeAllTables = function(clearDivSize)
	{
		if (_tblHeader) _removeTable(_tblHeader);
		if (_tblNumber) _removeTable(_tblNumber);
		if (_tblCorner) _removeTable(_tblCorner);
		_tblHeader = null;
		_tblNumber = null;
		_tblCorner = null;

		if (_isDivMode) {
			_removeEventListener(_divSource, "scroll", _onDivScroll);
			if (clearDivSize) {
				if (_autoWidth || (_divMaxWidth > 0)) _divSource.style.width = "";
				if (_autoHeight || (_divMaxHeight > 0)) _divSource.style.height = "";
			}
		}
	};

////////////////////////////////////////////////////////////////////////////////////////////////////
// div の自動リサイズ処理

	/*
	 * 全ての div の表示・非表示制御
	 */
	this.hideAllDivs = function(hide)
	{
		if (!_isDivMode) return false; // body モードでも呼ばれるのでチェック

		var divHeader = (_tblHeader == null) ? null : _tblHeader.parentNode;
		var divNumber = (_tblNumber == null) ? null : _tblNumber.parentNode;
		var divCorner = (_tblCorner == null) ? null : _tblCorner.parentNode;

		if (hide) {
			if (divHeader) _setClientWidth(divHeader, DIV_MIN_WIDTH);
			if (divNumber) _setClientHeight(divNumber, DIV_MIN_HEIGHT);
		}

		var display = hide ? "none" : "";
		if (divHeader) divHeader.style.display = display;
		if (divNumber) divNumber.style.display = display;
		if (divCorner) divCorner.style.display = display;

		if (_autoWidth || _autoHeight || (_divMaxWidth > 0) || (_divMaxHeight > 0)) {
			_hideSourceTable(hide); // (_isIE || _isFirefox || _isOpera)
		}

		if (hide) {
			if (_autoWidth) _setOffsetWidth(_divSource, DIV_MIN_WIDTH);
			if (_autoHeight) _setOffsetHeight(_divSource, DIV_MIN_HEIGHT);
		}

		return true;
	};


	/*
	 * div-max-width, div-max-height パラメータの処理
	 */
	this.resizeSourceDiv1 = function()
	{
		if (!_isDivMode) return; // body モードでも呼ばれるのでチェック
		var tblSource = (_tblDummy != null) ? _tblDummy : _tblSource;

		if (_divMaxWidth > 0) {
			_maxWidth = _body.clientWidth / 100 * _divMaxWidth;
			_maxWidth = Math.max(_maxWidth, DIV_MIN_WIDTH);
			var width = _offsetWidth(tblSource) + _divPaddingWidth;
			width = Math.min(width, _maxWidth);
			_setOffsetWidth(_divSource, width);
		}

		if (_divMaxHeight > 0) {
			_maxHeight = _body.clientHeight / 100 * _divMaxHeight;
			_maxHeight = Math.max(_maxHeight, DIV_MIN_HEIGHT);
			var height = _offsetHeight(tblSource) + _divPaddingHeight;
			height = Math.min(height, _maxHeight);
			_setOffsetHeight(_divSource, height);
		}
	};


	/*
	 * body からはみ出ないように元 div の高さを設定。
	 * (hideAllDivs(true) が呼ばれた状態で呼ばれる)
	 */
	this.resizeDivHeight = function()
	{
		if (!_isDivMode) return;
		if (!_autoHeight) return;

		//
		// body からはみ出ない最大の高さを取得
		//
		if (_body.scrollHeight > _body.clientHeight) {
			// 対象 div が最小の高さでも body をはみ出ている場合は、
			// 現在の body.scrollHeight を変えず、かつ
			// body.clientHeight を超えない最大サイズを求める。
			var over1 = Math.max(_body.scrollHeight - _body.clientHeight, 0);
			_divSource.style.height = Math.max(_body.clientHeight - 30, MIN_SIZE) + PX;
			var over2 = Math.max(_body.scrollHeight - _body.clientHeight, 0);
			_maxHeight = _offsetHeight(_divSource) - (over2 - over1);
		} else {
			_divSource.style.height = _body.clientHeight + PX;
			var over = Math.max(_body.scrollHeight - _body.clientHeight, 0);
			_maxHeight = _offsetHeight(_divSource) - over;
		}

		_maxHeight--; // 拡大表示対応
		_maxHeight = Math.max(_maxHeight, DIV_MIN_HEIGHT);

		//
		// div の高さを設定
		//
		var tblSource = (_tblDummy != null) ? _tblDummy : _tblSource;
		var height = _offsetHeight(tblSource) + _divPaddingHeight;
		height = Math.min(height, _maxHeight);
		if (_divFullMode) height = _maxHeight;
		_setOffsetHeight(_divSource, height);
	};


	/*
	 * body からはみ出ないように元 div の幅を設定。
	 * 高さ設定後に呼ばれる。
	 */
	this.resizeDivWidth = function()
	{
		if (!_isDivMode) return; // body モードでも呼ばれるのでチェック
		if (!_autoWidth) return;

		//
		// body からはみ出ない最大の幅を取得
		//
		if (_body.scrollWidth > _body.clientWidth) {
			// 対象 div が最小幅でも body をはみ出ている場合は、
			// 現在の body.scrollWidth を変えず、かつ
			// body.clientWidth を超えない最大サイズを求める。
			var over1 = Math.max(_body.scrollWidth - _body.clientWidth, 0);
			_divSource.style.width = Math.max(_body.clientWidth - 16, MIN_SIZE) + PX;
			var over2 = Math.max(_body.scrollWidth - _body.clientWidth, 0);
			_maxWidth = _offsetWidth(_divSource) - (over2 - over1);
		} else {
			_divSource.style.width = _body.clientWidth + PX;
			var over = Math.max(_body.scrollWidth - _body.clientWidth, 0);
			_maxWidth = _offsetWidth(_divSource) - over;
		}

		_maxWidth--; // 拡大表示対応
		_maxWidth = Math.max(_maxWidth, DIV_MIN_WIDTH);

		//
		// div の幅を設定
		//
		var tblSource = (_tblDummy != null) ? _tblDummy : _tblSource;
		var width = _offsetWidth(tblSource) + _divPaddingWidth;
		width = Math.min(width, _maxWidth);
		if (_divFullMode) width = _maxWidth;
		if (_isIE) width--;
		_setOffsetWidth(_divSource, width);
	};


	/*
	 * 元 div サイズの微調整
	 */
	this.resizeSourceDiv2 = function()
	{
		if (!_isDivMode) return; // body モードでも呼ばれるのでチェック

		var tblSource = (_tblDummy != null) ? _tblDummy : _tblSource;

		//
		// 幅調整 (div の overflow:auto 対応)。
		// 幅にまだ余裕があるのに水平スクロールバーがある場合はその分広げる。
		//
		if ((_autoWidth || (_divMaxWidth > 0)) &&
		    (_divSource.scrollWidth > _divSource.clientWidth) &&
		    (_offsetWidth(_divSource) < _maxWidth)) {
			var over = _divSource.scrollWidth - _divSource.clientWidth;
			var width = Math.min(_offsetWidth(_divSource) + over, _maxWidth);
			var o = _divSource.style.overflowY;
			if (_isIE) _divSource.style.overflowY = "hidden"; // おまじない
			_setOffsetWidth(_divSource, width);
			if (_isIE) _divSource.style.overflowY = o; // 戻す
		}

		//
		// 高さ調整 (div の overflow:auto 対応)。
		// 高さにまだ余裕があるのに垂直スクロールバーがある場合はその分広げる。
		//
		if ((_autoHeight || (_divMaxHeight > 0)) &&
		    (_divSource.scrollHeight > _divSource.clientHeight) &&
		    (_offsetHeight(_divSource) < _maxHeight)) {
			var over = _divSource.scrollHeight - _divSource.clientHeight;
			var height = Math.min(_offsetHeight(_divSource) + over, _maxHeight);
			_setOffsetHeight(_divSource, height);
		}

		//
		// div の幅変更により、table 幅が小さくなることがあるので、
		// 幅にまだ余裕がある場合は、小さくなった分大きくする。
		//
		if ((_autoWidth || (_divMaxWidth > 0)) &&
		    (_initTableWidth > _offsetWidth(tblSource)) &&
		    (_offsetWidth(_divSource) < _maxWidth)) {
			var diff = _initTableWidth - _offsetWidth(tblSource);
			var width = Math.min(_offsetWidth(_divSource) + diff, _maxWidth);
			_setOffsetWidth(_divSource, width);
		}
	};


	/*
	 * 元テーブルのサイズ、列幅、行の高さが変わっていたらリサイズ処理を呼び出す
	 * _isFirefox || _isSafari
	 */
	this.checkZoom = function()
	{
		if ((Math.abs(_offsetWidth(_tblSource) - _tableWidth) >= 1) ||
		    (Math.abs(_offsetHeight(_tblSource) - _tableHeight) >= 1) ||
		    (_colsWidthList(_cells(_tblSource.rows[0])) != _colsWidthInfo) ||
		    (_rowsHeightList(_tblSource.rows) != _rowsHeightInfo)) {
			return true;
		}
	};

////////////////////////////////////////////////////////////////////////////////////////////////////

	/*
	 * メイン処理
	 */
	this.main = function()
	{
		// (1) 前回実行時と現在の元テーブルのサイズ比較
		var w = _offsetWidth(_tblSource);
		var h = _offsetHeight(_tblSource);
		var widthChanged = (_tableWidth != w);
		var heightChanged = (_tableHeight != h);
		_tableWidth = w;
		_tableHeight = h;

		// 拡大表示すると、テーブル幅は同じでも列幅が変わる場合がある
		var cw = _colsWidthList(_cells(_tblSource.rows[0]));
		if (_colsWidthInfo != cw) {
			_colsWidthInfo = cw;
			widthChanged = true;
		}
		// 各行の高さも念のためチェック
		var rh = _rowsHeightList(_tblSource.rows);
		if (_rowsHeightInfo != rh) {
			_rowsHeightInfo = rh;
			heightChanged = true;
		}

		// (2) 固定テーブルの作成
		var created = _createTables(widthChanged || heightChanged);

		// (3) 固定テーブルのリサイズ
		if (created || widthChanged || heightChanged) {
			_resizeTables();
		}

		// (4) 固定テーブルの親 div のリサイズ
		if (_isDivMode) {
			_resizeDivs();
		}

		// (5) スクロール
		if (_isDivMode) {
			_scrollDivByFixed();
			_onDivScroll();
		} else {
			_onBodyScroll2(true);
		}
	};


	/*
	 * 各固定テーブルの作成。不要な場合は作成しない。
	 */
	function _createTables(sizeChanged)
	{
		//
		// 元テーブルの幅や高さが変わった場合、固定テーブルを一旦削除する。
		//
		if (sizeChanged) {
			if (_tblHeader != null) _removeTable(_tblHeader);
			if (_tblNumber != null) _removeTable(_tblNumber);
			if (_tblCorner != null) _removeTable(_tblCorner);
			_tblHeader = null;
			_tblNumber = null;
			_tblCorner = null;
		}

		var created = false;

		var srcTrList = _tblSource.rows;
		var parent = (_isDivMode ? _tblSource.parentNode : _body);
		var scrV = false;
		var scrH = false;

		if (_isDivMode) {
			scrV = ((_divSource.clientHeight < _divSource.scrollHeight) &&
				(_divSource.clientHeight != 0));
			scrH = ((_divSource.clientWidth < _divSource.scrollWidth) &&
				(_divSource.clientWidth != 0));
			// (_isIE) サイズを設定しないと、clientWidth, clientHeight は 0 が返る
		}

		// div モードの場合も、body のスクロールサイズをチェックするようにした
		if (!_isDivMode || (DIV_BODY_SCROLL != 0)) {
			if (!scrV) scrV = (_body.clientHeight < _body.scrollHeight);
			if (!scrH) scrH = (_body.clientWidth < _body.scrollWidth);

			if (!scrV && !scrH && (_body == document.documentElement)) {
				// 標準モードの場合に、documentElement の scrollHeight, scrollWidth の値が
				// 不正な場合があるため、document.body の値で再チェックする(暫定処置 ver1.7)。
				// 不正になる条件は不明。
				scrV = (document.body.clientHeight < document.body.scrollHeight);
				scrH = (document.body.clientWidth < document.body.scrollWidth);
			}
		}

		// 表示領域が狭く、固定部しか表示できない場合は固定テーブルを作成しない。
		if ((_fixedRows > 0) && scrV) {
			if (_rowsHeight(srcTrList, _fixedRows) + 30 >= parent.clientHeight) {
				scrV = false;
			}
		}
		if ((_fixedCols > 0) && scrH) {
			var srcTdList = _cells(srcTrList[0]);
			if (_colsWidth(srcTdList, _fixedCols) + 30 >= parent.clientWidth) {
				scrH = false;
			}
		}

		//
		// 横のスクロールバーがある場合、
		// 「行番号」のテーブルを作成し、元の行番号列の上に重ねて配置する。
		//
		if ((_fixedCols > 0) && scrH) {
			if (_tblNumber == null) {
				_tblNumber = _createTable(TID_NUMBER, srcTrList.length, _fixedCols);
				if (_tblNumber != null) created = true;
			}
		} else {
			if (_tblNumber != null) {
				_removeTable(_tblNumber);
				_tblNumber = null;
			}
		}

		//
		// 縦のスクロールバーがある場合、
		// 「見出し」のテーブルを作成し、元の見出し行の上に重ねて配置する。
		// CAPTION がある場合を考慮し、「行番号」より「見出し」を上に配置するため、後に作成。
		//
		if ((_fixedRows > 0) && scrV) {
			if (_tblHeader == null) {
				_tblHeader = _createTable(TID_HEADER, _fixedRows, _columnCount);
				if (_tblHeader != null) created = true;
			}
		} else {
			if (_tblHeader != null) {
				_removeTable(_tblHeader);
				_tblHeader = null;
			}
		}

		//
		// 見出しテーブルと行番号テーブルの両方がある場合、
		// 「左上隅」のテーブルを作成する。
		// 最後に作成するため zIndex の制御は不要。
		//
		if ((_tblHeader != null) && (_tblNumber != null)) {
			if (_tblCorner == null) {
				_tblCorner = _createTable(TID_CORNER, _fixedRows, _fixedCols);
				created = true;
			}
		} else {
			if (_tblCorner != null) {
				_removeTable(_tblCorner);
				_tblCorner = null;
			}
		}

		return created;
	}


	/*
	 * 各固定テーブルのサイズを元テーブルに合わせて設定。
	 */
	function _resizeTables()
	{
		// (1) 見出し
		if (_tblHeader != null) {
			_setOffsetWidth(_tblHeader, _offsetWidth(_tblSource));
			_checkHeaderHeight(_tblHeader);
		}

		// (2) 行番号
		if (_tblNumber != null) {
			var srcWidth = _colsWidth(_cells(_tblSource.rows[0]), _fixedCols);
			var dstWidth = _colsWidth(_cells(_tblNumber.rows[0]), _fixedCols);
			var diff = dstWidth - srcWidth;
			if (diff != 0) _setOffsetWidth(_tblNumber, _offsetWidth(_tblNumber) - diff);
			_checkHeaderHeight(_tblNumber);
			_setOffsetHeight(_tblNumber, _offsetHeight(_tblSource));
		}

		// (3) 左上隅
		if (_tblCorner != null) {
			_setOffsetWidth(_tblCorner, _offsetWidth(_tblNumber));
			_checkHeaderHeight(_tblCorner);
			_setOffsetHeight(_tblCorner, _offsetHeight(_tblHeader));
		}
	}


	/*
	 * 各固定テーブルの親 div のサイズを元テーブルに合わせて設定。
	 */
	function _resizeDivs()
	{
		var gap = 1; // 表示領域が小さい場合に、スクロールしていることがわかるようにするための隙間のサイズ

		var divHeader = (_tblHeader == null) ? null : _tblHeader.parentNode;
		var divNumber = (_tblNumber == null) ? null : _tblNumber.parentNode;
		var divCorner = (_tblCorner == null) ? null : _tblCorner.parentNode;

		if (divHeader != null) _checkDivPosition(divHeader);
		if (divNumber != null) _checkDivPosition(divNumber);
		if (divCorner != null) _checkDivPosition(divCorner);

		// (1) 見出し
		if (divHeader != null) {
			_setClientWidth(divHeader, _divSource.clientWidth);
			_setClientHeight(divHeader, _offsetHeight(_tblHeader));

			var marginTop = _rect(_tblHeader).bottom - _rect(divHeader).bottom;
			if (marginTop > 0) { // table の marginTop 分
				_setClientHeight(divHeader, _offsetHeight(_tblHeader) + marginTop);
			}

			// 元データの高さが見出しより小さい場合は、元データと同じ高さにする。
			if (divHeader.clientHeight >= _divSource.clientHeight) {
				_setClientHeight(divHeader, _divSource.clientHeight - gap);
			}
		}

		// (2) 行番号
		if (divNumber != null) {
			_setClientHeight(divNumber, _divSource.clientHeight);
			_setClientWidth(divNumber, _offsetWidth(_tblNumber));

			var marginLeft = _rect(_tblNumber).right - _rect(divNumber).right;
			if (marginLeft > 0) { // table の marginLeft 分
				_setClientWidth(divNumber, _offsetWidth(_tblNumber) + marginLeft);
			}

			// 元データの幅が行番号より小さい場合は、元データと同じ幅にする。
			if (divNumber.clientWidth >= _divSource.clientWidth) {
				_setClientWidth(divNumber, _divSource.clientWidth - gap);
			}
		}

		// (3) 左上隅
		if (divCorner != null) {
			_setClientHeight(divCorner, divHeader.clientHeight);
			_setClientWidth(divCorner, divNumber.clientWidth);
		}
	}

////////////////////////////////////////////////////////////////////////////////////////////////////

	/*
	 * (1) div の 位置を取得。
	 * (2) div の幅・高さが指定されているかをチェック。
	 * (3) div の border + padding + 子(table)の margin の合計サイズを取得。
	 * div の overflow が scroll の場合は、スクロールバーの幅も含む。
	 */
	function _getDivStyle()
	{
		// 元の div と table を一時的にコピーして調べる。
		var divTest = _divSource.cloneNode(false);
		var tblTest = _tblSource.cloneNode(false);

		divTest.style.position = "absolute";
		divTest.style.left = "0px";
		divTest.style.top = "0px";
		divTest.style.minWidth = "1px";
		divTest.style.minHeight = "1px";
		_divSource.parentNode.appendChild(divTest);

		// (DIV_BODY_SCROLL) div の margin 調査用
		var divTest2 = document.createElement("DIV");
		divTest2.style.position = "absolute";
		divTest2.style.left = "0px";
		divTest2.style.top = "0px";
		_divSource.parentNode.appendChild(divTest2);

		// (1) div の 位置(style.top, style.left に指定する値)
		// (table 作成時、スクロール時に参照)
		var srcRect = _rect(_divSource); // ***
		var testRect = _rect(divTest);
		_divTop = srcRect.top - testRect.top;
		_divLeft = srcRect.left - testRect.left;

		tblTest.style.width = "50px";
		tblTest.style.height = "50px";
		var tbody = document.createElement("TBODY");
		var tr = document.createElement("TR");
		var td = document.createElement("TD");
		td.appendChild(document.createTextNode("x"));
		tr.appendChild(td);
		tbody.appendChild(tr);
		tblTest.appendChild(tbody);

		// (2) div の幅・高さが指定されているかをチェック
		// div に table を追加し、サイズが変わらなければ固定とみなす。
		var width = divTest.offsetWidth;
		var height = divTest.offsetHeight;
		divTest.appendChild(tblTest);

		_autoWidth = (divTest.offsetWidth != width);
		_autoHeight = (divTest.offsetHeight != height);
		if (_autoWidth) {
			if ((_divMaxWidth > 0) || ((_divAutoSize != "both") && (_divAutoSize != "width"))) {
				_autoWidth = false;
			}
		} else {
			_divMaxWidth = -1;
		}
		if (_autoHeight) {
			if ((_divMaxHeight > 0) || ((_divAutoSize != "both") && (_divAutoSize != "height"))) {
				_autoHeight = false;
			}
		} else {
			_divMaxHeight = -1;
		}

		// (3) div の border + padding + 子(table)の margin の合計サイズを取得
		// (div のりサイズ処理で参照)
		_divPaddingWidth = _offsetWidth(divTest) - _offsetWidth(tblTest);
		_divPaddingHeight = _offsetHeight(divTest) - _offsetHeight(tblTest);

		// DIV_BODY_SCROLL
		var test2Rect = _rect(divTest2);
		_divMarginTop = testRect.top - test2Rect.top;
		_divMarginLeft = testRect.left - test2Rect.left;
		_tblMarginTop = _rect(tblTest).top - testRect.top;
		_tblMarginLeft = _rect(tblTest).left - testRect.left;

		divTest.parentNode.removeChild(divTest);
		divTest2.parentNode.removeChild(divTest2);
	}


	/*
	 * 見出しの全セルの padding と border サイズを取得
	 */
	function _getCellPadding()
	{
		if ((_fixedRows == 0) && (_fixedCols == 0)) return;

		// 調査用テーブルを作成(元テーブルをコピー)
		var tblTest = _tblSource.cloneNode(false);
		// (_isIE) absolute にしないと、先頭列の right が不正になる場合がある。
		tblTest.style.position = "absolute";
		tblTest.style.left = "0px";
		tblTest.style.top = "0px";
		tblTest.style.width = "auto";
		tblTest.style.height = "auto";
		tblTest.width = "";
		tblTest.height = "";

		var tbody;
		var srcTbody = _getElementByTagName(_tblSource, "THEAD");
		if (srcTbody == null) {
			srcTbody = _getElementByTagName(_tblSource, "TBODY");
		}
		if (srcTbody != null) {
			tbody = srcTbody.cloneNode(false);
		} else {
			tbody = document.createElement("TBODY");
		}
		tblTest.appendChild(tbody);
		_tblSource.parentNode.appendChild(tblTest);

		var spanInfo = _getSpanInfo();

		//
		// 固定する行の全セルのコピーを空の状態で作成
		//
		var rows = (_fixedRows > 0 ? _fixedRows : 1);
		var srcTrList = _tblSource.rows;
		var TEST_WIDTH = 0;
		for (var row = 0; row < rows; row++) { // 行のループ
			var srcTr = srcTrList[row];
			var dstTr = srcTr.cloneNode(false);
			dstTr.style.height = _trHeight(srcTr) + PX;
			var srcTdList = _cells(srcTr);

			for (var col = 0; col < srcTdList.length; col++) { // 列のループ
				var srcTd = srcTdList[col];
				var dstTd = srcTd.cloneNode(false);
				dstTd.width = "";
				if (srcTd.colSpan == 1) {
					dstTd.style.width = TEST_WIDTH + PX;
				} else {
					// 固定見出し内の他の行の同じ列に、colSpan=1 のセルが存在するか(ver1.7)
					var existsColumn = true;
					for (var i = 1; i < srcTd.colSpan; i++) {
						if (!spanInfo[srcTd.$FXH_COLINDEX + i]) {
							existsColumn = false;
							break;
						}
					}
					if (existsColumn) {
						dstTd.style.width = "auto";
					} else {
						dstTd.style.width = TEST_WIDTH + PX;
					}
				}
				// (_isIE) 空セルだと clientWidth が取得できない場合があるため(常に0)、
				// テキストを設定。colSpan, rowSpan が関係していそう。
				dstTd.appendChild(document.createTextNode(" "));
				dstTr.appendChild(dstTd);
			}
			tbody.appendChild(dstTr);

			// (IE7(互換表示) && 標準モード)の場合、縦方向のセルパディングを求める
			if (/*(row == 0) &&*/ _isIE && (_IEver <= 7) && !_isBackCompat) {
				_cellPaddingV = _trHeight(dstTr) - _trHeight(srcTr);
				// ↑通常は 0 のはずだが、(IE7(互換表示) && 標準モード)の場合は 0 にならない。
				_cellPaddingV += 2;
			}
		}

		//
		// 全列のボーダーとパディングのサイズを取得。
		// (_isIE8 && 標準モード)だとこの処理が遅い。
		//
		for (var row = 0; row < rows; row++) { // 行のループ
			var srcTdList = _cells(srcTrList[row]);
			var dstTdList = _cells(tbody.rows[row]);

			for (var col = 0; col < srcTdList.length; col++) {
				var srcTd = srcTdList[col];
				var dstTd = dstTdList[col];
				if (dstTd.style.width == "auto") continue;

				var cellId = row + "." + srcTd.cellIndex;
				_cellBorder[cellId] = _offsetWidth(dstTd) - dstTd.clientWidth;
				_cellPadding[cellId] = dstTd.clientWidth - TEST_WIDTH;

				var style = srcTd.currentStyle || document.defaultView.getComputedStyle(srcTd, '');
				var cssPadding = -1;
				if ((style.paddingLeft.match(/px$/) != null) &&
				    (style.paddingRight.match(/px$/) != null)) {
					cssPadding = _pixel(style.paddingLeft) + _pixel(style.paddingRight);
				} else
				if (style.padding.match(/px$/) != null) {
					cssPadding = _pixel(style.padding) * 2;
				}

				if (cssPadding > 0) {
					_cellPadding[cellId] = Math.min(_cellPadding[cellId], cssPadding);
				}

				_maxCellPadding = Math.max(_cellPadding[cellId], _maxCellPadding);
			}
		}

		tblTest.parentNode.removeChild(tblTest);
	}


	/*
	 * rowSpan, colSpan により存在しない列番号を調べる(ver1.7)。
	 */
	function _getSpanInfo()
	{
		var rowCount = (_fixedRows > 0 ? _fixedRows : 1);
		var srcTrList = _tblSource.rows;

		var tmpbuf = new Array(rowCount);
		for (var row = 0; row < rowCount; row++) {
			tmpbuf[row] = new Array(_columnCount);
			for (var col = 0; col < _columnCount; col++) {
				tmpbuf[row][col] = true;
			}
		}

		for (var row = 0; row < rowCount; row++) {
			var srcTr = srcTrList[row];
			var srcTdList = _cells(srcTr);
			var tdIndex = 0;

			for (var col = 0; col < _columnCount; ) {
				if (!tmpbuf[row][col]) {
					col++;
					continue;
				}
				if (tdIndex >= srcTdList.length) {
					break;
				}
				var srcTd = srcTdList[tdIndex];

				// _fixedRows が 0 で、先頭行に rowspan が 2 以上のセルがあるとこけていたのを修正(ver1.9)
				if ((srcTd.rowSpan >= 2) && (srcTd.colSpan >= 2)) {
					for (var i = 0; i < srcTd.rowSpan; i++) {
						if (row + i >= rowCount) break; // (ver1.9)
						for (var j = 0; j < srcTd.colSpan; j++) {
							if ((i == 0) && (j == 0)) continue;
							tmpbuf[row + i][col + j] = false;
						}
					}
				} else {
					if (srcTd.rowSpan >= 2) {
						for (var i = 1; i < srcTd.rowSpan; i++) {
							if (row + i >= rowCount) break; // (ver1.9)
							tmpbuf[row + i][col] = false;
						}
					}
					if (srcTd.colSpan >= 2) {
						for (var i = 1; i < srcTd.colSpan; i++) {
							tmpbuf[row][col + i] = false;
						}
					}
				}

				tdIndex++;
				col += srcTd.colSpan;
			}
		}

		var retbuf = new Array(_columnCount);
		for (var col = 0; col < _columnCount; col++) {
			retbuf[col] = false;
		}
		for (var row = 0; row < rowCount; row++) {
			var srcTr = srcTrList[row];
			var srcTdList = _cells(srcTr);
			var tdIndex = 0;

			for (var col = 0; col < _columnCount; col++) {
				if (!tmpbuf[row][col]) {
					continue;
				}
				if (tdIndex >= srcTdList.length) {
					break;
				}
				var srcTd = srcTdList[tdIndex];
				srcTd.$FXH_COLINDEX = col;
				if (srcTd.colSpan == 1) {
					retbuf[col] = true;
				}
				tdIndex++;
			}
		}

		return retbuf;
	}

////////////////////////////////////////////////////////////////////////////////////////////////////
// テーブルコピー

	/*
	 * 固定テーブルの作成
	 */
	function _createTable(tid, rows, cols)
	{
		var table = _tblSource.cloneNode(false);

		var srcCaption = _getElementByTagName(_tblSource, "CAPTION");
		var srcThead = _getElementByTagName(_tblSource, "THEAD");
		var srcTbody = _getElementByTagName(_tblSource, "TBODY");
		var caption = null;
		var thead = null;
		var tbody = null;

		// CAPTION
		if (srcCaption != null) {
			caption = srcCaption.cloneNode(true);
			caption.style.backgroundColor = _bgcolor;
			caption.style.overflow = "hidden";
			if (tid != TID_HEADER) {
				caption.innerHTML = "&nbsp;";
				caption.style.height = _offsetHeight(srcCaption) + PX;
				caption.style.backgroundColor = "transparent";
				// (_isMobile && (_isChrome || _isOpera))
				// TODO スマホの Chrome, Opera は transparent だとゴミが表示される
			}
			table.appendChild(caption);
		}

		// THEAD
		var theadRows = 0;
		if (srcThead != null) {
			thead = srcThead.cloneNode(false);
			table.appendChild(thead);
			theadRows = srcThead.rows.length;
		}

		// TBODY
		if ((srcTbody != null) && (theadRows < rows)) {
			tbody = srcTbody.cloneNode(false);
			table.appendChild(tbody);
		}

		// 固定するセルをコピー
		if (_copyCells(table, tbody, thead, tid, rows, cols, theadRows) == false) {
			return null;
		}

		_linkElement(table, _tblSource, tid, _copyId, false);
		if (caption != null) _linkElement(caption, srcCaption, tid, _copyId, true);
		if (thead != null) _linkElement(thead, srcThead, tid, _copyId, false);
		if (tbody != null) _linkElement(tbody, srcTbody, tid, _copyId, false);

		// right, bottom の margin, border, padding を 0 にする
		if (tid != TID_HEADER) { // NUMBER or CORNER
			table.style.marginRight      = "0px";
			table.style.borderRightWidth = "0px";
			table.style.paddingRight     = "0px";
		}
		if (tid != TID_NUMBER) { // HEADER or CORNER
			table.style.marginBottom      = "0px";
			table.style.borderBottomWidth = "0px";
			table.style.paddingBottom     = "0px";
		}
		table.style.minWidth  = "1px";
		table.style.minHeight = "1px";

		var boxShadow = null;
		if (_boxShadow != null) {
			// 1px 1px 1px 0 rgba(0,0,0,0.4);
			var buf = _boxShadow.split(" ");
			boxShadow = "";
			for (var i = 0; i < buf.length; i++) {
				if ((i == 0) && (tid == TID_HEADER)) {
					boxShadow += "0 "; // 横の影をなくす
				} else
				if ((i == 1) && (tid == TID_NUMBER)) {
					boxShadow += "0 "; // 縦の影をなくす
				} else {
					boxShadow += buf[i] + " ";
				}
			}
		}

		if (_isDivMode) {
			var div = _divSource.cloneNode(false);
			div.$FXH_PADDING_WIDTH = undefined;
			div.$FXH_PADDING_HEIGHT = undefined;
			_linkElement(div, _divSource, tid, _copyId, false);
			div.className += " fixed_header_display_none_at_print";
			div.style.overflowX = "hidden";
			div.style.overflowY = "hidden";
			_removeEventListener(div, "scroll", _onDivScroll);
			if (tid != TID_CORNER) {
				_addEventListener(div, "scroll", function() { _onFixedDivScroll(div, tid); } );
			}

			// right, bottom の margin, border, padding を 0 にする
			if (tid == TID_HEADER) {
				div.style.borderRightWidth = "0px";
			} else { // NUMBER or CORNER
				div.style.marginRight      = "0px";
				div.style.borderRightWidth = "0px";
				div.style.paddingRight     = "0px";
			}
			if (tid == TID_NUMBER) {
				div.style.borderBottomWidth = "0px";
			} else { // HEADER or CORNER
				div.style.marginBottom      = "0px";
				div.style.borderBottomWidth = "0px";
				div.style.paddingBottom     = "0px";
			}
			div.style.width  = "30px";
			div.style.height = "30px";
			div.style.minWidth  = "1px";
			div.style.minHeight = "1px";

			if (boxShadow != null) {
				div.style.boxShadow = boxShadow;
			}

			// 元のテーブルに重ねる
			div.style.position = "absolute";
			div.style.top = _divTop + PX;
			div.style.left = _divLeft + PX;

			div.style.backgroundColor = _bgcolor;
			div.appendChild(table);
			_divSource.parentNode.appendChild(div);
		} else {
			// body mode
			table.className += " fixed_header_display_none_at_print";
			if (boxShadow != null) {
				table.style.boxShadow = boxShadow;
			}

			switch (_positionMode) {
			case POS_FIXED:
				table.style.position = "fixed";
				break;
			case POS_ABSOLUTE:
				table.style.position = "absolute";
				break;
			default: // POS_MIX
				table.style.position = "fixed";
				break;
			}
			table.style.marginTop  = "0px";
			table.style.marginLeft = "0px";
			table.style.top  = _bodyTop + PX;
			table.style.left = _bodyLeft + PX;

			table.style.backgroundColor = _bgcolor;
			_tblSource.parentNode.appendChild(table);
		}

		if (_isOpera) { // 処理中のイメージがはっきり見えてしまうので隠しておく
			_showTable(_isDivMode ? div : table, false);
		}

		return table;
	}


	function _copyCells(table, tbody, thead, tid, rows, cols, theadRows)
	{
		var nextRows = new Array(cols); // rowSpan 制御用
		for (var i = 0; i < nextRows.length; i++) {
			nextRows[i] = 0;
		}

		var srcTrList = _tblSource.rows;
		for (var row = 0; row < rows; row++) { // 行のループ
			var srcTr = srcTrList[row];
			var dstTr = srcTr.cloneNode(false);
			_linkElement(dstTr, srcTr, tid, _copyId, false);

			// 行の高さ設定
			dstTr.style.height = (_trHeight(srcTr) - _cellPaddingV) + PX;
			// _cellPaddingV は (IE7(互換表示) && 標準モード) の場合のみ設定される。
			// この場合、パディングやボーダーを引いた値を指定しないと行が高くなってしまう。

			if ((row == _fixedRows - 1) && (tid != TID_NUMBER)) { // HEADER or CORNER
				if (_borderColor != null) dstTr.style.borderBottomColor = _borderColor;
				if (_borderStyle != null) dstTr.style.borderBottomStyle = _borderStyle;
				if (_borderWidth != null) dstTr.style.borderBottomWidth = _borderWidth;
			}
			if (tid != TID_HEADER) { // NUMBER or CORNER
				dstTr.style.borderRightWidth = "0px";
			}
			if (row < theadRows) {
				thead.appendChild(dstTr);
			} else
			if (tbody != null) {
				tbody.appendChild(dstTr);
			} else {
				// tbody がない場合は table に追加(ver1.9)
				table.appendChild(dstTr);
			}

			var tdList = _cells(srcTr);
			var tdIndex = 0;
			for (var col = 0; col < cols; ) { // 列のループ
				if (row < nextRows[col]) {
					col++;
					continue; // rowSpan により存在しない列
				}
				if (tdIndex >= tdList.length) break;
				var srcTd = tdList[tdIndex];
				tdIndex++;

				// 指定された固定行数が rowSpan で結合された行の途中の場合
				if (row + srcTd.rowSpan > rows) {
					return false;
				}

				nextRows[col] = row + srcTd.rowSpan; // この列の次の行
				if (srcTd.colSpan >= 2) {
					for (var i = 1; i < srcTd.colSpan; i++) {
						nextRows[col + i] = nextRows[col];
					}
					// 指定された固定列数が colSpan で結合された列の途中の場合
					if (col + srcTd.colSpan > cols) {
						return false;
					}
				}
				// (注意) (_isIE) dstTd の rowSpan は常に 1。
				// document に追加後は正常に取得できるようだ。

				_radioCtl(srcTd, "backup");
				var dstTd = srcTd.cloneNode(true); // 子もコピー
				_radioCtl(srcTd, "restore");
				_linkElement(dstTd, srcTd, tid, _copyId, true);
				dstTr.appendChild(dstTd);

				// セル幅設定(cellPadding 未調査のセルは行わない)
				try {
					var cellId = row + "." + srcTd.cellIndex;
					if (_cellPadding[cellId] != undefined) {
						var pad = _cellPadding[cellId] + _cellBorder[cellId];
						// fjn 20170613
						dstTd.style.width = (_offsetWidth(srcTd) - 0) + PX;
					} else
					if (_isIE && (_IEver <= 8) && (srcTd.colSpan >= 2)) {
						//
						// IE8 では、colSpan 指定のセルも、幅指定しないと微妙にずれる。
						// 指定する値は大きくなければ何でも良いようだ。
						// dstTd.style.width = "1px";
						// ↑でもよいが、↓のようにしておく。
						//
						dstTd.style.width = (srcTd.clientWidth - _maxCellPadding) + PX;
					}
				} catch (e) {
					// (_isIE) 標準モードだとcellIndexにアクセスできない場合がある
					// tr 配下に form がある場合？
				}

				var s = dstTd.style;
				if ((row + srcTd.rowSpan == _fixedRows) && (tid != TID_NUMBER)) {
					if (_borderColor != null) s.borderBottomColor = _borderColor;
					if (_borderStyle != null) s.borderBottomStyle = _borderStyle;
					if (_borderWidth != null) s.borderBottomWidth = _borderWidth;
				}
				if ((col + srcTd.colSpan == _fixedCols) && (tid != TID_HEADER)) {
					if (_borderColor != null) s.borderRightColor = _borderColor;
					if (_borderStyle != null) s.borderRightStyle = _borderStyle;
					if (_borderWidth != null) s.borderRightWidth = _borderWidth;
				}

				col += srcTd.colSpan;
			}

			// rowSpan により、カラムがない行の場合。
			// (_isIE) 「rowSpan が指定された行」に高さを足してやらないとずれる場合がある。
			// ver1.5で削除した処理を復活(ver1.6)
			if (_isIE && (_IEver <= 9) && (tdIndex == 0)) {
				dstTr.style.height = "0px";
				var bottom = _rect(srcTr).bottom;
				var parent = dstTr.parentNode;
				if (parent.tagName != "TABLE") parent = parent.parentNode;
				var dstTrList = parent.rows;
				var r = row - 1;
				while (dstTrList[r].style.height == "0px") r--;
				var height = bottom - _rect(srcTrList[r]).top;
				dstTrList[r].style.height = (height - _cellPaddingV) + PX; // rowSpan 指定行
			}
		}

		_radioCtl(table, "sync");

		return true;
	}


	/*
	 * 位置がずれていたら直す
	 */
	function _checkDivPosition(div)
	{
		if (div.style.position == "fixed") return; // DIV_BODY_SCROLL

		var srcRect = _rect(_divSource);
		var dstRect = _rect(div);
		var topDiff = dstRect.top - srcRect.top;
		var leftDiff = dstRect.left - srcRect.left;

		//
		// (_isIE) 拡大表示の場合に、本関数が呼ばれるたびに
		// topDiff が +1, -1, +1, -1, ... となる場合がある。
		// 本関数はスクロールのたびに呼ばれるため、この場合、
		// スクロール中に見出しが上下に揺れる。
		// 以下は、この状況に陥らないための処置。
		// left については遭遇したことがないが、念のため同じ処置をしておく。
		//
		if (_isIE) {
			if ((topDiff == -1) && (div.$TOP_DIFF == 1)) {
				topDiff = 0;
			} else {
				div.$TOP_DIFF = topDiff;
			}
			if ((leftDiff == -1) && (div.$LEFT_DIFF == 1)) {
				leftDiff = 0;
			} else {
				div.$LEFT_DIFF = leftDiff;
			}
		}

		if (Math.abs(topDiff) >= 1) {
			div.style.top = (_pixel(div.style.top) - topDiff) + PX;
		}
		if (Math.abs(leftDiff) >= 1) {
			div.style.left = (_pixel(div.style.left) - leftDiff) + PX;
		}
	}


	/*
	 * 作成したテーブルの削除
	 */
	function _removeTable(table)
	{
		if (_isDivMode) {
			table = table.parentNode; // div
		}
		_unlinkElement(table);
		if (table.parentNode) { // 元テーブルが破棄されている場合は null がありえる
			table.parentNode.removeChild(table);
		}
	}

////////////////////////////////////////////////////////////////////////////////////////////////////
// スクロール

	/*
	 * body の onScroll のイベントハンドラ。
	 * スクロールされた時に、見出し行と行番号列もスクロールさせる。
	 */
	this.onBodyScroll = function()
	{
		if (_bodyHeader) _bodyTop = Math.max(_rect(_bodyHeader).bottom, 0);
		if (_leftHeader) _bodyLeft = Math.max(_rect(_leftHeader).right, 0);

	//	if (_isDivMode) return; // div モードでも呼ばれるのでチェック
		if (_isDivMode) {
			_scrollDivByFixed(); // DIV_BODY_SCROLL
			return;
		}


		if (_positionMode == POS_ABSOLUTE) {
			// ちらつくのでスクロール中は非表示にしておく
			if (!_isMobile) {
				if (_getBodyScrollTop() != _bodyScrollTop) {
					_showTable(_tblHeader, false);
					_showTable(_tblCorner, false);
				}
				if (_getBodyScrollLeft() != _bodyScrollLeft) {
					_showTable(_tblNumber, false);
					_showTable(_tblCorner, false);
				}
			}
			if (_timerId_bodyScroll != null) clearTimeout(_timerId_bodyScroll);
			_timerId_bodyScroll = setTimeout(_onBodyScroll2, 200);
		} else {
			_onBodyScroll2();
		}
	};


	function _onBodyScroll2(isFirst)
	{
		_timerId_bodyScroll = null;

		var topChanged = (_getBodyScrollTop() != _bodyScrollTop);
		var leftChanged = (_getBodyScrollLeft() != _bodyScrollLeft);
		_bodyScrollTop = _getBodyScrollTop();
		_bodyScrollLeft = _getBodyScrollLeft();
		if (topChanged && leftChanged) isFirst = true;

		var srcRect = _rect(_tblSource);


		var showHeader = ((_tblHeader != null) &&
			(srcRect.top < _bodyTop) &&
			(srcRect.bottom >= _tblHeader.offsetHeight + _bodyTop));
		var showNumber = ((_tblNumber != null) &&
			(srcRect.left < _bodyLeft) &&
			(srcRect.right >= _tblNumber.offsetWidth + _bodyLeft));

		// 表示・非表示
		if (_tblHeader != null) _showTable(_tblHeader, showHeader);
		if (_tblNumber != null) _showTable(_tblNumber, showNumber);
		if (_tblCorner != null) _showTable(_tblCorner, (showHeader && showNumber));

		//
		// POS_MIX の場合
		// absolute と fixed の切り替え。
		//                     見出しTBL 行番号TBL
		// 上下にスクロール中  fixed     absolute
		// 左右にスクロール中  absolute  fixed
		//
		if ((_positionMode == POS_MIX) && (_tblHeader != null)) {
			var style = _tblHeader.style;
			if (isFirst || (topChanged && (style.position == "absolute"))) {
				// absolute --> fixed
				style.position = "fixed";
				style.left = srcRect.left + PX;
				style.top = _bodyTop + PX;
				if (_tblCorner != null) _tblCorner.style.top = style.top;
				_adjustHeaderLeft(false);
			} else
			if (!topChanged && leftChanged && (style.position == "fixed")) {
				// fixed --> absolute
				style.position = "absolute";
				style.left = (_bodyScrollLeft + srcRect.left) + PX;
				style.top = (_bodyScrollTop + _bodyTop) + PX;
				_adjustHeaderTop();
			}
		}
		if ((_positionMode == POS_MIX) && (_tblNumber != null)) {
			var style = _tblNumber.style;
			if (isFirst || (leftChanged && (style.position == "absolute"))) {
				// absolute --> fixed
				style.position = "fixed";
				style.left = _bodyLeft + PX;
				style.top = srcRect.top + PX;
				_adjustNumberTop(false);
			} else
			if (!leftChanged && topChanged && (style.position == "fixed")) {
				// fixed --> absolute
				style.position = "absolute";
				style.left = (_bodyScrollLeft + _bodyLeft) + PX;
				style.top = (_bodyScrollTop + srcRect.top) + PX;
				_adjustNumberLeft();
			}
		}

		//
		// POS_FIXED の場合
		//
		if (_positionMode == POS_FIXED) {
			if ((_tblHeader != null) && leftChanged) {
				_tblHeader.style.left = srcRect.left + PX;
			}
			if ((_tblNumber != null) && topChanged) {
				_tblNumber.style.top = srcRect.top + PX;
			}
			if ((_tblHeader != null) && (_tblHeader.style.top != (_bodyTop + PX))) {
				_tblHeader.style.top = _bodyTop + PX;
				if (_tblCorner != null) _tblCorner.style.top = _tblHeader.style.top;
			}
		}

		//
		// POS_ABSOLUTE の場合
		//
		if ((_positionMode == POS_ABSOLUTE) && isFirst) {
			if (_tblHeader != null) {
				_tblHeader.style.left = (_bodyScrollLeft + srcRect.left) + PX;
			}
			if (_tblNumber != null) {
				_tblNumber.style.top = (_bodyScrollTop + srcRect.top) + PX;
			}
		}
		if ((_positionMode == POS_ABSOLUTE) && (topChanged || leftChanged)) {
			// 少し前に位置づけて、徐々に本来の位置に移動させる
			var topMoveSize, leftMoveSize;
			var count = (_isMobile ? 1 : 4);
			if ((_tblHeader != null) && topChanged) {
				_tblHeader.style.top = (_bodyTop + _bodyScrollTop - _offsetHeight(_tblHeader)) + PX;
				if (_tblCorner != null) _tblCorner.style.top = _tblHeader.style.top;
				topMoveSize = _offsetHeight(_tblHeader) / count;
			}
			if ((_tblNumber != null) && leftChanged) {
				_tblNumber.style.left = (_bodyLeft + _bodyScrollLeft - _offsetWidth(_tblNumber)) + PX;
				if (_tblCorner != null) _tblCorner.style.left = _tblNumber.style.left;
				leftMoveSize = _offsetWidth(_tblNumber) / count;
			}
			if (_timerId_bodyScroll_abs != null) clearTimeout(_timerId_bodyScroll_abs);
			_onBodyScroll_absolute(topChanged, leftChanged, topMoveSize, leftMoveSize);
		}

		if ((_tblHeader != null) && leftChanged) _adjustHeaderLeft(false);
		if ((_tblNumber != null) && topChanged) _adjustNumberTop(false);
	}


	/*
	 * (_isIE) 少しずつ本来の位置に移動させる
	 */
	function _onBodyScroll_absolute(topChanged, leftChanged, topMoveSize, leftMoveSize)
	{
		_timerId_bodyScroll_abs = null;

		var bodyTop = _getBodyScrollTop() + _bodyTop;
		var bodyLeft = _getBodyScrollLeft() + _bodyLeft;
		var tableTop = bodyTop;
		var tableLeft = bodyLeft;

		if ((_tblHeader != null) && topChanged) {
			tableTop = _pixel(_tblHeader.style.top) + topMoveSize;
			if (topMoveSize > 0) {
				tableTop = Math.min(tableTop, bodyTop);
			} else {
				tableTop = Math.max(tableTop, bodyTop);
			}
			_tblHeader.style.top = tableTop + PX;
			if (_tblCorner != null) _tblCorner.style.top = _tblHeader.style.top;
		}

		if ((_tblNumber != null) && leftChanged) {
			tableLeft = _pixel(_tblNumber.style.left) + leftMoveSize;
			if (leftMoveSize > 0) {
				tableLeft = Math.min(tableLeft, bodyLeft);
			} else {
				tableLeft = Math.max(tableLeft, bodyLeft);
			}
			_tblNumber.style.left = tableLeft + PX;
			if (_tblCorner != null) _tblCorner.style.left = _tblNumber.style.left;
		}

		if ((bodyTop == tableTop) && (bodyLeft == tableLeft)) {
			if ((_tblHeader != null) && topChanged) {
				_adjustHeaderTop();
				if (_tblCorner != null) _tblCorner.style.top = _tblHeader.style.top;
			}
			if ((_tblNumber != null) && leftChanged) {
				_adjustNumberLeft();
				if (_tblCorner != null) _tblCorner.style.left = _tblNumber.style.left;
			}
			return; // スクロール処理終了
		}

		var func = function() { _onBodyScroll_absolute(topChanged, leftChanged, topMoveSize, leftMoveSize); };
		_timerId_bodyScroll_abs = setTimeout(func, 20);
	}


	/*
	 * div モードでも、body のスクロールに反応させる(DIV_BODY_SCROLL)。
	 *
	 * (_isSafari) fixed にしたり absolute にしたりすると、android の Safari が正常に表示できない。
	 */
	function _scrollDivByFixed()
	{
		if (DIV_BODY_SCROLL == 0) return;
		if (_positionMode == POS_ABSOLUTE) return; // "position: fixed" 未サポートの場合

		var divHeader = (_tblHeader == null) ? null : _tblHeader.parentNode;
		var divNumber = (_tblNumber == null) ? null : _tblNumber.parentNode;
		var divCorner = (_tblCorner == null) ? null : _tblCorner.parentNode;
		var divRect = _rect(_divSource);
		var tblRect = _rect(_tblSource);

		var n = 0;
		if (_isIE11) n = 1; // div.scrollHeight が 1 多い場合がある
		var headerFixed = false;
		var numberFixed = false;
		if (divHeader && ((DIV_BODY_SCROLL == 2) || (_divSource.scrollHeight - n <= _divSource.clientHeight))) {
			var bottom = Math.min(divRect.bottom, tblRect.bottom);
			if ((divRect.top < _bodyTop - _tblMarginTop) && (bottom >= _offsetHeight(_tblHeader))) {
				if (divHeader.style.position != "fixed") {
					divHeader.style.position = "fixed";
					divHeader.style.top = (_bodyTop - _divMarginTop - _tblMarginTop) + PX;
				}
				divHeader.style.left = (divRect.left - _divMarginLeft) + PX;
				_showTable(divHeader, true);
				headerFixed = true;
			} else {
				if (divHeader.style.position != "absolute") {
					divHeader.style.position = "absolute";
					divHeader.style.top = _divTop + PX;
					divHeader.style.left = _divLeft + PX;
				}
				_showTable(divHeader, (_divSource.scrollTop > 0));
				// (DIV_BODY_SCROLL == 1) の場合、scrollTop は常に0
				_checkDivPosition(divHeader);
			}

			if (_isVisible(divHeader)) _adjustHeaderLeft(false);
		}

		if (divNumber && ((DIV_BODY_SCROLL == 2) || (_divSource.scrollWidth <= _divSource.clientWidth))) {
			var right = Math.min(divRect.right, tblRect.right);
			if ((divRect.left < _bodyLeft - _tblMarginLeft) && (right >= _offsetWidth(_tblNumber))) {
				if (divNumber.style.position != "fixed") {
					divNumber.style.position = "fixed";
					divNumber.style.left = (_bodyLeft - _divMarginLeft - _tblMarginLeft) + PX;
				}
				divNumber.style.top = (divRect.top - _divMarginTop) + PX;
				_showTable(divNumber, true);
				numberFixed = true;
			} else {
				if (divNumber.style.position != "absolute") {
					divNumber.style.position = "absolute";
					divNumber.style.top = _divTop + PX;
					divNumber.style.left = _divLeft + PX;
				}
				_showTable(divNumber, (_divSource.scrollLeft > 0));
				// (DIV_BODY_SCROLL == 1) の場合、scrollLeft は常に0
				_checkDivPosition(divNumber);
			}

			if (_isVisible(divNumber)) _adjustNumberTop(false);
		}

		if (divCorner) {
			if (headerFixed || numberFixed) {
				if (divCorner.style.position != "fixed") {
					divCorner.style.position = "fixed";
				}
				if (headerFixed) {
					divCorner.style.top = divHeader.style.top;
				} else {
					divCorner.style.top = divNumber.style.top;
				}
				if (numberFixed) {
					divCorner.style.left = divNumber.style.left;
				} else {
					divCorner.style.left = divHeader.style.left;
				}
			} else {
				if (divCorner.style.position != "absolute") {
					divCorner.style.position = "absolute";
					divCorner.style.top = _divTop + PX;
					divCorner.style.left = _divLeft + PX;
					_checkDivPosition(divCorner);
				}
			}
			_showTable(divCorner, _isVisible(divHeader) && _isVisible(divNumber));
		}
	}


	/*
	 * divSource の onScroll のイベントハンドラ。
	 * スクロールされた時に、見出し行と行番号列もスクロールさせる。
	 * スマホの Safari, Firefox は、div のスクロール処理が重い。Chrome, Opera は速かった。(2014/6)
	 */
	function _onDivScroll()
	{

		var divHeader = (_tblHeader == null) ? null : _tblHeader.parentNode;
		var divNumber = (_tblNumber == null) ? null : _tblNumber.parentNode;
		var divCorner = (_tblCorner == null) ? null : _tblCorner.parentNode;
		var headerFixed = (divHeader && (divHeader.style.position == "fixed")); // DIV_BODY_SCROLL
		var numberFixed = (divNumber && (divNumber.style.position == "fixed")); // DIV_BODY_SCROLL

		// 表示・非表示制御(スクロールしていない場合は非表示)
		_showTable(divHeader, (_divSource.scrollTop > 0) || headerFixed);
		_showTable(divNumber, (_divSource.scrollLeft > 0) || numberFixed);
		_showTable(divCorner, _isVisible(divHeader) && _isVisible(divNumber));

		if (divHeader != null) _checkDivPosition(divHeader);
		if (divNumber != null) _checkDivPosition(divNumber);
		if (divCorner != null) _checkDivPosition(divCorner);

		//
		// (_isIE) 拡大表示すると、
		// 固定テーブルの列の幅や行の高さがずれることがあるので、
		// 少しでも見た目が揃う位置にスクロールする。
		// IE 以外でもやっておく。
		// ずれが発生しなければ、本処理は次の２行でよい。
		// if (divHeader != null) divHeader.scrollLeft = _divSource.scrollLeft;
		// if (divNumber != null) divNumber.scrollTop = _divSource.scrollTop;
		//

		var leftChanged = (_divSource.scrollLeft != _divScrollLeft);
		var topChanged = (_divSource.scrollTop != _divScrollTop);
		_divScrollLeft = _divSource.scrollLeft;
		_divScrollTop = _divSource.scrollTop;

		if ((divHeader != null) && leftChanged) {
			divHeader.$FXH_SCROLL_LEFT = _divSource.scrollLeft;
			divHeader.scrollLeft = _divSource.scrollLeft;
			if (divHeader.scrollLeft > 0) {
				_adjustHeaderLeft(false);
			}
		}

		if ((divNumber != null) && topChanged) {
			divNumber.$FXH_SCROLL_TOP = _divSource.scrollTop;
			divNumber.scrollTop = _divSource.scrollTop;
			if (divNumber.scrollTop > 0) {
				_adjustNumberTop(false);
			}
		}
	}


	/*
	 * 固定テーブルの親divは、スクロールバーを表示しないため、通常はスクロールイベントが
	 * 発生することはないが、<input> や <a> 等のフォーカスを受け取る要素がある場合は、
	 * TABキー操作によりスクロールが発生する。
	 */
	function _onFixedDivScroll(div, tid, atOnce)
	{
		//
		// JavaScript( で scrollLeft, scrollTop を操作すると onScroll イベントが発生するのでそれは無視。
		// 値が「同じ」だったら _onDivScroll() によるスクロール、とみなしたいところだが、
		// chrome で拡大表示すると、最大で 2px 程度の差がでるため、
		// 少し余裕を見て、値の差が「5px未満」なら _onDivScroll() によるスクロールとみなす。
		//
		if (tid == TID_HEADER) {
			if (Math.abs(div.scrollLeft - div.$FXH_SCROLL_LEFT) < 5) return;
		} else {
			if (Math.abs(div.scrollTop - div.$FXH_SCROLL_TOP) < 5) return;
		}

		if (tid == TID_HEADER) {
			_divSource.scrollLeft = div.scrollLeft;
		} else {
			_divSource.scrollTop = div.scrollTop;
		}
	}


	/*
	 * body モード、かつ position が absolute、かつ拡大表示の場合、
	 * わずかに隙間ができる場合があるので調整。
	 */
	function _adjustHeaderTop()
	{
		if (_tblHeader == null) return;
		var rect = _rect(_tblHeader);
		if (rect.top <= _bodyTop) return;

		_tblHeader.style.top = (_bodyTop + _pixel(_tblHeader.style.top) - rect.top) + PX;
	}


	function _adjustNumberLeft()
	{
		if (_tblNumber == null) return;
		var rect = _rect(_tblNumber);
		if (rect.left <= _bodyLeft) return;

		_tblNumber.style.left = (_bodyLeft + _pixel(_tblNumber.style.left) - rect.left) + PX;
	}


	/*
	 * (_isIE) 拡大表示すると、
	 * 固定テーブルの列の幅や行の高さがずれることがあるので、
	 * 少しでも見た目が揃う位置にスクロールする。
	 * IE 以外でもやっておく。
	 */
	function _adjustHeaderLeft(atOnce)
	{
		if (_tblHeader == null) return;
		if (!atOnce) {
			if (_timerId_adjustLeft != null) clearTimeout(_timerId_adjustLeft);
			_timerId_adjustLeft = setTimeout(function() { _adjustHeaderLeft(true); }, 200);
			return;
		}
		_timerId_adjustLeft = null;

		// 現在見えている列のずれを計算
		var minLeft;
		var maxLeft;
		if (_isDivMode) {
			minLeft = _rect(_divSource).left;
			maxLeft = minLeft + _divSource.clientWidth;
		} else {
			minLeft = 0;
			maxLeft = _body.clientWidth;
		}

		if (_isVisible(_tblNumber)) {
			minLeft += _offsetWidth(_tblNumber);
		}

		var srcTdList = _cells(_tblSource.rows[0]);
		var dstTdList = _cells(_tblHeader.rows[0]);
		var diff = 0;
		var cols = 0;

		for (var col = 0; col < srcTdList.length; col++) {
			var srcLeft = _rect(srcTdList[col]).left;
			if (srcLeft < minLeft) continue;
			if (srcLeft > maxLeft) break;
			var dstLeft = _rect(dstTdList[col]).left;
			diff += dstLeft - srcLeft;
			cols++;
		}
		if (cols == 0) return;
		if (diff == 0) return; // 差分なし

		diff = diff / cols;
		diff = Math.round(diff); // 四捨五入(しなくてもよい)
		if (diff == 0) return;

		if (_isDivMode) {
			var divHeader = _tblHeader.parentNode;
			if (divHeader.style.position == "fixed") { // DIV_BODY_SCROLL
				var left = _pixel(divHeader.style.left) - diff;
				divHeader.style.left = left + PX;
			} else {
				divHeader.$FXH_SCROLL_LEFT = divHeader.scrollLeft + diff;
				divHeader.scrollLeft += diff;
			}
		} else {
			var left = _pixel(_tblHeader.style.left) - diff;
			_tblHeader.style.left = left + PX;
		}
	}


	function _adjustNumberTop(atOnce)
	{
		if (_tblNumber == null) return;
		if (_isOpera) return; // ずれてしまう場合があるのでやらない
		if (!atOnce) {
			if (_timerId_adjustTop != null) clearTimeout(_timerId_adjustTop);
			_timerId_adjustTop = setTimeout(function() { _adjustNumberTop(true); }, 200);
			return;
		}
		_timerId_adjustTop = null;

		// 現在見えている行のずれを計算
		var minTop;
		var maxTop;
		if (_isDivMode) {
			minTop = _rect(_divSource).top;
			maxTop = minTop + _divSource.clientHeight;
		} else {
			minTop = 0;
			maxTop = _body.clientHeight;
		}

		if (_isVisible(_tblHeader)) {
			minTop += _offsetHeight(_tblHeader);
		}

		var srcTrList = _tblSource.rows;
		var dstTrList = _tblNumber.rows;
		var diff = 0;
		var rows = 0;

		for (var row = 0; row < srcTrList.length; row++) {
			var srcTop = _rect(srcTrList[row]).top;
			if (srcTop < minTop) continue;
			if (srcTop > maxTop) break;
			if (_cells(dstTrList[row]).length == 0) continue;
			var dstTop = _rect(dstTrList[row]).top;
			diff += (dstTop - srcTop);
			rows++;
		}
		if (rows == 0) return;
		if (diff == 0) return; // 差分なし

		diff = diff / rows;
		diff = Math.round(diff); // 四捨五入(しなくてもよい)
		if (diff == 0) return;

		if (_isDivMode) {
			var divNumber = _tblNumber.parentNode;
			if (divNumber.style.position == "fixed") { // DIV_BODY_SCROLL
				var top = _pixel(divNumber.style.top) - diff;
				divNumber.style.top = top + PX;
			} else {
				divNumber.$FXH_SCROLL_TOP = divNumber.scrollTop + diff;
				divNumber.scrollTop += diff;
			}
		} else {
			var top = _pixel(_tblNumber.style.top) - diff;
			_tblNumber.style.top = top + PX;
		}
	}


	function _isVisible(table)
	{
		return (table != null) && (table.style.visibility == "visible");
	}

	/*
	 * スクロールしていない場合は、固定テーブルは表示しない。
	 * パラメータ名は table だが、div モードの場合は div。
	 */
	function _showTable(table, show)
	{
		if (table == null) return;

		var value = (show ? "visible" : "hidden");
		if (table.style.visibility == value) return;
		table.style.visibility = value;

		if (_isIE) {
			//
			// 固定テーブルが非表示でも、下にある元テーブルをいじれない
			// (文字選択やリンクのクリックができない)。
			// スクロールして一旦固定テーブルを表示して非表示にするといじれるようになる。
			// table に caption があるとこの現象になる？
			// 原因が全くわからないので、非表示の場合は重なり順を下にする。(ver1.5)
			//
			table.style.zIndex = (show ? _zIndex : (_zIndex - 1));
		}

		if (_isOpera) {
			// Opera で visibility を操作するとボーダーが消えてしまう。
			// 次のおまじないで消えなくなる。
			table.style.opacity = (show ? _opacity : 0);
		}
		if (show && (_opacity < 1.0)) {
			// (_isIE) createTable()でfilterを設定すると、表示されなくなるのでここで行う
			table.style.opacity = _opacity;
			table.style.filter = "alpha(opacity=" + (_opacity * 100) + ")"; // _isIE
		}
	}

////////////////////////////////////////////////////////////////////////////////////////////////////

	/*
	 * (_isIE || _isFirefox || _isOpera)
	 * テーブルのデータ量が多いと div のリサイズ処理が非常に遅くなるので、
	 * div のリサイズ処理中は元 table と同じサイズのダミー table を作成し
	 * 元 table は非表示にしておく。
	 */
	function _hideSourceTable(hideFlag)
	{
		if (!_isIE && !_isIE11 && !_isFirefox && !_isOpera) return;

		if (hideFlag) {
			_tblDummy = _createDummyTable();
			_divSource.appendChild(_tblDummy);
			_setOffsetWidth(_tblDummy, _offsetWidth(_tblSource));
			_setOffsetHeight(_tblDummy, _offsetHeight(_tblSource));
			_tblSource.style.display = "none"; // ***
		} else {
			_tblSource.style.display = _tableDisplay; // ***
			_divSource.removeChild(_tblDummy);
			_tblDummy = null;
		}
	}


	function _createDummyTable()
	{
		var table = _tblSource.cloneNode(false);
		var tbody = document.createElement("TBODY");
		var tr = document.createElement("TR");
		var td = document.createElement("TD");
		td.appendChild(document.createTextNode("dummy"));
		tr.appendChild(td);
		tbody.appendChild(tr);
		table.appendChild(tbody);
		return table;
	}

////////////////////////////////////////////////////////////////////////////////////////////////////
// サイズ設定
// clientWidth, clientHeight でサイズ設定するのは、コピーdiv。
// offsetWidth, offsetHeight でサイズ設定するのは、コピーtableと元div。

	/* clientWidth を指定された長さ(px)にする */
	function _setClientWidth(element, width)
	{
		_setWidth(element, "CLIENT", width);
	}


	/* offsetWidth を指定された長さ(px)にする */
	function _setOffsetWidth(element, width)
	{
		_setWidth(element, "OFFSET", width);
	}


	/* clientHeight を指定された長さ(px)にする */
	function _setClientHeight(element, height)
	{
		_setHeight(element, "CLIENT", height);
	}


	/* offsetHeight を指定された長さ(px)にする */
	function _setOffsetHeight(element, height)
	{
		_setHeight(element, "OFFSET", height);
	}


	function _setWidth(element, sizeType, width)
	{
		var w = width;
		if (element.$FXH_PADDING_WIDTH != undefined) {
			w -= element.$FXH_PADDING_WIDTH;
		}
		var currentWidth;
		var diff;
		for (var retry = 0; retry < 2; retry++) {
			if (w < MIN_SIZE) w = MIN_SIZE;
			element.style.width = w + PX;
			currentWidth = ((sizeType == "CLIENT") ? element.clientWidth : _offsetWidth(element));
			// 差分(padding や border 分)を調整
			diff = currentWidth - width;
			if (element.$FXH_PADDING_WIDTH == undefined) {
				element.$FXH_PADDING_WIDTH = diff;
			}
			if ((diff == 0) || (w == MIN_SIZE)) break;
			w -= diff;
		}
	}


	function _setHeight(element, sizeType, height)
	{
		var h = height;
		if (element.$FXH_PADDING_HEIGHT != undefined) {
			h -= element.$FXH_PADDING_HEIGHT;
		}
		var currentHeight;
		var diff;
		for (var retry = 0; retry < 2; retry++) {
			if (h < MIN_SIZE) h = MIN_SIZE;
			element.style.height = h + PX;
			currentHeight = ((sizeType == "CLIENT") ? element.clientHeight : _offsetHeight(element));
			// 差分(padding や border 分)を調整
			diff = currentHeight - height;
			if (element.$FXH_PADDING_HEIGHT == undefined) {
				element.$FXH_PADDING_HEIGHT = diff;
			}
			if ((diff == 0) || (h == MIN_SIZE)) break;
			h -= diff;
		}
	}

////////////////////////////////////////////////////////////////////////////////////////////////////
// サイズ調整

	/*
	 * (_isIE) セル幅が元のセルより1ピクセルでも狭くなると、
	 * 文字列が折り返されて行が高くなってしまうことがある。
	 * この場合は、セルのパディングを減らすことにする。
	 * 拡大表示すると発生しやすいが、拡大表示でなくても発生することもある。
	 * 念のため、IE 以外のブラウザでもやっておく。
	 */
	function _checkHeaderHeight(table)
	{
		var srcTrList = _tblSource.rows;
		var dstTrList = table.rows;
		var rows = dstTrList.length;
		var srcHeight = _rowsHeight(srcTrList, rows) + HEIGHT_MARGIN; // 少し余裕を持たせてチェック

		if (_rowsHeight(dstTrList, rows) < srcHeight) {
			return; // OK
		}

		//
		// 各セルの左右の padding を 1px ずつ減らす。
		// padding 未調査のセルは 0 にする。
		//
		for (var row = 0; row < rows; row++) {
			var diff = _trHeight(dstTrList[row]) - _trHeight(srcTrList[row]);
			if (diff < HEIGHT_MARGIN) continue;
			var srcTdList = _cells(srcTrList[row]);
			var dstTdList = _cells(dstTrList[row]);

			for (var col = 0; col < dstTdList.length; col++) {
				var srcTd = srcTdList[col];
				var dstTd = dstTdList[col];

				var padding = 0;
				var cellId = row + "." + srcTd.cellIndex;
				if (_cellPadding[cellId] != undefined) {
					if (_cellPadding[cellId] <= 0) continue;
					padding = _cellPadding[cellId] - 2;
					if (padding < 0) padding = 0;
				}
				var paddingLeft = Math.ceil(padding / 2); // 繰上げ
				var paddingRight = Math.floor(padding / 2); // 繰下げ

				dstTd.style.paddingLeft = paddingLeft + PX;
				dstTd.style.paddingRight = paddingRight + PX;
				dstTd.style.width = (srcTd.clientWidth - padding) + PX;

				if (_rowsHeight(dstTrList, rows) < srcHeight) {
					return; // OK
				}
			}
		}
	}

} // End of _FixedHeader

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * _FixedElementList 内部クラス
 */
function _FixedElementList()
{
	var _fixedElements = new Array();

	this.add = function(obj)
	{
		_fixedElements.push(obj);
	};


	this.remove = function(obj)
	{
		for (var i = 0; i < _fixedElements.length; i++) {
			if (_fixedElements[i] == obj) {
				_fixedElements.splice(i, 1);
				return;
			}
		}
	};


	this.get = function(src)
	{
		for (var i = 0; i < _fixedElements.length; i++) {
			var obj = _fixedElements[i];
			if (obj.$SOURCE_ELEMENT == src) {
				return obj;
			}
		}
		return null;
	};


	this.getAll = function(src)
	{
		var array = null;
		for (var i = 0; i < _fixedElements.length; i++) {
			var obj = _fixedElements[i];
			if (obj.$SOURCE_ELEMENT == src) {
				if (array == null) array = new Array();
				array.push(obj);
			}
		}
		return array;
	};

} // End of _FixedElementList

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * src と dst のリンク付けとidやイベントハンドラの設定
 */
function _linkElement(dst, src, tid, copyId, linkChild)
{
	if (src.id && !copyId) {
		dst.removeAttribute("id");
	}
	if (src.name) {
		if ((src.tagName == "INPUT") && (src.type == "radio")) {
			dst.name = RADIO_PREFIX + tid + "_" + src.name; // tid追加(ver1.4)
		} else {
			dst.removeAttribute("name");
		}
	}

	dst.$SOURCE_ELEMENT = src;

	if (_fixedList != null) {
		// (_isIE8 && 標準モード)の場合、
		// src.$FIXED_ELEMENT = dst
		// が非常に遅いので _fixedList で管理する
		_fixedList.add(dst);
	} else {
		// 1つの src に対して dst は1つまたは3つ
		if (!src.$FXH_FIXED_ELEMENT) {
			// 未設定の場合
			src.$FXH_FIXED_ELEMENT = dst;
		} else {
			// 設定済みの場合
			var obj = src.$FXH_FIXED_ELEMENT;
			if (!obj.$IS_ARRAY) {
				var array = new Array();
				array.$IS_ARRAY = true;
				src.$FXH_FIXED_ELEMENT = array;
				array.push(obj);
			}
			src.$FXH_FIXED_ELEMENT.push(dst);
		}
	}

	_setEventHandler(dst, src);

	if (linkChild) { // 子ノードも
		for (var i = 0; i < dst.childNodes.length; i++) {
			var child = dst.childNodes[i];
			if (!child) continue; // (_isIE) 元の html が不正な場合、undefined の場合がある。
			if (!child.tagName) continue; // テキストノード等
			_linkElement(child, src.childNodes[i], tid, copyId, linkChild); // 再帰
		}
	}
}


/*
 * src から dst へのリンクを削除
 */
function _unlinkElement(dst)
{
	if (_fixedList != null) { // (_isIE8 && 標準モード)
		_fixedList.remove(dst);
		var src = dst.$SOURCE_ELEMENT;
		if (src && src.$FXH_ON_CHANGE_FUNC && (_fixedList.get(src) == null)) {
			_removeEventListener(src, "change", src.$FXH_ON_CHANGE_FUNC);
			src.$FXH_ON_CHANGE_FUNC = undefined;
		}
	} else {
		var src = dst.$SOURCE_ELEMENT;
		if (src && src.$FXH_FIXED_ELEMENT) {
			var obj = src.$FXH_FIXED_ELEMENT;
			if (!obj.$IS_ARRAY || (obj.length == 1)) {
				src.$FXH_FIXED_ELEMENT = undefined;
				if (src.$FXH_ON_CHANGE_FUNC) {
					_removeEventListener(src, "change", src.$FXH_ON_CHANGE_FUNC);
					src.$FXH_ON_CHANGE_FUNC = undefined;
				}
			} else {
				var array = new Array();
				for (var i = 0; i < obj.length; i++) {
					if (obj[i] != dst) {
						array.push(obj[i]);
					}
				}
				src.$FXH_FIXED_ELEMENT = array;
			}
		}
	}

	for (var i = 0; i < dst.childNodes.length; i++) {
		var child = dst.childNodes[i];
		if (!child) continue; // (_isIE) 元の html が不正な場合、undefined の場合がある。
		if (!child.tagName) continue; // テキストノード等
		_unlinkElement(child); // 再帰
	}
}


/*
 * コピーしたエレメントのエベントハンドラを設定。
 * コピー元エレメントのエベントハンドラを呼び出す。
 */
function _setEventHandler(dst, src)
{
	if (src.onclick)     dst.onclick     = function() { return src.onclick();     };
	if (src.ondblclick)  dst.ondblclick  = function() { return src.ondblclick();  };
	if (src.onkeydown)   dst.onkeydown   = function() { return src.onkeydown();   };
	if (src.onkeypress)  dst.onkeypress  = function() { return src.onkeypress();  };
	if (src.onkeyup)     dst.onkeyup     = function() { return src.onkeyup();     };
	if (src.onmousedown) dst.onmousedown = function() { return src.onmousedown(); };
	if (src.onmouseup)   dst.onmouseup   = function() { return src.onmouseup();   };
	if (src.onmouseover) dst.onmouseover = function() { return src.onmouseover(); };
	if (src.onmouseout)  dst.onmouseout  = function() { return src.onmouseout();  };
	if (src.onmousemove) dst.onmousemove = function() { return src.onmousemove(); };

	//
	// コピーした input, select, textarea の onClick, onChange 処理。
	//
	if ((src.tagName == "INPUT") ||
	    (src.tagName == "SELECT") ||
	    (src.tagName == "TEXTAREA")) {

		// (_isIE) cloneNode() で selected, checked がコピーされないのでコピー
		switch (src.type) {
		case "select-one":
		case "select-multiple":
		case "checkbox":
			_copyValue(src, dst);
		}

		switch (src.type) {
		case "checkbox":
		case "radio":
		case "select-one":
		case "select-multiple":
		case "text":
		case "password":
		case "textarea":
			dst.onclick = function() {
				_copyValue(dst, src);  // dst の値を src にコピー
				_copyValues(src, dst); // src の値を他の dst にコピー
				if (src.onclick) return src.onclick();
				return true;
			};
			dst.onchange = function() {
				_copyValue(dst, src);  // dst の値を src にコピー
				_copyValues(src, dst); // src の値を他の dst にコピー
				if (src.onchange) return src.onchange();
				return true;
			};
			if (src.$FXH_ON_CHANGE_FUNC == undefined) {
				src.$FXH_ON_CHANGE_FUNC = function() {
					_copyValues(src); // src の値を全 dst にコピー
				};
				_addEventListener(src, "change", src.$FXH_ON_CHANGE_FUNC);
			}
			break;

		// コピーした button, submit, image, reset ボタンの onclick にも反応するようにした(ver1.9) 
		case "button":
		case "submit":
		case "image":
		case "reset":
			dst.onclick = function() { src.click(); };
			break;

		case "hidden":
		case "file":
			break;
		}

		// form の reset イベント処理(ver1.9)
		if (src.form && (src.$FXH_ON_RESET_FUNC == undefined)) {
			src.$FXH_ON_RESET_FUNC = function() {
				// form の全要素の値をコピー要素にコピー
				for (var i = 0; i < src.form.elements.length; i++) {
					_copyValues(src.form.elements[i]);
				}
			};
			// resetイベントは、reset処理完了前に(各要素の値がリセットされる前に)通知されるようだ
		//	_addEventListener(src.form, "reset", src.$FXH_ON_RESET_FUNC);
			_addEventListener(src.form, "reset", function() { setTimeout(src.$FXH_ON_RESET_FUNC, 30); });
		}
	} else
	if (src.tagName == "FORM") {
		//
		// コピーした <form> の onSubmit を無効化。
		// 現状の createTable() の実装では、次のように各行が <form> で囲まれているような場合、
		// この <form> はコピーしないため、ここは通らない。
		// <form><tr><td>...</td></tr></from>
		// <tr><form><td>...</td></from></tr>
		//
		dst.onsubmit = function() { return false; };
	}

	// TODO onFocus, onBlur, onSubmit, onReset, onSelect
}


function _copyValues(src, excludeElement)
{
	if (excludeElement == undefined) excludeElement = null; // しなくてもよい

	var dst;
	if (_fixedList != null) { // (_isIE8 && 標準モード)
		dst = _fixedList.getAll(src);
		if (!dst) return;
		for (var i = 0; i < dst.length; i++) { // 最大3つ
			if (dst[i] == excludeElement) continue;
			_copyValue(src, dst[i]);
		}
		return;
	}

	dst = src.$FXH_FIXED_ELEMENT;
	if (!dst) return;
	if (!dst.$IS_ARRAY) {
		if (dst != excludeElement) {
			_copyValue(src, dst);
		}
	} else {
		for (var i = 0; i < dst.length; i++) { // 最大3つ
			if (dst[i] == excludeElement) continue;
			_copyValue(src, dst[i]);
		}
	}
}


function _copyValue(src, dst)
{
	switch (src.type) {
	case "checkbox":
	case "radio":
		dst.checked = src.checked;
		break;
	case "select-one":
	case "select-multiple":
		for (var i = 0; i < src.length; i++) {
			dst.options[i].selected = src.options[i].selected;
		}
		dst.selectedIndex = src.selectedIndex;
		break;
	case "text":
	case "password":
	case "textarea":
		dst.value = src.value;
		break;
	default:
		// button, submit, image, reset, hidden, file
		try { dst.value = src.value; } catch (e) {}
		break;
	}
}


function _copyStyle(src, dst, styleName)
{
	var buf = styleName.split(",");
	for (var i = 0; i < buf.length; i++) {
		var name = _trim(buf[i]);
		try {
			eval("dst.style." + name + " = src.style." + name);
		} catch (e) {
		}
	}
}


function _setStyle(element, styleName, value)
{
	try {
		eval("element.style." + styleName + " = value");
		return true;
	} catch (e) {
		return false;
	}
}


/*
 * ラジオボタンの checked のコピー
 * (1) cloneNode() する前に src の checked を記憶(backup)
 * (2) cloneNode() した後に src の checked を元に戻す(restore)
 * (3) 最後に (name を変更した) dst に src の checked をコピー(sync)
 */
function _radioCtl(parentNode, cmd)
{
	var elements = parentNode.getElementsByTagName("INPUT");
	for (var i = 0; i < elements.length; i++) {
		var obj = elements[i];
		if (obj.type != "radio") continue;
		switch (cmd) {
		case "backup":  obj.$FXH_CHECKED = obj.checked; break; // src
		case "restore": obj.checked = obj.$FXH_CHECKED; break; // src
		case "sync":    obj.checked = obj.$SOURCE_ELEMENT.checked; break; // dst
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * CSS をドキュメントに追加
 */
function _createCSS(media, cssText)
{
	var style = document.createElement("STYLE");
	style.setAttribute("type", "text/css");
	style.setAttribute("media", media);
	if (style.styleSheet) { // (_isIE)
		style.styleSheet.cssText = cssText;
	} else {
		style.appendChild(document.createTextNode(cssText));
	}
	document.body.appendChild(style);
}


/*
 * 指定された要素の背景色(透明の場合は null)を返す。
 */
function _getBackgroundColor(element)
{
	var style = element.currentStyle || document.defaultView.getComputedStyle(element, '');
	var bgcolor = style.backgroundColor;
	if (bgcolor == "transparent") return null;
	if (bgcolor.match(/^rgba\(/) == null) return bgcolor;

	// rgba(x,x,x,x)指定の場合
	var str = bgcolor.replace(/^rgba\(/, "").replace(/\)/, ""); // ()内の数字を取得
	var rgba = str.split(",");
	var alpha = Number(rgba[3]);
	if (alpha == 1) return bgcolor; // 不透明ならそのまま返す

	// 半透明の場合は不透明にして返す
	var rgb = "rgb("
		+ _color(Number(rgba[0]), alpha) + ", "
		+ _color(Number(rgba[1]), alpha) + ", "
		+ _color(Number(rgba[2]), alpha)
		+ ")";
	return rgb;
}


/* alpha の分、適当に色を薄くする */
function _color(num, alpha)
{
	var n = Math.round(num + (255 - num) * (1 - alpha));
	return Math.min(n, 255);
}

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
 * イベントリスナーの追加
 */
function _addEventListener(obj, event, func)
{
	if (obj.addEventListener) {
		obj.addEventListener(event, func, false);
	} else
	if (obj.attachEvent) { // (_isIE)
		obj.attachEvent("on" + event, func);
	}
}


/*
 * イベントリスナーの削除
 */
function _removeEventListener(obj, event, func)
{
	if (obj.removeEventListener) {
		obj.removeEventListener(event, func, false);
	} else
	if (obj.detachEvent) { // (_isIE)
		obj.detachEvent("on" + event, func);
	}
}


/*
ScrollTop, ScrollLeft の取得
            互換  標準
IE          body  documentElement
Firefox     body  documentElement
Opera       body  documentElement
Chrome      body  body
Safari      body  body
*/
function _getBodyScrollTop()
{
	if (_isBackCompat || _isSafari) return document.body.scrollTop;
	if (_isChrome) {
		/* 2019-09-27 Chrome 61 bug fix */
		if (browserInfo.versionNo > 60) {
			return document.scrollingElement.scrollTop;
		}
		return document.body.scrollTop;
	}
	return document.documentElement.scrollTop;
}

function _getBodyScrollLeft()
{
	if (_isBackCompat || _isSafari) return document.body.scrollLeft;
	if (_isChrome) {
		/* 2019-09-27 Chrome 61 bug fix */
		if (browserInfo.versionNo > 60) {
			return document.scrollingElement.scrollLeft;
		}
		return document.body.scrollLeft;
	}
	return document.documentElement.scrollLeft;
}


/*
 * getBoundingClientRect() は、
 * Firefox, Chrome, Safari では小数付きの値になるが、
 * IE, Opera では小数はない。(2012/05)
 */

function _offsetWidth(element)
{
//	if (_isIE) return element.offsetWidth;
	var rect = _rect(element);
	return rect.right - rect.left;
}


function _offsetHeight(element)
{
//	if (_isIE) return element.offsetHeight;
	var rect = _rect(element);
	return rect.bottom - rect.top;
}


/*
 * 先頭行から指定された行数分の高さを返す
 */
function _rowsHeight(trList, rows)
{
	return _rect(trList[rows - 1]).bottom - _rect(trList[0]).top;
}


/*
 * 先頭列から指定された列数分の幅を返す
 */
function _colsWidth(tdList, cols)
{
	var columnCount = 0;
	for (var i = 0; i < tdList.length; i++) {
		columnCount += tdList[i].colSpan;
		if (columnCount == cols) {
			cols = i + 1;
			break;
		}
	}
	return _rect(tdList[cols - 1]).right - _rect(tdList[0]).left;
}


/*
 * 各カラムの幅をカンマ区切りのリストにして返す
 */
function _colsWidthList(tdList)
{
	var str = "";
	for (var i = 0; i < tdList.length; i++) {
		if (i > 0) str += ",";
		str += tdList[i].offsetWidth;
	}
	return str;
}


/*
 * 各行の高さをカンマ区切りのリストにして返す
 */
function _rowsHeightList(trList)
{
	var str = "";
	for (var i = 0; i < trList.length; i++) {
		if (i > 0) str += ",";
		str += trList[i].offsetHeight;
	}
	return str;
}


/*
 * 行の高さを返す。
 * 「IE8 & 標準モード」の場合、rowSpan の指定があると、TR.offsetHeight は、
 * rowSpan 分の行の高さを返してくるので clientHeight を使用。
 * ちなみに、IE9 の場合は、TR.clientHeight は常に 0 で返ってくる。
 */
function _trHeight(tr)
{
	if (_isIE && (_IEver == 8) && !_isBackCompat) { // (IE8 & 標準モード)
		return tr.clientHeight;
	} else {
		return _offsetHeight(tr);
	}
}


/*
 * 指定された tr 配下の td または th タグの配列を返す。
 * td と th が混在している場合は両方返す。
 * (_isIE && 標準モード)で、<tr> と <td> の間に <form> があったりすると、
 * TR.cells で取得できないので、本関数で取得。
 */
function _cells(tr)
{
	// セル内にさらに table がある場合に、そのセルまで取得しないように、
	// tr.getElementsByTagName("TD") 方式はやめる(ver1.9)
	var childNodes = tr.childNodes;
	if (childNodes.length == 0) return childNodes;
	var list = new Array();
	for (var i = 0; i < childNodes.length; i++) {
		var el = childNodes[i];
		if ((el.tagName == "TD") || (el.tagName == "TH")) {
			list.push(el);
		}
	}
	if ((list.length == 0) && (childNodes[0].tagName == "FORM")) {
 		// <tr> と <td> の間に <form> がある場合に限り、<form> の下を取得
		return _cells(childNodes[0]); // 再帰呼出
	}
	return list;
}


function _pixel(str)
{
	if (str.match(/px$/) != null) {
		str = str.substring(0, str.length - 2); // 最後の "px" を削除
	}
	return Number(str);
}


function _percent(str)
{
	if (str.match(/%$/) != null) {
		str = str.substring(0, str.length - 1); // 最後の "%" を削除
	}
	return Number(str);
}


/*
function _round(n) { return (n > 0) ? Math.round(n) : -Math.round(-n); } // 四捨五入
function _floor(n) { return (n > 0) ? Math.floor(n) : -Math.floor(-n); } // 切捨て
function _ceil(n)  { return (n > 0) ? Math.ceil(n)  : -Math.ceil(-n);  } // 切上げ
*/


function _trim(str)
{
	return str.replace(/^[ 　]+/, "").replace(/[ 　]+$/, "");
}


function _rect(element)
{
	return element.getBoundingClientRect();
}


function _getElementByTagName(element, tagName)
{
	var elements = element.getElementsByTagName(tagName);
	if (elements.length == 0) return null;
	return elements[0];
}


/*
 * おまじない。
 * Firefox の場合、処理の最初にこれをやっておくと速い。
 */
function _createObjectForFirefox()
{
	var span = document.createElement("SPAN");
	span.style.display = "none";
	span.style.position = "absolute";
	span.style.top = "0px";
	span.style.left = "0px";
	document.body.appendChild(span);
	return span;
}


/*
 * おまじない。(ver1.7)
 * Chrome の場合、スクロール時に固定見出しが前後に揺れることがあるが、
 * これをやっておくと発生しない。position=fixed の DOM があればいいようだ。
 */
function _createObjectForChrome()
{
	var span = document.createElement("SPAN");
	span.style.position = "fixed";
	span.style.position = "absolute";
	span.style.top = "0";
	span.style.left = "0";
	span.style.height = "0";
	span.style.width = "0";
	document.body.appendChild(span);
}

////////////////////////////////////////////////////////////////////////////////////////////////////

} // End of FixedMidashi


/*
 * onLoad 時に自動実行
 * ⇒ わかりづらくなるので自動実行はやらない
 * ・イベントリスナーの実行順序が定かでなく、
 *   適切な場所で FixedMidashi.create() が呼ばれる保証がない
 * ・自動実行させたくない場合もある
 */
//if (window.addEventListener) {
//	window.addEventListener("load", FixedMidashi.create, false);
//} else
//if (window.attachEvent) { // (_isIE)
//	window.attachEvent("onload", FixedMidashi.create);
//}

////////////////////////////////////////////////////////////////////////////////////////////////////
