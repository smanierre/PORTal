import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Member } from "./api";
import { RootState } from "./store";
import { getEmptyMember } from "../lib/utils";

export interface AdminMemberState {
  selectedMember: Member | null;
  pendingSelectedMember: Member | null;
  updatedMemberDraft: Member | null;
  updatePending: boolean;
  updateError: string;
  newMember: boolean;
  memberChanges: boolean;
  showChangeWarning: boolean;
}

const initialState: AdminMemberState = {
  selectedMember: null,
  pendingSelectedMember: null,
  updatedMemberDraft: null,
  updatePending: false,
  updateError: "",
  newMember: false,
  memberChanges: false,
  showChangeWarning: false,
};

//TODO: Fix logic for switching from edited member to new member in all reducers

export const adminMemberSlice = createSlice({
  name: "adminMember",
  initialState,
  reducers: {
    selectMember: (state, action: PayloadAction<{ member: Member | null, force?: boolean }>) => {
      if (state.memberChanges && !action.payload.force) {
        state.pendingSelectedMember = action.payload.member
        state.showChangeWarning = true;
        return;
      }
      if (action.payload.force) {
        state.updatedMemberDraft = state.pendingSelectedMember;
        state.selectedMember = state.pendingSelectedMember;
        state.pendingSelectedMember = null;
        state.showChangeWarning = false;
        state.memberChanges = false;
        return
      }
      if (state.newMember && state.memberChanges) {
        state.newMember = true
      } else {
        state.newMember = false;
      }
      state.memberChanges = false;
      state.updatedMemberDraft = action.payload.member;
      state.selectedMember = action.payload.member;
    },
    updateLocalSelectedMember: (state, action: PayloadAction<Member>) => {
      state.memberChanges = true;
      state.updatedMemberDraft = action.payload;
    },
    newMember: (state) => {
      if (state.memberChanges) {
        state.showChangeWarning = true
        state.pendingSelectedMember = getEmptyMember();
        return
      }
      state.selectedMember = getEmptyMember();
      state.updatedMemberDraft = getEmptyMember();
      state.newMember = true
    },
    closeDialog: (state) => {
      state.showChangeWarning = false;
    },
    changesCommitted: (state) => {
      state.memberChanges = false
    }
  },
});

export const {
  selectMember,
  newMember,
  updateLocalSelectedMember,
  closeDialog,
  changesCommitted
} = adminMemberSlice.actions;
export const adminMemberSelector = (state: RootState) => state.adminMember;
export default adminMemberSlice.reducer;
