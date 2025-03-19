import { RequirementSelector } from "./RequirementSelector";
import { validateRequirementEditor } from "./helpers/RequirementHelper";

export class RequirementEditor extends HTMLElement {
  requirementSelector: RequirementSelector;
  form: HTMLFormElement;
  cancelButton: HTMLButtonElement | null;
  submitButton: HTMLButtonElement | null;
  dialog: HTMLDialogElement | null;
  submitTextElem: HTMLSpanElement | null;

  // Form inputs for later use
  idInput: HTMLInputElement | null;
  typeInput: HTMLSelectElement | null;
  nameInput: HTMLInputElement | null;
  notesInput: HTMLTextAreaElement | null;
  daysValidForInput: HTMLInputElement | null;
  referenceInput: HTMLInputElement | null;
  qualificationInput: HTMLSelectElement | null;

  // Labels for showing and hiding elements
  nameLabel: HTMLLabelElement | null;
  notesLabel: HTMLLabelElement | null;
  daysValidForLabel: HTMLLabelElement | null;
  qualificationLabel: HTMLLabelElement | null;

  constructor() {
    super();
    this.handleTypeChange = this.handleTypeChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
    this.handleFormUpdate = this.handleFormUpdate.bind(this);
    this.setQualificationDisplay = this.setQualificationDisplay.bind(this);
    this.setNonQualificationDisplay =
      this.setNonQualificationDisplay.bind(this);

    if (this.parentElement?.parentElement instanceof RequirementSelector) {
      this.requirementSelector = this.parentElement?.parentElement;
    } else {
      console.error(
        "This elements 2x parent needs to be a requirement selector!",
      );
      this.requirementSelector = new RequirementSelector();
    }

    const form = this.querySelector("form");
    if (form instanceof HTMLFormElement) {
      this.form = form;
    } else {
      console.error("There should be a form in this component!");
      this.form = new HTMLFormElement();
    }

    const cancelButton = this.querySelector(".danger");
    if (cancelButton instanceof HTMLButtonElement) {
      this.cancelButton = cancelButton;
    } else {
      this.cancelButton = null;
    }

    const idInput = this.querySelector("input[name=id]");
    if (idInput instanceof HTMLInputElement) {
      this.idInput = idInput;
    } else {
      this.idInput = null;
    }

    const typeInput = this.querySelector("select[name=type]");
    if (typeInput instanceof HTMLSelectElement) {
      this.typeInput = typeInput;
    } else {
      this.typeInput = null;
    }

    const nameInput = this.querySelector("input[name=name]");
    if (nameInput instanceof HTMLInputElement) {
      this.nameInput = nameInput;
    } else {
      this.nameInput = null;
    }

    const notesInput = this.querySelector("textarea[name=notes]");
    if (notesInput instanceof HTMLTextAreaElement) {
      this.notesInput = notesInput;
    } else {
      this.notesInput = null;
    }

    const daysValidForInput = this.querySelector("input[name=days_valid_for]");
    if (daysValidForInput instanceof HTMLInputElement) {
      this.daysValidForInput = daysValidForInput;
    } else {
      this.daysValidForInput = null;
    }

    const referenceInput = this.querySelector("input[name=reference]");
    if (referenceInput instanceof HTMLInputElement) {
      this.referenceInput = referenceInput;
    } else {
      this.referenceInput = null;
    }

    const qualificationInput = this.querySelector(
      "select[name=qualification_id]",
    );
    if (qualificationInput instanceof HTMLSelectElement) {
      this.qualificationInput = qualificationInput;
    } else {
      this.qualificationInput = null;
    }

    const submitTextElem = this.querySelector(".form-submit--text");
    if (submitTextElem instanceof HTMLSpanElement) {
      this.submitTextElem = submitTextElem;
    } else {
      this.submitTextElem = null;
    }

    if (this.parentElement instanceof HTMLDialogElement) {
      this.dialog = this.parentElement;
    } else {
      this.dialog = null;
    }

    const submitButton = this.querySelector(".form-submit");
    if (submitButton instanceof HTMLButtonElement) {
      this.submitButton = submitButton;
    } else {
      this.submitButton = null;
    }

    const nameLabel = this.querySelector("label[for=name]");
    if (nameLabel instanceof HTMLLabelElement) {
      this.nameLabel = nameLabel;
    } else {
      this.nameLabel = null;
    }

    const notesLabel = this.querySelector("label[for=notes]");
    if (notesLabel instanceof HTMLLabelElement) {
      this.notesLabel = notesLabel;
    } else {
      this.notesLabel = null;
    }

    const daysValidForLabel = this.querySelector("label[for=days_valid_for]");
    if (daysValidForLabel instanceof HTMLLabelElement) {
      this.daysValidForLabel = daysValidForLabel;
    } else {
      this.daysValidForLabel = null;
    }

    const qualificationLabel = this.querySelector(
      "label[for=qualification_id]",
    );
    if (qualificationLabel instanceof HTMLLabelElement) {
      this.qualificationLabel = qualificationLabel;
    } else {
      this.qualificationLabel = null;
    }
  }

