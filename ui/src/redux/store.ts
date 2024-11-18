import { configureStore } from "@reduxjs/toolkit";
import { api } from "./api";
import adminMemberSlice from "./adminMemberSlice";
import adminQualificationSlice from "./adminQualificationSlice"
import adminRequirementSlice from "./adminRequirementSlice"
import adminReferenceSlice from "./adminReferenceSlice"

export const store = configureStore({
  reducer: {
    [api.reducerPath]: api.reducer,
    adminMember: adminMemberSlice,
    adminQualification: adminQualificationSlice,
    adminRequirement: adminRequirementSlice,
    adminReference: adminReferenceSlice,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export type AppStore = typeof store;
