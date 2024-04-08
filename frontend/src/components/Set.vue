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
            <el-table-column prop="value" label="值">
                <template #default="scope">
                    <span v-show="scope.$index !== editIndex">{{ scope.row }}</span>
                    <el-input v-show="scope.$index === editIndex" v-model="scope.row">
                    </el-input>
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" width="120">
                <template #default="scope">
                    <el-button link type="primary" size="small" @click="handleEdit(scope.$index)">编辑</el-button>
                    <el-button link type="warning" size="small" @click="handleDelete(scope.row)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>个数：{{ value.length }}</template>
    </el-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import KeyInfo from './KeyInfo.vue';
import { ElMessage } from 'element-plus'
import { deleteKeyField } from '@/api/key';
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';

const connStore = useConnStore();
const dbStore = useDBStore();

let props = defineProps({
    type: { type: String, default: '' },
    name: { type: String, default: '' },
    ttl: { type: Number, default: 0 },
    value: { type: Array, default: [] },
})

const editIndex = ref(-1)

const handleEdit = (index: number) => {
    editIndex.value = index
}

function handleDelete(row: string) {
    let data = []
    data.push(row)
    deleteKeyField(connStore.current, dbStore.current, props.name, { "fields": data }).then(() => {
        ElMessage({
            type: 'success',
            message: '删除成功',
        })
    })
}
</script>