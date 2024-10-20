import React from "react"

interface LayoutProps {
    children: React.JSX.Element[]
}

export default function Layout(props: LayoutProps) {
    return (
        <div className={"h-screen w-screen grid grid-cols-1 grid-rows-main"}>
            {props.children}
        </div>
    )
}