package html

const AppTemplate = `<!DOCTYPE html>
<html lang="ru">

<head>
    <meta charset="utf-8">
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" integrity="sha384-JcKb8q3iqJ61gNV9KGb8thSsNjpSL0n8PARn9HuZOnIxN0hoP+VmmDGMN5t9UJ0Z" crossorigin="anonymous">
	<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
	<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js" integrity="sha384-9/reFTGAW83EW2RDu2S0VKaIzap3H66lZH81PoYlFhbGU+6BZp6G7niu735Sk7lN" crossorigin="anonymous"></script>
	<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js" integrity="sha384-B4gt1jrGC7Jh4AgTPSdUtOBvfO8shuf57BaghqFfPlYxofvL8/KUEfYiJOMMV+rV" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.14.0/css/all.min.css" integrity="sha512-1PKOgIY59xJ8Co8+NE6FZ+LOAZKjy+KY8iq0G4B3CyeY6wYHN3yt9PW0XpSriVlkMXe40PTKnXrLnZ9+fkDaog==" crossorigin="anonymous" />
    <link href="/static/css/app.css" rel="stylesheet" media="all">
    <script src="/static/js/app.js" async></script>
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
                    <a class="nav-link {{if eq .CurrentPageNum 1}} active {{end}}>" href="/">
                        <i class="fas fa-home"></i>&nbsp;Главная
                    </a>
                </li>
                <li class="nav-item">
                    <a class="nav-link {{if eq .CurrentPageNum 3}} active {{end}}" href="/kontakty">
                        <i class="fas fa-id-card"></i>&nbsp;Свяжитесь с нами
                    </a>
                </li>
				<li class="nav-item">
		            <a class="nav-link {{if eq .CurrentPageNum 4}} active {{end}}" href="/istoriya_prosmotrov">
		                <i class="fas fa-eye"></i>&nbsp;Вы смотрели
		            </a>
		        </li>
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
