import { Dispatch, SetStateAction } from "react";

interface LoginRes {
    member: Member
    qualifications: Qualification[]
    subordinates: Member[]
}

interface Member {
    id: string
    first_name: string
    last_name: string
    rank: string
}

interface Qualification {
    id: string
    name: string
    initial_requirements: Requirement[]
    recurring_requirements: Requirement[]
    notes: string
    expires: bool
    expiration_days: number
}

interface Requirement {
    id: string
    name: string
    reference: Reference
    description: string
    notes: string
    days_valid_for: number
}

interface Reference {
    id: string
    name: string
    volume: number
    paragraph: string
}

interface MemberContext {
    Member: Member | null,
    SetMember: Dispatch<SetStateAction<Member | null>>
}