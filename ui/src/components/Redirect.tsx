import React, { useEffect, useState } from "react";
import { AppCtx } from "../App";
import { useLocation, useNavigate } from "react-router-dom";
import FullPageSpinner from "./FullPageSpinner";

interface RedirectProps {
    children: React.ReactNode
    appCtx: AppCtx
    setAppCtx: React.Dispatch<React.SetStateAction<AppCtx>>
}

export default function Redirect({ children, setAppCtx }: RedirectProps) {
    const nav = useNavigate()
    const location = useLocation()
    const [waiting, setWaiting] = useState(true)
    useEffect(() => {
        const data = localStorage.getItem("data")
        if (data === null) {
            setAppCtx({ member: null, qualifications: [], subordinates: [] })
            nav("/login")
        } else {
            const parsedData = JSON.parse(data) as AppCtx
            setAppCtx({ ...parsedData })
            if (location.pathname === "/login" || location.pathname === "/") {
                nav("/dashboard")
            }
        }
        setWaiting(false)
    }, [])
    return (
        waiting ?
            <FullPageSpinner /> :
            <>
                {children}
            </>
    )
}