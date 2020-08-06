package html

var ProductList = `
<h1 class="pageHeader text-center">Журналы и выкройки для шитья</h1>
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

var Pagination = `
<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item">Журналы</a></li>
    <li class="breadcrumb-item active" aria-current="page">Страница {{add .Index 1}}</li>
  </ol>
</nav>
{{.ListHTML}}
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
		<a class="page-link" href="/zhurnaly/stranitsa/{{.Index}}" aria-label="предыдущая">
		    <i class="fas fa-angle-left"></i>
		</a>
	    </li>
	    {{end}}
	    {{range .PageNumbers}}
	    <li class="page-item {{if eq $.Index .}}active{{end}} ">
		<a class="page-link" href="/zhurnaly/stranitsa/{{add . 1}}">{{add . 1}}</a>
	    </li>
	    {{end}}
	    {{if not .Last}}
	    <li class="page-item">
		<a class="page-link" href="/zhurnaly/stranitsa/{{add .Index 2}}" aria-label="следующая">
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
