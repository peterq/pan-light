import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui'
import State from './util/state'
import RealTime from './realtime/realtime'

Vue.config.productionTip = false

// const $rt = new RealTime('ws://localhost:8001/demo/ws')
const $rt = new RealTime((location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host + '/demo/ws')
const $state = new Vue(State)
Vue.use(ElementUI)
Vue.mixin({
    created: function () {
        this.$state = this.$options.state || (this.$parent && this.$parent.$state) || {}
        this.$event = this.$options.event || (this.$parent && this.$parent.$event) || window.$event || {}
        this.$rt = this.$options.$rt || (this.$parent && this.$parent.$rt) || {}
    }
})
const $event = (function () {
    function e(e, t) {
        console.log('event', e, t)
        n[e] && n[e].map(function (e) {
            setTimeout(() => e(t), 0)
        })
    }

    function t(e, t) {
        n[e] || (n[e] = []), n[e].push(t)
    }
    var n = {}
    return {fire: e, on: t}
})()
window.$event = $event

new Vue({
    $event,
    $state,
    $rt,
    render: h => h(App),
}).$mount('#app')
