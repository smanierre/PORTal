import {BrowserRouter, Route, Routes} from "react-router-dom";
import Root from "./pages/Root.tsx";
import Login from "./pages/Login.tsx";
import Layout from "./Layout.tsx";
import Nav from "./components/Nav.tsx";
import {createContext, useEffect, useState} from "react";
import useValidateSession from "./hooks/useValidateSession.ts";
import {MemberContext} from "./index";
import {Progress} from "./components/ui/progress.tsx";
import Dashboard from "./pages/Dashboard.tsx";

export const MemberCtx = createContext<MemberContext>({Member: null, SetMember: null})

export default function App() {
    const [member, setMember] = useValidateSession()
    const [progress, setProgress] = useState(0)
    useEffect(() => {

        const interval = setInterval(() => setProgress((prev) => {
            return prev + Math.floor(Math.random() * 100) % 25
        }), 250)
        return () => {
            clearInterval(interval)
        }
    }, []);
    return (
        progress <= 100 ?
            <div className={"h-screen w-screen relative"}>
                <div className={"w-1/2 absolute top-1/2 left-1/2 -translate-x-1/2"}>
                    <Progress value={progress}/>
                </div>
            </div>
            :
            <BrowserRouter>
                <MemberCtx.Provider value={{Member: member, SetMember: setMember}}>
                    <Layout>
                        <Nav loggedIn={member !== null}/>
                        <Routes>
                            <Route path="/" element={<Root/>}/>
                            <Route path="/login" element={<Login/>}/>
                            <Route path="/dashboard" element={<Dashboard />}/>
                        </Routes>
                    </Layout>
                </MemberCtx.Provider>
            </BrowserRouter>
    )
}