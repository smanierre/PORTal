class MultiSelect extends HTMLElement {
  constructor() {
    super();

    // Set selected class on all the options that are selected
    this.options.forEach((option) => {
      if (option.selected) {
        option.classList.add("selected");
      }
    });
    this.addEventListener("mousedown", (e) => {
      e.preventDefault();
      if (this.options.indexOf(e.target) === -1) {
        return;
      }
      if (e.target.selected) {
        e.target.classList.remove("selected");
        e.target.selected = false;
      } else {
        e.target.classList.add("selected");
        e.target.selected = true;
      }
      this.dispatchEvent(new Event("change", { bubbles: true }));
    });
  }

  get options() {
    return Array.from(this.querySelectorAll("option"));
  }
}
customElements.define("multi-select", MultiSelect);
