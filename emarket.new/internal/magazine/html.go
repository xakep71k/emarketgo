package magazine

const appHTML = `<!DOCTYPE html>
<html lang="ru">

<head>
    <meta charset="utf-8">
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" integrity="sha384-JcKb8q3iqJ61gNV9KGb8thSsNjpSL0n8PARn9HuZOnIxN0hoP+VmmDGMN5t9UJ0Z" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js" integrity="sha384-9/reFTGAW83EW2RDu2S0VKaIzap3H66lZH81PoYlFhbGU+6BZp6G7niu735Sk7lN" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js" integrity="sha384-B4gt1jrGC7Jh4AgTPSdUtOBvfO8shuf57BaghqFfPlYxofvL8/KUEfYiJOMMV+rV" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.14.0/css/all.min.css" integrity="sha512-1PKOgIY59xJ8Co8+NE6FZ+LOAZKjy+KY8iq0G4B3CyeY6wYHN3yt9PW0XpSriVlkMXe40PTKnXrLnZ9+fkDaog==" crossorigin="anonymous" />
    <link href="/static/css/all.css" rel="stylesheet" media="all">
    <script src="/static/js/all.js" async></script>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>{{.Title}}</title>
</head>

<body id="{{.ID}}">
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
                        <i class="fas fa-home"></i>&nbsp;Главная
                    </a>
                </li>
		<li class="nav-item">
		    <a class="nav-link {{if eq .CurrentPage 4}} active {{end}}" href="/istoriya_prosmotrov">
		        <i class="fas fa-eye"></i>&nbsp;Вы смотрели
		    </a>
		</li>
	<!--
		<li class="nav-item">
            <span name="cart-counter" class="counter" style="display: none;"></span>
		    <a class="nav-link {{if eq .CurrentPage 5}} active {{end}}" href="/zakazy/novyy">
		        <i class="fas fa-shopping-cart"></i>&nbsp;Корзина
		    </a>
		</li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 2}} active {{end}}" href="/dostavka">
                        <i class="fas fa-shipping-fast"></i>&nbsp;Доставка
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPage 3}} active {{end}}" href="/kontakty">
                        <i class="fas fa-id-card"></i>&nbsp;Свяжитесь с нами
                    </a>
                </li>
	-->
            </ul>
        </div>
    </nav>
    {{.Body}}
</body>
<script>
	$(function() {
		setCartCounter();
		setProductsInCart();
	});
</script>
</html>
`

const contactHTML = `
<h1 class="pageHeader text-center">{{.Title}}</h1>
<div class="contact-page">
    <h4 style="text-align: center;"><p>Связаться с нами можете любым удобным для Вас способом.</p><p>Мы рады Вам всегда! <i class="far fa-smile"></i></p></h4>
    <ul class="list-group contacts">
        <li class="list-group-item"><a class="nolink telegram" href="https://t.me/torg4u" target="_blank"><i class="fab fa-telegram-plane"></i> Telegram</a></li>
        <li class="list-group-item"><a class="nolink whatsapp" href="https://api.whatsapp.com/send?phone=79057740885" target="_blank">
            <i class="fab fa-whatsapp"></i> WhatsApp</a></li>
        <li class="list-group-item"><a class="nolink viber" href="viber://add?number=79057740885" target="_blank"><i class="fab fa-viber"></i> Viber</a></li>
        <li class="list-group-item"><a class="nolink email" href="mailto:support@torg4u.ru?subject=Вопрос" target="_blank"><i class="fas fa-at"></i> support@torg4u.ru</a></li>
        <li class="list-group-item"><a class="nolink phone" href="tel:+79057740885"><i class="fas fa-phone"></i> +7(905)774-08-85</a></li>
    </ul>
</div>
`
const magazineListHTML = `
<div class="main-container" onload="">
    <div class="productsContainer" id="products">
        {{range magazinemagazines}}
        <div class="card productCard" card-product-id="{{.ID}}">
            <div class="productThumb">
                <a href="/zhurnaly/{{.ID}}">
                    <img class="card-img-top" alt="{{.Title}}" src="/product/image/{{.ID}}" />
                </a>
            </div>
			<!--
            <div class="card-body align-bottom">
                {{if ne .Quantity 0}}
                <i class="fas fa-ruble-sign"></i>
                <span>{{.Price}}</span>
                {{else}}
                <span>под заказ</span>
                {{end}}
                <i onclick="putInCart(this)" class="fas fa-shopping-cart product-cart" data-name="product-cart" data-product-id="{{.ID}}"></i>
            </div>
			-->
        </div>
        {{end}}
    </div>
</div>
`

