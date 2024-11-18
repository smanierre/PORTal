import { useGetAllReferencesQuery, Reference } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import AdminReferenceEditor from "./AdminReferenceEditor";
import SearchableList from "../../generic/SearchableList";
import { newReference, selectReference } from "../../../redux/adminReferenceSlice";
import { useAppDispatch } from "../../../redux/hooks";

export default function AdminReferencePane() {
    const { data: references, isLoading, isError } = useGetAllReferencesQuery()
    const dispatch = useAppDispatch();
    return (
        isLoading ?
            <FullPageSpinner />
            :
            isError ?
                <div>error</div>
                :
                (
                    <div className="grid grid-cols-adminPane">
                        <SearchableList
                            items={references as Reference[]}
                            displayFunc={ref => ref.name}
                            filterFunc={referenceFilter}
                            sortFunc={(ref1, ref2) => ref1.name.toLowerCase()[0] > ref2.name[0].toLowerCase()[0] ? 1 : -1}
                            newItemFunc={newReference}
                            selectItem={reference => { dispatch(selectReference({ reference })) }}
                            typeName="Reference"
                        />
                        <AdminReferenceEditor />
                    </div>
                )
    )
}

function referenceFilter(searchTerm: string): (ref: Reference) => boolean {
    return (ref: Reference) => ref.name.includes(searchTerm)
}