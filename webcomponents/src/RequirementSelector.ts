import { RequirementEditor } from "./RequirementEditor";

export class RequirementSelector extends HTMLElement {
    static #itemElement: HTMLButtonElement | null = null;
    static #editorElement: Element | null = null;
    #newRequirementCounter = 0;

    dialog: HTMLDialogElement
    addButton: HTMLButtonElement
    requirementList: HTMLSpanElement

    constructor() {
        super();

        // Get dialog for requirement editor and give it a unique ID for HTMX to target
        const dialog = this.querySelector("dialog");
        if (dialog instanceof HTMLDialogElement) {
            this.dialog = dialog
        } else {
            console.error("There should be a dialog element in this component!")
            this.dialog = new HTMLDialogElement();
        }
        const randLetter = String.fromCharCode(65 + Math.floor(Math.random() * 26));
        this.dialog.id = randLetter + Date.now();

        // Set the add button for setup during connectedCallback
        const addButton = this.querySelector("#add-button");
        if (addButton instanceof HTMLButtonElement) {
            this.addButton = addButton
        } else {
            console.error("There should be a button with the id 'add-button' in this component!")
            this.addButton = new HTMLButtonElement()
        }


        // Set the element that the list of requirements will be appended to
        const requirementList = this.querySelector(".requirement-items span");
        if (requirementList instanceof HTMLSpanElement) {
            this.requirementList = requirementList
        } else {
            console.error("There should be a span inside an element with class 'requirement-items' in this component!")
            this.requirementList = new HTMLSpanElement()
        }

    }

    get requirements(): Requirement[] {
        const reqText = this.getAttribute("requirements")
        return reqText ? JSON.parse(reqText) : []
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
            return new Promise<Node>((resolve, reject) => {
                fetch("/components/requirementItem", {
                    method: "GET",
                }).then((res) => {
                    if (res.ok) {
                        res.text().then((html) => {
                            const parser = new DOMParser();
                            const elem = parser.parseFromString(html, "text/html")
                            this.#itemElement = elem.querySelector("button");
                            if (this.#itemElement === null) {
                                reject("expected element wasn't retrieved from the server.")
                                return
                            }
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
            return new Promise<Node>((resolve, reject) => {
                fetch("/components/requirementEditor", {
                    method: "GET",
                }).then(res => {
                    if (res.ok) {
                        res.text().then(html => {
                            const parser = new DOMParser();
                            const elem = parser.parseFromString(html, "text/html")
                            this.#editorElement = elem.querySelector("requirement-editor")
                            if (this.#editorElement === null) {
                                reject("expected element wasn't retrieved from the server.")
                                return
                            }
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
                if (req instanceof HTMLButtonElement) {
                    const nameContainer = req?.querySelector(".name");
                    nameContainer!.textContent = requirement.name;
                    req?.setAttribute("data-id", requirement.id);
                    req?.addEventListener("click", this.updateRequirement.bind(this));
                    const cancelIcon = req.querySelector(
                        ".remove-requirement-button",
                    );
                    cancelIcon?.addEventListener("click", this.removeRequirement.bind(this));
                }
                this.requirementList.appendChild(req);
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
            cancelIcon?.removeEventListener(
                "click",
                this.removeRequirement.bind(this),
            );
        });
    }

    async newRequirement() {
        const editor = await RequirementSelector.editor;
        if (editor instanceof HTMLElement) {
            editor.setAttribute("data-parent-qualification-id", this.dataset.parentQualificationId || "")
        }
        this.dialog.innerHTML = '';
        this.dialog.appendChild(editor);
        this.dialog.showModal();

    }

    async updateRequirement(e: Event) {
        if (e === null || e.target === null) {
            return
        }
        let target: HTMLElement | null = e.target as HTMLElement;
        // If the inner span was clicked, set the new target to be the parent with the dataset
        if (e.target instanceof Element && e.target.tagName === "SPAN") {
            target = e.target.parentElement;
        }
        const editor = await RequirementSelector.editor as RequirementEditor;
        if (editor instanceof HTMLElement) {
            editor.setAttribute("data-parent-qualification-id", this.dataset.parentQualificationId || "")
        }
        this.dialog.innerHTML = '';
        this.dialog.appendChild(editor);
        const req = this.requirements.filter(req => req.id === target?.dataset.id)
        editor.setFormData(req[0])
        this.dialog.showModal();
    }

    async setAddedOrUpdatedQualificationRequirement(requirement: Requirement, display: string) {
        let currentReqs = this.requirements;
        if (requirement.id === "") {
            // This is a new requirement, give it a dummy ID so that it can be found when clicked on, add it to the attribute and create a display item for it.
            requirement.id = `${this.newRequirementCounter}`
            currentReqs.push(requirement);
            const req = await RequirementSelector.item;
            if (req instanceof HTMLElement) {
                const nameContainer = req.querySelector(".name");
                nameContainer!.textContent = display;
                req.setAttribute("data-id", requirement.id);
                req.addEventListener("click", this.updateRequirement.bind(this));
            }
            this.requirementList.appendChild(req);
        } else {
            // Existing requirement, update the attribute holding them and the name in the display item
            currentReqs = currentReqs.map(req => {
                if (req.id === requirement.id) {
                    req = requirement;
                }
                return req
            })
            const reqItems = this.requirementList.querySelectorAll("button");
            reqItems.forEach(item => {
                if (item.dataset.id === requirement.id) {
                    const nameSpan = item.querySelector(".name")
                    nameSpan!.textContent = requirement.name;
                }
            })
        }
        this.requirements = currentReqs;
        this.dispatchEvent(new Event("updated"))
    }

    async setAddedOrUpdatedRequirement(requirement: Requirement) {
        let currentReqs = this.requirements;
        if (requirement.id === "") {
            // This is a new requirement, give it a dummy ID so that it can be found when clicked on, add it to the attribute and create a display item for it.
            requirement.id = `${this.newRequirementCounter}`
            currentReqs.push(requirement);
            const req = await RequirementSelector.item;
            if (req instanceof HTMLElement) {
                const nameContainer = req.querySelector(".name");
                nameContainer!.textContent = requirement.name;
                req.setAttribute("data-id", requirement.id);
                req.addEventListener("click", this.updateRequirement.bind(this));
            }
            this.requirementList.appendChild(req);
        } else {
            // Existing requirement, update the attribute holding them and the name in the display item
            currentReqs = currentReqs.map(req => {
                if (req.id === requirement.id) {
                    req = requirement;
                }
                return req
            })
            const reqItems = this.requirementList.querySelectorAll("button");
            reqItems.forEach(item => {
                if (item.dataset.id === requirement.id) {
                    const nameSpan = item.querySelector(".name")
                    nameSpan!.textContent = requirement.name;
                }
            })
        }
        this.requirements = currentReqs;
        this.dispatchEvent(new Event("updated"))
    }

    async removeRequirement(e: Event) {
        e.stopPropagation()
        if (e.target === null) {
            return
        }
        const target = e.target as HTMLElement;
        target.parentElement?.remove();
        const remainingRequirements = this.requirements.filter(requirement => {
            target.parentElement?.dataset.id !== requirement.id
        })
        this.setAttribute("requirements", JSON.stringify(remainingRequirements))
    }
}