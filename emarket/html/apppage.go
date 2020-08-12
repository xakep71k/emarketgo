package html

import (
	"text/template"
)

var AppPage htmlTemplate

func init() {
	var err error
	AppPage.template, err = template.New("app page").Funcs(defaultTemplateFuncs()).Parse(appPage)

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
    <link href="/static/css/app.css" rel="stylesheet">
    <script src="/bootstrap/js/jquery-3.5.1.slim.min.js"></script>
    <script src="/bootstrap/js/popper.min.js"></script>
    <script src="/bootstrap/js/bootstrap.min.js"></script>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>{{ .Title }}</title>
</head>

<body id="{{.ID}}">
	<script>
		function setCartCounter() {
			if (typeof (Storage) !== "undefined") {
				let counters = document.getElementsByName("cart-counter")
				let inCart = localStorage.getItem("{{keyCart}}")
				if (inCart != null) {
					const counterVal = Object.keys(JSON.parse(inCart)).length
					if (counterVal != 0) {
						for(let i = 0; i < counters.length; i++) {
							counters[i].style.display = ""
							counters[i].innerHTML = counterVal
						}
						return
					}
				}
				for(let i = 0; i < counters.length; i++) {
					counters[i].style.display = "none"
					counters[i].innerHTML = ""
				}
			}
		}
		function setProductsInCart() {
			if (typeof (Storage) !== "undefined") {
				let counters = document.getElementsByName("cart-counter")
				let inCart = localStorage.getItem("{{keyCart}}")
				if (inCart != null) {
					inCart = JSON.parse(inCart)
					for (let pid in inCart) {
						let carts = document.querySelectorAll("[data-product-id='" + pid + "']");
						for (let i = 0; i < carts.length; i++) {
							carts[i].classList.remove("fa-shopping-cart")
							carts[i].classList.add("fa-cart-plus")
						}
					}
				}
			}
		}
    </script>
    <nav class="navbar navbar-expand-lg">
        <a class="navbar-brand" href="/"><i class="fas fa-child"></i></a>
        <button class="navbar-toggler navbar-toggler-right collapsed" type="button" data-toggle="collapse"
            data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false"
            aria-label="Toggle navigation">
			<span name="cart-counter" class="counter" style="display: none;"></span>
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
		    <a class="nav-link {{if eq .CurrentPage 4}} active {{end}}" href="/istoriya_prosmotrov">
		        <i class="fas fa-eye"></i> Вы смотрели
		    </a>
		</li>
		<li class="nav-item">
            <span name="cart-counter" class="counter" style="display: none;"></span>
		    <a class="nav-link {{if eq .CurrentPage 5}} active {{end}}" href="/zakazy/novyy">
		        <i class="fas fa-shopping-cart"></i> Корзина
		    </a>
		</li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 2}} active {{end}}" href="/dostavka">
                        <i class="fas fa-shipping-fast"></i> Доставка
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 3}} active {{end}}" href="/kontakty">
                        <i class="fas fa-id-card"></i> Свяжитесь с нами
                    </a>
                </li>
            </ul>
        </div>
    </nav>
    {{.Body}}
</body>
<script>
	setCartCounter()
	setProductsInCart()
</script>
</html>
`
