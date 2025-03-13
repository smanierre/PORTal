import { RequirementEditor } from "./RequirementEditor";
import { RequirementSelector } from "./RequirementSelector";
import { QualificationEditor } from "./QualificationEditor";
import { Dropdown } from "./Dropdown";
import { SearchableList } from "./SearchableList";

customElements.define("requirement-editor", RequirementEditor);
customElements.define("requirement-selector", RequirementSelector);
customElements.define("qualification-editor", QualificationEditor);
customElements.define("dropdown-component", Dropdown);
customElements.define("searchable-list", SearchableList);

document.addEventListener("htmx:beforeOnLoad", (evt: Event) => {
  const e = evt as HTMXEvent;
  if (e.detail.xhr.status === 401 || e.detail.xhr.status === 500) {
    e.detail.shouldSwap = true;
    e.detail.isError = false;
  }
});
