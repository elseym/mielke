import * as CleanWebpackPlugin from "clean-webpack-plugin";
import * as HtmlWebpackInlineSourcePlugin from "html-webpack-inline-source-plugin";
import * as webpack from "webpack";
import * as merge from "webpack-merge";
import baseConfig, {Path} from "./webpack.common.config";

/**
 * This configuration relays on webpack being executed with the `-p` (production) option.
 *
 * Though this can also be configured manually if necessary, just remember to remove the `-p` flag in the build process.
 * @see https://webpack.js.org/guides/production-build/#the-manual-way
 */

export default merge({
  output: {
    chunkFilename: "[chunkhash:8].js",
    filename: "[name].[chunkhash:8].js",
  },
  plugins: [
    new CleanWebpackPlugin([Path.output], {
      root: Path.projectRoot,
      verbose: true,
    }),
    new HtmlWebpackInlineSourcePlugin(),
  ],
}, baseConfig);
