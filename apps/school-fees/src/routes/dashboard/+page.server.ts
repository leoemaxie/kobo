import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { parentStudents, students } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import { kobo } from '$lib/server/kobo-client';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(302, '/login');
	}
	if (locals.user.role === 'superadmin') {
		throw redirect(302, '/admin/super');
	}
	if (locals.user.role === 'admin') {
		if (locals.user.status === 'pending') {
			throw redirect(302, '/admin/pending');
		}
		throw redirect(302, '/admin/students');
	}

	const linked = await db
		.select({ student: students })
		.from(parentStudents)
		.innerJoin(students, eq(parentStudents.studentId, students.id))
		.where(eq(parentStudents.parentId, locals.user.id));

	const studentsWithBalance = await Promise.all(linked.map(async ({ student }) => {
		try {
			const statement = await kobo.accounts.getStatement(student.koboIdentityId);
			return {
				id: student.id,
				name: student.name,
				class: student.className,
				balance: `₦ ${(statement.closing_balance_kobo / 100).toLocaleString()}`
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
