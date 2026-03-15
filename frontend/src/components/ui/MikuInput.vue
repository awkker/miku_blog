<template>
  <label class="block space-y-2">
    <span v-if="label" class="text-sm font-medium text-slate-800">{{ label }}</span>

    <div class="relative">
      <input
        :value="modelValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        :autocomplete="autocomplete"
        :aria-label="ariaLabel || label"
        class="w-full rounded-xl border bg-white/65 px-3 py-2.5 text-sm text-slate-900 placeholder:text-slate-500 transition duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70"
        :class="[
          hasTrailingSlot ? 'pr-11' : '',
          error
            ? 'border-red-300/85 focus-visible:border-red-300 focus-visible:ring-red-300/70'
            : 'border-slate-300/90 focus-visible:border-miku/85',
          disabled ? 'cursor-not-allowed opacity-70' : '',
        ]"
        @input="onInput"
        @blur="$emit('blur')"
        @focus="$emit('focus')"
      />

      <div
        v-if="hasTrailingSlot"
        class="absolute inset-y-0 right-2 flex items-center"
      >
        <slot name="trailing" />
      </div>
    </div>

    <p v-if="error" class="text-xs text-red-700">{{ error }}</p>
  </label>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

interface Props {
  modelValue: string
  label?: string
  type?: string
  placeholder?: string
  error?: string
  disabled?: boolean
  required?: boolean
  autocomplete?: string
  ariaLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  label: undefined,
  placeholder: '',
  error: '',
  disabled: false,
  required: false,
  autocomplete: 'off',
  ariaLabel: undefined,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'focus'): void
  (e: 'blur'): void
}>()

const slots = useSlots()
const hasTrailingSlot = computed(() => Boolean(slots.trailing))

function onInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>
