const esbuild = require('esbuild')

const isWatchMode = process.argv.includes('--watch')
const buildOptions = {
  entryPoints: ['src/index.tsx'],
  bundle: true,
  outfile: '../../static/webcomponents/floweditor.bundle.js',
  minify: true,
  sourcemap: false,
  loader: {
    '.tsx': 'tsx',
    '.css': 'css',
  },
  target: ['es6'],
}

const main = async () => {
  if (isWatchMode) {
    const ctx = await esbuild.context(buildOptions)
    await ctx.watch()
  }
  else {
    try {
      await esbuild.build(buildOptions)
    }
    catch (e) {
      console.error(e)
      process.exit(1)
    }
  }
}

main()
