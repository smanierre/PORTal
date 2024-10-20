import { Qualification } from ".."

interface QualificationListProps {
    qualifications: Qualification[]
}

export default function QualificationList({ qualifications }: QualificationListProps) {
    return <div>
        {qualifications.length > 0 ? qualifications.map(qualification => {
            return <p key={qualification.id}>{qualification.name} {qualification.notes} {qualification.expiration_days} {qualification.expires}</p>
        }) : <p>No qualifications</p>}
    </div>
}