import React from "react"

export default function AdminLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="h-full w-full grid grid-rows-admin">
            {children}
        </div>
    )
}