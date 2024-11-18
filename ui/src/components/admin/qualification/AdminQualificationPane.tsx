import { useGetAllQualificationsQuery, Qualification } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import AdminQualificationEditor from "./AdminQualificationEditor";
import SearchableList from "../../generic/SearchableList";
import { newQualification, selectQualification } from "../../../redux/adminQualificationSlice";
import { useAppDispatch } from "../../../redux/hooks";

export default function AdminQualificationPane() {
    const { data: qualifications, isLoading, isError } = useGetAllQualificationsQuery()
    const dispatch = useAppDispatch()
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
                            items={qualifications as Qualification[]}
                            displayFunc={qual => qual.name}
                            filterFunc={qualificationFilter}
                            sortFunc={(q1, q2) => q1.name.toLowerCase()[0] > q2.name.toLowerCase()[0] ? 1 : -1}
                            newItemFunc={newQualification}
                            selectItem={qualification => {
                                dispatch(selectQualification({ qualification }))
                            }}
                            typeName="Qualification"
                        />
                        <AdminQualificationEditor />
                    </div>
                )
    )
}

function qualificationFilter(searchTerm: string): (item: Qualification) => boolean {
    return item => item.name.toLowerCase().includes(searchTerm)
}