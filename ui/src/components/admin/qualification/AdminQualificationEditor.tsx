import { Qualification, Requirement, useCreateQualificationMutation, useGetAllRequirementsQuery, useUpdateQualificationMutation } from "../../../redux/api";
import { useState, useEffect } from "react";
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
    setChanges,
} from "../../../redux/adminQualificationSlice";
import ChangesPendingAlert from "../../generic/ChangesAlert";
import Picker from "../../generic/Picker";

export default function AdminQualificationEditor() {
    const {
        selectedQualification,
        updatePending,
        newQualification,
        qualificationChanges,
        showChangeWarning,
    } = useAppSelector(adminQualificationSelector);

    const dispatch = useAppDispatch();
    const [triggerUpdate, updateResponse] = useUpdateQualificationMutation();
    const [triggerCreate, createResponse] = useCreateQualificationMutation();
    const { data: requirements } = useGetAllRequirementsQuery();
    const [qualification, setQualification] = useState(selectedQualification)

    useEffect(() => {
        setQualification(selectedQualification)
    }, [selectedQualification])

    useEffect(() => {
        if (JSON.stringify(selectedQualification) !== JSON.stringify(qualification)) {
            dispatch(setChanges(true))
            return
        }
        dispatch(setChanges(false))
    }, [qualification])

    function commitUpdate() {
        if (!qualification) {
            return
        }
        switch (newQualification) {
            case true:
                triggerCreate(qualification)
                if (createResponse.error) {
                    console.log(typeof createResponse.error)
                }
                dispatch(changesCommitted())
                break
            case false:
                triggerUpdate(qualification)
                if (updateResponse.error) {
                    console.log(typeof updateResponse.error)
                }
        }
        dispatch(changesCommitted())
    }

    return qualification === null ? null : (
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
                        value={qualification.id}
                        disabled
                    />
                </label>
                <label>
                    Name:
                    <Input
                        className="inline w-48 bg-primary"
                        value={qualification.name}
                        onChange={(e) => {
                            setQualification({ ...qualification, name: e.target.value })
                        }}
                    />
                </label>
                <label>
                    Notes:
                    <Textarea
                        className="w-1/2 h-72 bg-primary"
                        value={qualification.notes}
                        onChange={(e) => {
                            setQualification({ ...qualification, notes: e.target.value })
                        }}
                    />
                </label>
                <label htmlFor="expires">
                    Expires:
                    <Checkbox
                        id="expires"
                        checked={qualification.expires}
                        onClick={() => {
                            setQualification({ ...qualification, expires: !qualification.expires })
                        }}
                    />
                </label>
                <label>
                    Expiration interval (days):
                    <Input
                        type="number"
                        disabled={!qualification.expires}
                        className="inline w-48 bg-primary"
                        value={qualification.expiration_days}
                        onChange={(e) => {
                            setQualification({ ...qualification, expiration_days: Number(e.target.value) })
                        }}
                    />
                </label>
                <label>
                    Initial Requirements:
                    <Picker
                        className="h-64"
                        pickedItems={[...qualification.initial_requirements]}
                        unpickedItems={determineUnpickedItems(qualification.initial_requirements, requirements ? requirements : [])}
                        setPicked={(q) => {
                            setQualification({ ...qualification, initial_requirements: qualification.initial_requirements.concat(q) })
                        }}
                        setUnpicked={(q) => {
                            setQualification({
                                ...qualification,
                                initial_requirements: qualification.initial_requirements.filter(
                                    req => req.id !== q.id
                                )
                            })
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

function determineUnpickedItems(picked: Requirement[], allReqs: Requirement[]): Requirement[] {
    if (picked.length === allReqs.length) {
        return [] as Requirement[]
    }
    if (picked.length === 0) {
        return allReqs
    }
    return allReqs.filter(
        req => picked.find(
            picked => picked.id !== req.id
        )
    ) as Requirement[]
}