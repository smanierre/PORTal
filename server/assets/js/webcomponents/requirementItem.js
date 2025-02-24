class RequirementItem extends HTMLElement {
    constructor() {
        super();
        this.removeItem = this.removeItem.bind(this);

        this.target = this.parentNode.parentNode.querySelector('select');
        this.removeIcon = this.querySelector('span');
        this.removeIcon.addEventListener('click', this.removeItem)
    }

    connectedCallback() {
        const option = document.createElement("option");
        option.selected = true;
        option.value = this.dataset.id
        this.target.appendChild(option);
    }

    disconnectedCallback() {
        const option = this.target.querySelector(`option[value=${this.dataset.id}]`);
        option.remove();
    }

    removeItem() {
        this.remove()
    }
}

customElements.define('requirement-item', RequirementItem);