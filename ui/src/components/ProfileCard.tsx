import {cn} from "../lib/utils.ts";
import {useContext} from "react";
import {MemberCtx} from "../App.tsx";

export default function ProfileCard({className}: {className?:string}) {
    const ctx = useContext(MemberCtx)
    return <article className={`${className !== undefined ? className : ""} w-1/2 h-1/2 bg-red-500`}>
        <p>First Name: {ctx.Member?.first_name}</p>
        <p>Last Name: {ctx.Member?.last_name}</p>
        <p>Rank: {ctx.Member?.rank}</p>
        <p>Qualifications: </p>
    </article>
}