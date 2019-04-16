import RFB from '../lib/vnc/core/rfb'

const div = document.createElement('div')
div.style.width = '100%'
div.style.height = '100%'

const rfb = new RFB(div, 'ws://loacalhost', {})

rfb.sendCredentials()

export default {
    div,
    rfb
}
