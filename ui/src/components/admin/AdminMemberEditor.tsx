import React, { useEffect, useState } from "react"
import { Member } from "../.."
import Selector from "../generic/Selector"
import { ranks } from "../../constants"
import { Input } from "../ui/input"
import { Checkbox } from "../ui/checkbox"
import { Button } from "../ui/button"
import { getBaseUrl } from "../../lib/utils"
import { LoadingSpinner } from "../ui/spinner"
import SubordinatePicker from "./SubordinatePicker"

interface AdminMemberEditorProps {
    selectedMember: Member
    newMember: boolean
    setAddedMember: React.Dispatch<React.SetStateAction<number>>
    addedMember: number,
    setSelectedMember: React.Dispatch<React.SetStateAction<Member>>
    setNewMember: React.Dispatch<React.SetStateAction<boolean>>
    members: Member[]
    setMembers: React.Dispatch<React.SetStateAction<Member[]>>
}


export default function AdminMemberEditor({ selectedMember, newMember, setAddedMember, addedMember, setSelectedMember, setNewMember, members, setMembers }: AdminMemberEditorProps) {
    const [id, setId] = useState(selectedMember ? selectedMember.id : "")
    const [rank, setRank] = useState(selectedMember ? selectedMember.rank : "")
    const [firstName, setFirstName] = useState(selectedMember ? selectedMember.first_name : "")
    const [lastName, setLastName] = useState(selectedMember ? selectedMember.last_name : "")
    const [username, setUsername] = useState(selectedMember ? selectedMember.username : "")
    const [admin, setAdmin] = useState(selectedMember ? selectedMember.admin : false)
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    // const [subordinates, setSubordinates] = useState<Member[]>([])
    const [waiting, setWaiting] = useState(false)
    useEffect(() => {
        if (selectedMember !== null) {
            setId(selectedMember.id)
            setRank(selectedMember.rank)
            setFirstName(selectedMember.first_name)
            setLastName(selectedMember.last_name)
            setUsername(selectedMember.username)
            setAdmin(selectedMember.admin)
        }
    }, [selectedMember])
    async function updateMember(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault()
        setWaiting(true)
        if (password !== "" && confirmPassword !== password) {
            console.log("Passwords dont match!")
            return
        }
        const m = {
            "id": id,
            "username": username,
            "rank": rank,
            "first_name": firstName,
            "last_name": lastName,
            "admin": admin,
            "password": password,
            supervisor_id: selectedMember.supervisor_id
        }
        if (newMember) {
            await addMember(m, setAddedMember, addedMember, setSelectedMember);
        } else {
            await updateUser();
        }
        setWaiting(false)
        setPassword("")
        setConfirmPassword("")
        setNewMember(false)
    }
    return (
        <div className="grid grid-cols-adminPane">
            <form className="p-4 flex flex-col gap-4 flex-wrap" onSubmit={e => { updateMember(e) }}>
                <label>
                    ID: <Input className="inline w-72 bg-primary" value={id} disabled />
                </label>
                <label >
                    Rank: <Selector options={ranks} value={rank} setValue={setRank} />
                </label>
                <label >
                    First name: <Input className="inline w-48 bg-primary" value={firstName} onChange={e => { setFirstName(e.target.value) }} />
                </label>
                <label >
                    Last name: <Input className="inline w-48 bg-primary" value={lastName} onChange={e => { setLastName(e.target.value) }} />
                </label>
                <label>
                    Username: <Input className="inline w-48 bg-primary" value={username} onChange={e => { setUsername(e.target.value) }} />
                </label>
                <label htmlFor="admin">
                    Admin: <Checkbox id="admin" checked={admin} onClick={() => { setAdmin(!admin) }} />

                </label>
                <label>
                    Password: <Input type="password" className="inline w-48 bg-primary" value={password} onChange={e => { setPassword(e.target.value) }} />
                </label>

                <label>
                    Confirm Password: <Input type="password" className="inline w-48 bg-primary" value={confirmPassword} onChange={e => { setConfirmPassword(e.target.value) }} />
                </label>
                <Button type="submit" className=" w-36 bg-background hover:bg-background-dark text-white">{waiting === true ? <LoadingSpinner className="h-6 w-6 inline-block" /> : newMember ? "Create Member" : "Update Member"}</Button>
            </form>
            <SubordinatePicker className="h-72" members={members} setMembers={setMembers} supervisorID={selectedMember.id} />
        </div>
    )
}

async function addMember(m: Member, setAddedMember: React.Dispatch<React.SetStateAction<number>>, addedMember: number, setSelectedMember: React.Dispatch<React.SetStateAction<Member>>) {
    const res = await fetch(`${getBaseUrl()}/api/member`, {
        method: "POST",
        credentials: "same-origin",
        body: JSON.stringify(m)
    })
    if (res.status !== 201) {
        console.log("it failed")
        return
    }
    const memberJson = await res.json() as Member
    setSelectedMember(memberJson)
    setAddedMember(++addedMember)

}

async function updateUser() {

}