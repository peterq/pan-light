import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui'
import {$event, $rt, $state} from "./app"

Vue.config.productionTip = false

Vue.use(ElementUI)
Vue.mixin({
    created: function () {
        this.$state = this.$options.state || (this.$parent && this.$parent.$state) || {}
        this.$event = this.$options.event || (this.$parent && this.$parent.$event) || window.$event || {}
        this.$rt = this.$options.$rt || (this.$parent && this.$parent.$rt) || {}
    }
})

window.$event = $event

new Vue({
    $event,
    $state,
    $rt,
    render: h => h(App),
}).$mount('#app')
