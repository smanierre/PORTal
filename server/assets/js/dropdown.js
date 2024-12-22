document.addEventListener("alpine:init", () => {
    Alpine.data("dropdown", () => ({
        open: false,
        selected: "Members",
        toggle() {
            this.open = ! this.open
        },
        selectHandler(e) {
            this.selected = e
            this.open = false
        }
    }))
})