export class SearchableList extends HTMLElement {
    items: HTMLLIElement[]
    selected: HTMLLIElement | null
    addBtn: HTMLButtonElement
    searchInput: HTMLInputElement
    clearBtn: HTMLButtonElement

    constructor() {
        super();

        this.handleItemClick = this.handleItemClick.bind(this);
        this.handleAddClick = this.handleAddClick.bind(this);
        this.handleFilter = this.handleFilter.bind(this);
        this.handleClearSearchClick = this.handleClearSearchClick.bind(this);

        const itemsList = this.querySelectorAll(".searchable-list-item");
        const items: HTMLLIElement[] = [];
        itemsList.forEach(item => {
            if (item instanceof HTMLLIElement) {
                items.push(item);
            }
        });
        this.items = items;
        this.selected = null;

        const addBtn = this.querySelector(".searchable-list-add")
        if (addBtn instanceof HTMLButtonElement) {
            this.addBtn = addBtn
        } else {
            console.error("There should be a button with the class 'searchable-list-add' in this component!")
            this.addBtn = new HTMLButtonElement
        }

        const input = this.querySelector("input")
        if (input instanceof HTMLInputElement) {
            this.searchInput = input
        } else {
            console.error("This element should have an input!")
            this.searchInput = new HTMLInputElement
        }

        const clearBtn = this.querySelector(".searchable-list-search button");
        if (clearBtn instanceof HTMLButtonElement) {
            this.clearBtn = clearBtn
        } else {
            console.error("There should be a button in the '.searchable-list-search' div!")
            this.clearBtn = new HTMLButtonElement
        }
    }

    connectedCallback() {
        this.items.forEach(item => {
            item.addEventListener("click", this.handleItemClick)
            if (item.dataset.selected) {
                item.classList.add("searchable-list-item--selected")
                this.selected = item;
            }
        })

        this.addBtn.addEventListener("click", this.handleAddClick)
        this.searchInput.addEventListener("keyup", this.handleFilter)
        this.clearBtn.addEventListener("click", this.handleClearSearchClick)
    }

    disconnectedCallback() {
        this.items.forEach(item => {
            item.removeEventListener("click", this.handleItemClick);
        })

        this.addBtn.removeEventListener("click", this.handleAddClick)
        this.searchInput.removeEventListener("keyup", this.handleFilter)
        this.clearBtn.removeEventListener("click", this.handleClearSearchClick)
    }

    handleItemClick(e: Event) {
        if (e.target === null) {
            return;
        }
        if (e.target instanceof HTMLLIElement) {
            if (this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        }
    }

    handleAddClick() {
        if (this.selected) {
            this.selected.classList.remove("searchable-list-item--selected");
            this.selected = null;
        }
    }

    handleFilter(e: Event) {
        if (e.target === null) {
            return;
        }
        if (e.target instanceof HTMLInputElement) {
            const input = e.target as HTMLInputElement
            if (input.value === "") {
                this.items.forEach(item => {
                    item.hidden = false;
                })
            }
            this.items.forEach(item => {
                if (!item.textContent?.toLowerCase().includes(input.value.toLowerCase())) {
                    item.hidden = true;
                } else {
                    item.hidden = false;
                }
            })
        }
    }

    handleClearSearchClick() {
        this.searchInput.value = "";
        this.items.forEach(item => {
            item.hidden = false;
        });
    }
}