import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Requirement } from "./api";
import { RootState } from "./store";
import { getEmptyRequirement } from "../lib/utils";

export interface AdminRequirementState {
    selectedRequirement: Requirement | null;
    pendingSelectedRequirement: Requirement | null;
    updatedRequirementDraft: Requirement | null;
    updatePending: boolean;
    updateError: string;
    newRequirement: boolean;
    requirementChanges: boolean;
    showChangeWarning: boolean;
}

const initialState: AdminRequirementState = {
    selectedRequirement: null,
    pendingSelectedRequirement: null,
    updatedRequirementDraft: null,
    updatePending: false,
    updateError: "",
    newRequirement: false,
    requirementChanges: false,
    showChangeWarning: false
}

export const adminRequirementSlice = createSlice({
    name: "adminRequirement",
    initialState,
    reducers: {
        selectRequirement: (state, action: PayloadAction<{ requirement: Requirement | null, force?: boolean }>) => {
            if (state.requirementChanges && !action.payload.force) {
                state.pendingSelectedRequirement = action.payload.requirement
                state.showChangeWarning = true;
                return;
            }
            if (action.payload.force) {
                state.updatedRequirementDraft = state.pendingSelectedRequirement;
                state.selectedRequirement = state.pendingSelectedRequirement;
                state.pendingSelectedRequirement = null;
                state.showChangeWarning = false;
                state.requirementChanges = false;
                return
            }
            state.newRequirement = false;
            state.requirementChanges = false;
            state.updatedRequirementDraft = action.payload.requirement;
            state.selectedRequirement = action.payload.requirement;
        },
        updateLocalSelectedRequirement: (state, action: PayloadAction<Requirement>) => {
            state.requirementChanges = true;
            state.updatedRequirementDraft = action.payload;
        },
        newRequirement: (state) => {
            if (state.requirementChanges) {
                state.showChangeWarning = true
                state.pendingSelectedRequirement = getEmptyRequirement();
                return
            }
            state.selectedRequirement = getEmptyRequirement();
            state.updatedRequirementDraft = getEmptyRequirement();
            state.newRequirement = true
        },
        closeDialog: (state) => {
            state.showChangeWarning = false;
        },
        changesCommitted: (state) => {
            state.requirementChanges = false
        }
    }
})

export const { selectRequirement, newRequirement, closeDialog, changesCommitted, updateLocalSelectedRequirement } = adminRequirementSlice.actions
export const adminRequirementSelector = (state: RootState) => state.adminRequirement
export default adminRequirementSlice.reducer