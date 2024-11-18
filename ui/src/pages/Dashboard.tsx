import { Accordion, AccordionItem } from "@radix-ui/react-accordion";
import ProfileCard from "../components/ProfileCard.tsx";
import useLoginRequired from "../hooks/useLoginRequired.ts";
import {
  AccordionContent,
  AccordionTrigger,
} from "../components/ui/accordion.tsx";
import {
  useGetLoggedInMemberQuery,
  useGetMemberSubordinatesQuery,
} from "../redux/api.ts";
import FullPageSpinner from "../components/FullPageSpinner.tsx";
import { Member } from "../redux/api.ts";

export default function Dashboard() {
  useLoginRequired();
  const { data: member, isLoading: memberIsLoading } =
    useGetLoggedInMemberQuery();
  const { data: subordinates, isLoading: subordinatesIsLoading } =
    useGetMemberSubordinatesQuery(member ? member.id : "");
  return (
    <div className={"w-full h-full px-4"}>
      <Accordion
        type="single"
        collapsible
        defaultValue="profile"
        className="m-4"
      >
        <AccordionItem
          value="profile"
          className="px-2 my-4 bg-background text-primary"
        >
          <AccordionTrigger>Profile</AccordionTrigger>
          <AccordionContent>
            {memberIsLoading || subordinatesIsLoading ? (
              <FullPageSpinner />
            ) : (
              <ProfileCard
                member={member as Member}
                subordinates={subordinates as Member[]}
              />
            )}
          </AccordionContent>
        </AccordionItem>
        <AccordionItem
          value="qualifications"
          className="px-2 my-4 bg-background text-primary"
        >
          <AccordionTrigger>Qualifications</AccordionTrigger>
          <AccordionContent>
            {/* <QualificationList qualifications={qualifications} /> */}
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
