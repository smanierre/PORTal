class RequirementSelector extends HTMLElement {
  static #itemElement = null;

  constructor() {
    super();

    // Get dialog for requirement editor and give it a unique ID for HTMX to target
    this.dialog = this.querySelector("dialog");
    const randLetter = String.fromCharCode(65 + Math.floor(Math.random() * 26));
    this.dialog.id = randLetter + Date.now();

    // Set the add button for setup during connectedCallback
    this.addButton = this.querySelector("#add-button");

    // Set the element that the list of requirements will be appended to
    this.requirementList = this.querySelector(".requirement-items span");

    // Set the select element to append options to
    this.select = this.querySelector("select");
  }

  // itemElement is an instance of the HTML returned when a requirement is added to a qualification.
  // For existing qualifications we need to get the itemElement from the server and set it here for instances to use.
  static get item() {
    UPDATE THIS TO RETURN A PROMISE THAT CAN BE AWAITED!
    if (this.#itemElement === null) {
      fetch("/components/requirementItem", {
        method: "GET",
      }).then((res) => {
        if (res.ok) {
          res.text().then((html) => {
            const item = document.createElement("button");
            item.outerHTML = html;
            this.#itemElement = item;
            return this.#itemElement;
          });
        } else {
          // TODO: Handle this
          console.error("Failed to load requirement item from server!");
        }
      });
    } else {
      return this.#itemElement.cloneNode(true);
    }
  }

  connectedCallback() {
    // Add the event listener for the add button
    this.addButton.addEventListener("click", this.addRequirement.bind(this));

    // For each of the existing requirements, add event listeners and add an <option> to the select with the id as a value.
    JSON.parse(this.getAttribute("requirements")).forEach(
      async (requirement) => {
        const req = await RequirementSelector.item;
        req.textContent = requirement.Name;
        req.setAttribute("data-id", requirement.ID);
        req.addEventListener("click", this.updateRequirement.bind(this));
        this.requirementList.appendChild(req);

        // const cancelIcon = requirement.querySelector(
        //   ".remove-requirement-button",
        // );
        // cancelIcon.addEventListener("click", this.removeRequirement.bind(this));

        const option = document.createElement("option");
        option.value = requirement.ID;
        option.selected = true;
        this.select.appendChild(option);
      },
    );
  }

  disconnectedCallback() {
    // Remove the event listener for the add button
    this.addButton.removeEventListener("click", this.addRequirement);

    const requirements = this.querySelectorAll(".requirement-item");

    requirements.forEach((requirement) => {
      requirement.removeEventListener(
        "click",
        this.updateRequirement.bind(this),
      );

      const cancelIcon = requirement.querySelector(
        ".remove-requirement-button",
      );
      cancelIcon.removeEventListener(
        "click",
        this.removeRequirement.bind(this),
      );
    });
  }

  async addRequirement() {
    try {
      await htmx.ajax("GET", "/admin/requirements/add", `#${this.dialog.id}`);
    } catch (e) {
      // TODO: Handle this
      console.log("caught error");
      return;
    }
    this.dialog.showModal();
  }

  async updateRequirement(e) {
    await htmx.ajax("GET", `/admin/requirements/${e.target.dataset.id}`, {
      target: `#${this.dialog.id}`,
      handler: (e, f) => {
        if (f.xhr.status !== 200) {
          htmx.swap("#toast", f.xhr.response, { swapStyle: "outerHTML" });
        } else {
          htmx.swap(`#${this.dialog.id}`, f.xhr.response, {
            swapStyle: "innerHTML",
          });
          this.dialog.showModal();
        }
      },
    });
  }

  async removeRequirement(e) {
    e.stopPropagation();
    const res = await fetch(
      `/admin/requirements/${e.target.parentElement.dataset.id}`,
      {
        method: "DELETE",
      },
    );
    // Remove requirement from select list
    const options = this.select.querySelectorAll("option");
    options.forEach((option) => {
      if (option.value === e.target.parentElement.dataset.id) {
        option.remove();
      }
    });
    e.target.parentElement.remove();
  }
}

customElements.define("requirement-selector", RequirementSelector);