const magazineHTML = `
<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item">Журналы</li>
        <li class="breadcrumb-item"><a href="/zhurnaly/stranitsa/{{.PageNum}}">Страница {{.PageNum}}</a></li>
    <li class="breadcrumb-item active" aria-current="page">{{.Title}}</li>
  </ol>
</nav>
<h1 class="at-center text-center respH1" style="height: auto;">{{.Title}}</h1>
<div class="show-info container" id="productDetails" productId="{{.ID}}">
        <div class="row show-info">
                <div class="col-md-auto showCardLeft">
                        <div class="card">
							<img class="img-fluid view overlay" alt="{{.Title}}" itemprop="image" src="/product/image/{{.ID}}" />
							<!--
							<div class="card-body">
								<p>
									<b>
					{{if ne .Quantity 0}}
										Цена:
										<span>{{.Price}}</span>
										<span>рублей</span>
					{{else}}
										под заказ
					{{end}}
									</b>
									<i onclick="putInCart(this)" class="fas fa-shopping-cart product-cart" data-name="product-cart" data-product-id="{{.ID}}"></i>
								</p>
							</div>
							-->
                        </div>
                </div>
                <div class="col-lg"><p class="pre infoColor" itemprop="description">{{.Description}}</p></div>
        </div>
</div>
<br>
<div class="at-center" style="height: auto;">
<a href="/zhurnaly/stranitsa/{{.PageNum}}" class="btn"><i class="fas fa-arrow-alt-circle-left"></i>&nbsp;&nbsp;К списку журналов</a>
</div>
<script>
    if (typeof (Storage) !== "undefined") {
        let viewed = localStorage.getItem("{{keyHistory}}");
        if (viewed == null) {
            viewed = new Array();
        } else {
            viewed = JSON.parse(viewed);
        }
        const pid = "{{.ID}}";
        let index = 0;
        do {
            index = viewed.indexOf(pid);
            if (index > -1) {
                viewed.splice(index, 1);
            } else {
                break;
            }
        } while(true);
        viewed.unshift(pid);

        const limit = 30;
        if (viewed.length > limit) {
            viewed.splice(limit, viewed.length - limit);
        }
        localStorage.setItem("{{keyHistory}}", JSON.stringify(viewed));
    }
</script>
`

const paginationHTML = `
<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item">Журналы</a></li>
    <li class="breadcrumb-item active" aria-current="page">Страница {{.PageNum}}</li>
  </ol>
</nav>
<h1 class="pageHeader text-center">{{.Title}}</h1>
{{.ListHTML}}
{{$PrevPageNum := add .PageNum -1}}
{{$NextPageNum := add .PageNum 1}}
<div id="emarketPagination">
    <nav aria-label="page product navigation">
	<ul class="pagination">
	    {{if not .First}}
	    {{if ne (index .PageNumbers 0) 0}}
	    <li class="page-item">
                 <a class="page-link" href="/zhurnaly/stranitsa/1">
                    <i class="fas fa-angle-double-left"></i>
                 </a>
            </li>
	    {{end}}
	    <li class="page-item">
		<a class="page-link" href="/zhurnaly/stranitsa/{{$PrevPageNum}}" aria-label="предыдущая">
		    <i class="fas fa-angle-left"></i>
		</a>
	    </li>
	    {{end}}
	    {{range .PageNumbers}}
	    <li class="page-item {{if eq $PrevPageNum .}}active{{end}} ">
		<a class="page-link" href="/zhurnaly/stranitsa/{{add . 1}}">{{add . 1}}</a>
	    </li>
	    {{end}}
	    {{if not .Last}}
	    <li class="page-item">
		<a class="page-link" href="/zhurnaly/stranitsa/{{$NextPageNum}}" aria-label="следующая">
		    <i class="fas fa-angle-right"></i>
		</a>
	    </li>
	    {{$latestIndex := add (len .PageNumbers) -1}}
	    {{if ne (index .PageNumbers $latestIndex) (add .MaxPages -1)}}
	    <li class="page-item">
                 <a class="page-link" href="/zhurnaly/stranitsa/{{.MaxPages}}">
                    <i class="fas fa-angle-double-right"></i>
                 </a>
            </li>
	    {{end}}
	    {{end}}
	</ul>
    </nav>
</div>
`
