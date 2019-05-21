import Vue from 'vue'
import {$event, $rt, $state} from "../app"

// 对话框管理器
const con = document.createElement('div')
document.body.appendChild(con)

// 复用组件实例
const instances = []
const comps = []

// 对话框组件混入
export const dialogMixin = {
    $event, $rt, $state,
    data: {
        visible: false,
        promise: null,
        param: {}
    },
    methods: {
        open(param) {
            this.visible = true
            this.param = param
            setTimeout(this.onOpen.bind(this, param))
            return new Promise((resolve, reject) => {
                this.promise = {resolve, reject}
            })
        },
        resolve(data) {
            this.visible = false
            this.$options.promise.resolve(data)
        },
        reject(reason) {
            this.visible = false
            this.$options.promise.reject(reason)
        },
        getPromise() {
            return this.$options.promise
        },
        onOpen() {
            throw new Error('请实现onOpen方法')
        }
    }
}

// 打开一个对话框组件并返回这个实例
export function openDialog(comp, param, recreate = true) {
    var ins
    // 构造Promise, 对话框操作完成销毁组件并移出dom
    const promise = new Promise((resolve, reject) => setTimeout(() => Object.assign(promise, {resolve, reject})))
        .finally(() => {
            if (recreate) {
                setTimeout(() => {
                    con.removeChild(ins.$el)
                    ins.$destroy()
                }, 1e3)
            }
        })
    if (!recreate && comps.indexOf(comp) > -1) {
        const ins = instances[comps.indexOf(comp)]
        ins.$options.promise = promise
        ins.open(param)
        return ins
    }
    // 强制混入
    comp.mixins = comp.mixins || []
    if (comp.mixins.indexOf(dialogMixin) < 0) {
        comp.mixins.push(dialogMixin)
    }
    // 构造实例插入容器
    ins = new Vue({promise, ...comp}).$mount()
    con.appendChild(ins.$el)
    ins.open(param)
    if (!recreate) {
        comps.push(comp)
        instances.push(ins)
    }
    return ins
}

