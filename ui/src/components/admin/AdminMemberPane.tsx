import { useState, useEffect } from "react";
import AdminMemberEditor from "./AdminMemberEditor";
import AdminMemberList from "./AdminMemberList";
import { Member } from "../..";
import { getEmptyMember, getBaseUrl } from "../../lib/utils";


export default function AdminMemberPane() {
    const [selectedMember, setSelectedMember] = useState<Member>(getEmptyMember)
    const [newMember, setNewMember] = useState(true)
    const [addedMember, setAddedMember] = useState(0)
    const [members, setMembers] = useState<Member[]>([])

    useEffect(() => {
        const fetchMembers = async () => {
            const res = await fetch(`${getBaseUrl()}/api/members`, {
                credentials: "include",
            })
            if (!res.ok) {
                console.error("error getting members")
            }
            const membersJson = await res.json()
            setMembers(membersJson)
        }
        fetchMembers()
    }, [addedMember])
    return (
        <div className="grid grid-cols-adminPane">
            <AdminMemberList members={members} setSelectedMember={setSelectedMember} selectedMember={selectedMember} setNewMember={setNewMember} addedMember={addedMember} />
            <AdminMemberEditor members={members} setMembers={setMembers} setNewMember={setNewMember} selectedMember={selectedMember} newMember={newMember} setAddedMember={setAddedMember} addedMember={addedMember} setSelectedMember={setSelectedMember} />
        </div>
    )
}