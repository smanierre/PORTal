import React, { SetStateAction } from "react";
import { Input } from "../ui/input";

interface SearchProps {
    searchTerm: string
    setSearchTerm: React.Dispatch<SetStateAction<string>>
    className?: string
}

export default function Search({ searchTerm, setSearchTerm, className }: SearchProps) {
    return (
        <Input className={`${className} bg-primary ml-auto`} placeholder="Search..." value={searchTerm} onChange={e => { setSearchTerm(e.target.value) }} />
    )
}