import { SetStateAction, useState } from "react"
import { Member } from "../.."
import { Button } from "../ui/button"
import { MoveLeftIcon, MoveRightIcon } from "lucide-react"


interface SubordinatePickerProps {
    members: Member[]
    setMembers: React.Dispatch<SetStateAction<Member[]>>
    supervisorID: string
    className?: string
}

export default function SubordinatePicker({ members, setMembers, supervisorID, className }: SubordinatePickerProps) {
    const [selected, setSelected] = useState<Member | null>(null)

    function handleAddSubordinate() {
        setMembers(members.map(member => {
            if (member.id === selected?.id) {
                member.supervisor_id = supervisorID
                return member
            }
            return member
        }))
    }

    function handleRemoveSubordinate() {
        setMembers(members.map(member => {
            if (member.id === selected?.id) {
                member.supervisor_id = ""
                return member
            }
            return member
        }))
    }
    return (
        <div className={`${className ? className : ""} grid grid-cols-subordinatePicker`}>
            <div className="h-full">
                <p>Subordinates:</p>
                <ul className="overflow-scroll border-black border h-full max-h-72">
                    {members.filter(member => member.supervisor_id !== "" && member.supervisor_id === supervisorID)
                        .map(member =>
                            <li key={member.id}
                                className={`${member === selected ? "bg-background text-white" : ""}`}
                                onClick={() => setSelected(member)}
                            >{member.rank} {member.first_name} {member.last_name}</li>
                        )}
                </ul>
            </div>
            <div className="flex flex-col items-center gap-4 justify-center">
                <Button className="bg-background hover:bg-background-dark" disabled={selected === null} onClick={handleAddSubordinate}><MoveLeftIcon color="white" /></Button>
                <Button className="bg-background  hover:bg-background-dark" disabled={selected === null} onClick={handleRemoveSubordinate}><MoveRightIcon color="white" /></Button>
            </div>
            <div className="h-full">
                <p>Available Members:</p>
                <ul className="overflow-scroll border-black border h-full max-h-72">
                    {members.filter(member => member.supervisor_id !== supervisorID && member.supervisor_id === "")
                        .sort((m1, m2) => m1.last_name.toLowerCase()[0] > m2.last_name[0].toLowerCase()[0] ? 1 : -1)
                        .map(member =>
                            <li key={member.id}
                                className={`${member === selected ? "bg-background text-white" : ""}`}
                                onClick={() => setSelected(member)}
                            >{member.rank} {member.first_name} {member.last_name}</li>
                        )}
                </ul>
            </div>
        </div>
    )
}