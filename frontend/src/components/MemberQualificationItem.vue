<script setup lang="ts">
import type {
  InitialMemberRequirement,
  MemberQualification,
  RecurringMemberRequirement,
} from '@/index';
import { computed, onMounted, ref } from 'vue';
import { useMemberStore } from '@/stores/member_store.ts';
import { useQualificationStore } from '@/stores/qualification_store.ts';
import { client } from '@/stores/client.ts';
import InitialMemberRequirementRow from '@/components/InitialMemberRequirementRow.vue';
import RecurringMemberRequirementRow from '@/components/RecurringMemberRequirementRow.vue';
import { EmptyInitialMemberRequirement, EmptyRecurringMemberRequirement } from '@/empties.ts';

const props = defineProps<{
  memberQualification: MemberQualification;
}>();

const memberStore = useMemberStore();
const qualificationStore = useQualificationStore();
const startingHeight = '50px';
const initialMemberRequirements = ref<InitialMemberRequirement[]>([]);
const recurringMemberRequirements = ref<RecurringMemberRequirement[]>([]);

onMounted(async () => {
  const { response, data, error } = await client.GET('/api/member/requirements/initial', {
    params: {
      query: {
        member_id: props.memberQualification.member_id,
        qualification_id: props.memberQualification.qualification_id,
      },
    },
  });
  if (!response.ok || error || !data) {
    console.error(error);
    return;
  }
  initialMemberRequirements.value = data;

  const res = await client.GET('/api/member/requirements/recurring', {
    params: {
      query: {
        member_id: props.memberQualification.member_id,
        qualification_id: props.memberQualification.qualification_id,
      },
    },
  });
  if (!res.response.ok || res.error || !res.data) {
    console.error(res.error);
    return;
  }
  recurringMemberRequirements.value = res.data;
});

const qualification = computed(() => {
  return qualificationStore.allQualifications.find(
    (q) => q.id === props.memberQualification.qualification_id,
  );
});

const assignedBy = computed(() => {
  return memberStore.allMembers.find((m) => m.id === props.memberQualification.assigned_by);
});
</script>

<template>
  <details :name="props.memberQualification.member_id">
    <summary class="header">
      {{ qualification?.name }}
    </summary>
    <div class="details-content">
      <p>Assigned by: {{ `${assignedBy?.first_name} ${assignedBy?.last_name}` }}</p>
      <p>
        Date Assigned:
        {{
          new Date(memberQualification.date_assigned).toLocaleString('en-US', {
            timeZone: 'America/New_York',
            dateStyle: 'short',
          })
        }}
      </p>
      Initial Requirements:
      <table>
        <thead>
          <tr class="table-heading">
            <th>Requirement</th>
            <th>Completed</th>
            <th>Completed Date</th>
            <th>Completed By</th>
          </tr>
        </thead>
        <InitialMemberRequirementRow
          v-for="ir in qualification?.initial_requirements"
          :requirement="ir"
          :member-requirement="
            initialMemberRequirements.find((i) => i.requirement_id === ir.id) ||
            EmptyInitialMemberRequirement
          "
          :key="ir.id"
        />
      </table>
      <p class="recurring-header">Recurring Requirements:</p>
      <table>
        <thead>
          <tr class="table-heading">
            <th>Requirement</th>
            <th>Current</th>
            <th>Last Completed</th>
            <th>Completed By</th>
            <th>Due</th>
          </tr>
        </thead>
        <RecurringMemberRequirementRow
          v-for="rr in qualification?.recurring_requirements"
          :requirement="rr"
          :member-requirement="
            recurringMemberRequirements.find((i) => i.requirement_id === rr.id) ||
            EmptyRecurringMemberRequirement
          "
          :key="rr.id"
        />
      </table>
    </div>
  </details>
</template>

<style scoped>
details {
  width: 100%;
  border: 1px solid rgb(var(--blue));
}

summary {
  background-color: rgba(var(--blue));
  color: white !important;
  padding-left: 1rem;
  &:focus {
    color: white !important;
  }
}

.details-content {
  padding-left: 1rem;
  display: flex;
  flex-direction: column;
}

.header {
  width: 100%;
  height: v-bind(startingHeight);
  margin-bottom: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

th {
  text-align: center;
}

.recurring-header {
  margin-top: 1rem;
}
</style>
