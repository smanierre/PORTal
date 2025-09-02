<script setup lang="ts">
import type { Member } from '@/index';
import { onMounted, ref, toRaw, watch } from 'vue';
import { useAppStore } from '@/stores/app_store.ts';
import SupervisorList from '@/components/SupervisorList.vue';
import { useMemberStore } from '@/stores/member_store.ts';

const props = defineProps<{
  selectedMember: Member;
  newMember: boolean;
}>();

const emit = defineEmits<{
  newMember: [member: Member];
  enabledChanged: [];
}>();

const appStore = useAppStore();
const memberStore = useMemberStore();
const workingMember = ref(structuredClone(toRaw(props.selectedMember)));
const password = ref('');
const confirmPassword = ref('');
const potentialSupervisors = ref<Member[]>([]);

watch(
  () => props.selectedMember,
  (newMember) => {
    workingMember.value = structuredClone(toRaw(newMember));
    updateSupervisors(newMember.grade);
  },
);

watch(
  () => workingMember.value.grade,
  (newGrade) => {
    updateSupervisors(newGrade);
  },
);

onMounted(() => {
  updateSupervisors(props.selectedMember.grade);
});

async function updateSupervisors(newGrade: string) {
  const res = await memberStore.getPotentialSupervisors(newGrade, props.selectedMember.id);

  if (res.length === 0) {
    workingMember.value.supervisor_id = '';
  }
  if (
    workingMember.value.supervisor_id === '' &&
    res.find((m) => m.id === props.selectedMember.supervisor_id)
  ) {
    workingMember.value.supervisor_id = props.selectedMember.supervisor_id;
  }
  potentialSupervisors.value = res;
}

async function updateMember() {
  // TODO: Handle an error when updating
  await memberStore.updateMember({
    ...workingMember.value,
    password: password.value,
  });
}

async function createMember() {
  const member = await memberStore.createMember({
    ...workingMember.value,
    password: password.value,
  });

  workingMember.value = member;
  emit('newMember', member);
}

function handleSubmit() {
  if (props.newMember) {
    createMember();
  } else {
    updateMember();
  }
}

async function disableMember() {
  await memberStore.disableMember(props.selectedMember.id);
  emit('enabledChanged');
}

async function enableMember() {
  await memberStore.enableMember(props.selectedMember.id);
  emit('enabledChanged');
}
</script>

<template>
  <form class="editor" @submit.prevent="handleSubmit">
    <label>ID: <input name="id" readonly class="id-input" :value="workingMember.id" /></label>
    <label
      >Rank:
      <select name="grade" v-model="workingMember.grade">
        <option
          v-for="grade in Object.entries(appStore.ranks || {})"
          :key="grade[0]"
          :value="grade[0]"
        >
          {{ grade[1] }}
        </option>
      </select>
    </label>
    <label
      >Username:
      <input name="username" :readonly="!props.newMember" v-model="workingMember.username"
    /></label>
    <label>First Name: <input name="first_name" v-model="workingMember.first_name" /></label>
    <label>Last Name: <input name="last_name" v-model="workingMember.last_name" /></label>
    <SupervisorList
      class="supervisor-list"
      v-if="potentialSupervisors.length > 0"
      :potential-supervisors="potentialSupervisors"
      :selected-id="workingMember.supervisor_id"
      @supervisor-updated="(id) => (workingMember.supervisor_id = id)"
    />
    <label>Admin: <input name="admin" type="checkbox" v-model="workingMember.admin" /></label>
    <label>Password: <input name="password" type="password" v-model="password" /></label>
    <label
      >Confirm Password: <input name="confirm_password" type="password" v-model="confirmPassword"
    /></label>
    <div class="editor-submit-container">
      <button
        v-if="workingMember.id && !workingMember.disabled"
        :disabled="!!password && !confirmPassword"
      >
        Update
      </button>
      <button
        v-if="workingMember.id && !workingMember.disabled"
        class="disable-button"
        @click="disableMember"
      >
        Disable
      </button>
      <button v-if="props.newMember">Create</button>
      <button v-if="workingMember.disabled" @click="enableMember">Enable</button>
    </div>
  </form>
</template>

<style scoped>
.editor {
  padding: 1rem;
  display: flex;
  gap: 0;
  flex-direction: column;
}

label {
  width: fit-content;
  display: flex;
  align-items: center;
  gap: 1rem;
  text-wrap: nowrap;
}

input {
  height: 2rem;
  border-radius: 0;
  font-size: 0.8rem;
  translate: 0 0.4rem;
}

.id-input {
  width: 19rem;
}

.supervisor-list {
  width: 18rem;
}

select {
  border-radius: 0;
  font-size: 0.8rem;
  translate: 0 0.4rem;
}

input[type='checkbox'] {
  width: 2rem;

  &:checked {
    background-color: rgb(var(--blue-dark));
    border: none;
  }
}

.editor-submit-container {
  display: flex;
  gap: 1rem;
}

.disable-button {
  background-color: rgb(var(--red));
  color: white;

  &:hover {
    background-color: rgb(var(--red-dark));
  }
}
</style>
