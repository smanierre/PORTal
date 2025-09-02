import type { Member, MemberQualification } from '@/index';
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { client } from '@/stores/client.ts';
import { EmptyMember, EmptyMemberQualification } from '@/empties.ts';
import { useAppStore } from '@/stores/app_store.ts';

export const useMemberStore = defineStore('member', () => {
  const appStore = useAppStore();

  const allMembers = ref<Member[]>([]);
  // Get members on store load
  getAllMembers();

  function getAllMembers() {
    client.GET('/api/members').then((res) => {
      if (!res.response.ok || res.error) {
        console.error(res.error);
        return;
      }
      allMembers.value = res.data;
    });
  }

  const allEnabledMembers = computed(() => {
    return allMembers.value.filter((m) => !m.disabled);
  });

  const allDisabledMembers = computed(() => {
    return allMembers.value.filter((m) => m.disabled);
  });

  async function createMember(member: Member & { password: string }): Promise<Member> {
    const { response, data, error } = await client.POST('/api/member', {
      body: {
        ...member,
      },
    });
    if (!response.ok || error) {
      console.error(error);
      return EmptyMember;
    }
    allMembers.value.push(data);
    return data;
  }

  async function updateMember(member: Member & { password: string }): Promise<Member> {
    const { response, data, error } = await client.PUT('/api/member', {
      body: {
        ...member,
      },
    });
    if (!response.ok || error) {
      console.error(error);
      return EmptyMember;
    }
    allMembers.value.forEach((member: Member, index) => {
      if (member.id === data?.id) {
        allMembers.value[index] = data;
      }
    });
    return data || EmptyMember;
  }

  async function getPotentialSupervisors(grade: string, memberID?: string): Promise<Member[]> {
    const { response, data, error } = await client.GET('/api/potentialSupervisors', {
      params: {
        query: {
          id: memberID,
          grade: grade,
        },
      },
    });
    if (!response.ok || error) {
      console.error(error);
      return [];
    }
    return data || [];
  }

  async function disableMember(memberID: string) {
    const { error } = await client.POST('/api/disableMember', {
      params: {
        query: {
          id: memberID,
        },
      },
    });
    if (error) {
      console.log(error);
    }
    getAllMembers();
  }

  async function enableMember(memberID: string) {
    const { error } = await client.POST('/api/enableMember', {
      params: {
        query: {
          id: memberID,
        },
      },
    });
    if (error) {
      console.log(error);
    }
    getAllMembers();
  }

  async function assignQualification(
    memberID: string,
    qualificationID: string,
  ): Promise<MemberQualification> {
    const { response, data, error } = await client.POST('/api/member/qualification', {
      params: {
        query: {
          member_id: memberID,
          qualification_id: qualificationID,
          assigned_by: appStore.member?.id || '',
        },
      },
    });
    if (!response.ok || error || !data) {
      console.error(error);
      return EmptyMemberQualification;
    }
    return data;
  }

  return {
    allMembers,
    allEnabledMembers,
    allDisabledMembers,
    createMember,
    updateMember,
    getPotentialSupervisors,
    disableMember,
    enableMember,
    assignQualification,
  };
});
