export class Dropdown extends HTMLElement {
    list: HTMLUListElement
    selected: string
    btn: HTMLButtonElement
    options: HTMLLIElement[]

    constructor() {
        super();

        this.handleClickInside = this.handleClickInside.bind(this);
        this.handleClickOutside = this.handleClickOutside.bind(this);
        this.handleItemClick = this.handleItemClick.bind(this);

        const list = this.querySelector("ul");
        if (list instanceof HTMLUListElement) {
            this.list = list
        } else {
            console.error("There should be a ul element in this component!")
            this.list = new HTMLUListElement
        }
        this.selected = ""

        const btn = this.querySelector("button");
        if (btn instanceof HTMLButtonElement) {
            this.btn = btn
        } else {
            console.error("There should be a button in this component!")
            this.btn = new HTMLButtonElement
        }

        const optionsQuery = this.querySelectorAll("li");
        const options: HTMLLIElement[] = [];
        optionsQuery.forEach(option => {
            if (option instanceof HTMLLIElement) {
                options.push(option)
            }
        })
        this.options = options

    }

    connectedCallback() {
        this.selected = this.dataset.selected || "Members"
        this.btn.textContent = this.selected;
        this.list.hidden = true;
        this.btn.addEventListener("click", this.handleClickInside)
        window.addEventListener("click", this.handleClickOutside)

        this.options.forEach(option => {
            option.addEventListener("click", this.handleItemClick)
        })
    }

    disconnectedCallback() {
        this.btn.removeEventListener("click", this.handleClickInside)
        window.removeEventListener("click", this.handleClickOutside)
        this.options.forEach(option => {
            option.removeEventListener("click", this.handleItemClick)
        })
    }

    handleClickInside(e: Event) {
        e.stopPropagation()
        this.list.hidden = !this.list.hidden;
    }

    handleClickOutside() {
        this.list.hidden = true;
    }

    handleItemClick(e: Event) {
        if (e.target === null) {
            return;
        }
        if (e.target instanceof HTMLLIElement) {
            this.btn.textContent = e.target.textContent
            this.list.hidden = true;
        }
    }
}