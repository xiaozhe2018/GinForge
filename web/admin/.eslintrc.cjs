/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution')

module.exports = {
  root: true,
  env: {
    node: true,
    browser: true,
    es2022: true
  },
  'extends': [
    'plugin:vue/vue3-essential'
  ],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module'
  },
  rules: {
    'vue/multi-word-component-names': 'off'
  },
  ignorePatterns: [
    'auto-imports.d.ts',
    'components.d.ts',
    '*.d.ts'
  ]
}




