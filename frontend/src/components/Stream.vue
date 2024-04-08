<template>
    <el-card>
        <template #header>
            <KeyInfo :name="name" :type="type" :ttl="ttl" />
        </template>

        <el-button @click="showEditDialog('添加', -1)" icon="plus" plain type="primary">添加新行</el-button>
        <el-table :data="value" style="width: 100%">
            <el-table-column type="index" label="#" width="50">
                <template #default="scope">
                    {{ scope.$index + 1 }}
                </template>
            </el-table-column>
            <el-table-column prop="ID" label="ID">
            </el-table-column>
            <el-table-column prop="Values" label="记录">
                <template #default="scope">
                    <span>{{ JSON.stringify(scope.row.Values) }}</span>
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" width="120">
                <template #default="scope">
                    <el-button link type="warning" size="small" @click="handleDelete(scope.$index)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>个数：{{ value.length }}</template>
        <el-dialog :title="editDialog.title" v-model="editDialog.visible" width="400px" :destroy-on-close="true"
            :close-on-click-modal="false">
            <el-form :model="form" label-width="auto">
                <el-form-item label="ID">
                    <el-input v-model="form.id"></el-input>
                </el-form-item>
                <el-form-item label="字段">
                    <el-input v-model="form.field"></el-input>
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
import { deleteKeyField, addKeyField } from '@/api/key';
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';

const connStore = useConnStore()
const dbStore = useDBStore()

let props = defineProps({
    type: { type: String, default: '' },
    name: { type: String, default: '' },
    ttl: { type: Number, default: 0 },
    value: { type: Array, default: [] },
})

let form = reactive({
    id: '*',
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
        form.id = props.value[index].ID;
        form.field = props.value[index].Values[0];
        form.value = props.value[index].Values[1];
    }
};

function handleDelete(index: number) {
    let data = []
    data.push(props.value[index].ID);
    deleteKeyField(connStore.current, dbStore.current, props.name, { "fields": data }).then(() => {
        props.value.splice(index, 1);
        ElMessage({
            type: 'success',
            message: '删除成功',
        })
    })
}
function handleEdit() {
    addKeyField(connStore.current, dbStore.current, props.name, form).then(() => {
        ElMessage({
            type: 'success',
            message: '添加成功',
        })
        form.value = undefined;
        form.field = undefined;
        editDialog.visible = false;
    })
}
</script>