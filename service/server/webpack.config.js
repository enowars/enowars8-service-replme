const path = require('path');

const config = {
  entry: {
    index: path.resolve(__dirname, 'client', 'index.ts'),
    term: path.resolve(__dirname, 'client', 'term.ts'),
  },
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
  },
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, 'static', 'js')
  },
  mode: 'development'
};
module.exports = config;
