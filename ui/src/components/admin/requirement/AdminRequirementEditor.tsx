import { Reference, useCreateRequirementMutation, useGetAllReferencesQuery, useUpdateRequirementMutation } from "../../../redux/api";
import { Input } from "../../ui/input";
import { Button } from "../../ui/button";
import { Textarea } from "../../ui/textarea"
import { LoadingSpinner } from "../../ui/spinner";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
    adminRequirementSelector,
    changesCommitted,
    closeDialog,
    selectRequirement,
    updateLocalSelectedRequirement,
} from "../../../redux/adminRequirementSlice";
import ChangesPendingAlert from "../../generic/ChangesAlert";

import Selector from "../../generic/Selector";
import FullPageSpinner from "../../FullPageSpinner";

export default function AdminRequirementEditor() {
    const {
        selectedRequirement,
        updatedRequirementDraft,
        updatePending,
        newRequirement,
        showChangeWarning,
        requirementChanges
    } = useAppSelector(adminRequirementSelector);
    const dispatch = useAppDispatch();
    const [triggerUpdate, updateResponse] = useUpdateRequirementMutation();
    const [triggerCreate, createResponse] = useCreateRequirementMutation();
    const { data: references } = useGetAllReferencesQuery();

    function commitUpdate() {
        if (!updatedRequirementDraft) {
            return
        }
        switch (newRequirement) {
            case true:
                triggerCreate(updatedRequirementDraft)
                if (createResponse.error) {
                    console.log(typeof createResponse.error)
                }
                dispatch(changesCommitted())
                break
            case false:
                triggerUpdate(updatedRequirementDraft)
                if (updateResponse.error) {
                    console.log(typeof updateResponse.error)
                }
        }
        dispatch(changesCommitted())
    }

    return updatedRequirementDraft === null || selectedRequirement === null ? null : (
        <>
            <form
                className="p-4 flex gap-4 flex-col"
                onSubmit={(e) => {
                    e.preventDefault();
                    commitUpdate()
                }}
            >
                <label>
                    ID:
                    <Input
                        className="inline w-80 bg-primary"
                        value={updatedRequirementDraft.id}
                        disabled
                    />
                </label>
                <label>
                    Name:
                    <Input
                        className="inline w-48 bg-primary"
                        value={updatedRequirementDraft.name}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedRequirement({
                                    ...updatedRequirementDraft,
                                    name: e.target.value,
                                }),
                            );
                        }}
                    />
                </label>
                <label>
                    Notes:
                    <Textarea
                        className="w-1/2 h-72 bg-primary"
                        value={updatedRequirementDraft.notes}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedRequirement({
                                    ...updatedRequirementDraft,
                                    notes: e.target.value,
                                }),
                            );
                        }}
                    />
                </label>
                <label>
                    Expiration interval (days):
                    <Input
                        type="number"
                        className="inline w-48 bg-primary"
                        value={updatedRequirementDraft.days_valid_for}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedRequirement({
                                    ...updatedRequirementDraft,
                                    days_valid_for: Number(e.target.value)
                                })
                            )
                        }}
                    />
                </label>
                <label>
                    Reference:
                    {references ?
                        <Selector
                            options={
                                references.map((reference) => ({ label: reference.name, value: reference.id }))
                            }
                            setValue={value => {
                                dispatch(updateLocalSelectedRequirement({
                                    ...updatedRequirementDraft,
                                    reference: references.find((val, index, obj) => value === val.id) as Reference
                                })
                                )
                            }
                            }
                            value={updatedRequirementDraft.reference.id}
                        /> :
                        <FullPageSpinner />
                    }
                </label>
                <Button
                    type="submit"
                    className=" w-36 bg-background hover:bg-background-dark text-white"
                    disabled={!newRequirement && !requirementChanges}
                >
                    {updatePending ? (
                        <LoadingSpinner className="h-6 w-6 inline-block" />
                    ) : newRequirement ? (
                        "Create Requirement"
                    ) : (
                        "Update Requirement"
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
                    dispatch(selectRequirement({ requirement: selectedRequirement, force: true }));
                }}
            />
        </>
    );
}
