import { Qualification, Requirement, useCreateQualificationMutation, useGetAllRequirementsQuery, useUpdateQualificationMutation } from "../../../redux/api";
import { Input } from "../../ui/input";
import { Checkbox } from "../../ui/checkbox";
import { Button } from "../../ui/button";
import { Textarea } from "../../ui/textarea"
import { LoadingSpinner } from "../../ui/spinner";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
    adminQualificationSelector,
    changesCommitted,
    closeDialog,
    selectQualification,
    updateLocalSelectedQualification,
} from "../../../redux/adminQualificationSlice";
import ChangesPendingAlert from "../../generic/ChangesAlert";
import Picker from "../../generic/Picker";

export default function AdminQualificationEditor() {
    const {
        selectedQualification,
        updatedQualificationDraft,
        updatePending,
        newQualification,
        showChangeWarning,
        qualificationChanges
    } = useAppSelector(adminQualificationSelector);
    const dispatch = useAppDispatch();
    const [triggerUpdate, updateResponse] = useUpdateQualificationMutation();
    const [triggerCreate, createResponse] = useCreateQualificationMutation();
    const { data: requirements } = useGetAllRequirementsQuery();


    function commitUpdate() {
        console.log(document.activeElement)
        if (!updatedQualificationDraft) {
            return
        }
        switch (newQualification) {
            case true:
                triggerCreate(updatedQualificationDraft)
                if (createResponse.error) {
                    console.log(typeof createResponse.error)
                }
                dispatch(changesCommitted())
                break
            case false:
                triggerUpdate(updatedQualificationDraft)
                if (updateResponse.error) {
                    console.log(typeof updateResponse.error)
                }
        }
        dispatch(changesCommitted())
    }

    return updatedQualificationDraft === null || selectedQualification === null ? null : (
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
                        className="inline w-72 bg-primary"
                        value={updatedQualificationDraft.id}
                        disabled
                    />
                </label>
                <label>
                    Name:
                    <Input
                        className="inline w-48 bg-primary"
                        value={updatedQualificationDraft.name}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedQualification({
                                    ...updatedQualificationDraft,
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
                        value={updatedQualificationDraft.notes}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedQualification({
                                    ...updatedQualificationDraft,
                                    notes: e.target.value,
                                }),
                            );
                        }}
                    />
                </label>
                <label htmlFor="expires">
                    Expires:
                    <Checkbox
                        id="expires"
                        checked={updatedQualificationDraft.expires}
                        onClick={() => {
                            dispatch(
                                updateLocalSelectedQualification({
                                    ...updatedQualificationDraft,
                                    expires: !updatedQualificationDraft.expires,
                                }),
                            );
                        }}
                    />
                </label>
                <label>
                    Expiration interval (days):
                    <Input
                        type="number"
                        disabled={!updatedQualificationDraft.expires}
                        className="inline w-48 bg-primary"
                        value={updatedQualificationDraft.expiration_days}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedQualification({
                                    ...updatedQualificationDraft,
                                    expiration_days: Number(e.target.value)
                                })
                            )
                        }}
                    />
                </label>
                <label>
                    Initial Requirements:
                    <Picker
                        className="h-64"
                        pickedItems={[...updatedQualificationDraft.initial_requirements]}
                        unpickedItems={requirements?.filter(req => updatedQualificationDraft.initial_requirements.findIndex(
                            (value, index, obj) => value.id == req.id
                        ) == -1) as Requirement[]}
                        setPicked={(q) => {
                            dispatch(updateLocalSelectedQualification({
                                ...updatedQualificationDraft,
                                initial_requirements: updatedQualificationDraft.initial_requirements.concat(q),
                            }))
                        }}
                        setUnpicked={(q) => {
                            dispatch(updateLocalSelectedQualification({
                                ...updatedQualificationDraft,
                                initial_requirements: updatedQualificationDraft.initial_requirements.filter(req => req.id !== q.id)
                            }))
                        }}
                    />
                </label>
                <br />
                <Button
                    type="submit"
                    className=" w-36 bg-background hover:bg-background-dark text-white"
                    disabled={!newQualification && !qualificationChanges}
                >
                    {updatePending ? (
                        <LoadingSpinner className="h-6 w-6 inline-block" />
                    ) : newQualification ? (
                        "Create Qualification"
                    ) : (
                        "Update Qualification"
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
                    dispatch(selectQualification({ qualification: selectedQualification, force: true }));
                }}
            />
        </>
    );
}
