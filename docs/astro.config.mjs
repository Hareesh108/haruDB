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
						{ label: 'Installation', slug: 'guides/installation' },
					],
				},
				{
					label: 'Core Features',
					items: [
						{ label: 'SQL Operations', slug: 'guides/sql-operations' },
						{ label: 'Indexes & Optimization', slug: 'guides/indexes' },
						{ label: 'Transactions & ACID', slug: 'guides/transactions' },
						{ label: 'Data Integrity & Recovery', slug: 'guides/data-integrity' },
					],
				},
				{
					label: 'Examples',
					items: [
						{ label: 'Real-world Examples', slug: 'guides/examples' },
					],
				},
				{
					label: 'Security & Management',
					items: [
						{ label: 'Authentication & Users', slug: 'guides/authentication' },
						{ label: 'Backup & Restore', slug: 'guides/backup-restore' },
						{ label: 'Data Storage & Encryption', slug: 'guides/data-storage' },
					],
				},
				{
					label: 'Deployment',
					items: [
						{ label: 'Docker', slug: 'guides/docker' },
						{ label: 'Connect', slug: 'guides/connect' },
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
					label: 'Support',
					items: [
						{ label: 'Troubleshooting', slug: 'guides/troubleshooting' },
						{ label: 'Roadmap', slug: 'reference/roadmap' },
					],
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
