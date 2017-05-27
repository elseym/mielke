import * as HtmlWebpackPlugin from "html-webpack-plugin";
import * as path from "path";
import * as merge from "webpack-merge";

export interface Path {
  output: string; // path to webpacks output directory
  projectRoot: string;
}

export const Path: Path = {
  output: path.join(__dirname, "..", "public"),
  projectRoot: path.join(__dirname, ".."),
};

export default merge({}, {
  context: Path.projectRoot,
  entry: {
    app: "./src/mielke.tsx",
  },
  module: {
    rules: [
      { loader: "awesome-typescript-loader", test: /\.tsx?$/ },
      { loader: "tslint-loader", test: /\.tsx?$/, enforce: "pre" },
    ],
  },
  output: {
    filename: "[name].js",
    path: Path.output,
  },
  plugins: [
    new HtmlWebpackPlugin({
      filename: "mielke.html",
      template: "src/mielke.html",
      inlineSource: ".(js|css)$"
    }),
  ],
  resolve: {
    extensions: [ ".ts", ".tsx", ".js", ".jsx" ],
  },
});
