const path = require('path');

const config = {
  entry: path.resolve(__dirname, 'client', 'session.ts'),
  devtool: 'inline-source-map',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      },
      {
        test: /\.js$/,
        use: ["source-map-loader"],
        enforce: "pre",
        exclude: /node_modules/
      }
    ]
  },
  resolve: {
    modules: [
      'node_modules',
    ],
    extensions: [ '.tsx', '.ts', '.js' ],
    alias: {
      common: path.resolve('./out/common'),
      browser: path.resolve('./out/browser')
    },
  },
  output: {
    filename: 'session.js',
    path: path.resolve(__dirname, 'static')
  },
  mode: 'development'
};
module.exports = config;
