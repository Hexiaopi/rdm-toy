<template>
    <el-row :gutter="20">
        <el-col :span="8">
            <el-row>
                <el-input v-model="filter" style="width: auto;" placeholder="match" clearable />
                <el-button-group>
                    <el-button :icon="Search" @click="getKeys" />
                    <el-button :icon="Refresh" @click="getKeys" />
                    <el-button :icon="Plus" @click="dialogFormVisible = true" />
                    <el-button :icon="Delete" @click="handleKeysDelete" />
                </el-button-group>
            </el-row>

            <el-scrollbar max-height="800px">
                <p v-for="key in keys" :key="key" class="scrollbar-demo-item">
                    <el-button @click.stop="handleKeyClick(key)" style="width: 85%;overflow: hidden;">
                        {{ key.name }}
                    </el-button>

                    <el-button-group class="ml-4" style="margin-left: auto;">
                        <el-button type="danger" title="删除" size="small" @click="handleDelete(key.name)"
                            :icon="Delete"></el-button>
                    </el-button-group>
                </p>
            </el-scrollbar>
        </el-col>
        <el-col :span="16">
            <el-tabs v-model="editableTabsValue" type="card" class="demo-tabs" closable @tab-remove="removeTab">
                <el-tab-pane v-for="item in editableTabs" :key="item.name" :label="item.name" :name="item.name">
                    <template v-if="item.type === 'string'">
                        <String :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else-if="item.type === 'list'">
                        <List :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else-if="item.type === 'hash'">
                        <Hash :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else-if="item.type === 'set'">
                        <Set :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else-if="item.type === 'zset'">
                        <ZSet :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else-if="item.type === 'stream'">
                        <Stream :name="item.name" :type="item.type" :ttl="item.ttl" :value="item.value" />
                    </template>
                    <template v-else>
                        {{ item.name }}
                    </template>
                </el-tab-pane>
            </el-tabs>
        </el-col>
    </el-row>
    <el-dialog v-model="dialogFormVisible" title="添加新键" width="500">
        <el-form :model="form" label-width="auto">
            <el-form-item label="键">
                <el-input v-model="form.name" autocomplete="off" required />
            </el-form-item>
            <el-form-item label="TTL">
                <el-input v-model.number="form.ttl" autocomplete="off" required />
            </el-form-item>
            <el-form-item label="类型">
                <el-select v-model="form.type">
                    <el-option label="String" value="string" />
                    <el-option label="List" value="list" />
                    <el-option label="Hash" value="hash" />
                    <el-option label="Set" value="set" />
                    <el-option label="ZSet" value="zset" />
                    <el-option label="Stream" value="stream" />
                </el-select>
            </el-form-item>
            <el-form-item label="值" v-show="form.type === 'string'">
                <el-input type="textarea" v-model="form.string"></el-input>
            </el-form-item>
            <el-form-item v-show="form.type === 'list' || form.type === 'set'">
                <el-table :data="form.list">
                    <el-table-column type="index" label="#" width="50">
                        <template #default="scope">
                            {{ scope.$index + 1 }}
                        </template>
                    </el-table-column>
                    <el-table-column label="值">
                        <template #default="scope">
                            <el-input v-model="form.list[scope.$index]"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column align="right">
                        <template #header>
                            <el-button :icon="Plus" @click="addList()" />
                        </template>
                        <template #default="scope">
                            <el-button size="small" type="danger" @click="deleteList(scope.$index)"
                                :icon="Delete"></el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-form-item>
            <el-form-item label="ID" v-show="form.type === 'stream'">
                <el-input v-model="form.id"></el-input>
            </el-form-item>
            <el-form-item v-show="form.type === 'hash' || form.type === 'stream'">
                <el-table :data="form.hash">
                    <el-table-column type="index" label="#" width="50">
                        <template #default="scope">
                            {{ scope.$index + 1 }}
                        </template>
                    </el-table-column>
                    <el-table-column label="字段">
                        <template #default="scope">
                            <el-input v-model="form.hash[scope.$index].field"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column label="值">
                        <template #default="scope">
                            <el-input v-model="form.hash[scope.$index].value"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column align="right">
                        <template #header>
                            <el-button :icon="Plus" @click="addHash()" />
                        </template>
                        <template #default="scope">
                            <el-button size="small" type="danger" @click="deleteHash(scope.$index)"
                                :icon="Delete"></el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-form-item>
            <el-form-item v-show="form.type === 'zset'">
                <el-table :data="form.zset">
                    <el-table-column type="index" label="#" width="50">
                        <template #default="scope">
                            {{ scope.$index + 1 }}
                        </template>
                    </el-table-column>
                    <el-table-column label="分数" prop="score">
                        <template #default="scope">
                            <el-input v-model="form.zset[scope.$index].score"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column label="成员">
                        <template #default="scope">
                            <el-input v-model="form.zset[scope.$index].member"></el-input>
                        </template>
                    </el-table-column>
                    <el-table-column align="right">
                        <template #header>
                            <el-button :icon="Plus" @click="addZSet()" />
                        </template>
                        <template #default="scope">
                            <el-button size="small" type="danger" @click="deleteZSet(scope.$index)"
                                :icon="Delete"></el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="dialogFormVisible = false">取消</el-button>
                <el-button type="primary" @click="handleAdd">确认</el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script lang="ts" setup>
