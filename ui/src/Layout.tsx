import react from "@vitejs/plugin-react-swc";
import {JSX} from "react";

export default function Layout(props) {
    return (
        <div className={"h-screen w-screen grid grid-cols-1 grid-rows-main"}>
            {props.Nav}
            {props.children}
        </div>
    )
}