import React, { useEffect, useState } from "react";
import { Member } from "../..";
import Selector from "../generic/Selector";
import { Input } from "../ui/input";
import { Checkbox } from "../ui/checkbox";
import { Button } from "../ui/button";
import {
  convertGrade,
  getBaseUrl,
  Grades,
  isHigherRank,
} from "../../lib/utils";
import { LoadingSpinner } from "../ui/spinner";

interface AdminMemberEditorProps {
  selectedMember: Member;
  newMember: boolean;
  setAddedMember: React.Dispatch<React.SetStateAction<number>>;
  addedMember: number;
  setSelectedMember: React.Dispatch<React.SetStateAction<Member>>;
  setNewMember: React.Dispatch<React.SetStateAction<boolean>>;
  members: Member[];
  setMembers: React.Dispatch<React.SetStateAction<Member[]>>;
}

export default function AdminMemberEditor({
  selectedMember,
  newMember,
  setAddedMember,
  addedMember,
  setSelectedMember,
  setNewMember,
  members,
  setMembers,
}: AdminMemberEditorProps) {
  const [member, setMember] = useState(selectedMember);
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [waiting, setWaiting] = useState(false);

  useEffect(() => {
    setMember(selectedMember);
  }, [selectedMember]);

  async function updateMember(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setWaiting(true);
    if (password !== "" && confirmPassword !== password) {
      console.log("Passwords dont match!");
      return;
    }

    if (newMember) {
      await addMember(
        member,
        setAddedMember,
        addedMember,
        setSelectedMember,
        password,
      );
    } else {
      await updateExistingMember(member);
      setMembers(members.map((m) => (m.id === member.id ? member : m)));
    }
    setWaiting(false);
    setPassword("");
    setConfirmPassword("");
    setNewMember(false);
    setMember(member);
  }
  return (
    <div className="" /*"grid grid-cols-adminPane*/>
      <form
        className="p-4 flex gap-4 flex-wrap"
        onSubmit={(e) => {
          updateMember(e);
        }}
      >
        <label>
          ID:{" "}
          <Input
            className="inline w-72 bg-primary"
            value={member.id}
            disabled
          />
        </label>
        <label>
          Rank:{" "}
          <Selector
            options={Grades.map((grade) => {
              return {
                label: convertGrade(grade),
                value: grade,
              };
            })}
            value={member.rank}
            setValue={(selectedValue) => {
              members.forEach((m) => {
                if (
                  m.id === member.supervisor_id &&
                  isHigherRank(member.rank, m.rank)
                ) {
                  console.log("blanking supervisor id");
                  setMember({
                    ...member,
                    rank: selectedValue,
                    supervisor_id: "",
                  });
                  return;
                }
              });
              setMember({ ...member, rank: selectedValue });
            }}
          />
        </label>
        <label>
          First name:{" "}
          <Input
            className="inline w-48 bg-primary"
            value={member.first_name}
            onChange={(e) => {
              setMember({ ...member, first_name: e.target.value });
            }}
          />
        </label>
        <label>
          Last name:{" "}
          <Input
            className="inline w-48 bg-primary"
            value={member.last_name}
            onChange={(e) => {
              setMember({ ...member, last_name: e.target.value });
            }}
          />
        </label>
        <label>
          Username:{" "}
          <Input
            className="inline w-48 bg-primary"
            value={member.username}
            onChange={(e) => {
              setMember({ ...member, username: e.target.value });
            }}
          />
        </label>
        {members.filter((m) => isHigherRank(m.rank, member.rank)).length > 0 ? (
          <label>
            Supervisor:
            <Selector
              value={member.supervisor_id}
              options={members
                .filter(
                  (m) =>
                    m.supervisor_id === "" &&
                    m.id !== member.id &&
                    !isHigherRank(member.rank, m.rank),
                )
                .map((item) => ({
                  label: `${item.first_name} ${item.last_name}`,
                  value: item.id,
                }))}
              setValue={(selectedId) => {
                setMember({ ...member, supervisor_id: selectedId });
              }}
            />
          </label>
        ) : null}
        <label htmlFor="admin">
          Admin:{" "}
          <Checkbox
            id="admin"
            checked={member.admin}
            onClick={() => {
              setMember({ ...member, admin: !member.admin });
            }}
          />
        </label>
        <label>
          Password:{" "}
          <Input
            type="password"
            className="inline w-48 bg-primary"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
            }}
          />
        </label>

        <label>
          Confirm Password:{" "}
          <Input
            type="password"
            className="inline w-48 bg-primary"
            value={confirmPassword}
            onChange={(e) => {
              setConfirmPassword(e.target.value);
            }}
          />
        </label>
        <Button
          type="submit"
          className=" w-36 bg-background hover:bg-background-dark text-white"
        >
          {waiting === true ? (
            <LoadingSpinner className="h-6 w-6 inline-block" />
          ) : newMember ? (
            "Create Member"
          ) : (
            "Update Member"
          )}
        </Button>
      </form>
    </div>
  );
}

async function addMember(
  m: Member,
  setAddedMember: React.Dispatch<React.SetStateAction<number>>,
  addedMember: number,
  setSelectedMember: React.Dispatch<React.SetStateAction<Member>>,
  password: string,
) {
  const res = await fetch(`${getBaseUrl()}/api/member`, {
    method: "POST",
    credentials: "same-origin",
    body: JSON.stringify({ ...m, password: password }),
  });
  if (res.status !== 201) {
    console.log("it failed");
    return;
  }
  const memberJson = (await res.json()) as Member;
  setSelectedMember(memberJson);
  setAddedMember(++addedMember);
}

async function updateExistingMember(member: Member) {
  const res = await fetch(`${getBaseUrl()}/api/member/${member.id}`, {
    method: "PUT",
    credentials: "same-origin",
    body: JSON.stringify(member),
  });
  if (res.status != 200) {
    console.log("it failed");
  }
}
