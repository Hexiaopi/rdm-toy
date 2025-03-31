<template>
  <div class="bg-main">
    <div ref="terminalRef" v-loading="loading" class="terminal" element-loading-text="拼命连接中"></div>
  </div>
</template>
<script lang="ts" setup>
import { ref, nextTick, onMounted, onBeforeUnmount, } from 'vue'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { useConnStore } from '@/store/Conn';
import 'xterm/css/xterm.css'

const props = defineProps({
  mountInit: {
    type: Boolean,
    default: true,
  }
})
const connStore = useConnStore();
const terminalRef: any = ref(null);
// 终端实例
let term: Terminal;
let socket: WebSocket;
let loading = ref(true)


function initSocket() {
  socket = new WebSocket(
    `ws://localhost:8080/ws/conn/` + connStore.current + '/cmd'
  )
  socket.onopen = () => {
    loading.value = false
    term.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
    term.focus()
  }
  socket.onerror = (e: Event) => {
    term.writeln('\r\n\x1b[31m提示: 连接错误...');
    console.log('连接错误', e);
  }
  socket.onclose = (e: CloseEvent) => {
    term.writeln('\r\n\x1b[31m提示: 连接中断...');
    console.log('socket close', e.reason);
  }
  socket.onmessage = (msg: any) => {
    term.write(msg.data);
    term.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
    //term.element && term.focus()
  }
}

// 初始化Terminal
const initTerm = () => {
  term = new Terminal({
    fontSize: 15,
    fontWeight: 'normal',
    fontFamily: 'JetBrainsMono, monaco, Consolas, Lucida Console, monospace',
    cursorBlink: true, //光标闪烁
    disableStdin: false, //是否应禁用输入
    rendererType: "canvas", //渲染类型
    rows: 30, //行数
    scrollback: 500, //终端中的回滚量
    windowsMode: true, // 根据窗口换行
    cursorStyle: "underline", //光标样式
    theme: {
      foreground: "#ECECEC", //字体
      background: "#000000", //背景色
      cursor: "help", //设置光标
    },
    bellStyle: 'sound',
  })
  term.open(terminalRef.value) //挂载dom窗口
  //注册自适应组件
  const fitAddon = new FitAddon();
  term.loadAddon(fitAddon) //自适应尺寸
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
  let history_lines: Array<string> = [];
  let last = 0
  let current = 0;
  term.onData((data) => {
    let code = data.charCodeAt(0);
    console.log("terminal:", code, data)
    switch (code) {
      case 127: //退格键
        if (curr_line.length > 0) {
          curr_line = curr_line.slice(0, -1);
          term.write('\x1b[2K\r');
          term.write('\x1b[32m' + connStore.current + ':' + '> \x1b[0m' + curr_line);
        } else {
          term.write("\x07");
        }
        break;
      case 13: //回车键
        if (curr_line === 'clear') {
          term.clear()
        }
        if (curr_line == "") {
          term.write('\r\n\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
        } else {
          socket.send(curr_line)
          history_lines.push(curr_line);
          last = history_lines.length - 1;
          current = last + 1;
          term.write('\r\n');
        }
        curr_line = ""
        break;
      case 27: //方向键
        if (data === '\u001b[A') {//up
          if (last >= 0 && current > 0) {
            current--
            //当前行有数据的时候进行删除掉在进行渲染上存储的历史数据
            if (curr_line.length > 0) {
              term.write('\x1b[2K\r');
              term.write('\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
            }
            let text = history_lines[current]
            term.write(text)
            //重点，一定要记住存储当前行命令保证下次up或down时不会光标错乱覆盖终端提示符
            curr_line = text
          } else {
            term.write("\x07");
          }
        } else if (data === '\u001b[B') {//down
          if (last >= 0 && current < last) {
            current++
            if (curr_line.length > 0) {
              term.write('\x1b[2K\r');
              term.write('\x1b[32m' + connStore.current + ':' + '> \x1b[0m');
            }
            let text = history_lines[current]
            term.write(text)
            curr_line = text
          } else {
            term.write("\x07");
          }
        } else {
          term.write("\x07");
        }
        break;
      default:
        curr_line += data;
        term.write(data);
        break;
    }
  })
}

onMounted(() => {
  if (props.mountInit) {
    init();
  }
})

function init() {
  if (term) {
    console.log('重新连接...');
    close();
  }
  nextTick(() => {
    initTerm();
    initSocket();
  })
}

onBeforeUnmount(() => {
  close();
})

function close() {
  console.log("terminal close");
  closeSocket();
  if (term) {
    term.dispose();
  }
}

function closeSocket() {
  socket && socket.readyState === 1 && socket.close();
}

</script>
<style lang="scss" scoped>
.terminal {
  width: 100%;
  height: calc(100% - 62px);
}
</style>