<template>
    <el-table :data="clients" style="width: 100%">
        <el-table-column prop="id" label="id" width="50" />
        <el-table-column prop="addr" label="addr" width="150" />
        <el-table-column prop="flags" label="flags" width="60" />
        <el-table-column prop="db" label="db" width="50" />
        <el-table-column prop="cmd" label="cmd" width="120" />
        <el-table-column prop="age" label="age" width="120" />
        <el-table-column prop="idle" label="idle" width="120" />
        <el-table-column prop="sub" label="sub" width="120" />
        <el-table-column prop="psub" label="psub" width="120" />
        <el-table-column fixed="right" label="操作" width="120">
            <template #default="scope">
                <el-button type="danger" size="small" @click="handleClick(scope.row.addr)" :icon="CircleClose">
                    关闭
                </el-button>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup>
import { reactive, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { CircleClose } from '@element-plus/icons-vue'
import { useConnStore } from '@/store/Conn';
import { listClients, deleteClient } from '@/api/conn';

const connStore = useConnStore()

const clients = reactive([])

function getData() {
    clients.length = 0
    listClients(connStore.current).then(result => {
        for (let i in result.data) {
            let client = result.data[i]
            clients.push(client)
        }
    })
}

const handleClick = (addr) => {
    deleteClient(connStore.current, { "addr": addr }).then(result => {
        ElMessage({
            message: '删除成功',
            type: 'success'
        })
        getData()
    })
}

onMounted(() => {
    getData()
})
</script>