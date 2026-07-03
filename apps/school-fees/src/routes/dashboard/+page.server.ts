import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parentStudents, students } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { koboFetch } from '$lib/server/kobo-client';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, '/login');
	}
	if (locals.user.isAdmin) {
		throw redirect(302, '/admin/students');
	}

	const linked = await db
		.select({ student: students })
		.from(parentStudents)
		.innerJoin(students, eq(parentStudents.studentId, students.id))
		.where(eq(parentStudents.parentId, locals.user.id));

	const studentsWithBalance = await Promise.all(linked.map(async ({ student }) => {
		try {
			const statement = await koboFetch(`/accounts/${student.koboIdentityId}/statement`);
			return {
				id: student.id,
				name: student.name,
				class: student.className,
				balance: `₦ ${(statement.balance / 100).toLocaleString()}`
			};
		} catch (e) {
			return {
				id: student.id,
				name: student.name,
				class: student.className,
				balance: 'Error fetching balance'
			};
		}
	}));

	return {
		students: studentsWithBalance,
		user: locals.user
	};
};
