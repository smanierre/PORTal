import { Reference } from "../../../redux/api"
import { useState } from "react";
import { Button } from "../../ui/button";
import Search from "../../generic/Search";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
    adminReferenceSelector,
    newReference,
    selectReference,
} from "../../../redux/adminReferenceSlice";

interface AdminReferenceListProps {
    references: Reference[];
}
export default function AdminReferenceList({ references }: AdminReferenceListProps) {
    const [searchTerm, setSearchTerm] = useState("");
    const { selectedReference } = useAppSelector(adminReferenceSelector);
    const dispatch = useAppDispatch();

    return (
        <div className="flex flex-col items-center align-middle">
            <div className="w-1/2 ml-auto">
                <Search
                    className="w-2/3 inline-block border-b-0"
                    searchTerm={searchTerm}
                    setSearchTerm={setSearchTerm}
                />
                <Button
                    className="inline-block bg-background hover:bg-background-dark text-white w-1/3"
                    onClick={() => {
                        setSearchTerm("");
                    }}
                >
                    Clear
                </Button>
            </div>
            <ul className="overflow-x-scroll ml-auto h-2/3 w-1/2 border-black border p-2">
                {references
                    .filter((reference) => {
                        return (
                            reference.name
                                .toLowerCase()
                                .includes(searchTerm.toLowerCase())
                        );
                    })
                    .sort((r1, r2) =>
                        r1.name.toLowerCase()[0] > r2.name[0].toLowerCase()[0]
                            ? 1
                            : -1,
                    )
                    .map((reference) => (
                        <ul
                            key={reference.id}
                            className={`${selectedReference?.id === reference.id ? "bg-background text-white" : ""} cursor-pointer`}
                            onClick={() => {
                                dispatch(selectReference({ reference }));
                            }}
                        >
                            {reference.name}
                        </ul>
                    ))}
            </ul>
            <Button
                className="bg-background text-white hover:bg-background-dark ml-auto mt-2"
                onClick={() => {
                    dispatch(newReference());
                }}
            >
                Add Reference
            </Button>
        </div>
    );
}
