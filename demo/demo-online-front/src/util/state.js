function data() {
    return {

    }
}

export default {
    data,
    watch: (function () {
        const map = {}

        function makeWatcher(key) {
            return function (val, old) {
                this.$event.fire('state.' + key, val)
                if (typeof val === 'string')
                    this.$event.fire('state.' + key + '.' + val, old)
            }
        }

        function addWatch(data, prefix = '') {
            for (let key in data) {
                if (data[key].__proto__ === {}.__proto__)
                    addWatch(data[key], prefix + key + '.')
                map[prefix + key] = makeWatcher(prefix + key)
            }
        }

        addWatch(data())
        return Object.assign({}, map, {
            username(val) {
                if (val) this.$event.fire('login-success', val)
            }
        })
    })(),
    created() {
    }
}
