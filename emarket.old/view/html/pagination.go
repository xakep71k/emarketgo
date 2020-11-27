package html

var Pagination = `
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
