import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Reference } from "./api";
import { RootState } from "./store";
import { getEmptyReference } from "../lib/utils";

export interface AdminReferenceState {
    selectedReference: Reference | null;
    pendingSelectedReference: Reference | null;
    updatedReferenceDraft: Reference | null;
    updatePending: boolean;
    updateError: string;
    newReference: boolean;
    referenceChanges: boolean;
    showChangeWarning: boolean;
}

const initialState: AdminReferenceState = {
    selectedReference: null,
    pendingSelectedReference: null,
    updatedReferenceDraft: null,
    updatePending: false,
    updateError: "",
    newReference: false,
    referenceChanges: false,
    showChangeWarning: false
}

export const adminReferenceSlice = createSlice({
    name: "adminReference",
    initialState,
    reducers: {
        selectReference: (state, action: PayloadAction<{ reference: Reference | null, force?: boolean }>) => {
            if (state.referenceChanges && !action.payload.force) {
                state.pendingSelectedReference = action.payload.reference
                state.showChangeWarning = true;
                return;
            }
            if (action.payload.force) {
                state.updatedReferenceDraft = state.pendingSelectedReference;
                state.selectedReference = state.pendingSelectedReference;
                state.pendingSelectedReference = null;
                state.showChangeWarning = false;
                state.referenceChanges = false;
                return
            }
            if (state.newReference && state.referenceChanges) {
                state.newReference = true
            } else {
                state.newReference = false;
            }
            state.newReference = false;
            state.referenceChanges = false;
            state.updatedReferenceDraft = action.payload.reference;
            state.selectedReference = action.payload.reference;
        },
        updateLocalSelectedReference: (state, action: PayloadAction<Reference>) => {
            state.referenceChanges = true;
            state.updatedReferenceDraft = action.payload;
        },
        newReference: (state) => {
            if (state.referenceChanges) {
                state.showChangeWarning = true
                state.pendingSelectedReference = getEmptyReference();
                return
            }
            state.selectedReference = getEmptyReference();
            state.updatedReferenceDraft = getEmptyReference();
            state.newReference = true
        },
        closeDialog: (state) => {
            state.showChangeWarning = false;
        },
        changesCommitted: (state) => {
            state.referenceChanges = false
        }
    }
})

export const { selectReference, newReference, closeDialog, changesCommitted, updateLocalSelectedReference } = adminReferenceSlice.actions
export const adminReferenceSelector = (state: RootState) => state.adminReference
export default adminReferenceSlice.reducer