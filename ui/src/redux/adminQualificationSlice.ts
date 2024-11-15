import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Qualification } from "./api";
import { RootState } from "./store";
import { getEmptyQualification } from "../lib/utils";

export interface AdminQualificationState {
    selectedQualification: Qualification | null;
    pendingSelectedQualification: Qualification | null;
    updatedQualificationDraft: Qualification | null;
    updatePending: boolean;
    updateError: string;
    newQualification: boolean;
    qualificationChanges: boolean;
    showChangeWarning: boolean;
}

const initialState: AdminQualificationState = {
    selectedQualification: null,
    pendingSelectedQualification: null,
    updatedQualificationDraft: null,
    updatePending: false,
    updateError: "",
    newQualification: false,
    qualificationChanges: false,
    showChangeWarning: false
}

export const adminQualificationSlice = createSlice({
    name: "adminQualification",
    initialState,
    reducers: {
        selectQualification: (state, action: PayloadAction<{ qualification: Qualification | null, force?: boolean }>) => {
            if (state.qualificationChanges && !action.payload.force) {
                state.pendingSelectedQualification = action.payload.qualification
                state.showChangeWarning = true;
                return;
            }
            if (action.payload.force) {
                state.updatedQualificationDraft = state.pendingSelectedQualification;
                state.selectedQualification = state.pendingSelectedQualification;
                state.pendingSelectedQualification = null;
                state.showChangeWarning = false;
                state.qualificationChanges = false;
                return
            }
            state.newQualification = false;
            state.qualificationChanges = false;
            state.updatedQualificationDraft = action.payload.qualification;
            state.selectedQualification = action.payload.qualification;
        },
        updateLocalSelectedQualification: (state, action: PayloadAction<Qualification>) => {
            state.qualificationChanges = true;
            state.updatedQualificationDraft = action.payload;
        },
        newQualification: (state) => {
            if (state.qualificationChanges) {
                state.showChangeWarning = true
                state.pendingSelectedQualification = getEmptyQualification();
                return
            }
            state.selectedQualification = getEmptyQualification();
            state.updatedQualificationDraft = getEmptyQualification();
            state.newQualification = true
        },
        closeDialog: (state) => {
            state.showChangeWarning = false;
        },
        changesCommitted: (state) => {
            state.qualificationChanges = false
        }
    }
})

export const { selectQualification, newQualification, closeDialog, changesCommitted, updateLocalSelectedQualification } = adminQualificationSlice.actions
export const adminQualificationSelector = (state: RootState) => state.adminQualification
export default adminQualificationSlice.reducer