<template>
    <el-tabs v-model="activeName" @tab-click="handleClick">
        <el-tab-pane label="基本" name="info">
            <el-tabs tab-position="left">
                <el-tab-pane label="Server">
                    <el-table :data="server" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="CPU">
                    <el-table :data="cpu" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Clients">
                    <el-table :data="clients" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Cluster">
                    <el-table :data="cluster" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Keyspace">
                    <el-table :data="keyspace" style="width: 100%">
                        <el-table-column prop="db" label="DB" width="180" />
                        <el-table-column prop="avg_ttl" label="AVG" width="180" />
                        <el-table-column prop="expires" label="Expires" />
                        <el-table-column prop="keys" label="Keys" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Memory">
                    <el-table :data="memory" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Persistence">
                    <el-table :data="persistence" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Replication">
                    <el-table :data="replication" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="Stats">
                    <el-table :data="stats" style="width: 100%">
                        <el-table-column prop="name" label="键" />
                        <el-table-column prop="value" label="值" />
                    </el-table>
                </el-tab-pane>
            </el-tabs>
        </el-tab-pane>
        <el-tab-pane label="配置" name="config">
            <el-table :data="config" style="width: 100%">
                <el-table-column prop="name" label="键" />
                <el-table-column prop="value" label="值" />
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="命令行" name="command">
            <Terminal />
        </el-tab-pane>
        <el-tab-pane label="慢日志" name="slowlog">
            <el-table :data="Slowlog" style="width: 100%">
                <el-table-column prop="ID" label="ID" width="50" />
                <el-table-column prop="ClientAddr" label="客户端" width="150" />
                <el-table-column prop="Time" label="时间" width="200">
                    <template #default="scope">{{ moment(scope.row.Time).format('YYYY-MM-DD HH:mm:ss') }}</template>
                </el-table-column>
                <el-table-column prop="Duration" label="耗时" width="100">
                    <template #default="scope">{{ scope.row.Duration / 1000 }}ms</template>
                </el-table-column>
                <el-table-column prop="Args" label="命令">
                    <template #default="scope">{{ scope.row.Args.join(" ") }}</template>
                </el-table-column>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="客户端" name="clients">
            <Clients />
        </el-tab-pane>
        <el-tab-pane label="图表" name="charts">
            <Charts />
        </el-tab-pane>
    </el-tabs>
</template>
<script lang="ts" setup>
import { onMounted, ref, } from 'vue'
import type { TabsPaneContext } from 'element-plus'
import { getConn } from '@/api/conn';
import { useConnStore } from '@/store/Conn';
import moment from "moment";

const connStore = useConnStore()
const activeName = ref('info')

let info = ref({})
let server = ref([])
let cpu = ref([])
let clients = ref([])
let cluster = ref([])
let keyspace = ref([])
let memory = ref([])
let persistence = ref([])
let replication = ref([])
let stats = ref([])
let config = ref([])
let Slowlog = ref([])

function getInfoArray(info) {
    let a = []
    for (let obj in info) {
        a.push({
            "name": obj,
            "value": info[obj]
        })
    }
    return a
}

function getKeySpace(info) {
    let a = []
    for (let obj in info) {
        a.push({
            "db": obj,
            "avg_ttl": info[obj].avg_ttl,
            "expires": info[obj].expires,
            "keys": info[obj].keys,
        })
    }
    return a
}

const handleClick = (tab: TabsPaneContext, event: Event) => {
    console.log(tab, event)
}

onMounted(() => {
    getConn(connStore.current, '').then(result => {
        info.value = result.data.Info
        //info
        server.value = getInfoArray(info.value.Server)
        cpu.value = getInfoArray(info.value.CPU)
        clients.value = getInfoArray(info.value.Clients)
        cluster.value = getInfoArray(info.value.Cluster)
        keyspace.value = getKeySpace(info.value.Keyspace)
        memory.value = getInfoArray(info.value.Memory)
        persistence.value = getInfoArray(info.value.Persistence)
        replication.value = getInfoArray(info.value.Replication)
        stats.value = getInfoArray(info.value.Stats)
        //config
        config.value = getInfoArray(result.data.Config)
        Slowlog.value = result.data.Slowlog
    })
})
</script>