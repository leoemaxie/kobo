import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { db } from '$lib/server/db';
import { students, parents, parentStudents } from '$lib/server/db/schema';
import { eq, and, or, ilike } from 'drizzle-orm';

export const GET: RequestHandler = async ({ url, locals }) => {
    if (!locals.user) {
        return json({ error: 'Unauthorized' }, { status: 401 });
    }

    const q = url.searchParams.get('q') || '';
    if (!q || q.length < 2) {
        return json({ students: [], parents: [] });
    }

    const searchPattern = `%${q}%`;

    try {
        if (locals.user.role === 'admin' || locals.user.role === 'superadmin') {
            // Admins can search all students and parents
            const matchedStudents = await db.select().from(students)
                .where(
                    or(
                        ilike(students.name, searchPattern),
                        ilike(students.id, searchPattern)
                    )
                )
                .limit(10);
            
            const matchedParents = await db.select().from(parents)
                .where(
                    or(
                        ilike(parents.name, searchPattern),
                        ilike(parents.email, searchPattern)
                    )
                )
                .limit(10);
            
            return json({
                students: matchedStudents,
                parents: matchedParents.map(p => ({
                    id: p.id,
                    name: p.name,
                    email: p.email,
                    role: p.role
                }))
            });
        } else if (locals.user.role === 'parent') {
            // Parents can only search their linked students
            // Since they are a parent, we find their linked students
            const parentLinks = await db.select().from(parentStudents)
                .where(eq(parentStudents.parentId, locals.user.id));
            
            if (parentLinks.length === 0) {
                return json({ students: [], parents: [] });
            }

            const linkedStudentIds = parentLinks.map(link => link.studentId);
            
            const matchedStudents = await db.select().from(students)
                .where(
                    and(
                        or(
                            ilike(students.name, searchPattern),
                            ilike(students.id, searchPattern)
                        ),
                        // Manually filtering below if "inArray" is tricky, but Drizzle has "inArray"
                        // Or we can just fetch all linked students and filter in memory if few, but let's just filter in memory for linked students as there are usually < 10
                    )
                );

            const filteredLinkedStudents = matchedStudents.filter(s => linkedStudentIds.includes(s.id));
            
            return json({
                students: filteredLinkedStudents,
                parents: [] // Parents don't search other parents
            });
        }

        return json({ students: [], parents: [] });
    } catch (e: any) {
        return json({ error: 'Search failed' }, { status: 500 });
    }
};
