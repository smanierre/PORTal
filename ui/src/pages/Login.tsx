import {useContext, useEffect, useState} from "react";
import {Button} from "../components/ui/button.tsx";
import {useNavigate} from "react-router-dom";
import {MemberCtx} from "../App.tsx";
import useValidateSession from "../hooks/useValidateSession.ts";

export default function Login() {
    const nav = useNavigate()
    const memberContext = useContext(MemberCtx)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")

    useEffect(() => {
        if(memberContext.Member !== null) {
            const nav = useNavigate()
            nav("/dashboard")
        }
    }, []);

    async function login() {
        const res = await fetch("http://localhost:8080/api/login",
            {
                body: JSON.stringify({username: username, password: password}),
                credentials: "include",
                method: "POST"
            })
        if(res.ok) {
            const member = await res.json()
            localStorage.setItem("member", JSON.stringify(member))
            memberContext.SetMember(member)
            nav("/dashboard")
        }
    }
    return (
        <>
            <input value={username} onChange={e => {
                setUsername(e.target.value)
            }}/>
            <input value={password} onChange={e => {
                setPassword(e.target.value)
            }}/>
            <Button onClick={login}>Login</Button>
        </>
    )
}