<script setup lang="ts">
import type { Member } from '@/index';
import { ref } from 'vue';

const props = defineProps<{
  potentialSupervisors: Member[];
  selectedId?: string;
}>();

const emit = defineEmits<{
  supervisorUpdated: [id: string | undefined];
}>();

const supervisorID = ref(props.selectedId);
</script>

<template>
  <label for="supervisor">
    Supervisor:
    <select
      name="supervisor"
      @change="emit('supervisorUpdated', supervisorID)"
      v-model="supervisorID"
    >
      <option v-for="member in props.potentialSupervisors" :key="member.id" :value="member.id">
        {{ `${member.first_name} ${member.last_name}` }}
      </option>
    </select>
  </label>
</template>
