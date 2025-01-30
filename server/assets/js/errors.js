document.addEventListener("htmx:beforeOnLoad", (e) => {
    if(e.detail.xhr.status === 401 || e.detail.xhr.status === 500) {
        e.detail.shouldSwap = true;
        e.detail.isError = false;
    }
})