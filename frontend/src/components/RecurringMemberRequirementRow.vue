<script setup lang="ts">
import { useQualificationStore } from '@/stores/qualification_store.ts';
import type { RecurringMemberRequirement, Requirement } from '@/index';
import { computed } from 'vue';
import { useMemberStore } from '@/stores/member_store.ts';

const props = defineProps<{
  memberRequirement: RecurringMemberRequirement;
  requirement: Requirement;
}>();

const qualificationStore = useQualificationStore();
const memberStore = useMemberStore();

const formattedDate = computed(() => {
  if (!props.memberRequirement.completed_date) {
    return '';
  }
  return new Date(props.memberRequirement.completed_date).toLocaleString('en-US', {
    timeZone: 'America/New_York',
    dateStyle: 'short',
  });
});

const completedBy = computed(() => {
  const m = memberStore.allEnabledMembers.find(
    (m) => m.id === props.memberRequirement.completed_by,
  );
  if (m) {
    return `${m.first_name} ${m.last_name}`;
  }
  return '';
});

function displayRequirement(r: Requirement): string {
  switch (r.type) {
    case 'Qualification':
      const qual = qualificationStore.allQualifications.find((q) => q.id === r.qualification_id);
      return `Qualification: ${qual?.name}`;
    case 'Grade':
      return `Grade: ${r.grade}`;
    default:
      return r.name;
  }
}
</script>

<template>
  <tr>
    <td class="requirement-name">{{ displayRequirement(props.requirement) }}</td>
    <td>
      {{ memberRequirement.completed_date ? '✅' : '❌' }}
    </td>
    <td>
      {{ formattedDate }}
    </td>
    <td>
      {{ completedBy }}
    </td>
  </tr>
</template>

<style scoped>
td {
  text-align: center;
}

.requirement-name {
  text-align: left;
}
</style>
