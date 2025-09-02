<script setup lang="ts">
import { ref } from 'vue';
import AdminToggle from '@/components/AdminToggle.vue';
import type { DropdownValue } from '..';
import AdminMember from './AdminMember.vue';
import AdminQualification from '@/views/AdminQualification.vue';

const selectedPane = ref<DropdownValue>('Members');

function toggle(value: DropdownValue) {
  selectedPane.value = value;
}
</script>

<template>
  <div class="admin-layout">
    <AdminToggle @toggle-selected="toggle" :selected-value="selectedPane" />
    <div class="content">
      <AdminMember v-if="selectedPane === 'Members'" />
      <AdminQualification v-if="selectedPane === 'Qualifications'" />
    </div>
  </div>
</template>

<style scoped>
.admin-layout {
  display: grid;
  height: 100%;
  width: 100%;
  grid-template-rows: 75px 1fr;
}

.content {
  display: grid;
  grid-template-columns: 25% 1fr;
}
</style>
