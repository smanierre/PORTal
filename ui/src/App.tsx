import { BrowserRouter, Route, Routes } from "react-router-dom";
import Login from "./pages/Login.tsx";
import Layout from "./Layout.tsx";
import Nav from "./components/Nav.tsx";

import Dashboard from "./pages/Dashboard.tsx";
import Qualifications from "./pages/Qualifications.tsx";
import { createContext, useContext, useState } from "react";
import { Member, Qualification } from "./index";
import Admin from "./pages/Admin.tsx"
import Redirect from "./components/Redirect.tsx";


export interface AppCtx {
  member: Member | null
  qualifications: Qualification[]
  subordinates: Member[]
}
export const AppContext = createContext<AppCtx>(Object.create({ member: null }))

export default function App() {
  const [appCtx, setAppCtx] = useState(useContext(AppContext))
  return (
    <BrowserRouter>
      <AppContext.Provider value={appCtx}>
        <Redirect appCtx={appCtx} setAppCtx={setAppCtx}>
          <Layout>
            <Nav loggedIn={appCtx.member !== null} admin={appCtx.member ? appCtx.member.admin ? true : false : false} setContext={setAppCtx} />
            <Routes>
              <Route path="/login" element={<Login setContext={setAppCtx} />} />
              <Route path="/dashboard" element={<Dashboard member={appCtx.member} qualifications={appCtx.qualifications} subordinates={appCtx.subordinates} />} />
              <Route path="/qualifications" element={<Qualifications />} />
              <Route path="/admin" element={<Admin />} />
            </Routes>
          </Layout>
        </Redirect>
      </AppContext.Provider>
    </BrowserRouter>
  )
}