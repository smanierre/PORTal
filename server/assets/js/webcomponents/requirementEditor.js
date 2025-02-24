class RequirementEditor extends HTMLElement {
  constructor() {
    super();
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCancel = this.handleCancel.bind(this);

    this.requirementSelector = this.parentElement.parentElement;
    this.form = this.querySelector("form");
    this.cancelButton = this.querySelector(".danger");
    this.idInput = this.querySelector("input[name=id]");

    this.addEventListener("submit", this.handleSubmit);
    this.cancelButton.addEventListener("click", this.handleCancel);
  }

  async handleSubmit(e) {
    e.preventDefault();
    const data = new FormData(this.form);
    if (this.idInput.value === "") {
      const res = await fetch("/admin/requirements/add", {
        method: "POST",
        body: new URLSearchParams(data),
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      });
      if (res.status !== 200) {
        // TODO: Handle this
        console.log("ERROR");
        return;
      }

      const select = this.requirementSelector.querySelector("select");
      htmx.swap(".requirement-items span", await res.text(), {
        swapStyle: "beforeend",
      });
    } else {
      const res = await fetch(`/admin/requirements/${this.idInput.value}`, {
        method: "POST",
        body: new URLSearchParams(data),
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      });
      if (res.status !== 200) {
        //TODO: Handle this
        console.log("ERROR");
        return;
      }
    }
    this.parentElement.close();
  }

  handleCancel() {
    this.parentElement.close();
  }
}

customElements.define("requirement-editor", RequirementEditor);
