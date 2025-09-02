<script setup lang="ts">
import type { Qualification, Requirement } from '@/index';
import { ref, toRaw } from 'vue';
import { useQualificationStore } from '@/stores/qualification_store.ts';
import { Grades } from '@/empties.ts';

const props = defineProps<{
  requirement: Requirement;
  initial: boolean;
  parentQualificationId: string;
  assignedQualificationIds: string[];
}>();

const emit = defineEmits<{
  closeModal: [];
  changed: [];
  addedRequirement: [initial: boolean, requirement: Requirement];
  updatedRequirement: [requirement: Requirement];
}>();

const qualificationStore = useQualificationStore();

const workingRequirement = ref<Requirement>(structuredClone(toRaw(props.requirement)));
const potentialQualifications = ref<Qualification[]>(
  filterPotentialQualifications(qualificationStore.allQualifications),
);

function filterPotentialQualifications(quals: Qualification[]): Qualification[] {
  return quals.filter((q) => {
    return (
      workingRequirement.value.qualification_id === q.id ||
      (!props.assignedQualificationIds.includes(q.id) && props.parentQualificationId !== q.id)
    );
  });
}
async function handleSubmit() {
  if (!workingRequirement.value.id) {
    workingRequirement.value.initial = props.initial;
    emit('addedRequirement', props.initial, workingRequirement.value);
  } else {
    emit('updatedRequirement', workingRequirement.value);
  }
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="editor">
    <label
      >ID
      <input name="requirement_id" class="editor-input" v-model="workingRequirement.id" readonly
    /></label>
    <label
      >Type
      <select @change.once="emit('changed')" v-model="workingRequirement.type" name="type">
        <option v-if="props.initial" value="Qualification">Qualification</option>
        <option value="WBT">WBT</option>
        <option v-if="props.initial" value="Grade">Grade</option>
        <option v-if="!props.initial" value="Proficiency">Proficiency</option>
      </select>
    </label>
    <label
      for="name"
      v-if="workingRequirement.type === 'WBT' || workingRequirement.type === 'Proficiency'"
      >Name
      <input @input.once="emit('changed')" v-model="workingRequirement.name" name="name" />
    </label>
    <label for="notes"
      >Notes:
      <textarea
        @input.once="emit('changed')"
        v-model="workingRequirement.notes"
        name="notes"
      ></textarea>
    </label>
    <label
      v-if="
        (workingRequirement.type === 'WBT' && !props.initial) ||
        workingRequirement.type === 'Proficiency'
      "
      for="days_valid_for"
    >
      Days Valid For
      <input
        @input.once="emit('changed')"
        v-model="workingRequirement.days_valid_for"
        name="days_valid_for"
        type="number"
      />
    </label>
    <label v-if="workingRequirement.type === 'Grade'" for="grade">
      Grade
      <select @change.once="emit('changed')" v-model="workingRequirement.grade" name="grade">
        <option v-for="grade in Grades" :value="grade" :key="grade">{{ grade }}</option>
      </select>
    </label>
    <label v-if="workingRequirement.type === 'Qualification'" for="qualification_id">
      Qualification
      <select
        @change.once="emit('changed')"
        v-model="workingRequirement.qualification_id"
        name="qualification_id"
        :disabled="potentialQualifications.length === 0"
      >
        <option
          v-for="qualification in potentialQualifications"
          :key="qualification.id"
          :value="qualification.id"
        >
          {{ qualification.name }}
        </option>
      </select>
    </label>
    <label
      >Reference
      <input
        @input.once="emit('changed')"
        name="reference"
        v-model="workingRequirement.reference"
        placeholder="e.g. AFI 24-605V2 2.8"
      />
    </label>
    <button class="submit-create-button" v-if="props.requirement.id !== ''">Update</button>
    <button class="submit-create-button" v-else>Create</button>
    <button class="cancel-button" type="button" @click="emit('closeModal')">Cancel</button>
  </form>
</template>

<style scoped>
.editor {
  width: 75%;
  height: 75%;
  padding: 2rem 2rem 0 2rem;
  background-color: white;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr 1fr;
  align-items: center;
  justify-items: center;
  label {
    height: 2rem;
    width: 80%;
  }
}

.submit-create-button {
  grid-column-start: 1;
  justify-self: end;
  margin-right: 1rem;

  &:hover {
    background-color: rgb(var(--blue-dark));
  }
}

.cancel-button {
  margin-left: 1rem;
  grid-column-start: 2;
  justify-self: start;
  margin: 0;

  background-color: rgb(var(--red));
  border: 1px solid rgb(var(--red));
  &:hover {
    background-color: rgb(var(--red-dark));
  }
}
</style>
