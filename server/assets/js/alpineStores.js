document.addEventListener("alpine:init", () => {
    Alpine.data("dropdown", () => ({
        open: false,
        selected: "Members",
        toggle() {
            this.open = ! this.open
        },
        selectHandler(selectedItem) {
            this.selected = selectedItem
            this.open = false
        }
    }));

    Alpine.data("adminMemberList", () => ({
        selected: null,
        selectHandler(e) {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedMember() {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
                this.selected = null;
            }
        }
    }));

    Alpine.data("adminMemberEditor", () => ({
      gradeOnChange(e, memberId) {
          htmx.ajax("GET", `/admin/members/${memberId}/potentialSupervisors/${e.target.value}`, {target: "#supervisor", swap: "outerHTML"})
      }
    }));

    Alpine.data("adminReferenceList", () => ({
        selected: null,
        selectHandler(e) {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedReference() {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
                this.selected = null;
            }
        }
    }));

    Alpine.data("adminRequirementList", () => ({
        selected: null,
        selectHandler(e) {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedRequirement() {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
                this.selected = null;
            }
        }
    }));

    Alpine.data("adminQualificationList", () => ({
        selected: null,
        selectHandler(e) {
            if(this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedQualification() {
            if (this.selected) {
                this.selected.classList.remove("searchable-list-item--selected");
                this.selected = null;
            }
        }
    }));

    Alpine.data("adminQualificationEditor", () => ({}));
})

