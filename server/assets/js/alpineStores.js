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
            if(this.selected !== null) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedMember() {
            this.selected.classList.remove("searchable-list-item--selected");
            this.selected = null;
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
            if(this.selected !== null) {
                this.selected.classList.remove("searchable-list-item--selected");
            }
            this.selected = e.target;
            this.selected.classList.add("searchable-list-item--selected");
        },
        clearSelectedReference() {
            this.selected.classList.remove("searchable-list-item--selected");
            this.selected = null;
        }
    }))
})

