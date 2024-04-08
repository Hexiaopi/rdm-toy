<template>
  <div class="bg-main">
    <div ref="terminal" v-loading="loading" class="terminal" element-loading-text="拼命连接中"></div>
  </div>
</template>
<script setup>
import { ref, onMounted, } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { useConnStore } from '@/store/Conn';
import 'xterm/css/xterm.css'

const connStore = useConnStore()

const terminal = ref(null)
const fitAddon = new FitAddon()

let socket = ref(null)
let loading = ref(true)
let term = ref(null)

// 初始化WS
const initWS = () => {
  if (!socket.value) {
    createWS()
  }
  if (socket.value && socket.value.readyState > 1) {
    socket.value.close()
    createWS()
  }
}

// let socket = new WebSocket(`ws://localhost:8080/ws/conn/` + connStore.current + '/cmd')

// socket.onopen = () => {
//   term.value.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
// }

// socket.onmessage = () => {
//   const data = message.data
//   term.value.write(data)
//   term.value.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
//   term.value.element && term.value.focus()
// }

// socket.onclose = () => {
//   console.log('socket close')
// }

// socket.onerror = (er) => {
//   console.log('socket err', ex)
// }


// 创建WS
const createWS = () => {
  socket.value = new WebSocket(
    `ws://localhost:8080/ws/conn/` + connStore.current + '/cmd'
  )
  socket.value.onopen = runRealTerminal //WebSocket 连接已建立
  socket.value.onmessage = onWSReceive //收到服务器消息
  socket.value.onclose = closeRealTerminal //WebSocket 连接已关闭
  socket.value.onerror = errorRealTerminal //WebSocket 连接出错
}

//WebSocket 连接已建立
const runRealTerminal = () => {
  loading.value = false
  term.value.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
  term.value.focus()
}
//WebSocket收到服务器消息
const onWSReceive = (message) => {
  const data = message.data
  term.value.write(data)
  term.value.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
  term.value.element && term.value.focus()
}
//WebSocket 连接出错
const errorRealTerminal = (ex) => {
  console.log('socket err', ex)
}
//WebSocket 连接已关闭
const closeRealTerminal = () => {
  console.log('socket close')
}

// 初始化Terminal
const initTerm = () => {
  term.value = new Terminal({
    rendererType: "canvas", //渲染类型
    rows: 30, //行数
    scrollback: 500, //终端中的回滚量
    disableStdin: false, //是否应禁用输入
    windowsMode: true, // 根据窗口换行
    cursorStyle: "underline", //光标样式
    cursorBlink: true, //光标闪烁
    theme: {
      foreground: "#ECECEC", //字体
      background: "#000000", //背景色
      cursor: "help", //设置光标
      lineHeight: 30,
    },
    bellStyle: 'sound',
  })
  term.value.open(terminal.value) //挂载dom窗口
  term.value.loadAddon(fitAddon) //自适应尺寸
  // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
  setTimeout(() => {
    fitAddon.fit()
  }, 5)
  termData() //Terminal 事件挂载
}

// 终端输入触发事件
const termData = () => {
  // 输入与粘贴的情况,onData不能重复绑定,不然会发送多次
  let curr_line = '';
  let history_lines = [];
  let last = 0
  let current = 0;
  term.value.onData((data) => {
    let code = data.charCodeAt(0);
    console.log("terminal:", code, data)
    switch (code) {
      case 127: //退格键
        if (curr_line.length > 0) {
          curr_line = curr_line.slice(0, -1);
          term.value.write('\x1b[2K\r');
          term.value.write('\x1b[32m' + connStore.current + ':' + '> \x1b[0m' + curr_line);
        } else {
          term.value.write("\x07");
        }
        break;
      case 13: //回车键
        if (curr_line === 'clear') {
          term.value.clear()
        }
        if (curr_line == "") {
          term.value.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
        } else {
          socket.value.send(curr_line)
          history_lines.push(curr_line);
          last = history_lines.length - 1;
          current = last + 1;
          term.value.write('\r\n');
        }
        curr_line = ""
        break;
      case 27: //方向键
        if (data === '\u001b[A') {//up
          if (last >= 0 && current > 0) {
            current--
            //当前行有数据的时候进行删除掉在进行渲染上存储的历史数据
            if (curr_line.length > 0) {
              term.value.write('\x1b[2K\r');
              term.value.write('\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
            }

            let text = history_lines[current]
            term.value.write(text)
            //重点，一定要记住存储当前行命令保证下次up或down时不会光标错乱覆盖终端提示符
            curr_line = text
          } else {
            term.value.write("\x07");
          }

        } else if (data === '\u001b[B') {//down
          if (last >= 0 && current < last) {
            current++
            if (curr_line.length > 0) {
              term.value.write('\x1b[2K\r');
              term.value.write('\x1b[32m' + name + ':' + db + '> \x1b[0m');
            }

            let text = history_lines[current]
            term.value.write(text)
            curr_line = text
          } else {
            term.value.write("\x07");
          }

        } else {
          term.value.write("\x07");
        }
        break;
      default:
        curr_line += data;
        term.value.write(data);
        break;
    }
  })
}

onMounted(() => {
  initWS()
  initTerm()
})

</script>
<style lang="scss" scoped>
.terminal {
  width: 100%;
  height: calc(100% - 62px);
}
</style>