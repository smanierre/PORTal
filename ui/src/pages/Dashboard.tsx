import { Accordion, AccordionItem } from "@radix-ui/react-accordion";
import ProfileCard from "../components/ProfileCard.tsx";
import QualificationList from "../components/QualificationList.tsx";
import useLoginRequired from "../hooks/useLoginRequired.ts";
import { Member, Qualification } from "../index";
import { AccordionContent, AccordionTrigger } from "../components/ui/accordion.tsx";

interface DashboardProps {
    member: Member | null
    qualifications: Qualification[]
    subordinates: Member[]
}

export default function Dashboard({ member, qualifications, subordinates }: DashboardProps) {
    useLoginRequired()
    if (member === null) {
        return <></>
    }
    return (
        <div className={"w-full h-full px-4"}>
            <Accordion type="single" collapsible defaultValue="profile" className="m-4">
                <AccordionItem value="profile" className="px-2 my-4 bg-background text-primary">
                    <AccordionTrigger>Profile</AccordionTrigger>
                    <AccordionContent>
                        <ProfileCard member={member} subordinates={subordinates} />
                    </AccordionContent>
                </AccordionItem>
                <AccordionItem value="qualifications" className="px-2 my-4 bg-background text-primary">
                    <AccordionTrigger>Qualifications</AccordionTrigger>
                    <AccordionContent>
                        <QualificationList qualifications={qualifications} />
                    </AccordionContent>
                </AccordionItem>
            </Accordion>
        </div>
    )
}