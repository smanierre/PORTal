<script setup lang="ts">
import FilteredList from '@/components/FilteredList.vue';
import { computed, ref } from 'vue';
import { useAppStore } from '@/stores/app_store.ts';
import AdminMemberEditor from '@/components/AdminMemberEditor.vue';
import type { Member } from '@/index';
import { EmptyMember } from '@/empties.ts';
import { useMemberStore } from '@/stores/member_store.ts';

const disabled = ref(false);
const memberStore = useMemberStore();
const appStore = useAppStore();

const filterValue = ref('');
const selectedMember = ref<Member | null>(null);
const newMember = ref(false);

const displayMembers = computed(() => {
  return memberStore.allMembers.filter((member) => {
    const r = Object.entries(appStore.ranks || {}).find((rank) => {
      return rank[0] === member.grade;
    });
    if (!r) {
      return (
        member.disabled === disabled.value &&
        `${member.first_name} ${member.last_name}`
          .toLowerCase()
          .includes(filterValue.value.toLowerCase())
      );
    }
    return (
      member.disabled === disabled.value &&
      `${r[1]} ${member.first_name} ${member.last_name}`
        .toLowerCase()
        .includes(filterValue.value.toLowerCase())
    );
  });
});

function selectMember(m: Member) {
  newMember.value = false;
  selectedMember.value = m;
}

function setNewMember() {
  newMember.value = true;
  selectedMember.value = EmptyMember;
}
</script>

<template>
  <div class="list-container">
    <div class="toggle-container">
      <button
        :class="{ selected: disabled === false }"
        class="enabled-button"
        @click="disabled = false"
      >
        Enabled
      </button>
      <button
        :class="{ selected: disabled === true }"
        class="disabled-button"
        @click="disabled = true"
      >
        Disabled
      </button>
    </div>
    <FilteredList
      @filter-change="
        (v) => {
          filterValue = v;
        }
      "
      @select-item="selectMember"
      :items="displayMembers"
      :selected-item="selectedMember || EmptyMember"
    />
    <button class="new-member" @click="setNewMember">New Member</button>
  </div>
  <AdminMemberEditor
    v-if="selectedMember || newMember"
    :selected-member="selectedMember || EmptyMember"
    :new-member="newMember"
    @new-member="(m) => (selectedMember = m)"
    @enabled-changed="selectedMember = null"
  />
</template>

<style scoped>
.list-container {
  margin-left: auto;
  display: flex;
  flex-direction: column;
  align-items: end;
}

.toggle-container {
  margin-top: 1rem;
  margin-left: 1rem;
  display: flex;
  flex-direction: row;
  justify-content: right;
}

.enabled-button {
  border-top-left-radius: 0.25rem;
  border-bottom-left-radius: 0.25rem;
}

.disabled-button {
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
