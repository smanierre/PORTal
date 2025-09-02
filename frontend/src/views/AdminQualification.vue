<script setup lang="ts">
import FilteredList from '@/components/FilteredList.vue';
import { computed, ref } from 'vue';
import { useQualificationStore } from '@/stores/qualification_store.ts';
import type { Qualification } from '@/index';
import { EmptyQualification } from '@/empties.ts';
import QualificationEditor from '@/components/QualificationEditor.vue';

const qualificationStore = useQualificationStore();

const filterValue = ref('');
const selectedQualification = ref<Qualification | null>(null);
const newQualification = ref(false);

const displayQualifications = computed(() => {
  return qualificationStore.allQualifications.filter((qualification) => {
    return qualification.name.toLowerCase().includes(filterValue.value.toLowerCase());
  });
});

function selectQualification(q: Qualification) {
  newQualification.value = false;
  selectedQualification.value = q;
}

function setNewQualification() {
  newQualification.value = true;
  selectedQualification.value = EmptyQualification;
}
</script>

<template>
  <div class="list-container">
    <FilteredList
      @filter-change="
        (v) => {
          filterValue = v;
        }
      "
      @select-item="selectQualification"
      :items="displayQualifications"
      :selected-item="selectedQualification || EmptyQualification"
    />
    <button class="new-member" @click="setNewQualification">New Qualification</button>
  </div>
  <QualificationEditor
    v-if="selectedQualification || newQualification"
    :selected-qualification="selectedQualification || EmptyQualification"
    :new-qualification="newQualification"
    @qualification-added="(q) => (selectedQualification = q)"
  />
</template>

<style scoped>
.list-container {
  margin-left: auto;
  display: flex;
  flex-direction: column;
  align-items: end;
}

button {
  color: rgb(var(--blue));
  border-radius: 0;
  background-color: white;
  opacity: 50%;
  border: 1px solid #646b79;
  transition: all 500ms ease;

  &:hover {
    opacity: 100%;
  }
}

.new-member {
  background-color: rgb(var(--blue));
  opacity: 100%;
  color: white;
  border-radius: 0.25rem;
  &:hover {
    background-color: rgb(var(--blue-dark));
  }
}
</style>
