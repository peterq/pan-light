module.exports = {
    devServer: {
        proxy: {
            '/demo': {
                target: 'http://localhost:8081',
                ws: true,
                changeOrigin: true,
            }

        }
    }
}