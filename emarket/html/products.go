package html

var ProductList = `
<div class="main-container">
    <div class="productsContainer" id="products">
        {{range .}}
        <div class="card productCard">
            <div class="productThumb">
                <a href="/zhurnaly/{{.ID}}">
                    <img class="card-img-top" alt="{{.Title}}" src="/product/image/{{.ID}}" />
                </a>
            </div>
            <div class="card-body align-bottom">
                {{if ne .Quantity 0}}
                <i class="fas fa-ruble-sign"></i>
                <span>{{.Price}}</span>
                {{else}}
                <span>под заказ</span>
                {{end}}
                <i class="fas fa-shopping-cart product-cart" data-name="product-cart" data-product-id="{{.ID}}"></i>
            </div>
        </div>
        {{end}}
    </div>
</div>
`

var Products = `
<h1 class="pageHeader text-center">Журналы и выкройки для шитья</h1>
` + ProductList

var Product = `
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
                                <i class="fas fa-shopping-cart product-cart" data-name="product-cart" data-product-id="{{.ID}}"></i>
                            </p>
                        </div>
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
    const key = "emarket.history.v1"
    if (typeof (Storage) !== "undefined") {
        let viewed = localStorage.getItem(key)
        if (viewed == null) {
            viewed = new Array()
        } else {
            viewed = JSON.parse(viewed)
        }
        const pid = "{{.ID}}"
        let index = 0
        do {
            index = viewed.indexOf(pid)
            if (index > -1) {
                viewed.splice(index, 1)
            } else {
                break
            }
        } while (true)
        viewed.unshift(pid)

        const limit = 30
        if (viewed.length > limit) {
            viewed.splice(limit, viewed.length - limit)
        }
        localStorage.setItem(key, JSON.stringify(viewed))
    }
</script>
`
