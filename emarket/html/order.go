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
	<form id="orderForm" action="/api/order" accept-charset="UTF-8" method="post" class="was-validated">
		<div class="form-group">
			<label class="form-control-label" for="customer_name">Ваше имя</label>
			<input type="text" class="form-control" id="customer_name" name="customer_name" required/>
			<div class="invalid-feedback">имя не заполнено</div>
		</div>
		<div class="form-group">
			<label class="form-control-label" for="phone_number">Ваш телефон</label>
			<input type="tel" class="form-control" id="phone_number" name="phone_number" required/>
			<div class="invalid-feedback">телефон не заполнен</div>
		</div>
		<div class="form-group">
			<label class="form-control-label" for="email">Ваш Email</label>
			<input type="email" class="form-control" id="email" name="email" required/>
			<div class="invalid-feedback">Email не заполнен</div>
		</div>
		<div class="form-group">
			<label>Как с Вами связаться?</label>
			<select id="contact_type" name="contact_type" class="custom-select" required>
				<option value="">Выберите способ связи</option>
				<option value="telephone">по телефону</option>
				<option value="viber">через Viber</option>
				<option value="whatsapp">через WhatsApp</option>
				<option value="email">по Email</option>
				<option value="telegram">через Telegram</option>
			</select>
			<div class="invalid-feedback">способ связи не выбран</div>
		</div>
		<div class="custom-control custom-checkbox mb-3">
			<input type="checkbox" class="custom-control-input" id="agreement" name="agreement" autocomplete="off" required/>
			<label class="custom-control-label" for="agreement">Я соглашаюсь на обработку персональных данных</label>
			<div class="invalid-feedback">не дано согласие</div>
		</div>
        <div class="text-center">
            <button id="sendOrder" name="sendOrder" type="submit" class="btn form-control">Отправить заказ</button>
        </div>
    </form>
</div>
<div id="confirmOrder" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-body">
				<div class="text-center">
					<p>Отправить заказ?</p>
				</div>
            </div>
			<div class="modal-footer">
				<button type="button" data-dismiss="modal" class="btn btn-primary" id="orderYes">да</button>
				<button type="button" data-dismiss="modal" class="btn" id="orderNo">нет</button>
			</div>
        </div>
    </div>
</div>
<script>
function confirmOrder(e) {
	var form = $("#orderForm")
	e.preventDefault()
	e.stopPropagation()
	if(form[0].checkValidity() === false) {
		return
	}
	$('#confirmOrder').modal({
		backdrop: 'static',
		keyboard: false
	}).on('click', '#orderYes', function(e) {
		form.trigger('submit')
	})
	$("#cancel").on('click',function(e){
		e.preventDefault()
		e.stopPropagation()
		$('#confirmOrder').modal.model('hide')
	})
}

$("#sendOrder").click("on", confirmOrder)
</script>
{{end}}
`
