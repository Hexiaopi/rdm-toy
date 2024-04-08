<template>
    <el-card>
        <template #header>
            <KeyInfo :name="name" :type="type" :ttl="ttl" />
        </template>

        <el-table :data="list" style="width: 100%">
            <el-table-column type="index" label="#" width="50">
                <template #default="scope">
                    {{ scope.$index + 1 }}
                </template>
            </el-table-column>
            <el-table-column prop="key" label="字段">
            </el-table-column>
            <el-table-column prop="value" label="值">

            </el-table-column>
            <el-table-column fixed="right" label="操作" width="120">
                <template #default="scope">
                    <el-button link type="primary" size="small" @click="handleEdit(scope.$index)">编辑</el-button>
                    <el-button link type="warning" size="small">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>个数：{{ list.length }}</template>
    </el-card>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import KeyInfo from './KeyInfo.vue';

let props = defineProps(['type', 'name', 'ttl', 'value'])

const list = reactive([])

for (let i in props.value) {
    let obj = { "key": i, "value": props.value[i] }
    list.push(obj)
}

const editIndex = ref(-1)

const handleEdit = (index) => {
    editIndex.value = index
}
</script>