<template>
  <router-view v-slot="{ Component, route }">
    <template v-if="route.meta.layout === 'full'">
      <component :is="Component" />
    </template>
    <template v-else>
      <Layout>
        <component :is="Component" />
      </Layout>
    </template>
  </router-view>
</template>

<script setup lang="ts">
import Layout from './layout/index.vue'
</script>

<style lang="scss">
@import './styles/index.scss';
</style>
