import type {
  InitialMemberRequirement,
  Member,
  MemberQualification,
  Qualification,
  RecurringMemberRequirement,
  Requirement,
} from '@/index';

export const EmptyMember: Member = {
  id: '',
  first_name: '',
  last_name: '',
  username: '',
  supervisor_id: '',
  grade: '',
  admin: false,
  disabled: false,
};

export const EmptyQualification: Qualification = {
  id: '',
  name: '',
  notes: '',
  initial_requirements: [],
  recurring_requirements: [],
};

export const EmptyRequirement: Requirement = {
  id: '',
  name: '',
  notes: '',
  grade: 'E1',
  type: 'WBT',
  initial: false,
  days_valid_for: 0,
  reference: '',
  qualification_id: '',
};

export const EmptyMemberQualification: MemberQualification = {
  member_id: '',
  qualification_id: '',
  date_assigned: '',
  assigned_by: '',
};

export const EmptyInitialMemberRequirement: InitialMemberRequirement = {
  member_id: '',
  assigned_by: '',
  requirement_id: '',
};

export const EmptyRecurringMemberRequirement: RecurringMemberRequirement = {
  member_id: '',
  requirement_id: '',
  id: '',
  assigned_by: '',
  completion_history: [],
};

export const Grades = ['E1', 'E2', 'E3', 'E4', 'E5', 'E6', 'E7', 'E8', 'E9'];
