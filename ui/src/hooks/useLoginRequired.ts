import { useNavigate } from "react-router-dom";
import { AppContext } from "../App";
import { useContext, useEffect } from "react";

export default function useLoginRequired() {
    const memberCtx = useContext(AppContext)
    const nav = useNavigate()
    if (memberCtx.member !== null) {
        return
    }
    useEffect(() => {
        nav("/login")
    }, [])
}