import '@/main.css';
import Admin from '@/views/Admin.vue';

import Dashboard from '@/views/Dashboard.vue';
import Login from '@/views/Login.vue';
import Members from '@/views/Members.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Dashboard,
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    {
      path: '/members',
      name: 'members',
      component: Members,
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin,
    },
  ],
});

export default router;
