<template>
    <el-row :gutter="20">

        <el-tag v-if="type === 'string'" size="large" class="string-tag">String</el-tag>
        <el-tag v-else-if="type === 'list'" size="large" class="list-tag">List</el-tag>
        <el-tag v-else-if="type === 'hash'" size="large" class="hash-tag">Hash</el-tag>
        <el-tag v-else-if="type === 'set'" size="large" class="set-tag">Set</el-tag>
        <el-tag v-else-if="type === 'zset'" size="large" class="zset-tag">ZSet</el-tag>
        <el-tag v-else-if="type === 'stream'" size="large" class="stream-tag">Stream</el-tag>
        <el-input :value="name" disabled style="width: auto;margin-left: 10px;margin-right: 10px;"></el-input>
        <el-button-group>
            <el-button type="warning" @click="dialogNameVisible = true"><el-icon>
                    <Edit />
                </el-icon>修改</el-button>
            <el-button type="primary"><el-icon>
                    <Refresh />
                </el-icon>刷新</el-button>
            <el-button type="danger" @click="handleDelete"><el-icon>
                    <Delete />
                </el-icon>删除</el-button>
            <el-button type="info" @click="dialogTTLVisible = true"><el-icon>
                    <AlarmClock />
                </el-icon>
                <span v-if="ttl == 0">永久</span>
                <span v-else>TTL:{{ ttl }}</span>
            </el-button>
        </el-button-group>
    </el-row>

    <el-dialog v-model="dialogNameVisible" title="重命名键" width="500">
        <el-form>
            <el-form-item>
                <el-input v-model="originName" autocomplete="off" />
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="dialogNameVisible = false">取消</el-button>
                <el-button type="primary" @click="updateName()">
                    确认
                </el-button>
            </div>
        </template>
    </el-dialog>

    <el-dialog v-model="dialogTTLVisible" title="设置过期时间" width="500">
        <el-form>
            <el-form-item>
                <el-input v-model.number="ttl" autocomplete="off" type="number" />
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="dialogTTLVisible = false">取消</el-button>
                <el-button type="primary" @click="updateTTL()">
                    确认
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import { defineProps, ref, toRef } from 'vue';
import { patchKeyName, patchKeyTTL, deleteKey } from '@/api/key';
import { ElMessage, ElMessageBox } from 'element-plus'
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';

let props = defineProps(['type', 'name', 'ttl'])

const connStore = useConnStore()
const dbStore = useDBStore()

const dialogNameVisible = ref(false)
const dialogTTLVisible = ref(false)

const originName = toRef(props.name)
const ttl = toRef(props.ttl)

function updateName() {
    dialogNameVisible.value = false
    patchKeyName(connStore.current, dbStore.current, props.name, originName.value).then(result => {
        console.log(result)
    })
}

function updateTTL() {
    dialogTTLVisible.value = false
    console.log(ttl)
    patchKeyTTL(connStore.current, dbStore.current, props.name, ttl.value).then(result => {
        console.log(result)
    })
}

const handleDelete = () => {
    ElMessageBox.confirm(
        '确定删除?',
        'Warning',
        {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
        }
    ).then(() => {
        deleteKey(connStore.current, dbStore.current, props.name).then(result => {
            ElMessage({
                type: 'success',
                message: '删除成功',
            })
        }).catch(() => {
            ElMessage({
                type: 'error',
                message: '删除失败',
            })
        })
    }).catch(() => {
        ElMessage({
            type: 'info',
            message: '取消删除',
        })
    })
}

</script>

<style scoped>
.string-tag {
    background-color: #95d475;
    color: white;
    border-color: #95d475;
}

.list-tag {
    background-color: #e3e058;
    color: white;
    border-color: #e3e058;
}

.hash-tag {
    background-color: #c1df26;
    color: white;
    border-color: #c1df26;
}


.set-tag {
    background-color: #58d7e3;
    color: white;
    border-color: #58d7e3;
}

.zset-tag {
    background-color: #1ba5e0;
    color: white;
    border-color: #1ba5e0;
}

.stream-tag {
    background-color: #e22c2c;
    color: white;
    border-color: #e22c2c;
}
</style>