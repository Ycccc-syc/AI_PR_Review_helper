import { ref } from 'vue'
import type { ReviewResponse } from '../api/review'

const result = ref<ReviewResponse | null>(null)
const prUrl = ref('')

export function useReview() {
  function setReview(data: ReviewResponse, url: string) {
    result.value = data
    prUrl.value = url
  }

  function clearReview() {
    result.value = null
    prUrl.value = ''
  }

  return { result, prUrl, setReview, clearReview }
}
