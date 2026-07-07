declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: {
				id: string;
				email: string;
				name: string;
				role: string;
			} | null;
			session: any | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
