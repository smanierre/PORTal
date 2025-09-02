import { defineStore } from 'pinia';
import { client } from '@/stores/client.ts';
import { ref } from 'vue';
import type { MemberQualification, Qualification } from '@/index';
import { EmptyQualification } from '@/empties.ts';

export const useQualificationStore = defineStore('qualification', () => {
  const allQualifications = ref<Qualification[]>([]);
  // Get qualifications on store load
  getQualifications();

  function getQualifications() {
    client.GET('/api/qualifications').then((res) => {
      if (!res.response.ok || res.error) {
        console.error(res.error);
        return;
      }
      allQualifications.value = res.data;
    });
  }

  async function createQualification(q: Qualification): Promise<Qualification> {
    const { response, data, error } = await client.POST('/api/qualification', {
      body: q,
    });
    if (!response.ok || error || !data) {
      console.error(error);
      return EmptyQualification;
    }
    allQualifications.value.push(data);
    return data;
  }

  async function updateQualification(q: Qualification): Promise<Qualification> {
    const { response, data, error } = await client.PUT('/api/qualification', {
      body: q,
    });
    if (!response.ok || error || !data) {
      console.error(error);
      return EmptyQualification;
    }
    allQualifications.value.forEach((qual, index) => {
      if (qual.id === q.id) {
        allQualifications.value[index] = data;
      }
    });
    return data;
  }

  async function getMemberQualifications(memberID: string): Promise<MemberQualification[]> {
    const { response, data, error } = await client.GET('/api/member/qualifications', {
      params: {
        query: {
          member_id: memberID,
        },
      },
    });
    if (!response.ok || error || !data) {
      console.error(error);
      return [];
    }
    return data;
  }

  return { allQualifications, createQualification, updateQualification, getMemberQualifications };
});
