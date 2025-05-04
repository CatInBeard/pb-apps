const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin');
const ejs = require('ejs');
const fs = require('fs');

const languages = [
    { lang: 'en', filename: 'index.html' },
    { lang: 'ru', filename: 'ru.html' },
    { lang: 'fr', filename: 'fr.html' },
    { lang: 'de', filename: 'de.html' },
    { lang: 'es', filename: 'es.html' },
];

const repoData = JSON.parse(fs.readFileSync('./src/repo.json', 'utf8'));


module.exports = {
    entry: './src/index.js',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: [MiniCssExtractPlugin.loader, 'css-loader']
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'style.min.css',
        }),
        new CssMinimizerPlugin({
            minimizerOptions: {
                preset: ['default', {
                    discardComments: { removeAll: true },
                }],
            },
        }),
        ...languages.map((lang) => {
            const translations = JSON.parse(fs.readFileSync(`./src/${lang.lang}.json`, 'utf8'));
            return new HtmlWebpackPlugin({
                template: './src/index.ejs',
                filename: lang.filename,
                templateParameters: {
                    ...translations,
                    repositories: repoData.repositories
                  },
                inject: true,
                link: true
            });
        })
    ]
};
