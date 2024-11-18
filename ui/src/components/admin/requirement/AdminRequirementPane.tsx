import { Requirement, useGetAllRequirementsQuery } from "../../../redux/api"
import FullPageSpinner from "../../FullPageSpinner";
import AdminRequirementEditor from "./AdminRequirementEditor";
import SearchableList from "../../generic/SearchableList";
import { newRequirement, selectRequirement } from "../../../redux/adminRequirementSlice";
import { useAppDispatch } from "../../../redux/hooks";

export default function AdminRequirementPane() {
    const { data: requirements, isLoading, isError } = useGetAllRequirementsQuery();
    const dispatch = useAppDispatch();
    return (
        isLoading ?
            <FullPageSpinner />
            :
            isError ?
                <div>error</div> :
                <div className="grid grid-cols-adminPane">
                    <SearchableList
                        items={requirements as Requirement[]}
                        displayFunc={req => req.name}
                        newItemFunc={newRequirement}
                        selectItem={requirement => dispatch(selectRequirement({ requirement }))}
                        sortFunc={(req1, req2) => req1.name.toLowerCase()[0] > req2.name.toLowerCase()[0] ? 1 : -1}
                        filterFunc={requirementFilter}
                        typeName="Requirement"
                    />
                    <AdminRequirementEditor />
                </div>
    )
}

function requirementFilter(searchTerm: string): (req: Requirement) => boolean {
    return (req: Requirement): boolean => {
        return req.name.toLowerCase().includes(searchTerm)
    }
}