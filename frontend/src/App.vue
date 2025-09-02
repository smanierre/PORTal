<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
import { useRouter } from 'vue-router';
import { useAppStore } from './stores/app_store.ts';

const store = useAppStore();
const router = useRouter();
if (!store.member) {
  router.replace('/login');
}

function logout() {
  store.logout();
  router.push('/login');
}
</script>

<template>
  <div class="layout">
    <nav v-if="store.member" :hidden="!store.member">
      <RouterLink class="nav-item" to="/">Dashboard</RouterLink>
      <RouterLink class="nav-item" to="/members" v-if="store.subordinates.length > 0"
        >Members</RouterLink
      >
      <RouterLink v-if="store.member.admin" class="nav-item" to="/admin">Admin</RouterLink>
      <button @click="logout" class="logout">Logout</button>
    </nav>
    <div v-else></div>
    <RouterView />
  </div>
</template>

<style scoped>
.layout {
  height: 100vh;
  width: 100vw;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: var(--nav-height) 1fr;
}

nav {
  position: sticky;
  display: flex;
  align-items: center;
  width: 100vw;
  top: 0;
  background-color: rgb(var(--blue));
  justify-content: left;
}

.nav-item {
  display: flex;
  width: fit-content;
  color: white;
  cursor: pointer;
  height: 100%;
  padding: 0 1rem;
  align-items: center;
  text-decoration: none;
  font-size: 0.8rem;

  &:hover {
    background-color: rgb(var(--blue-dark));
  }
}

.logout {
  padding: 0 1rem;
  height: 100%;
  margin-left: auto;
  border: none;
  font-size: 0.8rem;
  border-radius: 0;

  &:hover {
    background-color: rgb(var(--blue-dark));
  }
}
</style>
