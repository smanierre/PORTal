<script setup lang="ts">
import type { Qualification, Requirement } from '@/index';
import { ref, toRaw, watch } from 'vue';
import RequirementsSelector from '@/components/RequirementsSelector.vue';
import { useQualificationStore } from '@/stores/qualification_store.ts';

const props = defineProps<{
  selectedQualification: Qualification;
  newQualification: boolean;
}>();

const emit = defineEmits<{
  qualificationAdded: [qualification: Qualification];
}>();

watch(
  () => props.selectedQualification,
  (newQual) => {
    workingQualification.value = structuredClone(toRaw(newQual));
  },
);

const qualificationStore = useQualificationStore();
const workingQualification = ref(structuredClone(toRaw(props.selectedQualification)));

function handleAddedRequirement(initial: boolean, requirement: Requirement) {
  if (initial) {
    workingQualification.value.initial_requirements = [
      ...(workingQualification.value.initial_requirements || []),
      requirement,
    ];
  } else {
    workingQualification.value.recurring_requirements = [
      ...(workingQualification.value.recurring_requirements || []),
      requirement,
    ];
  }
}

function handleUpdatedRequirement(requirement: Requirement) {
  if (requirement.initial) {
    workingQualification.value.initial_requirements =
      workingQualification.value.initial_requirements?.map((req) => {
        if (req.id === requirement.id) {
          return requirement;
        } else {
          return req;
        }
      }) || [];
  } else {
    workingQualification.value.recurring_requirements =
      workingQualification.value.recurring_requirements?.map((req) => {
        if (req.id === requirement.id) {
          return requirement;
        } else {
          return req;
        }
      }) || [];
  }
}

function handleRemovedRequirement(requirement: Requirement) {
  if (requirement.initial) {
    workingQualification.value.initial_requirements =
      workingQualification.value.initial_requirements?.filter((req) => req.id !== requirement.id);
  } else {
    workingQualification.value.recurring_requirements =
      workingQualification.value.recurring_requirements?.filter((req) => req.id !== requirement.id);
  }
}

async function handleSubmit() {
  if (props.newQualification) {
    const data = await qualificationStore.createQualification(workingQualification.value);
    workingQualification.value = data;
    emit('qualificationAdded', data);
  } else {
    workingQualification.value = await qualificationStore.updateQualification(
      workingQualification.value,
    );
  }
}
</script>

<template>
  <form class="editor" @submit.prevent="handleSubmit">
    <label
      >ID: <input name="id" readonly class="id-input" :value="workingQualification.id"
    /></label>
    <label
      >Name:
      <input name="name" v-model="workingQualification.name" />
    </label>
    <span class="requirements">
      <label class="requirements-label"
        >Initial Requirements
        <RequirementsSelector
          :requirements="workingQualification.initial_requirements || []"
          :initial="true"
          :parent-qualification-id="workingQualification.id"
          :assigned-qualification-ids="
            workingQualification.initial_requirements?.map((q) => q.qualification_id || '') || []
          "
          @added-requirement="handleAddedRequirement"
          @updated-requirement="handleUpdatedRequirement"
          @removed-requirement="handleRemovedRequirement"
        />
      </label>
      <label class="requirements-label"
        >Recurring Requirements
        <RequirementsSelector
          :requirements="workingQualification.recurring_requirements || []"
          :initial="false"
          :parent-qualification-id="workingQualification.id"
          :assigned-qualification-ids="
            workingQualification.recurring_requirements?.map((q) => q.id) || []
          "
          @added-requirement="handleAddedRequirement"
          @updated-requirement="handleUpdatedRequirement"
          @removed-requirement="handleRemovedRequirement"
        />
      </label>
    </span>
    <label>Notes: <textarea name="notes" v-model="workingQualification.notes" /></label>
    <div class="editor-submit-container">
      <button v-if="workingQualification.id">Update</button>
      <button v-if="props.newQualification">Create</button>
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

.requirements {
  display: flex;
  gap: 2rem;
}

.requirements-label {
  display: flex;
  flex-direction: column;
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

  button:hover {
    background-color: rgb(var(--blue-dark));
  }
}
</style>
