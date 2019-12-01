package core

import "strings"

var defaultCol = `<colgroup>
		  <col width="40%">
		  <col width="60%">
		</colgroup>`
var tableHtml = `<html>
<html><head>
   <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0;">
 
  <meta name="Generator" content="EditPlus">
  <meta name="Author" content="">
  <meta name="Keywords" content="">
  <meta name="Description" content="">
  <title>Dotweb</title>
    <style>
    .overtable {
      width: 100%;
      overflow: hidden;
      overflow-x: auto;
    }
    body {
      max-width: 780px;
       margin:0 auto;
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      font-size: 1rem;
      color: #444;
    }
    table {
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      *border-collapse: collapse;
      /* IE7 and lower */
      border-spacing: 0;
      width: 100%;
      border-collapse: collapse;
      overflow-x: auto
    }
    caption {
      font-family: 'Microsoft Yahei', 'trebuchet MS', 'Lucida sans', Arial;
    text-align: left;
    padding: .5rem;
    font-weight: bold;
    font-size: 110%;
    color: #666;
    }
    tr {
      border-top: 1px solid #dfe2e5
    }
    tr:nth-child(2n) {
      background-color: #f6f8fa
    }
    td,
    th {
      border: 1px solid #dfe2e5;
      padding: .6em 1em;
    }
    .bordered tr:hover {
      background: #fbf8e9;
    }
    .bordered td,
    .bordered th {
      border: 1px solid #ccc;
      padding: 10px;
      text-align: left;
    }
  </style>
  <script>
  
(function(doc, win) {
    window.MPIXEL_RATIO = (function () {
        var Mctx = document.createElement("canvas").getContext("2d"),
            Mdpr = window.devicePixelRatio || 1,
            Mbsr = Mctx.webkitBackingStorePixelRatio ||
                Mctx.mozBackingStorePixelRatio ||
                Mctx.msBackingStorePixelRatio ||
                Mctx.oBackingStorePixelRatio ||
                Mctx.backingStorePixelRatio || 1;
    
        return Mdpr/Mbsr;
    })();

    function addEventListeners(ele,type,callback){
    
        try{  // Chrome、FireFox、Opera、Safari、IE9.0及其以上版本
            ele.addEventListener(type,callback,false);
        }catch(e){
            try{  // IE8.0及其以下版本
                ele.attachEvent('on' + type,callback);
            }catch(e){  // 早期浏览器
                ele['on' + type] = callback;
            }
        }
    }

    var docEl = doc.documentElement,
        resizeEvt = 'orientationchange' in window ? 'orientationchange' : 'resize';
    window.recalc = function() {
            var clientWidth = docEl.clientWidth < 768 ? docEl.clientWidth : 768;
            if (!clientWidth) return;
            docEl.style.fontSize = 10 * (clientWidth / 320) *  window.MPIXEL_RATIO + 'px';
        };
    window.recalc();
    
    addEventListeners(win, resizeEvt, recalc);
})(document, window);

</script>
</head>
<body>
<div class="overtable">
{{tableBody}}
</div>
</body>
</html>
`

// CreateTablePart create a table part html by replacing flags
func CreateTablePart(col, title, header, body string) string {
	template := `<br><table class="bordered">
		{{col}}
		<caption>{{title}}</caption>
	  <thead>
	 {{header}}
	  </thead>
		{{body}}
	</table>`
	if col == "" {
		col = defaultCol
	}
	data := strings.Replace(template, "{{col}}", col, -1)
	data = strings.Replace(data, "{{title}}", title, -1)
	data = strings.Replace(data, "{{header}}", header, -1)
	data = strings.Replace(data, "{{body}}", body, -1)
	return data
}

// CreateTableHtml create a complete page html by replacing {{tableBody}} and table part html
func CreateTableHtml(col, title, header, body string) string {
	template := `<br><table class="bordered">
		{{col}}
		<caption>{{title}}</caption>
	  <thead>
	 {{header}}
	  </thead>
		{{body}}
	</table>`

	if col == "" {
		col = defaultCol
	}
	data := strings.Replace(template, "{{col}}", col, -1)
	data = strings.Replace(data, "{{title}}", title, -1)
	data = strings.Replace(data, "{{header}}", header, -1)
	data = strings.Replace(data, "{{body}}", body, -1)
	return CreateHtml(data)
}

// CreateHtml create a complete page html by replacing {{tableBody}}
func CreateHtml(tableBody string) string {
	return strings.Replace(tableHtml, "{{tableBody}}", tableBody, -1)
}
