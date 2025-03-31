<template>
    <el-card>
        <template #header>
            <KeyInfo :name="name" :type="type" :ttl="ttl" />
        </template>

        <el-button @click="showEditDialog('添加', -1)" icon="plus" plain type="primary">添加新行</el-button>
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
                    <el-button link type="primary" size="small"
                        @click="showEditDialog('修改', scope.$index)">编辑</el-button>
                    <el-button link type="warning" size="small" @click="handleDelete(scope.$index)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>个数：{{ list.length }}</template>
        <el-dialog :title="editDialog.title" v-model="editDialog.visible" width="400px" :destroy-on-close="true"
            :close-on-click-modal="false">
            <el-form :model="form" label-width="auto">
                <el-form-item label="字段">
                    <el-input v-model="form.field" :disabled="editDialog.index >= 0"></el-input>
                </el-form-item>
                <el-form-item label="值">
                    <el-input v-model="form.value"></el-input>
                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="editDialog.visible = false">取 消</el-button>
                    <el-button type="primary" @click="handleEdit">确 定</el-button>
                </div>
            </template>
        </el-dialog>
    </el-card>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import KeyInfo from './KeyInfo.vue';
import { ElMessage } from 'element-plus'
import { deleteKeyField, addKeyField, pathcKeyField } from '@/api/key';
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';

const connStore = useConnStore()
const dbStore = useDBStore()

let props = defineProps({
    type: { type: String, default: '' },
    name: { type: String, default: '' },
    ttl: { type: Number, default: 0 },
    value: { type: Object, default: {} },
})

const list = reactive([])

for (let i in props.value) {
    let obj = { "key": i, "value": props.value[i] }
    list.push(obj)
}

let form = reactive({
    field: undefined,
    value: undefined,
})

const editDialog = reactive({
    index: -1,
    visible: false,
    title: '',
})

const showEditDialog = (title: string, index: number) => {
    editDialog.index = index;
    editDialog.title = title ? title : '';
    editDialog.visible = true;
    if (index >= 0) {
        form.value = list[index].value;
        form.field = list[index].key;
    }
};


function handleDelete(index: number) {
    let data = [];
    data.push(list[index].key);
    deleteKeyField(connStore.current, dbStore.current, props.name, { "fields": data }).then(() => {
        list.splice(index, 1);
        ElMessage({
            type: 'success',
            message: '删除成功',
        })
    })
}

function handleEdit() {
    console.log(editDialog.index)
    if (editDialog.index < 0) {
        addKeyField(connStore.current, dbStore.current, props.name, form).then(() => {
            list.push({ "key": form.field, "value": form.value });
            ElMessage({
                type: 'success',
                message: '添加成功',
            })
            form.value = undefined;
            form.field = undefined;
            editDialog.visible = false;
        })
    } else {
        pathcKeyField(connStore.current, dbStore.current, props.name, form).then(() => {
            list[editDialog.index].value = form.value
            ElMessage({
                type: 'success',
                message: '修改成功',
            })
            form.value = undefined;
            form.field = undefined;
            editDialog.visible = false;
        })
    }
}
</script>