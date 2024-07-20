import {useEffect} from "react";
import useValidateSession from "../hooks/useValidateSession.ts";
import useLoginRequired from "../hooks/useLoginRequired.ts";

export default function Root() {
    useLoginRequired()
    return <h1>Root</h1>
}