import React from "react";
import { useLocation, useNavigate } from "react-router-dom";
import FullPageSpinner from "./FullPageSpinner";
import { useGetLoggedInMemberQuery } from "../redux/api";

interface RedirectProps {
  children: React.ReactNode;
}

export default function Redirect({ children }: RedirectProps) {
  const nav = useNavigate();
  const location = useLocation();
  const { data, isLoading, error } = useGetLoggedInMemberQuery();
  function handleRedirect(children: React.ReactNode) {
    if (error !== null && location.pathname !== "/login") {
      nav("/login");
    } else if (location.pathname === "/login" && data != null) {
      nav("/dashboard");
    }
    return <>{children}</>;
  }
  return isLoading ? <FullPageSpinner /> : handleRedirect(children);
}
