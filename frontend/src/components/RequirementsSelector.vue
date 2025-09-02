<script lang="ts" setup>
import type { Requirement } from '@/index';
import { ref } from 'vue';
import RequirementEditor from '@/components/RequirementEditor.vue';
import { EmptyRequirement } from '@/empties.ts';
import { useQualificationStore } from '@/stores/qualification_store.ts';

const props = defineProps<{
  requirements: Requirement[];
  initial: boolean;
  parentQualificationId: string;
  assignedQualificationIds: string[];
}>();

const emit = defineEmits<{
  addedRequirement: [initial: boolean, requirement: Requirement];
  updatedRequirement: [requirement: Requirement];
  removedRequirement: [requirement: Requirement];
}>();

const qualificationStore = useQualificationStore();

const open = ref(false);
const requirement = ref<Requirement>(EmptyRequirement);
const changes = ref(false);

function handleNewRequirement() {
  requirement.value = EmptyRequirement;
  open.value = true;
  changes.value = false;
}

function selectRequirement(r: Requirement) {
  requirement.value = r;
  open.value = true;
  changes.value = false;
}

function handleModalClose() {
  if (changes.value) {
    const res = confirm('Changes are pending, are you sure you want to exit?');
    if (!res) {
      return;
    }
  }
  open.value = false;
}

function handleAddedRequirement(initial: boolean, requirement: Requirement) {
  emit('addedRequirement', initial, requirement);
  open.value = false;
}

function handleUpdatedRequirement(requirement: Requirement) {
  emit('updatedRequirement', requirement);
  open.value = false;
}
</script>

<template>
  <div class="requirement-items-container">
    <ul class="requirement-items">
      <li
        v-for="requirement in props.requirements"
        @click.self="selectRequirement(requirement)"
        :key="requirement.id"
        class="requirement-item"
      >
        {{
          requirement.type === 'Qualification'
            ? 'Qualification: ' +
                qualificationStore.allQualifications.filter((q) => {
                  return q.id === requirement.qualification_id;
                })[0].name || 'Unknown Qualification'
            : requirement.type === 'Grade'
              ? 'Grade: ' + requirement.grade
              : requirement.name
        }}
        <span @click="emit('removedRequirement', requirement)">X</span>
      </li>
      <li @click="handleNewRequirement">Add Requirement...</li>
    </ul>
    <Teleport to="body">
      <div v-if="open" @click.self="open = false" class="modal-container">
        <RequirementEditor
          @changed="changes = true"
          @close-modal="handleModalClose"
          :requirement="requirement"
          :initial="props.initial"
          :parentQualificationId="props.parentQualificationId"
          :assignedQualificationIds="props.assignedQualificationIds"
          @added-requirement="handleAddedRequirement"
          @updated-requirement="handleUpdatedRequirement"
        />
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.requirement-items-container {
  min-width: 250px;
  border: 1px solid black;
  min-height: 250px;
  overflow-y: scroll;
  height: fit-content;
  display: flex;
  gap: 3px;
  align-items: center;
  padding-left: 5px;
  padding-right: 5px;
}

.requirement-items {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: start;
  padding-top: 1rem;
  padding-left: 0;
  li {
    list-style: none;
    &:hover {
      background-color: rgb(var(--blue-light));
      cursor: pointer;
    }
  }
}

.requirement-item {
  display: flex;
  justify-content: space-between;

  span {
    margin-left: 0.5rem;
    margin-right: 0.5rem;
    &:hover {
      color: rgb(var(--red));
    }
  }
}

.modal-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
