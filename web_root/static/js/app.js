function setCartCounter() {
    if (typeof (Storage) !== "undefined") {
        let counters = document.getElementsByName("cart-counter");
        let inCart = localStorage.getItem("{{keyCart}}");
        if (inCart != null) {
            const counterVal = Object.keys(JSON.parse(inCart)).length;
            if (counterVal != 0) {
                for (let i = 0; i < counters.length; i++) {
                    counters[i].style.display = "";
                    counters[i].innerHTML = counterVal;
                }
                return;
            }
        }
        for (let i = 0; i < counters.length; i++) {
            counters[i].style.display = "none";
            counters[i].innerHTML = "";
        }
    }
}

function setProductsInCart() {
    if (typeof (Storage) !== "undefined") {
        let counters = document.getElementsByName("cart-counter");
        let inCart = localStorage.getItem("{{keyCart}}");
        if (inCart != null) {
            inCart = JSON.parse(inCart);
            for (let pid in inCart) {
                let carts = document.querySelectorAll("[data-product-id='" + pid + "']");;
                for (let i = 0; i < carts.length; i++) {
                    carts[i].classList.remove("fa-shopping-cart");
                    carts[i].classList.add("fa-cart-plus");
                }
            }
        }
    }
}

function putInCart(cart) {
    if (typeof (Storage) !== "undefined") {
        let inCart = localStorage.getItem("{{keyCart}}");
        if (inCart == null) {
            inCart = {};
        } else {
            inCart = JSON.parse(inCart);
        }

        const pid = cart.getAttribute("data-product-id");

        cart.classList.remove("fa-shopping-cart");
        cart.classList.remove("fa-cart-plus");
        $("#alertCart").remove();
        if (inCart[pid]) {
            delete inCart[pid];
            cart.classList.add("fa-shopping-cart");
            $('{{alertCartRemove}}').appendTo("body");
        } else {
            inCart[pid] = true;
            cart.classList.add("fa-cart-plus");
            $('{{alertCartPutIn}}').appendTo("body");
            $("#alertPutInCartCounter").html(Object.keys(inCart).length);
        }

        localStorage.setItem("{{keyCart}}", JSON.stringify(inCart));
        setCartCounter();
        $("#alertCart").slideDown("slow").delay(4000).fadeOut("slow");
    }
}

