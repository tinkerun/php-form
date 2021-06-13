const { build } = require('esbuild')

const isProduction = process.env.NODE_ENV === 'production'

const wasmPlugin = {
  name: 'wasm',
  setup(build) {
    let path = require('path')
    let fs = require('fs')

    // Resolve ".wasm" files to a path with a namespace
    build.onResolve({ filter: /\.wasm$/ }, args => {
      // If this is the import inside the stub module, import the
      // binary itself. Put the path in the "wasm-binary" namespace
      // to tell our binary load callback to load the binary file.
      if (args.namespace === 'wasm-stub') {
        return {
          path: args.path,
          namespace: 'wasm-binary',
        }
      }

      // Otherwise, generate the JavaScript stub module for this
      // ".wasm" file. Put it in the "wasm-stub" namespace to tell
      // our stub load callback to fill it with JavaScript.
      //
      // Resolve relative paths to absolute paths here since this
      // resolve callback is given "resolveDir", the directory to
      // resolve imports against.
      if (args.resolveDir === '') {
        return // Ignore unresolvable paths
      }
      return {
        path: path.isAbsolute(args.path) ? args.path : path.join(args.resolveDir, args.path),
        namespace: 'wasm-stub',
      }
    })

    // Virtual modules in the "wasm-stub" namespace are filled with
    // the JavaScript code for compiling the WebAssembly binary. The
    // binary itself is imported from a second virtual module.
    build.onLoad({ filter: /.*/, namespace: 'wasm-stub' }, async (args) => ({
      contents: `import wasm from ${JSON.stringify(args.path)}
          export default (imports) =>
            WebAssembly.instantiate(wasm, imports).then(
              result => result.instance)`,
    }))

    // Virtual modules in the "wasm-binary" namespace contain the
    // actual bytes of the WebAssembly file. This uses esbuild's
    // built-in "binary" loader instead of manually embedding the
    // binary data inside JavaScript code ourselves.
    build.onLoad({ filter: /.*/, namespace: 'wasm-binary' }, async (args) => ({
      contents: await fs.promises.readFile(args.path),
      loader: 'binary',
    }))
  },
}

const config = {
 entryPoints: [
   'index.ts'
 ],
 format: 'esm',
 bundle: true,
 outdir: 'lib/esm',
 plugins: [wasmPlugin],
 external: [
   'fs',
   'crypto',
   'util',
   'os'
 ],
 watch: true,
 ...isProduction && {
   minify: true,
   watch: false,
 }
}

const cjsConfig = {
  ...config,
  format: 'cjs',
  outdir: 'lib/cjs',
}

build(config).catch(() => process.exit(1))
build(cjsConfig).catch(() => process.exit(1))