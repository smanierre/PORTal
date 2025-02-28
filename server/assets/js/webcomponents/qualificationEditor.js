class QualificationEditor extends HTMLElement {
    constructor() {
        super();

        this.submitBtn = this.querySelector(".form-submit");
        this.form = this.querySelector("form")
        this.initialRequirementsSelector = this.querySelector('requirement-selector[data-reqtype="initial_requirements"]')
        this.recurringRequirementsSelector = this.querySelector('requirement-selector[data-reqtype="recurring_requirements"]')
    }

    get initialRequirements() {
        return JSON.parse(this.initialRequirementsSelector.getAttribute("requirements"))
    }

    get recurringRequirements() {
        return JSON.parse(this.recurringRequirementsSelector.getAttribute("requirements"))
    }

    // TODO: Make expiration days field visible based on whether or not expires is checked
    // TODO: Update update requirement button enable/disable

    connectedCallback() {
        this.submitBtn.addEventListener("click", this.handleSubmit.bind(this));
    }

    disconnectedCallback() {
        this.submitBtn.removeEventListener("click", this.handleSubmit.bind(this));
    }

    handleSubmit(e) {
        e.preventDefault();
        const formData = new FormData(this.form)
        console.log(formData)
        console.log(this.initialRequirements)
    }

}

customElements.define("qualification-editor", QualificationEditor);