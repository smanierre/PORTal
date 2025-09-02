<script lang="ts" setup>
import { useAppStore } from '@/stores/app_store.ts';
import FilteredList from '@/components/FilteredList.vue';
import { EmptyMember } from '@/empties.ts';
import { computed, ref, watch } from 'vue';
import MemberEditor from '@/components/MemberEditor.vue';
import type { InitialMemberRequirement, Member, MemberQualification } from '@/index';
import { useMemberStore } from '@/stores/member_store.ts';
import { useQualificationStore } from '@/stores/qualification_store.ts';

const appStore = useAppStore();
const memberStore = useMemberStore();
const qualificationStore = useQualificationStore();

const filterValue = ref('');
const selectedMember = ref<Member | null>(null);
const memberQualifications = ref<MemberQualification[] | null>(null);
const showModal = ref(false);
const qualificationToAssign = ref('');

watch(
  () => selectedMember.value,
  async (newValue) => {
    memberQualifications.value = await qualificationStore.getMemberQualifications(
      newValue?.id || '',
    );
  },
);

const displayItems = computed(() => {
  return appStore.subordinates.filter((member) => {
    const r = Object.entries(appStore.ranks || {}).find((rank) => {
      return rank[0] === member.grade;
    });
    if (!r) {
      return `${member.first_name} ${member.last_name}`
        .toLowerCase()
        .includes(filterValue.value.toLowerCase());
    }
    return `${r[1]} ${member.first_name} ${member.last_name}`
      .toLowerCase()
      .includes(filterValue.value.toLowerCase());
  });
});

async function assignQualification() {
  const data = await memberStore.assignQualification(
    selectedMember.value?.id || '',
    qualificationToAssign.value,
  );
  if (memberQualifications.value) {
    memberQualifications.value.push(data);
  } else {
    memberQualifications.value = [data];
  }
  showModal.value = false;
}
</script>

<template>
  <div class="layout">
    <div class="list-container">
      <FilteredList
        @filter-change="(v) => (filterValue = v)"
        @select-item="(item) => (selectedMember = item)"
        :items="displayItems"
        :selected-item="selectedMember || EmptyMember"
      />
    </div>
    <MemberEditor
      v-if="selectedMember"
      :member="selectedMember"
      :memberQualifications="memberQualifications || []"
      @assign-qualification="showModal = true"
    />
    <Teleport to="body" v-if="showModal">
      <div class="modal">
        <form @submit.prevent="assignQualification" class="assign-qualification-form">
          <label class="form-content">
            Qualification:
            <select v-model="qualificationToAssign">
              <option
                v-for="qualification in qualificationStore.allQualifications.filter((q) => {
                  let found = false;
                  memberQualifications?.forEach((mq) => {
                    if (mq.qualification_id === q.id) {
                      found = true;
                    }
                  });
                  return !found;
                })"
                :value="qualification.id"
                :key="qualification.id"
              >
                {{ qualification.name }}
              </option>
            </select>
            <button class="submit">Assign</button>
            <button class="cancel" @click="showModal = false" type="button">Cancel</button>
          </label>
        </form>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: 25% 1fr;
  grid-template-rows: 1fr;
}

.list-container {
  display: flex;
  flex-direction: column;
  align-items: end;
  margin-left: auto;
}

.modal {
  position: absolute;
  height: 100vh;
  width: 100vw;
  background-color: rgba(0, 0, 0, 0.5);
  top: 0;
  left: 0;
}

.assign-qualification-form {
  position: relative;
  top: 50%;
  left: 50%;
  height: 25%;
  width: 50%;
  translate: -50% -50%;
  background-color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.submit:hover {
  background-color: rgb(var(--blue-dark));
}
.cancel {
  background-color: rgb(var(--red));
  margin-bottom: 0;
  &:hover {
    background-color: rgb(var(--red-dark));
  }
}
</style>
