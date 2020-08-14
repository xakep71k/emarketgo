package html

var NewOrder = `
<h1 class="pageHeader">{{.Title}}</h1>
<div class="cart-page" id="cart-page">
</div>
<div id="loading" class="modal" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
	<div class="modal-header">
        <h5 class="modal-title">Загружаем...</h5>
      </div>
      <div class="modal-body">
		<img id="cart-loader" style="display: block; margin-right: auto; margin-left: auto;" alt="подождите, пожалуйста..." src="/static/loader100x100.gif" />
      </div>
    </div>
  </div>
</div>
<script>
    function setNoContent() {
		$("#loading").modal("hide")
        $("#cart-page").html('<div class="text-center">здесь ничего нет</div>')
		$("#cart-page").addClass('at-center')
    }

    function setError() {
		$("#loading").modal("hide")
        $("#cart-page").html('<div class="text-center">возникли неполадки :-( обратитесь, пожалуйста, для заказа <a href="/kontakty">к нам напрямую<a/></div>')
    }

    function reloadOrder() {
        if (typeof (Storage) !== "undefined") {
		    $("#loading").modal("show")
            let inCart = localStorage.getItem("{{keyCart}}")
            if (inCart == null) {
                inCart = {}
            } else {
                inCart = JSON.parse(inCart)
            }
            let products = new Array()
            for (let key in inCart) {
                if (inCart[key]) {
                    products.push(key)
                }
            }
            if (products.length != 0) {
                fetch("/api/cart",
                    {
                        headers: {
                            "Accept": "application/json",
                            "Content-Type": "application/json"
                        },
                        method: "POST",
                        body: JSON.stringify(products)
                    }).then(function (res) {
                        if (res.status == 200) {
                            res.text().then(function (text) {
                                if (text.length == 0) {
                                    setNoContent()
                                } else {
                                    $("#cart-page").html(text)
                                }
                            })
                        } else {
                            setError()
                        }
						$("#loading").modal("hide")
                    }).catch(function (res) {
                        setError()
                    });
            } else {
                setNoContent()
            }
        }
    }

	function removeProduct(pid) {
        if (typeof (Storage) !== "undefined") {
            let inCart = localStorage.getItem("{{keyCart}}")
            if (inCart == null) {
                inCart = {}
            } else {
                inCart = JSON.parse(inCart)
            }

		    delete inCart[pid]
			localStorage.setItem("{{keyCart}}", JSON.stringify(inCart))
			setCartCounter()
			reloadOrder()
		}
	}

	reloadOrder()
</script>
`

var OrderedProducts = `
<div id="cart-products">
	{{range .Products}}
	<div class="card productInCart" id="{{.ID}}" name="product">
		<div class="thumb-wrapper">
			<a href="/zhurnaly/{{.ID}}">
				<img src="/product/image/{{.ID}}" class="productImg" name="product-thumb">
			</a>
		</div>
		<div>
			<ul class="list-group noBorder">
				<li class="list-group-item product-title" name="title">{{.Title}}</li>
				<li class="list-group-item" data-name="in-stock"><i class="fas fa-ruble-sign"></i> <i style="font-style: normal;" name="price">{{.Price}}</i></li>
			</ul>
		</div>
		<button type="button" class="close close-button" aria-label="Close" name="close-button" onclick="removeProduct('{{.ID}}')">
			<span>&times;</span>
		</button>
	</div>
	{{end}}
</div>
{{if .Empty}}
<div class="text-center" id="emptyMessage">
    здесь ничего нет
</div>
{{else}}
<div class="card" id="cartSummary">
	<ul class="list-group summary">
		<li class="list-group-item">Итого стоимость заказа: <i style="font-style: normal;" id="totalPrice">{{.TotalPrice}}</i> руб.</li>
	</ul>
</div>
<div id="formWrapper" class="commonMargin">
    <form id="orderForm" action="/api/order" accept-charset="UTF-8" method="post">
        <label>Ваше имя</label>
        <input type="text" name="customer_name" />

        <label>Ваш телефон</label>
        <input type="text" name="phone_number" />

        <label>Ваш Email</label>
        <input type="text" name="email" />

        <label>Как с Вами связаться?</label>
        <select name="contact_type">
            <option selected="selected" value="telephone">по телефону</option>
            <option value="viber">через Viber</option>
            <option value="whatsapp">через WhatsApp</option>
            <option value="email">по Email</option>
            <option value="telegram">через Telegram</option>
        </select>
        <label for="agreement">
            <input type="checkbox" name="agreement" id="agreement" autocomplete="off" />
            Я соглашаюсь на обработку персональных данных
        </label>
        <div class="text-center">
            <input type="submit" name="commit" value="Отправить заказ" class="btn form-control btn-success"
                data-disable-with="Отправить заказ" />
        </div>
    </form>
</div>
{{end}}
`
