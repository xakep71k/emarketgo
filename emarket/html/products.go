package html

var ProductList = `
<h1 class="pageHeader">Журналы и выкройки для шитья</h1>
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
