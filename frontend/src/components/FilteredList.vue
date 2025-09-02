<script setup lang="ts" generic="T extends Member | Qualification">
import type { Member, Qualification } from '@/index';
import { ref } from 'vue';
import { useAppStore } from '@/stores/app_store.ts';

const props = defineProps<{
  items: T[];
  selectedItem: T;
}>();

const emit = defineEmits<{
  filterChange: [term: string];
  selectItem: [item: T];
}>();

const appStore = useAppStore();

const filterValue = ref('');

function display(item: Member | Qualification) {
  if ('name' in item) {
    return item.name;
  } else {
    const r = Object.entries(appStore.ranks || {}).find((rank) => {
      return rank[0] === item.grade;
    });
    if (!r) {
      return `${item.first_name} ${item.last_name}`;
    }
    return `${r[1]} ${item.first_name} ${item.last_name}`;
  }
}

function clearFilter() {
  emit('filterChange', '');
  filterValue.value = '';
}
</script>

<template>
  <div class="filter-container">
    <input
      @input="emit('filterChange', filterValue)"
      v-model="filterValue"
      class="list-filter"
      name="filter"
      type="text"
    />
    <button class="filter-clear" @click="clearFilter">Clear</button>
  </div>
  <ul class="searchable-list-items">
    <li
      v-for="item in props.items"
      :key="item.id"
      :class="{ selected: props.selectedItem.id === item.id }"
      class="searchable-list-item"
      @click="emit('selectItem', item)"
    >
      {{ display(item) }}
    </li>
  </ul>
</template>

<style scoped>
.filter-container {
  display: flex;
  justify-content: end;
  margin-top: 1rem;
}

.filter-clear {
  height: 2rem;
  padding: 0 1rem;
  font-size: 0.75rem;
  vertical-align: middle;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;

  &:hover {
    background-color: rgb(var(--blue-dark));
  }
}

.list-filter {
  width: 50%;
  height: 2rem;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.searchable-list-items {
  overflow: scroll;
  height: 20rem;
  width: 75%;
  border: 1px solid black;
  padding: 0.5rem;
  font-size: 0.9rem;
}

.searchable-list-item {
  list-style: none;

  &:hover {
    cursor: pointer;
  }
  &:hover:not(.selected) {
    background-color: rgb(var(--blue-light));
  }
}

.selected {
  background-color: rgb(var(--blue-dark));
  color: white;
}
</style>
