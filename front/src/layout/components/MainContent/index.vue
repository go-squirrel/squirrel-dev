<template>
  <main class="main-content">
    <slot />
  </main>
</template>

<style scoped>
.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}
</style>
