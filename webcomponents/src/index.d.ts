interface Requirement {
    id: string
    name: string
    notes: string
    type: "Qualification" | "Grade" | "WBT"
    days_valid_for: number
    reference: string
    qualification_id: string
}

interface HTMXEvent extends Event {
    detail: {
        xhr: {
            status: number
        },
        shouldSwap: boolean,
        isError: boolean
    }
}

interface adminMemberListStore {
    selected: HTMLLIElement | null
    selectHandler: (e: Event) => void
    clearSelectedMember: () => void
}

interface adminQualificationListStore {
    selected: HTMLLIElement | null
    selectHandler: (e: Event) => void
    clearSelectedQualification: () => void
}