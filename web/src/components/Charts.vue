<template>
    <div>
        <el-row :gutter="2">
            <el-col :span="4">
                <el-form-item label="自动更新">
                    <el-switch @change="autoChange" v-model="autoupdate" inline-prompt active-text="是"
                        inactive-text="否" />
                </el-form-item>
            </el-col>
            <el-col :span="4">
                <el-form-item label="更新间隔">
                    <el-select @change="intervalChange" v-model="interval" placeholder="间隔">
                        <el-option label="5s" value="5" />
                        <el-option label="10s" value="10" />
                        <el-option label="30s" value="30" />
                        <el-option label="1min" value="60" />
                        <el-option label="5min" value="300" />
                    </el-select>
                </el-form-item>
            </el-col>
        </el-row>
        <el-row :gutter="2" style="height: 400px;">
            <el-col :span="8">
                <v-chart :option="clientsOption" autoresize />
            </el-col>
            <el-col :span="8">
                <v-chart :option="cmdsOption" autoresize />
            </el-col>
            <el-col :span="8">
                <v-chart :option="keysOption" autoresize />
            </el-col>
        </el-row>
        <el-row :gutter="2" style="height: 400px;">
            <el-col :span="8">
                <v-chart :option="memoryOption" autoresize />
            </el-col>
            <el-col :span="8">
                <v-chart :option="cpuOption" autoresize />
            </el-col>
            <el-col :span="8">
                <v-chart :option="networkOption" autoresize />
            </el-col>
        </el-row>
    </div>
</template>

<script setup>
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, } from 'echarts/charts';
import { TitleComponent, TooltipComponent, LegendComponent, GridComponent } from 'echarts/components';
import VChart, { THEME_KEY } from 'vue-echarts';
import { ref, provide, onUnmounted, reactive } from 'vue';
import { getConn } from '@/api/conn';
import moment from 'moment';
import { useConnStore } from '@/store/Conn';
const connStore = useConnStore()

use([
    GridComponent,
    CanvasRenderer,
    LineChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
]);

provide(THEME_KEY, 'light');

let info = ref({})
function getInfo() {
    getConn(connStore.current, '').then(result => {
        console.log(result.data)
        info.value = result.data.Info
        let now = moment().format("HH:mm:ss")
        console.log(now)
        if (xtime.length < 10) {
            xtime.push(now)
        } else {
            xtime.shift()
            xtime.push(now)
        }
        clientChartChange(result.data.Info.Clients.connected_clients)
        cmdsChartChange(result.data.Info.Stats.instantaneous_ops_per_sec)
        let count = 0;
        for (let obj in result.data.Info.Keyspace) {
            count += result.data.Info.Keyspace[obj].keys
        }
        keysChartChange(count)
        memoryChartChange(result.data.Info.Memory.used_memory)
        cpuChartChange(result.data.Info.CPU.used_cpu_sys)
        networkChartChange(result.data.Info.Stats.total_net_input_bytes, result.data.Info.Stats.total_net_output_bytes)
    })
}


let timer = null
let xtime = []
let autoupdate = ref(false)
let interval = ref(5)

const autoChange = (val) => {
    if (val == false) {
        clearInterval(timer)
        timer = null
    } else {
        timer = setInterval(() => {
            getInfo()
        }, interval.value * 1000)
    }
}

if (autoupdate.value) {
    timer = setInterval(() => {
        getInfo()
    }, interval.value * 1000)
}
const intervalChange = (val) => {
    if (autoupdate.value) {
        clearInterval(timer)
        timer = setInterval(() => {
            getInfo()
        }, val * 1000)
    }
}

const clientsOption = reactive({
    title: {
        text: '连接数',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})
//客户端连接数
let clients = reactive([])
const clientChartChange = (num) => {
    if (clients.length < 10) {
        clients.push(num)
    } else {
        clients.shift()
        clients.push(num)
    }
    clientsOption.xAxis.data = xtime
    clientsOption.series[0].data = clients
}


const cmdsOption = reactive({
    title: {
        text: '每秒执行数',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})
//每秒处理数
let cmds = reactive([])
const cmdsChartChange = (num) => {
    if (cmds.length < 10) {
        cmds.push(num)
    } else {
        cmds.shift()
        cmds.push(num)
    }
    cmdsOption.xAxis.data = xtime
    cmdsOption.series[0].data = cmds
}

const keysOption = reactive({
    title: {
        text: '键值数',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})
//键值数
let keys = reactive([])
const keysChartChange = (num) => {
    if (keys.length < 10) {
        keys.push(num)
    } else {
        keys.shift()
        keys.push(num)
    }
    keysOption.xAxis.data = xtime
    keysOption.series[0].data = keys
}

const memoryOption = reactive({
    title: {
        text: '内存',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    grid: {
        left: 60,// 调整这个属性
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value',
        axisLabel: {
            formatter: (value) => {
                return trans(value)
            }
        }
    },
    series: [
        {
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})

const trans = (size) => {
    let kb = 1024
    let unit = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    let i = Math.floor(Math.log(size) / Math.log(kb))
    if (size > 0) {
        return (size / Math.pow(kb, i)).toPrecision(3) + ' ' + unit[i];
    } else {
        return '0B'
    }
}

let memorys = reactive([])
const memoryChartChange = (num) => {
    if (memorys.length < 10) {
        memorys.push(num)
    } else {
        memorys.shift()
        memorys.push(num)
    }
    memoryOption.xAxis.data = xtime
    memoryOption.series[0].data = memorys
}


const cpuOption = reactive({
    title: {
        text: 'CPU',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value'
    },
    series: [
        {
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})
let cpus = reactive([])
const cpuChartChange = (num) => {
    if (cpus.length < 10) {
        cpus.push(num)
    } else {
        cpus.shift()
        cpus.push(num)
    }
    cpuOption.xAxis.data = xtime
    cpuOption.series[0].data = cpus
}


const networkOption = reactive({
    title: {
        text: '网络流量',
        left: 'center'
    },
    tooltip: {
        trigger: 'axis',
        axisPointer: {
            type: 'cross',
        }
    },
    grid: {
        left: 60,// 调整这个属性
    },
    xAxis: {
        type: 'category',
        data: []
    },
    yAxis: {
        type: 'value',
        axisLabel: {
            formatter: (value) => {
                return trans(value)
            }
        }
    },
    series: [
        {
            name: '入口流量',
            data: [],
            type: 'line',
            areaStyle: {}
        },
        {
            name: '出口流量',
            data: [],
            type: 'line',
            areaStyle: {}
        }
    ]
})
let networkInput = reactive([])
let networkOutput = reactive([])
const networkChartChange = (input, output) => {
    if (networkInput.length < 10) {
        networkInput.push(input)
    } else {
        networkInput.shift()
        networkInput.push(input)
    }
    if (networkOutput, length < 10) {
        networkOutput.push(output)
    } else {
        networkOutput.shift()
        networkOutput.push(output)
    }
    networkOption.xAxis.data = xtime
    networkOption.series[0].data = networkInput
    networkOption.series[1].data = networkOutput
}

onUnmounted(() => {
    clearInterval(timer)
    timer = null
})

</script>

<style scoped lang="scss"></style>