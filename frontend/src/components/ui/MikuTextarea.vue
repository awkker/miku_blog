<template>
  <label class="block space-y-2">
    <span v-if="label" class="text-sm font-medium text-white/85">{{ label }}</span>

    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :required="required"
      :aria-label="ariaLabel || label"
      :rows="rows"
      class="w-full rounded-2xl border bg-white/20 px-3 py-2.5 text-sm text-white placeholder:text-white/55 transition duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-miku/70"
      :class="[
        error
          ? 'border-red-300/85 focus-visible:border-red-300 focus-visible:ring-red-300/70'
          : 'border-white/30 focus-visible:border-miku/85',
        disabled ? 'cursor-not-allowed opacity-70' : '',
      ]"
      @input="onInput"
      @blur="$emit('blur')"
      @focus="$emit('focus')"
    />

    <p v-if="error" class="text-xs text-red-200">{{ error }}</p>
  </label>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string
  label?: string
  placeholder?: string
  rows?: number
  error?: string
  disabled?: boolean
  required?: boolean
  ariaLabel?: string
}

withDefaults(defineProps<Props>(), {
  label: undefined,
  placeholder: '',
  rows: 5,
  error: '',
  disabled: false,
  required: false,
  ariaLabel: undefined,
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'focus'): void
  (e: 'blur'): void
}>()

function onInput(event: Event) {
  const target = event.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
}
</script>