  connectedCallback() {
    this.addEventListener("submit", this.handleSubmit);
    this.cancelButton?.addEventListener("click", this.handleCancel);
    this.form.addEventListener("change", this.handleFormUpdate);
    this.typeInput?.addEventListener("change", this.handleTypeChange);
    if (this.idInput!.value === "") {
      this.setQualificationDisplay();
    }
    const qualificationOptions =
      this.qualificationInput?.querySelectorAll("option");
    qualificationOptions?.forEach((option) => {
      if (option.value === this.dataset.parentQualificationId) {
        option.hidden = true;
        if (qualificationOptions.length === 1) {
          this.qualificationLabel!.hidden = true;
          this.qualificationInput!.disabled = true;
        }
      }
    });
    if (qualificationOptions?.length === 0) {
      this.qualificationLabel!.hidden = true;
      this.qualificationInput!.disabled = true;
    }
  }

  disconnectedCallback() {
    this.cancelButton?.removeEventListener("click", this.handleCancel);
    this.removeEventListener("submit", this.handleSubmit);
    this.form.removeEventListener("change", this.handleFormUpdate);
    this.typeInput?.removeEventListener("change", this.handleTypeChange);
  }

  handleFormUpdate() {
    if (!this.submitButton) {
      return;
    }
    const valid = validateRequirementEditor(new FormData(this.form));
    if (valid) {
      this.submitButton.disabled = false;
    } else {
      this.submitButton.disabled = true;
    }
  }

  async handleSubmit(e: Event) {
    e.preventDefault();
    if (
      this.typeInput?.value !== "WBT" &&
      this.typeInput?.value !== "Qualification" &&
      this.typeInput?.value !== "Qualification"
    ) {
      return;
    }
    const req: Requirement = {
      id: this.idInput?.value || "",
      days_valid_for: this.daysValidForInput?.value
        ? +this.daysValidForInput.value
        : 0,
      name: this.nameInput?.value || "",
      notes: this.notesInput?.value || "",
      reference: this.referenceInput?.value || "",
      type: this.typeInput?.value || "",
      qualification_id: this.qualificationInput?.value || "",
    };
    if (req.type === "Qualification" && req.qualification_id) {
      this.requirementSelector.setAddedOrUpdatedQualificationRequirement(
        req,
        this.qualificationInput![this.qualificationInput!.selectedIndex]
          .textContent!,
      );
      this.dialog?.close();
      return;
    }
    this.requirementSelector.setAddedOrUpdatedRequirement(req);
    this.dialog?.close();
    return;
  }

  setFormData(requirement: Requirement) {
    this.idInput!.value = requirement.id;
    this.referenceInput!.value = requirement.reference;
    this.typeInput!.value = requirement.type;
    if (requirement.type === "Qualification") {
      this.qualificationInput!.value = requirement.qualification_id;
      this.nameInput!.disabled = true;
      this.notesInput!.disabled = true;
      this.daysValidForInput!.disabled = true;
      this.setQualificationDisplay();
    } else {
      this.setNonQualificationDisplay();
      this.nameInput!.value = requirement.name;
      this.notesInput!.value = requirement.notes;
      this.daysValidForInput!.value = requirement.days_valid_for.toString();
      this.qualificationInput!.disabled = true;
    }
    // Since this is an existing requirement, update the button to say "Update" instead of "Create"
    this.submitTextElem!.textContent = "Update";
  }

  handleTypeChange(e: Event) {
    if (e.target === null) {
      return;
    }
    if (e.target instanceof HTMLSelectElement) {
      switch (e.target.value) {
        case "WBT":
          this.setNonQualificationDisplay();
          break;
        case "Grade":
          this.setNonQualificationDisplay();
          break;
        case "Qualification":
          this.nameLabel!.hidden = true;
          this.nameInput!.disabled = true;

          this.notesLabel!.hidden = true;
          this.notesInput!.disabled = true;

          this.daysValidForLabel!.hidden = true;
          this.daysValidForInput!.disabled = true;

          this.qualificationLabel!.hidden = false;
          this.qualificationInput!.disabled = false;
          this.setQualificationDisplay();
      }
    }
  }

  setQualificationDisplay() {
    this.nameLabel!.hidden = true;
    this.nameInput!.disabled = true;

    this.notesLabel!.hidden = true;
    this.notesInput!.disabled = true;

    this.daysValidForLabel!.hidden = true;
    this.daysValidForInput!.disabled = true;

    this.qualificationLabel!.hidden = false;
    this.qualificationInput!.disabled = false;
  }

  setNonQualificationDisplay() {
    this.nameLabel!.hidden = false;
    this.nameInput!.disabled = false;

    this.notesLabel!.hidden = false;
    this.notesInput!.disabled = false;

    this.daysValidForLabel!.hidden = false;
    this.daysValidForInput!.disabled = false;

    this.qualificationLabel!.hidden = true;
    this.qualificationInput!.disabled = true;
  }

  handleCancel() {
    this.dialog?.close();
  }
}
