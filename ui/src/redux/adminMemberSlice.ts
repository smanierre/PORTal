import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Member } from "./api";
import { RootState } from "./store";
import { getBaseUrl, getEmptyMember } from "../lib/utils";

export interface AdminState {
  selectedMember: Member | null;
  pendingSelectedMember: Member | null;
  updatedMemberDraft: Member | null;
  updatePending: boolean;
  updateError: string;
  newMember: boolean;
  memberChanges: boolean;
  showChangeWarning: boolean;
}

const initialState: AdminState = {
  selectedMember: null,
  pendingSelectedMember: null,
  updatedMemberDraft: null,
  updatePending: false,
  updateError: "",
  newMember: false,
  memberChanges: false,
  showChangeWarning: false,
};

export const adminMemberSlice = createSlice({
  name: "admin",
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
      state.newMember = false;
      state.memberChanges = false;
      state.updatedMemberDraft = action.payload.member;
      state.selectedMember = action.payload.member;
    },
    updateLocalSelectedMember: (state, action: PayloadAction<Member>) => {
      state.memberChanges = true;
      state.updatedMemberDraft = action.payload;
    },
    newMember: (state) => {
      state.selectedMember = getEmptyMember();
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
export const adminSelector = (state: RootState) => state.admin;
export default adminMemberSlice.reducer;
