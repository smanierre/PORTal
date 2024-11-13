import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getBaseUrl } from "../lib/utils";

export default function useAdminRequired() {
  const [waiting, setWaiting] = useState(true);
  const nav = useNavigate();
  useEffect(() => {
    const checkPrivleges = async () => {
      try {
        const res = await fetch(`${getBaseUrl()}/api/checkAdmin`, {
          credentials: import.meta.env.DEV ? "include" : "same-origin",
        });
        if (res.ok) {
          setWaiting(false);
          return;
        }
      } catch (err) {
        nav("/dashboard");
      }
    };
    checkPrivleges();
  }, []);
  return waiting;
}
