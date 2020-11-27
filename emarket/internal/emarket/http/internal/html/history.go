package html

const HistoryTemplate = `
<h1 class="pageHeader">{{.Title}}</h1>
<div class="at-center" id='historyPage'>
    <img id="loader" alt="подождите, пожалуйста..." src="/static/loader100x100.gif" />
</div>
<script>
  function getHistoryPage() {
    return $("#historyPage")[0];
  }
  function setNoContent() {
    getHistoryPage().innerHTML = '<div class="text-center">здесь ничего нет</div>';
  }
  $(function() {
	  if (typeof (Storage) !== "undefined") {
		let viewed = localStorage.getItem("{{keyHistory}}");
		viewed = JSON.parse(viewed);
		if (viewed != null && viewed.length != 0) {
		  fetch("/api/magazines",
			{
			  headers: {
				"Accept": "application/json",
				"Content-Type": "application/json"
			  },
			  method: "POST",
			  body: JSON.stringify(viewed)
			}).then(function (res) {
			  if (res.status == 200) {
				res.text().then(function (text) {
				  getHistoryPage().innerHTML = text;
				  let products = document.querySelectorAll("[card-product-id]");
				  if (products.length == 0) {
					setNoContent();
				  } else {
					getHistoryPage().classList.remove("at-center");
					setProductsInCart();
				  }
				  let dict = {};
				  for (let i = 0, max = products.length; i < max; i++) {
					dict[products[i].getAttribute("card-product-id")] = true;
				  }
				  let updateViewed = new Array();
				  for (let i = 0, max = viewed.length; i < max; i++) {
					if (dict[viewed[i]]) {
					  updateViewed.push(viewed[i]);
					}
				  }
				  localStorage.setItem("{{keyHistory}}", JSON.stringify(updateViewed));
				})
			  } else {
				setNoContent();
			  }
			}).catch(function (res) {
			  setNoContent();
			});
		} else {
		  setNoContent();
		}
	  }
  });
</script>
`
