<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <Header />
      </el-header>
      <el-container>
        <el-aside width="200px">
          <Sider />
        </el-aside>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>

</template>

<script setup lang="ts">
import Sider from './components/Sider.vue';
import Header from './components/Header.vue';

</script>

<style scoped></style>