import {Dispatch, SetStateAction} from "react";

interface Member {
    id: string
    first_name: string
    last_name: string
    rank: string
}

interface MemberContext {
    Member: Member | null,
    SetMember: Dispatch<SetStateAction<Member | null>>
}