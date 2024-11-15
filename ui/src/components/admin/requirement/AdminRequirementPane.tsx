import { Requirement, useGetAllRequirementsQuery } from "../../../redux/api"
import AdminRequirementList from "./AdminRequirementList";
import FullPageSpinner from "../../FullPageSpinner";
import AdminRequirementEditor from "./AdminRequirementEditor";

export default function AdminRequirementPane() {
    const { data: requirements, isLoading, isError } = useGetAllRequirementsQuery();
    return (
        isLoading ?
            <FullPageSpinner />
            :
            isError ?
                <div>error</div> :
                <div className="grid grid-cols-adminPane">
                    <AdminRequirementList requirements={requirements as Requirement[]} />
                    <AdminRequirementEditor />
                </div>


    )
}