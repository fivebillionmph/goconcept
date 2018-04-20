var webpack = require("webpack");
var path = require("path");

var base = {
	module: {
		loaders: [
			{
				test: /\.(js|jsx)$/,
				exclude: /node_modules/,
				loader: 'babel-loader',
				query: {
					presets: ["es2015"]
				}
			},
			{
				test: /\.less$/,
				loaders: ["style-loader", "css-loader", "less-loader"]
			},
			{
				test: /\.css$/,
				loaders: ["style-loader", "css-loader"]
			}
		]
	}
}

var main = Object.assign({}, base, {
	entry: [
		path.join(__dirname, "/index.js")
	],
	output: {
		path: __dirname + "/public",
		filename: "index.js",
		libraryTarget: "umd"
	}
});

module.exports = [main];
