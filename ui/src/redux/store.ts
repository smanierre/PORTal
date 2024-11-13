import { configureStore } from "@reduxjs/toolkit";
import { api } from "./api";
import adminMemberSlice from "./adminMemberSlice";

export const store = configureStore({
  reducer: {
    [api.reducerPath]: api.reducer,
    admin: adminMemberSlice,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(api.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export type AppStore = typeof store;
