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
  tagTypes: ["AllMembers", "Subordinates"],
  endpoints: (builder) => ({
    getMemberById: builder.query<Member, string>({
      query: (id) => ({
        url: `member/${id}`,
        method: "GET",
      }),
    }),
    getAllMembers: builder.query<Member[], void>({
      query: () => ({
        url: "members",
        method: "GET",
      }),
      providesTags: ["AllMembers"]
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
      providesTags: ["Subordinates"]
    }),
    updateMember: builder.mutation<Member, Member & { password: string }>({
      query: (member) => ({
        url: `member/${member.id}`,
        method: "PUT",
        body: member,
      }),
      invalidatesTags: ["AllMembers", "Subordinates"],

    })
  }),
});

export const {
  useGetMemberByIdQuery,
  useGetAllMembersQuery,
  useGetLoggedInMemberQuery,
  useLoginMutation,
  useGetMemberSubordinatesQuery,
  useUpdateMemberMutation
} = api;
