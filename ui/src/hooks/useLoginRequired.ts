import {useLocation, useNavigate} from "react-router-dom";
import useValidateSession from "./useValidateSession.ts";
import {useContext} from "react";
import {MemberCtx} from "../App.tsx";

export default function useLoginRequired() {
    const ctx = useContext(MemberCtx)
    const currentRoute = useLocation()
    const nav = useNavigate()
    if(ctx.Member !== null) {
        if(currentRoute.pathname === "/") {
            nav("/dashboard")
        }
        return
    }
    const [member, setMember] = useValidateSession()
    if(member === null) {
        nav("/login")
    }
    ctx.SetMember(member)
}