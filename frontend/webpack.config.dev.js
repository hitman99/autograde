const path = require('path');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const webpack = require('webpack');

const config = {
    context: path.resolve(__dirname, 'src/react'),

    entry: {
        app: ['./index.js']
    },

    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: './assets/js/[name].[chunkhash].bundle.js',
        publicPath: "/"
    },

    module: {
        rules: [
            {
                test: /(\.jsx$)|(\.js$)/,
                include: /src/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            },
            // HTML
            {test: /\.html$/, use: ['html-loader']},
            // CSS
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader']
            },
            // Images
            {
                test: /\.(png|jpg|gif|svg)$/,
                loader: 'url-loader',
                options: {
                    limit: 10000,
                    name: 'assets/images/[hash].[ext]'
                }
            },
            // Fonts
            {
                test: /\.(eot|ttf|woff|woff2)$/,
                loader: 'url-loader',
                options: {
                    limit: 10000,
                    name: 'assets/fonts/[hash].[ext]'
                }
            }
        ]
    },

    plugins: [
        new CleanWebpackPlugin({
            cleanOnceBeforeBuildPatterns: ['dist']
        }),
        new HtmlWebpackPlugin({
            template: 'index.html'
        })
        // ,
        // new CopyWebpackPlugin([{ from: 'static' }])
    ],
    optimization: {
        splitChunks: {
            chunks: 'async',
            minSize: 30000,
            maxSize: 0,
            minChunks: 1,
            maxAsyncRequests: 5,
            maxInitialRequests: 3,
            automaticNameDelimiter: '~',
            automaticNameMaxLength: 30,
            name: true,
            cacheGroups: {
                vendors: {
                    test: /[\\/]node_modules[\\/]/,
                    priority: -10
                },
                default: {
                    minChunks: 2,
                    priority: -20,
                    reuseExistingChunk: true
                }
            }
        }
    },
    devtool: 'inline-source-map',
    resolve: {
        extensions: ['.js', '.jsx']
    },
    devServer: {
        historyApiFallback: true,
        proxy: {
            '/signup': {
                target: 'http://localhost:80',
                pathRewrite: {'^/signup': ''}
            },
            "/lab/scenario": {
                target: 'http://localhost:8081',
                //pathRewrite: {'^/signup': ''}
            },
            "/lab/deps": {
                target: 'http://localhost:8081',
                //pathRewrite: {'^/signup': ''}
            },
            "/state": {
                target: 'http://localhost:8081',
                //pathRewrite: {'^/signup': ''}
            },
            "/lab/scenario/state": {
                target: 'http://localhost:8081',
                //pathRewrite: {'^/signup': ''}
            }
        }
    }
};

module.exports = config;