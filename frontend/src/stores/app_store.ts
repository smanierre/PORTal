import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import type { Member, Ranks, Error } from '@/index';
import { MEMBER_KEY, ORGANIZATION_KEY, RANKS_KEY } from '@/stores/local_storage.ts';
import { client } from '@/stores/client.ts';
import { useMemberStore } from '@/stores/member_store.ts';

export const useAppStore = defineStore('app', () => {
  const memberStore = useMemberStore();
  /************************ Logged in member ************************/
  const member = ref<Member | null>(null);
  // If a member is logged in on load, it will be stored in Local storage
  const lsMember = localStorage.getItem(MEMBER_KEY);
  if (lsMember) {
    member.value = JSON.parse(lsMember);
  }

  /************************ Subordinates ************************/
  const subordinates = computed(() => {
    if (member.value === null) {
      return [];
    }
    return memberStore.allEnabledMembers.filter((m) => m.supervisor_id === member.value.id);
  });

  /************************ Organization ************************/
  const organization = ref(localStorage.getItem(ORGANIZATION_KEY) || '');
  if (!organization.value) {
    client.GET('/api/organization').then((res) => {
      if (!res.data) {
        console.error(res.error);
        return;
      }
      // Sync organization to local storage
      localStorage.setItem(ORGANIZATION_KEY, res.data.organization);
      organization.value = res.data.organization;
    });
  }

  /************************ Ranks ************************/
  const ranks = ref<Ranks | null>(null);
  const lsRanks = localStorage.getItem(RANKS_KEY);
  if (lsRanks) {
    ranks.value = JSON.parse(lsRanks);
  } else {
    client.GET('/api/ranks').then((res) => {
      if (!res.response.ok || res.error) {
        console.error(res.error);
        return;
      }
      ranks.value = res.data;
      localStorage.setItem(RANKS_KEY, JSON.stringify(res.data));
    });
  }

  /************************ Logout ************************/
  async function login(username: string, password: string): Promise<Error | null> {
    const { data, error } = await client.POST('/api/login', {
      body: {
        username: username,
        password: password,
      },
    });
    if (error) {
      console.error(error);
      return error;
    }
    localStorage.setItem(MEMBER_KEY, JSON.stringify(data));
    member.value = data;
    return null;
  }

  /************************ Logout ************************/
  async function logout() {
    //Clear out members local storage
    localStorage.removeItem(MEMBER_KEY);
    member.value = null;
    const { response, error } = await client.POST('/api/logout');
    if (!response.ok || error) {
      console.error(error);
    }
  }
  return { member, subordinates, organization, ranks, login, logout };
});
