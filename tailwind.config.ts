import type { Config } from 'tailwindcss'
import typography from '@tailwindcss/typography'

const config: Config = {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        miku: 'rgb(102 205 170)',
        'miku-soft': 'rgba(102, 205, 170, 0.12)',
        'miku-strong': 'rgba(102, 205, 170, 0.8)',
      },
      fontFamily: {
        sans: ['"LXGW WenKai Screen"', '"LXGW WenKai"', '"PingFang SC"', '"Hiragino Sans GB"', '"Microsoft YaHei"', 'sans-serif'],
        mono: ['"JetBrains Mono"', '"Fira Code"', '"SFMono-Regular"', 'Menlo', 'Monaco', 'Consolas', '"Liberation Mono"', '"Courier New"', 'monospace'],
      },
      boxShadow: {
        glass: '0 18px 38px rgba(15, 23, 42, 0.22), inset 0 1px 0 rgba(255,255,255,0.55)',
      },
      keyframes: {
        blink: {
          '0%, 49%': { opacity: '1' },
          '50%, 100%': { opacity: '0' },
        },
        heartbeat: {
          '0%': { transform: 'scale(1)' },
          '25%': { transform: 'scale(1.25)' },
          '55%': { transform: 'scale(0.96)' },
          '100%': { transform: 'scale(1)' },
        },
        'dock-float': {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-6px)' },
        },
        'fade-up': {
          '0%': { opacity: '0', transform: 'translateY(24px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
      animation: {
        blink: 'blink 1s step-end infinite',
        heartbeat: 'heartbeat 320ms ease-out',
        'dock-float': 'dock-float 2.4s ease-in-out infinite',
        'fade-up': 'fade-up 620ms ease-out both',
      },
      typography: {
        DEFAULT: {
          css: {
            maxWidth: 'none',
          },
        },
      },
    },
  },
  plugins: [typography],
}

export default config
