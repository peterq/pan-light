module.exports = {
    "root": true,
    "env": {
        "node": true
    },
    "extends": [
        "plugin:vue/essential",
        "eslint:recommended"
    ],
    "rules": {
        "no-console": 0,
        'semi': [2, 'never'],
        'semi-spacing': [2, {
            'before': false,
            'after': true
        }],
    },
    "parserOptions": {
        "parser": "babel-eslint"
    }
}
