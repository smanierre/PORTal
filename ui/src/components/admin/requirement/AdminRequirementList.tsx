import { Requirement } from "../../../redux/api";
import { useState } from "react";
import { Button } from "../../ui/button";
import Search from "../../generic/Search";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
    adminRequirementSelector,
    newRequirement,
    selectRequirement,

} from "../../../redux/adminRequirementSlice";

interface AdminRequirementListProps {
    requirements: Requirement[]
}
export default function AdminRequirementList({ requirements }: AdminRequirementListProps) {
    const [searchTerm, setSearchTerm] = useState("");
    const { selectedRequirement } = useAppSelector(adminRequirementSelector);
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
                {requirements !== null ? requirements
                    .filter((requirement) => {
                        return (
                            requirement.name
                                .toLowerCase()
                                .includes(searchTerm.toLowerCase())
                        );
                    })
                    .sort((q1, q2) =>
                        q1.name.toLowerCase()[0] > q2.name.toLowerCase()[0]
                            ? 1
                            : -1,
                    )
                    .map((requirement) => (
                        <ul
                            key={requirement.id}
                            className={`${selectedRequirement?.id === requirement.id ? "bg-background text-white" : ""} cursor-pointer`}
                            onClick={() => {
                                dispatch(selectRequirement({ requirement }));
                            }}
                        >
                            {requirement.name}
                        </ul>
                    )) : null}
            </ul>
            <Button
                className="bg-background text-white hover:bg-background-dark ml-auto mt-2"
                onClick={() => {
                    dispatch(newRequirement());
                }}
            >
                Add Requirement
            </Button>
        </div>
    );
}