import { onMounted, watchEffect, reactive, ref, } from 'vue'
import { Plus, Refresh, Search, Delete } from '@element-plus/icons-vue'
import { listKeys, getKey, deleteKey, deleteKeys, addKey } from '@/api/key';
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';
import { ElMessage } from 'element-plus';

const keys = reactive([])
const filter = ref("*")
const connStore = useConnStore()
const dbStore = useDBStore()

const handleKeyClick = (key: any) => {
    console.log("click key", key)
    addTab(key.name)
}

const editableTabsValue = ref('')
const editableTabs = ref([])
const addTab = (targetName: string) => {
    const tabs = editableTabs.value
    for (let i in tabs) {
        if (tabs[i].name == targetName) {
            return
        }
    }
    getKey(connStore.current, dbStore.current, targetName).then(result => {
        let obj = { name: result.data.name, type: result.data.type, value: result.data.value, ttl: result.data.ttl }
        editableTabs.value.push(obj)
    })
    editableTabsValue.value = targetName
}
const removeTab = (targetName: string) => {
    const tabs = editableTabs.value
    let activeName = editableTabsValue.value
    if (activeName === targetName) {
        tabs.forEach((tab, index) => {
            if (tab.name === targetName) {
                const nextTab = tabs[index + 1] || tabs[index - 1]
                if (nextTab) {
                    activeName = nextTab.name
                }
            }
        })
    }

    editableTabsValue.value = activeName
    editableTabs.value = tabs.filter((tab) => tab.name !== targetName)
}

function getKeys() {
    keys.length = 0
    listKeys(connStore.current, dbStore.current, { "filter": filter.value }).then(result => {
        for (let i in result.data) {
            keys.push(result.data[i])
        }
    })
}


const dialogFormVisible = ref(false)
let form = reactive({
    name: '',
    ttl: -1,
    type: 'string',
    string: '',
    list: [''],
    hash: [{ field: '', value: '' }],
    zset: [{ score: '', member: '' }],
    id: '*',
})
function addList() {
    form.list.push('')
}
function deleteList(index: number) {
    form.list.splice(index, 1)
}
function addHash() {
    form.hash.push({ field: '', value: '' })
}
function deleteHash(index: number) {
    form.hash.splice(index, 1)
}
function addZSet() {
    form.zset.push({ score: '', member: '' })
}
function deleteZSet(index: number) {
    form.zset.splice(index, 1)
}

function handleKeysDelete() {
    deleteKeys(connStore.current, dbStore.current).then(result => {
        ElMessage({
            message: 'flush成功',
            type: 'success',
        })
        getKeys()
    })
}

function handleAdd() {
    let value = undefined
    switch (form.type) {
        case "string":
            value = form.string
            break;
        case "list":
            value = form.list
            break;
        case "hash":
            value = form.hash
            break;
        case "set":
            value = form.list
            break;
        case "zset":
            value = form.zset
            break;
        case "stream":
            value = { id: form.id, value: form.hash }
            break;
    }
    let key = { name: form.name, ttl: form.ttl, type: form.type, value: value }
    form = Object.assign(form, { name: '', ttl: -1, type: 'string', list: [''], hash: [{ field: '', value: '' }], zset: [{ score: '', member: '' }], id: '*', })
    addKey(connStore.current, dbStore.current, key).then(result => {
        ElMessage({
            message: '添加成功',
            type: 'success',
        })
        getKeys()
        dialogFormVisible.value = false
    })
}

function handleDelete(key: string) {
    deleteKey(connStore.current, dbStore.current, key).then(result => {
        ElMessage({
            message: '删除成功',
            type: 'success'
        })
        getKeys()
    }).catch(err => {
        ElMessage({
            message: err,
            type: 'error'
        })
    })
}

// onMounted(() => {
//     //查询db下的key
//     getKeys()
// })

watchEffect(() => {
    getKeys()
})
</script>

<style scoped></style>