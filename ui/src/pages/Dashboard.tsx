import useLoginRequired from "../hooks/useLoginRequired.ts";
import {Button} from "../components/ui/button.tsx";
import ProfileCard from "../components/ProfileCard.tsx";

export default function Dashboard() {
    useLoginRequired()
    return (
        <div className={"grid grid-rows-2 grid-cols-2 w-full h-full"}>
            <ProfileCard/>
        </div>
    )
}