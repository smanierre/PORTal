import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../components/ui/button.tsx";
import { LoadingSpinner } from "../components/ui/spinner.tsx";
import { useGetLoggedInMemberQuery, useLoginMutation } from "../redux/api.ts";
import FullPageSpinner from "../components/FullPageSpinner.tsx";

const errorHiddenClass =
  "h-0 w-4/5 p-0 mx-auto duration-500 transition-all text-xs";
const errorShownClass =
  "w-4/5 bg-red-600 p-6 mt-6 mx-auto duration-500 transition-all text-base";

export default function Login() {
  const nav = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showError, setShowError] = useState(false);
  const [errorText, setErrorText] = useState("");
  const [loading, setLoading] = useState(false);
  const [loginMutation] = useLoginMutation();
  const { data: member, refetch, isLoading } = useGetLoggedInMemberQuery();

  useEffect(() => {
    if (member !== undefined) {
      nav("/dashboard");
      return;
    }
    if (window.location.pathname !== "/login") {
      window.history.replaceState(null, "Login Page", "/login");
    }
  }, [member]);

  async function login(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setShowError(false);
    setErrorText("");
    const { error } = await loginMutation({
      username: username,
      password: password,
    });
    if (error && "status" in error) {
      switch (error.status) {
        case 401:
          setErrorText("Invalid Credentials");
          break;
        case 500:
          setErrorText("Server error");
          break;
        default:
          setErrorText("Unexpected error");
      }
      setShowError(true);
      setLoading(false);
      return;
    }
    await refetch();
    nav("/dashboard");
  }
  if (isLoading) {
    return <FullPageSpinner />;
  }
  return (
    <div className="flex flex-row w-full h-full justify-center items-center">
      <div className="w-96 h-fit pb-8 bg-background">
        <p className="text-center text-xl font-bold mt-4 text-white">
          103d LRS PORTal Login
        </p>
        <form autoComplete="on" onSubmit={(e) => login(e)}>
          <input
            name="username"
            disabled={loading}
            className="w-4/5 h-10 block mx-auto mt-10"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            name="password"
            disabled={loading}
            className="w-4/5 h-10 block mx-auto mt-5"
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button
            className="w-4/5 mx-auto mt-8 block whitespace-nowrap"
            type="submit"
            disabled={loading}
          >
            {loading ? (
              <LoadingSpinner className="h-6 w-6 inline-block" />
            ) : (
              "Login"
            )}
          </Button>
        </form>
        <div className={showError ? errorShownClass : errorHiddenClass}>
          {errorText}
        </div>
      </div>
    </div>
  );
}
