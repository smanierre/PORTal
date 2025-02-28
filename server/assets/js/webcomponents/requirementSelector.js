class RequirementSelector extends HTMLElement {
  static #itemElement = null;
  static #editorElement = null;
  #newRequirementCounter = 0;

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
  }

  get requirements() {
    return JSON.parse(this.getAttribute("requirements"))
  }

  set requirements(requirements) {
    this.setAttribute("requirements", JSON.stringify(requirements))
  }

  get newRequirementCounter() {
    this.#newRequirementCounter++
    return this.#newRequirementCounter
  }

  // itemElement is an instance of the HTML returned when a requirement is added to a qualification.
  // For existing qualifications we need to get the itemElement from the server and set it here for instances to use.
  static get item() {
    if (this.#itemElement === null) {
      return new Promise((resolve, reject) => {
        fetch("/components/requirementItem", {
          method: "GET",
        }).then((res) => {
          if (res.ok) {
            res.text().then((html) => {
              const parser = new DOMParser();
              const elem = parser.parseFromString(html, "text/html")
              this.#itemElement = elem.querySelector("button");
              resolve(this.#itemElement.cloneNode(true)); // Resolve the promise with the item
            }).catch(reject); // Handle errors in the inner `text` promise
          } else {
            reject(new Error("Failed to load requirement item from server!")); // Reject if the fetch fails
          }
        }).catch(reject); // Handle errors in the fetch itself
      });
    } else {
      return Promise.resolve(this.#itemElement.cloneNode(true)); // Return a resolved promise with the cloned item
    }
  }

  static get editor() {
    if (this.#editorElement == null) {
      return new Promise((resolve, reject) => {
        fetch("/components/requirementEditor", {
          method: "GET",
        }).then(res => {
          if (res.ok) {
            res.text().then(html => {
              const parser = new DOMParser();
              const elem = parser.parseFromString(html, "text/html")
              this.#editorElement = elem.querySelector("requirement-editor")
              resolve(this.#editorElement.cloneNode(true));
            }).catch(reject);
          } else {
            reject(new Error("Failed to load requirement editor from server!"));
          }
        }).catch(reject);
      });
    } else {
      return Promise.resolve(this.#editorElement.cloneNode(true));
    }
  }

  connectedCallback() {
    // Add the event listener for the add button
    this.addButton.addEventListener("click", this.newRequirement.bind(this));

    // For each of the existing requirements, add event listeners and add an <option> to the select with the id as a value.
    this.requirements.forEach(
      async (requirement) => {
        const req = await RequirementSelector.item;
        const nameContainer = req.querySelector(".name");
        nameContainer.textContent = requirement.Name;
        req.setAttribute("data-id", requirement.ID);
        req.addEventListener("click", this.updateRequirement.bind(this));
        this.requirementList.appendChild(req);

        const cancelIcon = req.querySelector(
          ".remove-requirement-button",
        );
        cancelIcon.addEventListener("click", this.removeRequirement.bind(this));
      },
    );
  }

  disconnectedCallback() {
    // Remove the event listener for the add button
    this.addButton.removeEventListener("click", this.newRequirement);

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

  async newRequirement() {
    const editor = await RequirementSelector.editor;
    this.dialog.innerHTML = '';
    this.dialog.appendChild(editor);
    this.dialog.showModal();
  }

  async updateRequirement(e) {
    let target = e.target;
    // If the inner span was clicked, set the new target to be the parent with the dataset
    if (e.target.tagName === "SPAN") {
      target = e.target.parentNode;
    }
    const editor = await RequirementSelector.editor;
    this.dialog.innerHTML = '';
    this.dialog.appendChild(editor);
    const req = this.requirements.filter(req => req.ID === target.dataset.id)
    editor.setFormData(req[0])
    this.dialog.showModal();
  }

  async setAddedOrUpdatedRequirement(requirement) {
    let currentReqs = this.requirements;
    if (requirement.ID === undefined) {
      // This is a new requirement, give it a dummy ID so that it can be found when clicked on, add it to the attribute and create a display item for it.
      requirement.ID = `${this.newRequirementCounter}`
      currentReqs.push(requirement);
      const req = await RequirementSelector.item;
      const nameContainer = req.querySelector(".name");
      nameContainer.textContent = requirement.Name;
      req.setAttribute("data-id", requirement.ID);
      req.addEventListener("click", this.updateRequirement.bind(this));
      this.requirementList.appendChild(req);
    } else {
      // Existing requirement, update the attribute holding them and the name in the display item
      currentReqs = currentReqs.map(req => {
        if (req.ID === requirement.ID) {
          req = requirement;
        }
        return req
      })
      const reqItems = this.requirementList.querySelectorAll("button");
      reqItems.forEach(item => {
        if (item.dataset.id === requirement.ID) {
          const nameSpan = item.querySelector(".name")
          nameSpan.textContent = requirement.Name;
        }
      })
    }
    this.requirements = currentReqs;
  }

  async removeRequirement(e) {
    e.stopPropagation()
    e.target.parentElement.remove();
    const remainingRequirements = JSON.parse(this.getAttribute("requirements")).filter(requirement => {
      e.target.parentElement.dataset.id !== requirement.ID
    })
    this.setAttribute("requirements", JSON.stringify(remainingRequirements))
  }
}

customElements.define("requirement-selector", RequirementSelector);
