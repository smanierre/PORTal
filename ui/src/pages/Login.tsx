import { useContext, useEffect, useState } from "react";
import { Button } from "../components/ui/button.tsx";
import { useNavigate } from "react-router-dom";
import { MemberCtx } from "../App.tsx";
import useValidateSession from "../hooks/useValidateSession.ts";
import { Input } from "../components/ui/input.tsx";

export default function Login() {
    const nav = useNavigate()
    const memberContext = useContext(MemberCtx)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")

    useEffect(() => {
        if (memberContext.Member !== null) {
            nav("/dashboard")
        }
    }, []);

    async function login() {
        const res = await fetch("http://localhost:8080/api/login",
            {
                body: JSON.stringify({ username: username, password: password }),
                credentials: "include",
                method: "POST"
            })
        if (res.ok) {
            const loginData = await res.json()
            localStorage.setItem("member", JSON.stringify(member))
            memberContext.SetMember(member)
            nav("/dashboard")
        }
    }
    return (
        <div className="flex flex-row w-full h-full justify-center items-center">
            <div className="w-96 h-80 bg-slate-100 rounded-2xl">
                <p className="text-center text-xl font-bold">Login</p>
                <input className="w-30 h-10 m-5" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
                <br />
                <input className="w-30 h-10 m-5" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
                <button onClick={login}>Login</button>
            </div>
        </div>
    )
}