document.querySelector(".login-form-container").addEventListener("htmx:beforeOnLoad", (e) => {
    if(e.detail.xhr.status === 401) {
        e.detail.shouldSwap = true;
        e.detail.isError = false;
    }
})

document.querySelector(".login-submit").addEventListener("click", (e) => {
    document.querySelector("#loginError").className = "login-error--hidden"
})