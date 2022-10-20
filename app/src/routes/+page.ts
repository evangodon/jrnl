import { api } from '$lib/api';
import dayjs from 'dayjs';
import type { PageLoad } from './$types';

type Response = {
	daily: {
		id: string;
		content: string;
		createdAt: Date;
		updatedAt: Date;
	};
};

export const load: PageLoad = async ({ fetch, depends }) => {
	depends('daily-entry:new');
	const date = dayjs().format('YYYY-MM-D');
	try {
		const response = await api.get(`daily/${date}`, { fetch });

		return await response.json<Response>();
	} catch (err: unknown) {
		// TODO: if error is 404
		const response = await api.get(`daily/template`, { fetch });
		const data = await response.json<{ template: string }>();
		return {
			daily: {
				id: '-1',
				content: data.template
			}
		};
	}
};

export const ssr = false;
export const csr = true;
