import AdminMemberEditor from "./AdminMemberEditor";
import AdminMemberList from "./AdminMemberList";
import { useGetAllMembersQuery } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import { Member } from "../../../redux/api";

export default function AdminMemberPane() {
  const { data: members, isLoading, isError } = useGetAllMembersQuery();
  return isLoading ? (
    <FullPageSpinner />
  ) : isError ? (
    /* TODO: Create an error page*/
    <div>
      "Error loading content"
    </div >
  ) :
    (
      <div className="grid grid-cols-adminPane">
        <AdminMemberList members={members as Member[]} />
        <AdminMemberEditor members={members as Member[]} />
      </div>
    );
}
