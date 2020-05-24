import { Message, MessageBox } from 'element-ui'
import { parseTime } from '@/utils'

const state = {
  ws: null,
  message: {},
  againConnectTime: 0 // 规定时间重连
}

const mutations = {
  SET_MESSAGE: (state, message) => {
    state.message = message
  },
  SET_WS: (state, ws) => {
    state.ws = ws
  },
  CLOSE_WS: (state) => {
    if (state.ws) {
      state.ws.close()
    }
  },
  SET_AGAINCONNECTTIME: (state, time) => {
    state.againConnectTime = time
  }
}

const actions = {
  init({ dispatch, commit, state }) {
    return new Promise((resolve, reject) => {
      const websocket = new WebSocket('ws://' + window.location.host + process.env.VUE_APP_BASE_API + '/ws/connect')
      websocket.onopen = () => {
        console.log('websocket连接成功, 时间：' + parseTime(new Date()))
        // 连接成功，当成重连次数0 置0
        commit('SET_AGAINCONNECTTIME', 0)
        resolve(this.webSocket)
      }

      websocket.onerror = (err) => {
        console.log('websocket连接发生错误, 时间：' + parseTime(new Date()))
        reject(err)
      }

      websocket.onmessage = (event) => {
        const responseData = JSON.parse(event.data)
        console.log(responseData)
        commit('SET_MESSAGE', responseData)
      }
      websocket.onclose = (e) => {
      // 1005 主动断开
      // websocket close code https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent
        // 顶号，后台发送关闭帧
        commit('SET_WS', null)

        if (e.code !== 1005) {
          if (state.againConnectTime === 0) {
            // 第一次连接断开进来 有一次重连机会
            commit('SET_AGAINCONNECTTIME', new Date().getTime())
            setTimeout(() => {
              console.log('首次断开，重新主动连接, 时间：' + parseTime(new Date()))
              dispatch('init')
            }, 60000)
          } else {
            if (new Date().getTime() - state.againConnectTime >= 60000) {
              console.log('主动连接失败，再次尝试, 时间：' + parseTime(new Date()))
              // 一分钟后的连接 一次重连机会已用完 还是连接失败,就弹窗询问用户
              MessageBox.confirm('检测到与服务器断开连接, 请重连！', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
              }).then(() => {
                dispatch('init')
              }).catch(() => {
                Message({
                  type: 'info',
                  message: '取消重连'
                })
              })
            }
          }
        }

        console.log('connection closed (' + e.code + ') ' + e.reason + ', 时间：' + parseTime(new Date()))
      }
      commit('SET_WS', websocket)
    })
  },

  close({ commit }) {
    return new Promise(resolve => {
      commit('CLOSE_WS')
      commit('SET_WS', null)
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

