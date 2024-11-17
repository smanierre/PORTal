import { useState } from "react";
import { Member, useCreateMemberMutation, useUpdateMemberMutation } from "../../../redux/api";
import Selector from "../../generic/Selector";
import { Input } from "../../ui/input";
import { Checkbox } from "../../ui/checkbox";
import { Button } from "../../ui/button";
import {
  convertGrade,
  validatePassword,
  Grades,
  isHigherRank,
  getEmptyMember,
} from "../../../lib/utils";
import { LoadingSpinner } from "../../ui/spinner";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
  adminMemberSelector,
  changesCommitted,
  closeDialog,
  selectMember,
  updateLocalSelectedMember,
} from "../../../redux/adminMemberSlice";
import ChangesPendingAlert from "../../generic/ChangesAlert";

interface AdminMemberEditorProps {
  members: Member[];
}

export default function AdminMemberEditor({ members }: AdminMemberEditorProps) {
  const {
    selectedMember,
    updatedMemberDraft,
    updatePending,
    newMember,
    showChangeWarning,
    memberChanges
  } = useAppSelector(adminMemberSelector);
  const dispatch = useAppDispatch();
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [triggerUpdate, updateResponse] = useUpdateMemberMutation()
  const [triggerCreate, createResponse] = useCreateMemberMutation()

  function handleSetGrade(selectedValue: string) {
    if (selectedMember === null || updatedMemberDraft === null) {
      return;
    }
    let currentSupervisorRank = members.reduce((_, cur) =>
      cur.id === selectedMember?.supervisor_id ? cur : getEmptyMember(),
    ).grade;
    // New rank is higher than supervisor, blank out current supervisor
    if (isHigherRank(selectedValue, currentSupervisorRank)) {
      dispatch(
        updateLocalSelectedMember({
          ...updatedMemberDraft,
          grade: selectedValue,
          supervisor_id: "",
        }),
      );
      return;
    }
    // New rank is lower than supervisor, and supervisor hasn't been changed, keep original supervisor
    else if (
      !isHigherRank(selectedValue, currentSupervisorRank) &&
      selectedMember.supervisor_id === updatedMemberDraft.supervisor_id
    ) {
      dispatch(
        updateLocalSelectedMember({
          ...updatedMemberDraft,
          grade: selectedValue,
        }),
      );
    }
    // New rank is lower than previous supervisor, but supervisor was updated to be empty, restore supervisor
    else if (
      !isHigherRank(selectedValue, currentSupervisorRank) &&
      updatedMemberDraft.supervisor_id === ""
    ) {
      dispatch(
        updateLocalSelectedMember({
          ...updatedMemberDraft,
          grade: selectedValue,
          supervisor_id: selectedMember.supervisor_id,
        }),
      );
    }
  }

  function commitUpdate() {
    if (password !== "" && password !== confirmPassword) {
      return;
    }
    if (!updatedMemberDraft) {
      return
    }
    switch (newMember) {
      case true:
        triggerCreate({ ...updatedMemberDraft, password: password })
        if (createResponse.error) {
          console.log(typeof createResponse.error)
        }
        break
      case false:
        triggerUpdate({ ...updatedMemberDraft, password: password })
        if (updateResponse.error) {
          console.log(typeof updateResponse.error)
        }
    }
    dispatch(changesCommitted())
  }

  return updatedMemberDraft === null || selectedMember === null ? null : (
    <div>
      <form
        className="p-4 flex gap-4 flex-col"
        onSubmit={(e) => {
          e.preventDefault();
          if (password !== "" && !validatePassword(password, confirmPassword)) {
            //TODO: handle this in UI
            return;
          }
          commitUpdate()
        }}
      >
        <label>
          ID:
          <Input
            className="inline w-80 bg-primary"
            value={updatedMemberDraft.id}
            disabled
          />
        </label>
        <label>
          Rank:
          <Selector
            options={Grades.map((grade) => {
              return {
                label: convertGrade(grade),
                value: grade,
              };
            })}
            value={updatedMemberDraft.grade}
            setValue={handleSetGrade}
          />
        </label>
        <label>
          First name:
          <Input
            className="inline w-48 bg-primary"
            value={updatedMemberDraft.first_name}
            onChange={(e) => {
              dispatch(
                updateLocalSelectedMember({
                  ...updatedMemberDraft,
                  first_name: e.target.value,
                }),
              );
            }}
          />
        </label>
        <label>
          Last name:
          <Input
            className="inline w-48 bg-primary"
            value={updatedMemberDraft.last_name}
            onChange={(e) => {
              dispatch(
                updateLocalSelectedMember({
                  ...updatedMemberDraft,
                  last_name: e.target.value,
                }),
              );
            }}
          />
        </label>
        <label>
          Username:
          <Input
            disabled={!newMember}
            className="inline w-48 bg-primary"
            value={updatedMemberDraft.username}
            onChange={(e) => {
              dispatch(
                updateLocalSelectedMember({
                  ...updatedMemberDraft,
                  username: e.target.value,
                }),
              );
            }}
          />
        </label>
        {members.filter((m) => isHigherRank(m.grade, updatedMemberDraft.grade))
          .length > 0 ? (
          <label>
            Supervisor:
            <Selector
              optional
              value={updatedMemberDraft.supervisor_id}
              options={members
                .filter(
                  (m) =>
                    m.supervisor_id === "" &&
                    m.id !== updatedMemberDraft.id &&
                    !isHigherRank(updatedMemberDraft.grade, m.grade),
                )
                .map((item) => ({
                  label: `${convertGrade(item.grade)} ${item.first_name} ${item.last_name}`,
                  value: item.id,
                }))}
              setValue={(selectedId) => {
                dispatch(
                  updateLocalSelectedMember({
                    ...updatedMemberDraft,
                    supervisor_id: selectedId,
                  }),
                );
              }}
            />
          </label>
        ) : null}
        <label htmlFor="admin">
          Admin:
          <Checkbox
            id="admin"
            checked={updatedMemberDraft.admin}
            onClick={() => {
              dispatch(
                updateLocalSelectedMember({
                  ...updatedMemberDraft,
                  admin: !updatedMemberDraft.admin,
                }),
              );
            }}
          />
        </label>
        <label>
          Password:
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
          Confirm Password:
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
          disabled={!newMember && !memberChanges}
        >
          {updatePending ? (
            <LoadingSpinner className="h-6 w-6 inline-block" />
          ) : newMember ? (
            "Create Member"
          ) : (
            "Update Member"
          )}
        </Button>
      </form>
      <ChangesPendingAlert
        open={showChangeWarning}
        onOpenChange={() => { }}
        accept={() => {
          dispatch(closeDialog());
        }}
        cancel={() => {
          dispatch(selectMember({ member: selectedMember, force: true }));
        }}
      />
    </div>
  );
}
