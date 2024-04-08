<template>
    <el-card>
        <template #header>
            <KeyInfo :name="name" :type="type" :ttl="ttl" />
        </template>

        <el-table :data="value" style="width: 100%">
            <el-table-column type="index" label="#" width="50">
                <template #default="scope">
                    {{ scope.$index + 1 }}
                </template>
            </el-table-column>
            <el-table-column prop="ID" label="字段">
            </el-table-column>
            <el-table-column prop="Values" label="分值">
                <template #default="scope">
                    <span>{{ JSON.stringify(scope.row.Values) }}</span>
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" width="120">
                <template #default="scope">
                    <el-button link type="primary" size="small" @click="handleEdit(scope.$index)">编辑</el-button>
                    <el-button link type="warning" size="small">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>个数：{{ value.length }}</template>
    </el-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import KeyInfo from './KeyInfo.vue';

defineProps(['type', 'name', 'ttl', 'value'])
const editIndex = ref(-1)

const handleEdit = (index) => {
    editIndex.value = index
}
</script>