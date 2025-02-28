class RequirementEditor extends HTMLElement {
  constructor() {
    super();
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCancel = this.handleCancel.bind(this);

    this.requirementSelector = this.parentElement.parentElement;
    this.form = this.querySelector("form");
    this.cancelButton = this.querySelector(".danger");
    this.idInput = this.querySelector("input[name=ID]");

    this.addEventListener("submit", this.handleSubmit);
    this.cancelButton.addEventListener("click", this.handleCancel);
  }

  async handleSubmit(e) {
    e.preventDefault();
    const data = new FormData(this.form);
    const req = {}
    data.entries().forEach(entry => {
      req[entry[0]] = entry[1]
    })
    if (this.idInput.value !== "") {
      req.ID = this.idInput.value
    }
    this.requirementSelector.setAddedOrUpdatedRequirement(req)
    this.parentElement.close()
    return
  }

  setFormData(requirement) {
    const idInput = this.form.querySelector('input[name="ID"]');
    const typeInput = this.form.querySelector('select[name="Type"]');
    const nameInput = this.form.querySelector('input[name="Name"]');
    const notesInput = this.form.querySelector('textarea[name="Notes"]');
    const daysValidForInput = this.form.querySelector('input[name="DaysValidFor"]');
    const referenceInput = this.form.querySelector('input[name="Reference"]');
    idInput.value = requirement.ID;
    typeInput.value = requirement.Type;
    nameInput.value = requirement.Name;
    notesInput.value = requirement.Notes;
    daysValidForInput.value = requirement.DaysValidFor;
    referenceInput.value = requirement.Reference;

    // Since this is an existing requirement, update the button to say "Update" instead of "Create"
    const submitBtnText = this.querySelector(".form-submit--text")
    submitBtnText.textContent = "Update";
  }

  handleCancel() {
    this.parentElement.close();
  }
}

customElements.define("requirement-editor", RequirementEditor);
