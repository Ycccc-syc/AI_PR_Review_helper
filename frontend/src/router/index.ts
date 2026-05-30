import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomePage },
    { path: '/auth', component: () => import('../pages/AuthPage.vue') },
    { path: '/history', component: () => import('../pages/HistoryPage.vue') },
    { path: '/ask', component: () => import('../pages/AskPage.vue') },
  ],
})

export default router
