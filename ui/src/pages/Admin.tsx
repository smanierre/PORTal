import AdminLayout from "../components/layouts/AdminLayout"
import FullPageSpinner from "../components/FullPageSpinner"
import useAdminRequired from "../hooks/useAdminRequired"
import AdminMemberPane from "../components/admin/member/AdminMemberPane"
import AdminQualificationPane from "../components/admin/qualification/AdminQualificationPane"
import AdminRequirementPane from "../components/admin/requirement/AdminRequirementPane"
import AdminReferencePane from "../components/admin/reference/AdminReferencePane"
import Selector from "../components/generic/Selector"
import React, { useState } from "react"

const options = [
    {
        value: "members",
        label: "Members",
    },
    {
        value: "qualifications",
        label: "Qualifications",
    },
    {
        value: "requirements",
        label: "Requirements",
    },
    {
        value: "references",
        label: "References"
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
                {choosePane(selectorValue)}
            </AdminLayout>
    )
}

function choosePane(selectorValue: string): React.ReactNode {
    switch (selectorValue) {
        case "members":
            return <AdminMemberPane />
        case "qualifications":
            return <AdminQualificationPane />
        case "requirements":
            return <AdminRequirementPane />
        case "references":
            return <AdminReferencePane />
    }
}