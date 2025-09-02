<script setup lang="ts">
import type { DropdownValue } from '..';

const props = defineProps<{
  selectedValue: DropdownValue;
}>();

const emit = defineEmits<{
  toggleSelected: [value: DropdownValue];
}>();
</script>

<template>
  <div class="toggle-container">
    <button
      :class="{ selected: props.selectedValue === 'Members' }"
      class="members-button"
      @click="emit('toggleSelected', 'Members')"
    >
      Members
    </button>
    <button
      :class="{ selected: props.selectedValue === 'Qualifications' }"
      class="qualifications-button"
      @click="emit('toggleSelected', 'Qualifications')"
    >
      Qualifications
    </button>
  </div>
</template>

<style scoped>
.toggle-container {
  margin-top: 1rem;
  margin-left: 1rem;
  display: flex;
  flex-direction: row;
  justify-content: left;
}

.members-button {
  border-top-left-radius: 0.25rem;
  border-bottom-left-radius: 0.25rem;
}

.qualifications-button {
  border-top-right-radius: 0.25rem;
  border-bottom-right-radius: 0.25rem;
}

.selected {
  background-color: rgb(var(--blue));
  opacity: 100%;
  color: white;
}

button {
  color: rgb(var(--blue));
  border-radius: 0;
  background-color: white;
  opacity: 50%;
  transition: all 500ms ease;
  font-size: 0.9rem;

  &:hover {
    opacity: 100%;
  }
}
</style>
