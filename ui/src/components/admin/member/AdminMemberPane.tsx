import AdminMemberEditor from "./AdminMemberEditor";
import { useGetAllMembersQuery } from "../../../redux/api";
import FullPageSpinner from "../../FullPageSpinner";
import { Member } from "../../../redux/api";
import SearchableList from "../../generic/SearchableList";
import { newMember, selectMember } from "../../../redux/adminMemberSlice";
import { convertGrade } from "../../../lib/utils";
import { useAppDispatch } from "../../../redux/hooks";

export default function AdminMemberPane() {
  const { data: members, isLoading, isError } = useGetAllMembersQuery();
  const dispatch = useAppDispatch()
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
        <SearchableList
          items={members as Member[]}
          newItemFunc={newMember}
          displayFunc={displayMember}
          filterFunc={memberFilter}
          selectItem={(member: Member) => {
            dispatch(selectMember({ member }))
          }}
          sortFunc={memberSort}
          typeName="Member"
        />
        <AdminMemberEditor members={members as Member[]} />
      </div>
    );
}

function displayMember(member: Member): string {
  return `${convertGrade(member.grade)} ${member.first_name} ${member.last_name}`
}

function memberFilter(searchTerm: string): (item: Member) => boolean {
  return (member: Member): boolean => {
    return member.first_name
      .toLowerCase()
      .includes(searchTerm.toLowerCase()) ||
      member.last_name
        .toLowerCase()
        .includes(searchTerm.toLowerCase()) ||
      convertGrade(member.grade)
        .toLowerCase()
        .includes(searchTerm.toLowerCase())
      ;
  }
}

function memberSort(i1: Member, i2: Member): number {
  return i1.last_name.toLowerCase()[0] > i2.last_name[0].toLowerCase()[0] ? 1 : -1
}