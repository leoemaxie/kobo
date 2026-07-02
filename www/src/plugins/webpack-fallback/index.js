module.exports = function (context, options) {
  return {
    name: 'webpack-fallback-plugin',
    configureWebpack(config, isServer, utils) {
      return {
        resolve: {
          fallback: {
            fs: false,
            path: false,
            os: false,
            crypto: false,
            stream: false,
            util: false,
            assert: false,
            http: false,
            https: false,
            zlib: false,
            url: false,
            buffer: false,
            timers: false,
            tls: false,
            net: false,
            child_process: false,
          },
        },
      };
    },
  };
};
