import { Member } from "../../redux/api";
import { useState } from "react";
import { convertGrade } from "../../lib/utils";
import { Button } from "../ui/button";
import Search from "../generic/Search";
import { useAppDispatch, useAppSelector } from "../../redux/hooks";
import {
  adminSelector,
  newMember,
  selectMember,
} from "../../redux/adminMemberSlice";

interface AdminMemberListProps {
  members: Member[];
}
export default function AdminMemberList({ members }: AdminMemberListProps) {
  const [searchTerm, setSearchTerm] = useState("");
  const { selectedMember } = useAppSelector(adminSelector);
  const dispatch = useAppDispatch();

  return (
    <div className="flex flex-col items-center align-middle">
      <div className="w-1/2 ml-auto">
        <Search
          className="w-2/3 inline-block border-b-0"
          searchTerm={searchTerm}
          setSearchTerm={setSearchTerm}
        />
        <Button
          className="inline-block bg-background hover:bg-background-dark text-white w-1/3"
          onClick={() => {
            setSearchTerm("");
          }}
        >
          Clear
        </Button>
      </div>
      <ul className="overflow-x-scroll ml-auto h-2/3 w-1/2 border-black border p-2">
        {members
          .filter((member) => {
            return (
              member.first_name
                .toLowerCase()
                .includes(searchTerm.toLowerCase()) ||
              member.last_name
                .toLowerCase()
                .includes(searchTerm.toLowerCase()) ||
              convertGrade(member.grade)
                .toLowerCase()
                .includes(searchTerm.toLowerCase())
            );
          })
          .sort((m1, m2) =>
            m1.last_name.toLowerCase()[0] > m2.last_name[0].toLowerCase()[0]
              ? 1
              : -1,
          )
          .map((member) => (
            <ul
              key={member.id}
              className={`${selectedMember?.id === member.id ? "bg-background text-white" : ""} cursor-pointer`}
              onClick={() => {
                dispatch(selectMember({ member: member }));
              }}
            >
              {convertGrade(member.grade)} {member.first_name}{" "}
              {member.last_name}
            </ul>
          ))}
      </ul>
      <Button
        className="bg-background text-white hover:bg-background-dark ml-auto mt-2"
        onClick={() => {
          dispatch(newMember());
        }}
      >
        Add Member
      </Button>
    </div>
  );
}
