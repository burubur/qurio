import { createRouter, createWebHistory } from 'vue-router'
import SourcesView from '../views/SourcesView.vue'
import SettingsView from '../views/SettingsView.vue'
import SourceDetailView from '../views/SourceDetailView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'sources',
      component: SourcesView
    },
    {
      path: '/sources/:id',
      name: 'source-detail',
      component: SourceDetailView
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView
    }
  ]
})

export default router
