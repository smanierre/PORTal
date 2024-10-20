import AdminLayout from "../components/layouts/AdminLayout"
import FullPageSpinner from "../components/FullPageSpinner"
import useAdminRequired from "../hooks/useAdminRequired"
import AdminMemberPane from "../components/admin/AdminMemberPane"
import AdminQualificationPane from "../components/admin/AdminQualificationPane"
import Selector from "../components/generic/Selector"
import { useState } from "react"

const options = [
    {
        value: "members",
        label: "Members",
    },
    {
        value: "qualifications",
        label: "Qualifications",
    }
]

export default function Admin() {
    const [selectorValue, setSelectorValue] = useState("members")
    const waiting = useAdminRequired()
    return (
        waiting ?
            <FullPageSpinner /> :
            <AdminLayout>
                <Selector value={selectorValue} setValue={setSelectorValue} options={options} />
                {selectorValue === "members" ?
                    <AdminMemberPane /> :
                    selectorValue === "qualifications" ?
                        <AdminQualificationPane /> : null
                }
            </AdminLayout>
    )
}