import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "./store";
import { Member } from "./api";

export interface MemberState {
  member: Member | null;
  loading: boolean;
  error: string | null;
}

const initialState: MemberState = {
  member: null,
  loading: false,
  error: null,
};

export const fetchMember = createAsyncThunk(
  "member/fetchMember",
  async (id: string) => {
    const res = await fetch(`http://localhost:8080/api/member/${id}`, {
      method: "GET",
      credentials: "same-origin",
    });
    const data = (await res.json()) as Member;
    return data;
  },
);

const memberSlice = createSlice({
  name: "member",
  initialState,
  reducers: {
    update: (state, action: PayloadAction<Member>) => {
      state.member = action.payload;
    },
    logout: (state) => {
      state.member = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchMember.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchMember.fulfilled, (state, action) => {
        state.member = action.payload;
        state.loading = false;
      })
      .addCase(fetchMember.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message as string;
      });
  },
});

export const { update, logout } = memberSlice.actions;
export const memberSelector = (state: RootState) => state.member;
export default memberSlice.reducer;
