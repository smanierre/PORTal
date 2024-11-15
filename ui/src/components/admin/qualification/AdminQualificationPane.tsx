import { useGetAllQualificationsQuery, Qualification } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import AdminQualificationList from "./AdminQualificationList";
import AdminQualificationEditor from "./AdminQualificationEditor";

export default function AdminQualificationPane() {
    const { data: qualifications, isLoading, isError } = useGetAllQualificationsQuery()
    return (
        isLoading ?
            <FullPageSpinner />
            :
            isError ?
                <div>error</div>
                :
                (
                    <div className="grid grid-cols-adminPane">
                        <AdminQualificationList qualifications={qualifications as Qualification[]} />
                        <AdminQualificationEditor />
                    </div>
                )
    )
}