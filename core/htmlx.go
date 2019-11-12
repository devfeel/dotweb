package core

import "strings"

func CreateTableHtml(title, header, body string) string {
	template := `<br><table class="bordered">
	<colgroup>
		  <col width="40%">
		  <col width="60%">
		</colgroup>
		<caption>{{title}}</caption>
	  <thead>
	 {{header}}
	  </thead>
		{{body}}
	</table>`
	html := strings.Replace(template, "{{title}}", title, -1)
	html = strings.Replace(html, "{{header}}", header, -1)
	html = strings.Replace(html, "{{body}}", body, -1)
	return html
}
