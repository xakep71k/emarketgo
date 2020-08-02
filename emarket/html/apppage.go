package html

import (
	"text/template"
)

var AppPage htmlTemplate

func init() {
	var err error
	AppPage.template, err = template.New("AppPage").Parse(appPage)

	if err != nil {
		panic(err)
	}
}

const appPage = `
<!DOCTYPE html>
<html lang="ru">

<head>
    <link href="/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="/fontawesome/css/all.min.css" rel="stylesheet">
    <link href="/static/css/custom_bootstrap.css" rel="stylesheet">
    <link href="/static/css/close_btn.css" rel="stylesheet">
    <link href="/static/css/app.css" rel="stylesheet">
    <script src="/bootstrap/js/jquery-3.5.1.slim.min.js"></script>
    <script src="/bootstrap/js/popper.min.js"></script>
    <script src="/bootstrap/js/bootstrap.min.js"></script>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>{{ .Title }}</title>
</head>

<body {{if .ID}} id="{{.ID}}" {{end}}>
    <nav class="navbar navbar-expand-lg">
        <a class="navbar-brand" href="/"><i class="fas fa-child"></i></a>
        <button class="navbar-toggler navbar-toggler-right collapsed" type="button" data-toggle="collapse"
            data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false"
            aria-label="Toggle navigation">
            <span class="my-1 mx-2 close">X</span>
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mr-auto">
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 1}} active {{end}}>" href="/">
                        <i class="fas fa-home"></i> Главная
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 2}} active {{end}} "
                        href="/dostavka">
                        <i class="fas fa-shipping-fast"></i> Доставка
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 3}} active {{end}} "
                        href="/kontakty">
                        <i class="fas fa-id-card"></i> Свяжитесь с нами
                    </a>
                </li>
            </ul>
    </nav>
    {{.Body}}
</body>

</html>
`
