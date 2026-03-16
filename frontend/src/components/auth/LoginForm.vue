<template>
  <LiquidGlassCard width="100%" maxWidth="28rem" padding="28px" class="mx-auto">
    <div class="text-center">
      <p class="text-xs uppercase tracking-[0.28em] text-slate-700">Miku Blog Console</p>
      <img
        src="/picture/author.jpg"
        alt="管理员头像"
        class="mx-auto mt-5 h-20 w-20 rounded-full border border-miku/70 object-cover shadow-[0_0_24px_rgba(57,197,187,0.45)] transition duration-300 hover:scale-105"
      />
      <h1 class="mt-4 text-2xl font-semibold text-slate-900">欢迎回来</h1>
      <p class="mt-1 text-sm text-slate-700">仅管理员可访问后台系统</p>
      <div class="mt-4 flex flex-wrap justify-center gap-2">
        <a
          href="/blog"
          class="rounded-xl border border-slate-300/80 bg-white/55 px-3 py-1 text-xs text-slate-900 transition hover:border-miku/40 hover:text-miku"
          aria-label="前往博客首页"
        >
          博客首页
        </a>
        <a
          href="/guestbook"
          class="rounded-xl border border-slate-300/80 bg-white/55 px-3 py-1 text-xs text-slate-900 transition hover:border-miku/40 hover:text-miku"
          aria-label="前往留言板"
        >
          留言板
        </a>
        <a
          href="/friends"
          class="rounded-xl border border-slate-300/80 bg-white/55 px-3 py-1 text-xs text-slate-900 transition hover:border-miku/40 hover:text-miku"
          aria-label="前往友链页面"
        >
          友情链接
        </a>
      </div>
    </div>

    <form
      class="mt-6 space-y-4 transition duration-300"
      :class="hasAuthError ? 'rounded-2xl ring-2 ring-red-300/80 ring-offset-2 ring-offset-transparent' : ''"
      @submit.prevent="handleSubmit"
    >
      <MikuInput
        v-model="identifier"
        label="用户名"
        placeholder="admin"
        autocomplete="username"
        :error="errors.identifier"
        aria-label="用户名"
        @focus="clearAuthError"
        required
      />

      <MikuInput
        v-model="password"
        :type="showPassword ? 'text' : 'password'"
        label="密码"
        placeholder="请输入密码"
        autocomplete="current-password"
        :error="errors.password"
        aria-label="密码"
        @focus="clearAuthError"
        required
      >
        <template #trailing>
          <button
            type="button"
            class="rounded-lg p-1 text-slate-600 transition duration-300 hover:bg-white/30 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70"
            :aria-label="showPassword ? '隐藏密码' : '显示密码'"
            @click="showPassword = !showPassword"
          >
            <svg v-if="showPassword" viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
              <path d="M3 3l18 18" />
              <path d="M10.58 10.58A3 3 0 0012 15a3 3 0 001.42-.36" />
              <path d="M9.88 5.09A10.94 10.94 0 0112 5c6 0 10 7 10 7a18.47 18.47 0 01-4.45 5.01" />
              <path d="M6.61 6.61A18.81 18.81 0 002 12s1.62 2.84 4.7 4.77" />
            </svg>
            <svg v-else viewBox="0 0 24 24" class="h-4 w-4 fill-none stroke-current stroke-[1.8]">
              <path d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7-10-7-10-7z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
          </button>
        </template>
      </MikuInput>

      <p
        v-if="formError"
        class="flex items-center gap-2 rounded-xl border border-red-300/90 bg-red-100/80 px-3 py-2 text-sm font-medium text-slate-900"
        role="status"
        aria-live="polite"
      >
        <svg viewBox="0 0 24 24" class="h-4 w-4 shrink-0 fill-none stroke-current stroke-[2]">
          <path d="M12 8v4" />
          <path d="M12 16h.01" />
          <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3l-8.47-14.14a2 2 0 00-3.42 0z" />
        </svg>
        {{ formError }}
      </p>
      <p
        v-if="successMessage"
        class="flex items-center gap-2 rounded-xl border border-emerald-300/85 bg-emerald-100/80 px-3 py-2 text-sm font-medium text-slate-900"
        role="status"
        aria-live="polite"
      >
        <svg viewBox="0 0 24 24" class="h-4 w-4 shrink-0 fill-none stroke-current stroke-[2]">
          <path d="M20 6L9 17l-5-5" />
        </svg>
        {{ successMessage }}
      </p>

      <MikuButton
        type="submit"
        :loading="submitting"
        :disabled="submitting"
        :full-width="true"
        aria-label="登录"
      >
        {{ submitting ? '登录中...' : '登录' }}
      </MikuButton>

      <p class="text-center text-xs text-slate-700" aria-live="polite">
        {{ submitting ? '正在验证账号信息…' : '请使用有效管理员账号登录后台' }}
      </p>
      <p class="text-center text-[11px] text-slate-600">
        测试账号：admin / admin123
      </p>
    </form>
  </LiquidGlassCard>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, reactive, ref } from 'vue'

import { authState, hydrateAuth, loginWithPassword } from '../../stores/auth'
import { getScopeLoading, setScopeStatus } from '../../stores/loading'
import MikuButton from '../ui/MikuButton.vue'
import MikuInput from '../ui/MikuInput.vue'
import LiquidGlassCard from '../ui/LiquidGlassCard.vue'

const identifier = ref('')
const password = ref('')
const showPassword = ref(false)
const formError = ref('')
const successMessage = ref('')

const errors = reactive({
  identifier: '',
  password: '',
})

const submitting = useStore(getScopeLoading('loginSubmit'))
const auth = useStore(authState)
const hasAuthError = computed(() => Boolean(formError.value))

onMounted(() => {
  hydrateAuth()

  const latestAuth = authState.get()

  if (latestAuth.status === 'authenticated') {
    window.location.replace('/admin')
  }
})

function validate() {
  errors.identifier = identifier.value.trim() ? '' : '用户名不能为空'
  errors.password = password.value.trim() ? '' : '密码不能为空'
  return !errors.identifier && !errors.password
}

function clearAuthError() {
  formError.value = ''
  successMessage.value = ''

  if (errors.identifier === '账号或密码错误，请重新输入') {
    errors.identifier = ''
  }

  if (errors.password === '账号或密码错误，请重新输入') {
    errors.password = ''
  }
}

async function handleSubmit() {
  formError.value = ''
  successMessage.value = ''

  if (!validate()) {
    return
  }

  setScopeStatus('loginSubmit', 'loading')

  try {
    await loginWithPassword(identifier.value, password.value)
    setScopeStatus('loginSubmit', 'success')
    successMessage.value = '登录成功，正在进入后台...'
    window.setTimeout(() => {
      window.location.assign('/admin')
    }, 250)
  } catch (error) {
    setScopeStatus('loginSubmit', 'error')
    formError.value = error instanceof Error ? error.message : '登录失败，请稍后重试。'
    errors.identifier = '账号或密码错误，请重新输入'
    errors.password = '账号或密码错误，请重新输入'
  }
}
</script>
