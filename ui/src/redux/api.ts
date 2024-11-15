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

export interface Qualification {
  id: string;
  name: string;
  initial_requirements: Requirement[];
  recurring_requirements: Requirement[];
  notes: string;
  expires: boolean;
  expiration_days: number;
}

export interface Requirement {
  id: string;
  name: string;
  reference: Reference;
  description: string;
  notes: string;
  days_valid_for: number;
}

export interface Reference {
  id: string;
  name: string;
  volume: number;
  paragraph: string;
}

export const api = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({
    baseUrl: `${getBaseUrl()}/api/`,
    credentials: import.meta.env.DEV ? "include" : "same-origin",
  }),
  tagTypes: ["AllMembers", "Subordinates", "AllQuals", "AllRequirements", "AllReferences"],
  endpoints: (builder) => ({
    // Member Queries
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
    getMemberSubordinates: builder.query<Member[], string>({
      query: (id) => ({
        url: `member/${id}/subordinates`,
        method: "GET",
      }),
      providesTags: ["Subordinates"]
    }),

    //Member Mutations
    login: builder.mutation<void, { username: string; password: string }>({
      query: (creds) => ({
        url: "login",
        method: "POST",
        body: creds,
      }),
    }),
    updateMember: builder.mutation<Member, Member & { password: string }>({
      query: (member) => ({
        url: `member/${member.id}`,
        method: "PUT",
        body: member,
      }),
      invalidatesTags: ["AllMembers", "Subordinates"],
    }),
    createMember: builder.mutation<Member, Member & { password: string }>({
      query: (member) => ({
        url: "member",
        method: "POST",
        body: member
      }),
      invalidatesTags: ["AllMembers", "Subordinates"],
    }),

    // Qualification Queries
    getAllQualifications: builder.query<Qualification[], void>({
      query: () => ({
        url: "qualifications",
        method: "GET"
      }),
      providesTags: ["AllQuals"]
    }),
    getQualificationById: builder.query<Qualification, string>({
      query: (id) => ({
        url: `qualification/${id}`,
        method: "GET"
      })
    }),

    // Qualificaiton Mutations
    updateQualification: builder.mutation<Qualification, Qualification>({
      query: (qualification) => ({
        url: `qualification/${qualification.id}`,
        method: "PUT",
        body: qualification
      }),
      invalidatesTags: ["AllQuals"]
    }),
    createQualification: builder.mutation<Qualification, Qualification>({
      query: (qualification) => ({
        url: "qualification",
        method: "POST",
        body: qualification
      }),
      invalidatesTags: ["AllQuals"]
    }),

    // Requirement Queries
    getAllRequirements: builder.query<Requirement[], void>({
      query: () => ({
        url: "requirements",
        method: "GET",
      }),
      providesTags: ["AllRequirements"]
    }),
    getRequirementById: builder.query<Requirement, string>({
      query: (id) => ({
        url: `requirement/${id}`,
        method: "GET",
      })
    }),

    // Requirement Mutations
    updateRequirement: builder.mutation<Requirement, Requirement>({
      query: (requirement) => ({
        url: `requirement/${requirement.id}`,
        method: "PUT",
        body: requirement,
      }),
      invalidatesTags: ["AllRequirements"]
    }),
    createRequirement: builder.mutation<Requirement, Requirement>({
      query: (requirement) => ({
        url: "requirement",
        method: "POST",
        body: requirement,
      }),
      invalidatesTags: ["AllRequirements"],
    }),

    // Reference Queries
    getAllReferences: builder.query<Reference[], void>({
      query: () => ({
        url: "references",
        method: "GET",
      }),
      providesTags: ["AllReferences"],
    }),
    getReferenceById: builder.query<Reference, string>({
      query: (id) => ({
        url: `reference/${id}`,
        method: "GET",
      })
    }),

    // Reference Mutations
    createReference: builder.mutation<Reference, Reference>({
      query: (reference) => ({
        url: "reference",
        method: "POST",
        body: reference,
      }),
      invalidatesTags: ["AllReferences"],
    }),
    updateReference: builder.mutation<Reference, Reference>({
      query: (reference) => ({
        url: `reference/${reference.id}`,
        method: "PUT",
        body: reference,
      }),
      invalidatesTags: ["AllReferences"]
    }),
  }),
});

export const {
  useGetMemberByIdQuery,
  useGetAllMembersQuery,
  useGetLoggedInMemberQuery,
  useGetMemberSubordinatesQuery,

  useLoginMutation,
  useCreateMemberMutation,
  useUpdateMemberMutation,

  useGetAllQualificationsQuery,
  useGetQualificationByIdQuery,

  useUpdateQualificationMutation,
  useCreateQualificationMutation,

  useGetAllRequirementsQuery,
  useGetRequirementByIdQuery,

  useUpdateRequirementMutation,
  useCreateRequirementMutation,

  useGetAllReferencesQuery,
  useGetReferenceByIdQuery,

  useCreateReferenceMutation,
  useUpdateReferenceMutation,
} = api;
