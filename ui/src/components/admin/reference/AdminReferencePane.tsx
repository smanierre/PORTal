import { useGetAllReferencesQuery, Reference } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import AdminReferenceList from "./AdminReferenceList";
import AdminReferenceEditor from "./AdminReferenceEditor";

export default function AdminReferencePane() {
    const { data: references, isLoading, isError } = useGetAllReferencesQuery()
    return (
        isLoading ?
            <FullPageSpinner />
            :
            isError ?
                <div>error</div>
                :
                (
                    <div className="grid grid-cols-adminPane">
                        <AdminReferenceList references={references as Reference[]} />
                        <AdminReferenceEditor />
                    </div>
                )
    )
}