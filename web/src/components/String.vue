<template>
    <el-card>
        <template #header>
            <KeyInfo :name="name" :type="type" :ttl="ttl" />
        </template>
        <el-input :modelValue="value" :autosize="{ minRows: 6, maxRows: 10 }" type="textarea" />
        <template #footer>
            <span>长度：{{ value.length }}</span>
        </template>
    </el-card>
</template>

<script lang="ts" setup>
import KeyInfo from './KeyInfo.vue';
defineProps(['type', 'name', 'ttl', 'value'])
</script>