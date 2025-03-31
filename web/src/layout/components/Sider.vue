<template>
    <el-menu class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose">
        <el-sub-menu v-for="conn, key in conns" :key="key" :index="conn.key" @click="handleConnClick(conn)">
            <template #title :id="'conn#' + conn.key">
                <el-icon>
                    <Connection />
                </el-icon>
                <span>{{ conn.title }}</span>
            </template>
            <el-menu-item v-for="db, index in conn.children" :key="index" :index="db.key">
                <template #title>
                    <el-icon>
                        <Setting />
                    </el-icon>
                    <el-badge :value="db.total" style="margin-top:10px">
                        <el-button @click.stop="handleDBClick(db)" style="margin-bottom:15px">
                            {{ db.title }}
                        </el-button>
                    </el-badge>
                    <!-- <el-button type="danger" size="small" :icon="Delete" style="margin-left: auto;" /> -->
                </template>
            </el-menu-item>
        </el-sub-menu>
    </el-menu>
</template>

<script lang="ts" setup>

import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useConnStore } from '@/store/Conn';
import { useDBStore } from '@/store/DB';
import { ElMessage } from 'element-plus'
import { Delete } from '@element-plus/icons-vue';
import { listConns, deleteConn } from '@/api/conn';
import { listDBs } from '@/api/db';

const connStore = useConnStore()
const dbStore = useDBStore()
const conns = reactive([])
const router = useRouter()


const handleOpen = (key: string, keyPath: string[]) => {
    console.log(key, keyPath)
    //查询连接下的db
    listDBs(key).then(result => {
        let dbs = []
        for (let i in result.data) {
            const obj = { key: result.data[i].id, title: result.data[i].name, total: result.data[i].total, children: [] }
            dbs.push(obj)
        }
        for (let i in conns) {
            if (conns[i].key === key) {
                conns[i].children = dbs
            }
        }
    })

}

const handleConnClick = (conn: any) => {
    console.log("conn click", conn)
    connStore.current = conn.key
    router.push(`/conn`)
}

const handleConnDelete = (conn: string) => {
    deleteConn(conn).then(result => {
        ElMessage({
            message: '删除成功',
            type: 'success'
        })
        if (connStore.current === conn) {
            connStore.current = ''
            router.push("/")
        } else {
            getConns()
        }
    })
}

const handleDBClick = (db: any) => {
    console.log("db click", db)
    dbStore.current = db.key
    router.push(`/db`)
}


const handleSelect = (key: string, keyPath: string[]) => {
    console.log("handle select", key, keyPath)
    router.push("/conn")
}

const handleClose = (key: string, keyPath: string[]) => {
    console.log(key, keyPath)
}

function getConns() {
    conns.length = 0
    listConns('').then(result => {
        console.log(result)
        for (let i in result.data) {
            const obj = { key: result.data[i].name, title: result.data[i].name, children: [] }
            conns.push(obj)
        }
    })
}

onMounted(() => {
    getConns()
})
</script>

<style scoped></style>