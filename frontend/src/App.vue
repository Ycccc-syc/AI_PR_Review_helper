<script setup lang="ts">
import { ref } from 'vue'
import ReviewInput from './components/ReviewInput.vue'
import ReviewResult from './components/ReviewResult.vue'
import { analyzePR, type ReviewResponse } from './api/review'

const reviewRef = ref<InstanceType<typeof ReviewInput>>()
const result = ref<ReviewResponse | null>(null)
const error = ref('')
const activeTab = ref('risks')

async function onSubmit(url: string) {
  error.value = ''
  result.value = null
  reviewRef.value?.setLoading(true)

  try {
    result.value = await analyzePR(url)
  } catch (e: any) {
    error.value = e.message || '分析失败，请稍后重试'
  } finally {
    reviewRef.value?.setLoading(false)
  }
}
</script>

<template>
  <div class="app">
    <ReviewInput ref="reviewRef" @submit="onSubmit" />

    <div v-if="error" class="error">{{ error }}</div>

    <ReviewResult
      v-if="result"
      :data="result"
      v-model:tab="activeTab"
    />
  </div>
</template>

<style scoped>
.app {
  max-width: 900px;
  margin: 0 auto;
}

.error {
  max-width: 600px;
  margin: 0 auto 24px;
  padding: 12px 20px;
  background: #fef2f2;
  color: #dc2626;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
}
</style>
