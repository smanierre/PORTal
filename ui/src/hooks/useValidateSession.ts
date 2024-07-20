import {Dispatch, SetStateAction, useEffect, useState} from "react";
import {Member} from "../index";

export default function useValidateSession() {
    const [member, setMember] = useState<Member | null>(null)
    useEffect(() => {
        async function validateSession(m: Member) {
            const res = await fetch("http://localhost:8080/api/validateSession",
                {
                    body: `{"id":"${member.id}"}`,
                    method: "POST",
                    credentials: "include"
                })
            if (!res.ok) {
                localStorage.removeItem("member")
                return null
            } else {
                setMember(m)
            }
        }

        const memberString = localStorage.getItem("member") as string | null
        if (memberString === null) {
            return undefined
        }
        const member = JSON.parse(memberString) as Member
        validateSession(member)
    }, []);
    return [member, setMember]
}