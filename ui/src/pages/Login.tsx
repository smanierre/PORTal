import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LoginRes } from "../index";
import { AppCtx } from "../App.tsx";
import { Button } from "../components/ui/button.tsx"
import { LoadingSpinner } from "../components/ui/spinner.tsx";
import { getBaseUrl } from "../lib/utils.ts";

interface LoginProps {
    setContext: React.Dispatch<React.SetStateAction<AppCtx>>
}

const errorHiddenClass = "h-0 w-4/5 p-0 mx-auto duration-500 transition-all text-xs"
const errorShownClass = "w-4/5 bg-red-600 p-6 mt-6 mx-auto duration-500 transition-all text-base"

export default function Login({ setContext }: LoginProps) {
    const nav = useNavigate()
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [showError, setShowError] = useState(false)
    const [waiting, setWaiting] = useState(false)
    const [errorText, setErrorText] = useState("")

    async function login(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        setWaiting(true)
        setShowError(false)
        setErrorText("")
        try {
            const res = await fetch(`${getBaseUrl()}/api/login`,
                {
                    body: JSON.stringify({ username: username, password: password }),
                    credentials: "include",
                    method: "POST"
                })
            if (res.ok) {
                const loginData = await res.json() as LoginRes
                localStorage.setItem("data", JSON.stringify(loginData))
                setContext({ ...loginData });
                nav("/dashboard")
            }
            switch (res.status) {
                case 401:
                    setErrorText("Invalid Credentials")
                    break
                case 500:
                    setErrorText("Server error")
                    break
                default:
                    setErrorText("Unexpected error")
            }
        } catch (err) {
            setErrorText("Unexpected error")
        }
        setShowError(true)
        setWaiting(false)
    }
    return (
        <div className="flex flex-row w-full h-full justify-center items-center">
            <div className="w-96 h-fit pb-8 bg-background">
                <p className="text-center text-xl font-bold mt-4 text-white">103d LRS PORTal Login</p>
                <form autoComplete="on" onSubmit={e => login(e)}>
                    <input name="username" disabled={waiting} className="w-4/5 h-10 block mx-auto mt-10" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
                    <input name="password" disabled={waiting} className="w-4/5 h-10 block mx-auto mt-5" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
                    <Button className="w-4/5 mx-auto mt-8 block" type="submit" disabled={waiting}>{waiting ? <LoadingSpinner className="h-6 w-6 inline-block" /> : "Login"}</Button>
                </form>
                <div className={showError ? errorShownClass : errorHiddenClass}>{errorText}</div>
            </div>
        </div>
    )
}