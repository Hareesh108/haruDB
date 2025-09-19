// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	integrations: [
		starlight({
			title: 'HaruDB Docs',
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/Hareesh108/haruDB' },
			],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Overview', slug: 'index' },
						{ label: 'Quick Start', slug: 'guides/quick-start' },
					],
				},
				{
					label: 'Features',
					items: [
						{ label: 'SQL Operations', slug: 'guides/sql-operations' },
						{ label: 'Indexes & Optimization', slug: 'guides/indexes' },
						{ label: 'Transactions & ACID', slug: 'guides/transactions' },
						{ label: 'Data Integrity & Recovery', slug: 'guides/data-integrity' },
					],
				},
				{
					label: 'Architecture',
					items: [
						{ label: 'WAL', slug: 'reference/wal' },
						{ label: 'Storage Engine', slug: 'reference/storage' },
						{ label: 'SQL Parser', slug: 'reference/sql-parser' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Installation', slug: 'guides/installation' },
						{ label: 'Docker', slug: 'guides/docker' },
						{ label: 'Connect', slug: 'guides/connect' },
						{ label: 'Troubleshooting', slug: 'guides/troubleshooting' },
					],
				},
				{
					label: 'Roadmap',
					items: [{ label: 'Planned Features', slug: 'reference/roadmap' }],
				},
				{
					label: 'About',
					items: [
						{ label: 'Vision', slug: 'reference/vision' },
						{ label: 'Contributing', slug: 'reference/contributing' },
						{ label: 'Author', slug: 'reference/author' },
						{ label: 'Disclaimer', slug: 'reference/disclaimer' },
					],
				},
			],
		}),
	],
});
