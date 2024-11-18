import { useCreateReferenceMutation, useUpdateReferenceMutation } from "../../../redux/api";
import { Input } from "../../ui/input";
import { Button } from "../../ui/button";
import { LoadingSpinner } from "../../ui/spinner";
import { useAppDispatch, useAppSelector } from "../../../redux/hooks";
import {
    adminReferenceSelector,
    changesCommitted,
    closeDialog,
    selectReference,
    updateLocalSelectedReference,
} from "../../../redux/adminReferenceSlice";
import ChangesPendingAlert from "../../generic/ChangesAlert";

export default function AdminReferenceEditor() {
    const {
        selectedReference,
        updatedReferenceDraft,
        updatePending,
        newReference,
        showChangeWarning,
        referenceChanges
    } = useAppSelector(adminReferenceSelector);
    const dispatch = useAppDispatch();
    const [triggerUpdate, updateResponse] = useUpdateReferenceMutation();
    const [triggerCreate, createResponse] = useCreateReferenceMutation();

    function commitUpdate() {
        if (!updatedReferenceDraft) {
            return
        }
        switch (newReference) {
            case true:
                triggerCreate(updatedReferenceDraft)
                if (createResponse.error) {
                    console.log(typeof createResponse.error)
                }
                dispatch(changesCommitted())
                break
            case false:
                triggerUpdate(updatedReferenceDraft)
                if (updateResponse.error) {
                    console.log(typeof updateResponse.error)
                }
        }
        dispatch(changesCommitted())
    }

    return updatedReferenceDraft === null || selectedReference === null ? null : (
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
                        value={updatedReferenceDraft.id}
                        disabled
                    />
                </label>
                <label>
                    Name:
                    <Input
                        className="inline w-48 bg-primary"
                        value={updatedReferenceDraft.name}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedReference({
                                    ...updatedReferenceDraft,
                                    name: e.target.value,
                                }),
                            );
                        }}
                    />
                </label>
                <label>
                    Volume (Optional, 0 for no volume):
                    <Input
                        className="inline w-48 bg-primary"
                        value={updatedReferenceDraft.volume}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedReference({
                                    ...updatedReferenceDraft,
                                    volume: Number(e.target.value),
                                }),
                            );
                        }}
                    />
                </label>
                <label>
                    Paragraph:
                    <Input
                        className="inline w-48 bg-primary"
                        value={updatedReferenceDraft.paragraph}
                        onChange={(e) => {
                            dispatch(
                                updateLocalSelectedReference({
                                    ...updatedReferenceDraft,
                                    paragraph: e.target.value,
                                }),
                            );
                        }}
                    />
                </label>
                <Button
                    type="submit"
                    className=" w-36 bg-background hover:bg-background-dark text-white"
                    disabled={!newReference && !referenceChanges}
                >
                    {updatePending ? (
                        <LoadingSpinner className="h-6 w-6 inline-block" />
                    ) : newReference ? (
                        "Create Reference"
                    ) : (
                        "Update Reference"
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
                    dispatch(selectReference({ reference: selectedReference, force: true }));
                }}
            />
        </>
    );
}
