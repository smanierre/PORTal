import { RequirementSelector } from "./RequirementSelector";
import { validateQualificationEditor } from "./helpers/QualificationHelper";

export class QualificationEditor extends HTMLElement {
  submitBtn: HTMLButtonElement;
  form: HTMLFormElement;
  initialRequirementsSelector: RequirementSelector;
  recurringRequirementsSelector: RequirementSelector;
  expirationInput: HTMLInputElement;
  expirationIntervalLabel: HTMLLabelElement;
  expirationIntervalInput: HTMLInputElement;
  initialRequirementsInput: HTMLInputElement;
  recurringRequirementsInput: HTMLInputElement;

  constructor() {
    super();

    const submitBtn = this.querySelector(".form-submit");
    if (submitBtn instanceof HTMLButtonElement) {
      this.submitBtn = submitBtn;
    } else {
      console.error(
        "This component should have a button with class 'form-submit'!",
      );
      this.submitBtn = new HTMLButtonElement();
    }

    const form = this.querySelector("form");
    if (form instanceof HTMLFormElement) {
      this.form = form;
    } else {
      console.error("This component should have a form element!");
      this.form = new HTMLFormElement();
    }

    const initialRequirementsSelector = this.querySelector(
      'requirement-selector[data-reqtype="initial_requirements"]',
    );
    if (initialRequirementsSelector !== null) {
      this.initialRequirementsSelector =
        initialRequirementsSelector as RequirementSelector;
    } else {
      console.error(
        "This component should have a requirement-selector element with data-reqtype='initial_requirements'",
      );
      this.initialRequirementsSelector = new RequirementSelector();
    }

    const recurringRequirementsSelector = this.querySelector(
      'requirement-selector[data-reqtype="recurring_requirements"]',
    );
    if (recurringRequirementsSelector !== null) {
      this.recurringRequirementsSelector =
        recurringRequirementsSelector as RequirementSelector;
    } else {
      console.error(
        "This component should have a requirement-selector element with data-reqtype='recurring_requirements'",
      );
      this.recurringRequirementsSelector = new RequirementSelector();
    }

    const expirationInput = this.querySelector('input[name="expires"]');
    if (expirationInput instanceof HTMLInputElement) {
      this.expirationInput = expirationInput;
    } else {
      console.error(
        "This component should have an input element with name 'expires'!",
      );
      this.expirationInput = new HTMLInputElement();
    }

    const expirationIntervalInput = this.querySelector(
      'input[name="expiration_interval"]',
    );
    if (expirationIntervalInput?.parentElement instanceof HTMLLabelElement) {
      this.expirationIntervalLabel = expirationIntervalInput.parentElement;
    } else {
      console.error(
        "This component should have an input element with name 'expiration_interval' with a parent label!",
      );
      this.expirationIntervalLabel = new HTMLLabelElement();
    }
    if (expirationIntervalInput instanceof HTMLInputElement) {
      this.expirationIntervalInput = expirationIntervalInput;
    } else {
      console.error(
        "This component should have an input element with name 'expiration_interval'!",
      );
      this.expirationIntervalInput = new HTMLInputElement();
    }

    const initialRequirementsInput = this.querySelector(
      'input[name="initial_requirements"]',
    );
    if (initialRequirementsInput instanceof HTMLInputElement) {
      this.initialRequirementsInput = initialRequirementsInput;
    } else {
      console.error(
        "This component should have an input element with name 'initial_requirements'!",
      );
      this.initialRequirementsInput = new HTMLInputElement();
    }

    const recurringRequirementsInput = this.querySelector(
      'input[name="recurring_requirements"]',
    );
    if (recurringRequirementsInput instanceof HTMLInputElement) {
      this.recurringRequirementsInput = recurringRequirementsInput;
    } else {
      console.error(
        "This component should have an input element with name 'recurring_requirements'!",
      );
      this.recurringRequirementsInput = new HTMLInputElement();
    }
  }

  get initialRequirements(): Requirement[] {
    const data = this.initialRequirementsSelector.getAttribute("requirements");
    return data === null ? [] : JSON.parse(data);
  }

  get recurringRequirements(): Requirement[] {
    const data =
      this.recurringRequirementsSelector.getAttribute("requirements");
    return data === null ? [] : JSON.parse(data);
  }

  connectedCallback() {
    this.submitBtn.addEventListener("click", this.handleSubmit.bind(this));
    this.expirationInput.addEventListener(
      "change",
      this.handleExpirationChange.bind(this),
    );
    this.form.addEventListener("change", this.handleFormUpdate.bind(this));
    this.initialRequirementsSelector.addEventListener(
      "updated",
      this.handleFormUpdate.bind(this),
    );
  }

  disconnectedCallback() {
    this.submitBtn.removeEventListener("click", this.handleSubmit.bind(this));
    this.expirationInput.removeEventListener(
      "change",
      this.handleExpirationChange.bind(this),
    );
    this.form.addEventListener("change", this.handleFormUpdate.bind(this));
    this.initialRequirementsSelector.removeEventListener(
      "updated",
      this.handleFormUpdate.bind(this),
    );
  }

  handleExpirationChange(e: Event) {
    if (e.target instanceof HTMLInputElement) {
      if (e.target.checked) {
        this.expirationIntervalLabel.hidden = false;
        this.expirationIntervalInput.disabled = false;
      } else {
        this.expirationIntervalLabel.hidden = true;
        this.expirationIntervalInput.disabled = true;
      }
    }
  }

  handleFormUpdate() {
    const valid = validateQualificationEditor(
      new FormData(this.form),
      this.initialRequirements,
      this.recurringRequirements,
    );
    if (valid) {
      this.submitBtn.disabled = false;
    } else {
      this.submitBtn.disabled = true;
    }
  }

  handleSubmit(e: Event) {
    e.preventDefault();
    // Clean up any IDs that were temporarily set
    const cleanInitialReqs = this.initialRequirements.map((requirement) => {
      if (requirement.id.length < 20) {
        // It would take a lottttttttt of temporary ID's for the dummy IDs to be as long as a UUID
        requirement.id = "";
      }
      return requirement;
    });
    const cleanRecurringReqs = this.recurringRequirements.map((requirement) => {
      if (requirement.id.length < 20) {
        // It would take a lottttttttt of temporary ID's for the dummy IDs to be as long as a UUID
        requirement.id = "";
      }
      return requirement;
    });
    this.initialRequirementsInput.value = JSON.stringify(cleanInitialReqs);
    this.recurringRequirementsInput.value = JSON.stringify(cleanRecurringReqs);
    this.form.dispatchEvent(new Event("submit", { cancelable: true }));
  }
}
