<script setup lang="ts">
import type { Member, MemberQualification } from '@/index';
import { useMemberStore } from '@/stores/member_store.ts';
import { computed } from 'vue';
import MemberQualificationItem from '@/components/MemberQualificationItem.vue';
import { useQualificationStore } from '@/stores/qualification_store.ts';

const props = defineProps<{
  member: Member;
  memberQualifications: MemberQualification[];
}>();

const emit = defineEmits<{
  assignQualification: [];
}>();

const memberStore = useMemberStore();
const qualificationStore = useQualificationStore();

const supervisorName = computed(() => {
  const s = memberStore.allMembers.find((m) => m.id === props.member.supervisor_id);
  if (s) {
    return `${s.first_name} ${s.last_name}`;
  }
  return '';
});
</script>

<template>
  <div class="member-layout">
    <div class="member-info">
      <span class="member-info-header">{{
        `${props.member.first_name} ${props.member.last_name}`
      }}</span>
      <span class="member-info-content">
        Grade: {{ props.member.grade }} Supervisor:
        {{ supervisorName }}
      </span>
    </div>
    <div class="qualifications">
      Assigned Qualifications:
      <ul class="qualification-list">
        <MemberQualificationItem
          v-for="mq in props.memberQualifications"
          :key="mq.qualification_id"
          :member-qualification="mq"
        />
      </ul>
      <button
        v-if="memberQualifications.length < qualificationStore.allQualifications.length"
        class="submit"
        @click="emit('assignQualification')"
      >
        Assign Qualification
      </button>
    </div>
  </div>
</template>

<style scoped>
.member-layout {
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: 100px 50vh;
  margin: 1rem;
}

.member-info {
  display: grid;
  grid-template-rows: 50px 1fr;
  align-items: center;
}

.member-info-header {
  font-size: 2rem;
  font-weight: bold;
}

.member-info-content {
  display: flex;
}

.submit:hover {
  background-color: rgb(var(--blue-dark));
}

.qualification-list {
  margin-top: 1rem;
  height: 100%;
  overflow: scroll;
}
</style>
