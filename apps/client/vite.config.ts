import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	optimizeDeps: { exclude: ['@facile/muse', '@facile/client', '@facile/mise-en-scene'] },
	server: { proxy: { '/api': 'http://127.0.0.1:4000' } },
	test: {
		include: ['src/**/*.test.ts'],
		environment: 'node'
	}
});
