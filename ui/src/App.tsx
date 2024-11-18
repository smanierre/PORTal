import { Route, Routes, useLocation } from "react-router-dom";
import Login from "./pages/Login.tsx";
import Layout from "./Layout.tsx";
import Nav from "./components/Nav.tsx";

import Dashboard from "./pages/Dashboard.tsx";
import Admin from "./pages/Admin.tsx";
import { useGetLoggedInMemberQuery } from "./redux/api.ts";
import FullPageSpinner from "./components/FullPageSpinner.tsx";

export default function App() {
  const { data: member, isLoading } = useGetLoggedInMemberQuery();
  const location = useLocation();
  return isLoading ? (
    <FullPageSpinner />
  ) : (
    <Layout>
      <Nav
        loggedIn={!!member}
        admin={member ? (member.admin ? true : false) : false}
      />
      <Routes location={member === undefined ? "/login" : location.pathname}>
        <Route path="/login" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
        {/*<Route path="/qualifications" element={<Qualifications />} />*/}
        <Route path="/admin" element={<Admin />} />
      </Routes>
    </Layout>
  );
}
