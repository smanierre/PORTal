import { Member } from "..";

interface ProfileCardProps {
    className?: string
    member: Member | null
    subordinates: Member[]
}

export default function ProfileCard({ className, member, subordinates }: ProfileCardProps) {
    return <article className={`${className !== undefined ? className : "p-4"}`}>
        <p>{member?.rank} {member?.first_name} {member?.last_name}</p>
        <p className="py-2">Subordinates:</p>
        <ul className="overflow-scroll">
            {subordinates && subordinates.map(subordinate => <li>{subordinate.rank} {subordinate.first_name} {subordinate.last_name}</li>)}
        </ul>

    </article>
}