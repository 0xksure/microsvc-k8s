import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { NodeGlobalsPolyfillPlugin } from '@esbuild-plugins/node-globals-polyfill';
import inject from '@rollup/plugin-inject';
import nodePolyfills from 'rollup-plugin-node-polyfills';
import path from 'path';


export default defineConfig({
	plugins: [sveltekit()],
	optimizeDeps: {
		include: ['@project-serum/anchor', '@solana/web3.js'],
		esbuildOptions: {
			target: 'esnext',
			plugins: [NodeGlobalsPolyfillPlugin({ buffer: true })]
		}
	},
	resolve: {
		alias: {
			$utils: path.resolve('src/utils/'),
			$static: path.resolve('static/'),
			stream: 'rollup-plugin-node-polyfills/polyfills/stream'
		}
	},
	build: {
		target: 'esnext',
		commonjsOptions: {
			transformMixedEsModules: true
		},
		rollupOptions: {
			plugins: [inject({ Buffer: ['buffer', 'Buffer'] }), nodePolyfills({ crypto: true })]
		}
	}
});
