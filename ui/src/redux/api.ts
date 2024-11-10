import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { getBaseUrl } from "../lib/utils";

export interface Member {
  id: string;
  first_name: string;
  last_name: string;
  username: string;
  grade: string;
  admin: boolean;
  supervisor_id: string;
}

export const api = createApi({
  reducerPath: "member",
  baseQuery: fetchBaseQuery({
    baseUrl: `${getBaseUrl()}/api/`,
    credentials: import.meta.env.DEV ? "include" : "same-origin",
  }),
  endpoints: (builder) => ({
    getMemberById: builder.query<Member, string>({
      query: (id) => ({
        url: `member/${id}`,
        method: "GET",
      }),
    }),
    addMember: builder.mutation<Member, Member>({
      query: (member) => ({
        url: "member",
        method: "POST",
        body: member,
      }),
    }),
    getAllMembers: builder.query<Member[], void>({
      query: () => ({
        url: "members",
        method: "GET",
      }),
    }),
    updateMember: builder.query<Member, Member>({
      query: (member) => ({
        url: `member/${member.id}`,
        method: "PUT",
        body: member,
      }),
    }),
    deleteMember: builder.query<void, string>({
      query: (id) => ({
        url: `member/${id}`,
        method: "DELETE",
      }),
    }),
    getLoggedInMember: builder.query<Member, void>({
      query: () => ({
        url: "member",
        method: "GET",
      }),
    }),
    login: builder.mutation<void, { username: string; password: string }>({
      query: (creds) => ({
        url: "login",
        method: "POST",
        body: creds,
      }),
    }),
    getMemberSubordinates: builder.query<Member[], string>({
      query: (id) => ({
        url: `member/${id}/subordinates`,
        method: "GET",
      }),
    }),
  }),
});

export const {
  useGetMemberByIdQuery,
  useAddMemberMutation,
  useGetAllMembersQuery,
  useUpdateMemberQuery,
  useDeleteMemberQuery,
  useGetLoggedInMemberQuery,
  useLoginMutation,
  useGetMemberSubordinatesQuery,
} = api;
