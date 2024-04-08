<template>
    <el-button type="primary" :icon="Plus" @click="dialogFormVisible = true">连接</el-button>
    <el-dialog v-model="dialogFormVisible" title="配置Redis连接" width="500">
        <el-form :model="form">
            <el-form-item label="名称">
                <el-input v-model="form.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="主机">
                <el-input v-model="form.host" autocomplete="off" />
            </el-form-item>
            <el-form-item label="端口">
                <el-input v-model="form.port" autocomplete="off" />
            </el-form-item>
            <el-form-item label="用户">
                <el-input v-model="form.username" autocomplete="off" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button type="warning" @click="handleTest" style="float: left;">测试</el-button>
                <el-button @click="dialogFormVisible = false">取消</el-button>
                <el-button type="primary" @click="handleCreate">
                    确认
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import { Plus } from "@element-plus/icons-vue"
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router';
import { createConn, testConn } from '@/api/conn';

const dialogFormVisible = ref(false)
const form = reactive({
    name: '',
    host: '',
    port: '',
    username: '',
    password: '',
})

const router = useRouter()

function handleCreate() {
    createConn(form).then(result => {
        ElMessage({
            message: '创建成功',
            type: 'success'
        })
    })
    dialogFormVisible.value = false
    router.push("/")
}

function handleTest() {
    testConn(form).then(result => {
        ElMessage({
            message: '连接成功',
            type: 'success',
        })
    })
}
</script>